package l1

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ten-protocol/go-ten/contracts/generated/DataAvailabilityRegistry"
	"github.com/ten-protocol/go-ten/go/common/stopcontrol"

	"github.com/ten-protocol/go-ten/go/ethadapter/contractlib"

	"github.com/ten-protocol/go-ten/go/host/storage"

	"github.com/ten-protocol/go-ten/go/common/gethutil"

	"github.com/ten-protocol/go-ten/go/enclave/crosschain"

	"github.com/ten-protocol/go-ten/go/common/subscription"

	"github.com/ten-protocol/go-ten/go/common/host"

	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/retry"
	"github.com/ten-protocol/go-ten/go/enclave/events"
	"github.com/ten-protocol/go-ten/go/ethadapter"
)

const (
	// minimum gap in block height to trigger eager fetching of historical blocks - 5 blocks should not happen in normal operation
	_minGapForEagerFetch = 5
	eagerFetchHeightKey  = "eager_fetch_height"
)

var (
	l1BlockTime      = 12 * time.Second
	_timeoutNoBlocks = 2 * l1BlockTime // after this timeout we assume the subscription to the L1 node is not working
	one              = big.NewInt(1)
	_headBlock       = big.NewInt(-1) // special value to indicate fetching up to head block (rather than to a target)
	ErrNoNextBlock   = errors.New("no next block")
)

// DataService is a host service for subscribing to new blocks and looking up L1 data
type DataService struct {
	blockSubscribers *subscription.Manager[host.L1BlockHandler]
	// this eth client should only be used by the repository, the repository may "reconnect" it at any time and don't want to interfere with other processes
	ethClient        ethadapter.EthClient
	contractRegistry contractlib.ContractRegistryLib
	blobResolver     BlobResolver
	blockResolver    storage.BlockResolver
	l1StartHash      gethcommon.Hash

	eagerFetchHeight atomic.Uint64 // the next block height to eagerly fetch historical blocks from
	head             gethcommon.Hash

	running       atomic.Bool
	hostInterrupt *stopcontrol.StopControl
	logger        gethlog.Logger
}

func NewL1DataService(ethClient ethadapter.EthClient, blockResolver storage.BlockResolver, contractRegistry contractlib.ContractRegistryLib, blobResolver BlobResolver, l1StartHash gethcommon.Hash, hostInterrupt *stopcontrol.StopControl, logger gethlog.Logger) *DataService {
	return &DataService{
		blockSubscribers: subscription.NewManager[host.L1BlockHandler](),
		ethClient:        ethClient,
		contractRegistry: contractRegistry,
		blobResolver:     blobResolver,
		blockResolver:    blockResolver,
		l1StartHash:      l1StartHash,
		running:          atomic.Bool{},
		hostInterrupt:    hostInterrupt,
		logger:           logger,
	}
}

func (r *DataService) Start() error {
	r.running.Store(true)
	r.hostInterrupt.OnStop(func() {
		r.logger.Info("Stopping L1 DataService as host is stopping")
		r.running.Store(false)
	})

	r.initEagerFetchHeight()

	// Repository constantly streams new blocks and forwards them to subscribers
	go r.streamLiveBlocks()

	return nil
}

func (r *DataService) Stop() error {
	r.running.Store(false)
	return nil
}

func (r *DataService) initEagerFetchHeight() {
	if r.blockResolver == nil {
		r.logger.Error("blockResolver not set, cannot load eager fetch height")
		return
	}

	height, err := r.blockResolver.GetMetadata(eagerFetchHeightKey)
	if err != nil {
		if r.l1StartHash == (gethcommon.Hash{}) {

			// note: this is expected in some test scenarios with dummy L1s
			r.logger.Warn("l1StartHash not set, blocks will be fetched from L1 genesis")

			r.updateEagerFetchHeight(1)
			return
		}

		startBlock, err := r.ethClient.HeaderByHash(r.l1StartHash)
		if err != nil {
			r.logger.Crit("Could not get l1StartHash block for eager fetch height initialization", log.ErrKey, err, "l1StartHash", r.l1StartHash)
			return
		}
		height = startBlock.Number.Uint64()
		r.updateEagerFetchHeight(height)
		r.logger.Info("Initialized eager fetch height from l1StartHash", "height", height)
		return
	}
	r.eagerFetchHeight.Store(height)
	r.logger.Info("Loaded eager fetch height from database", "height", height)

}

func (r *DataService) updateEagerFetchHeight(blockHeight uint64) {
	current := r.eagerFetchHeight.Load()
	if blockHeight > current {
		r.eagerFetchHeight.Store(blockHeight)

		err := r.blockResolver.SetMetadata(eagerFetchHeightKey, blockHeight)
		if err != nil {
			// this is not a disaster if it succeeds on future updates, we just use the in-mem value
			r.logger.Error("Failed to persist eager fetch height", log.ErrKey, err)
		}
	} else {
		r.logger.Debug("Skipping eager fetch height update for older block", "current", current, "blockHeight", blockHeight)
	}
}

func (r *DataService) HealthStatus(context.Context) host.HealthStatus {
	// todo (@matt) do proper health status based on last received block or something
	errMsg := ""
	if !r.running.Load() {
		errMsg = "not running"
	}
	return &host.BasicErrHealthStatus{ErrMsg: errMsg}
}

// Subscribe will register a new block handler to receive new blocks as they arrive, returns unsubscribe func
func (r *DataService) Subscribe(handler host.L1BlockHandler) func() {
	return r.blockSubscribers.Subscribe(handler)
}

// FetchNextBlock calculates the next canonical block that should be sent to requester after a given hash.
// It returns the block and a bool for whether it is the latest known head
func (r *DataService) FetchNextBlock(remoteHead gethcommon.Hash) (*types.Header, bool, error) {
	if remoteHead == r.head {
		// remoteHead is the latest known head
		return nil, false, ErrNoNextBlock
	}

	if remoteHead == gethutil.EmptyHash {
		// remoteHead is empty, so we are starting from genesis
		blk, err := r.ethClient.HeaderByNumber(big.NewInt(0))
		if err != nil {
			return nil, false, fmt.Errorf("could not find genesis block - %w", err)
		}
		return blk, false, nil
	}

	// the latestCanonAncestor will usually return the remoteHead itself but this step is necessary to walk back if there was a fork
	fork, err := r.latestCanonAncestor(remoteHead)
	if err != nil {
		return nil, false, err
	}

	// and send the canonical block at the height after that
	// (which may be a fork, or it may just be the next on the same branch if we are catching-up)
	blk, err := r.ethClient.HeaderByNumber(increment(fork.CommonAncestor.Number))
	if err != nil {
		if errors.Is(err, ethereum.NotFound) {
			return nil, false, ErrNoNextBlock
		}
		return nil, false, fmt.Errorf("could not find block after latest canon ancestor, height=%s - %w", increment(fork.CommonAncestor.Number), err)
	}

	return blk, blk.Hash() == r.head, nil
}

// FetchBlock - BlockResolver interface
func (r *DataService) FetchBlock(_ context.Context, blockHash common.L1BlockHash) (*types.Header, error) {
	h, err := r.blockResolver.ReadBlock(&blockHash)
	if err != nil {
		return r.ethClient.HeaderByHash(blockHash)
	}
	return h, nil
}

func (r *DataService) latestCanonAncestor(remote gethcommon.Hash) (*common.ChainFork, error) {
	ctx := context.Background()

	remoteHead, err := r.FetchBlock(ctx, remote)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch L1 block with hash=%s - %w", remote, err)
	}

	searchBackFrom, err := r.FetchBlockByHeight(remoteHead.Number)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch L1 block by height=%s - %w", remoteHead.Number, err)
	}

	fork, err := gethutil.LCA(ctx, searchBackFrom, remoteHead, r)
	if err != nil {
		return nil, fmt.Errorf("unable to calculate LCA - %w", err)
	}
	return fork, nil
}

// GetTenRelevantTransactions processes logs in their natural order without grouping by transaction hash.
func (r *DataService) GetTenRelevantTransactions(block *types.Header) (*common.ProcessedL1Data, error) {
	processed := &common.ProcessedL1Data{
		BlockHeader: block,
		Events:      []common.L1Event{},
	}
	networkConfigAddress := r.contractRegistry.NetworkConfigLib().GetContractAddr()
	allAddresses := r.contractRegistry.GetContractAddresses()

	if err := r.processNetworkConfigLogs(block, *networkConfigAddress, processed); err != nil {
		return nil, err
	}
	if err := r.processNetworkUpgradeLogs(block, *networkConfigAddress, processed); err != nil {
		return nil, err
	}
	if err := r.processMessageBusLogs(block, allAddresses.L1MessageBus, processed); err != nil {
		return nil, err
	}
	if err := r.processEnclaveRegistryLogs(block, allAddresses.EnclaveRegistry, processed); err != nil {
		return nil, err
	}
	if err := r.processRollupLogs(block, allAddresses.DataAvailabilityRegistry, processed); err != nil {
		return nil, err
	}

	return processed, nil
}

func (r *DataService) getContractLogs(block *types.Header, contractAddr gethcommon.Address) ([]types.Log, error) {
	blkHash := block.Hash()
	logs, err := r.ethClient.GetLogs(ethereum.FilterQuery{
		BlockHash: &blkHash,
		Addresses: []gethcommon.Address{contractAddr},
	})
	if err != nil {
		return nil, fmt.Errorf("unable to fetch logs for contract %s: %w", contractAddr.Hex(), err)
	}
	return logs, nil
}

func (r *DataService) processMessageBusLogs(block *types.Header, contractAddr gethcommon.Address, processed *common.ProcessedL1Data) error {
	logs, err := r.getContractLogs(block, contractAddr)
	if err != nil {
		return err
	}

	for _, l := range logs {
		if len(l.Topics) == 0 {
			r.logger.Error("Log has no topics. Should not happen", "txHash", l.TxHash)
			return errors.New("log has no topics")
		}
		txData, err := r.fetchTxAndReceipt(l.TxHash)
		if err != nil {
			r.logger.Error("Error creating transaction data. Should not happen", "txHash", l.TxHash, "error", err)
			return fmt.Errorf("error creating transaction data. Should not happen: %w", err)
		}
		switch l.Topics[0] {
		case ethadapter.CrossChainEventID:
			err = r.processCrossChainLogs(l, txData, processed)
		}
		if err != nil {
			r.logger.Error("Error processing log", "txHash", l.TxHash, "error", err)
			return fmt.Errorf("error processing log: %w", err)
		}
	}
	return nil
}

func (r *DataService) processNetworkUpgradeLogs(block *types.Header, contractAddr gethcommon.Address, processed *common.ProcessedL1Data) error {
	logs, err := r.getContractLogs(block, contractAddr)
	if err != nil {
		return err
	}

	// Filter logs for network upgrade events
	var upgradeLogs []types.Log
	for _, l := range logs {
		if len(l.Topics) > 0 && l.Topics[0] == ethadapter.UpgradedEventID {
			upgradeLogs = append(upgradeLogs, l)
		}
	}

	// Process each upgrade log
	for _, l := range upgradeLogs {
		txData, err := r.fetchTxAndReceipt(l.TxHash)
		if err != nil {
			r.logger.Error("Error creating transaction data", "txHash", l.TxHash, "error", err)
			continue
		}

		// Convert the log to network upgrade data
		networkUpgrades, err := events.ConvertLogsToNetworkUpgrades([]types.Log{l}, ethadapter.UpgradedEventName, ethadapter.NetworkConfigABI)
		if err != nil {
			r.logger.Error("Error converting logs to network upgrades", "txHash", l.TxHash, "error", err)
			return fmt.Errorf("error converting logs to network upgrades: %w", err)
		}

		txData.NetworkUpgrades = networkUpgrades
		processed.AddEvent(common.NetworkUpgradedTx, txData)
	}
	return nil
}

func (r *DataService) processEnclaveRegistryLogs(block *types.Header, contractAddr gethcommon.Address, processed *common.ProcessedL1Data) error {
	logs, err := r.getContractLogs(block, contractAddr)
	if err != nil {
		return err
	}

	for _, l := range logs {
		if len(l.Topics) == 0 {
			continue
		}
		txData, err := r.fetchTxAndReceipt(l.TxHash)
		if err != nil {
			r.logger.Error("Error creating transaction data", "txHash", l.TxHash, "error", err)
			continue
		}
		switch l.Topics[0] {
		case ethadapter.NetworkSecretInitializedEventID:
			if err = r.processEnclaveRegistrationTx(txData, processed); err == nil {
				err = r.processSequencerLogs(l, txData, processed, common.SequencerAddedTx)
			}
		case ethadapter.SequencerEnclaveGrantedEventID:
			err = r.processSequencerLogs(l, txData, processed, common.SequencerAddedTx)
		case ethadapter.SequencerEnclaveRevokedEventID:
			err = r.processSequencerLogs(l, txData, processed, common.SequencerRevokedTx)
		case ethadapter.NetworkSecretRequestedID:
			processed.AddEvent(common.SecretRequestTx, txData)
		case ethadapter.NetworkSecretRespondedID:
			processed.AddEvent(common.SecretResponseTx, txData)
		default:
			// there are known events that we don't care about here
			r.logger.Trace("Unknown log topic", "topic", l.Topics[0], "txHash", l.TxHash)
		}

		if err != nil {
			r.logger.Error("Error processing log", "txHash", l.TxHash, "error", err)
			return fmt.Errorf("error processing log: %w", err)
		}
	}
	return nil
}

func (r *DataService) processNetworkConfigLogs(block *types.Header, contractAddr gethcommon.Address, processed *common.ProcessedL1Data) error {
	logs, err := r.getContractLogs(block, contractAddr)
	if err != nil {
		return err
	}
	for _, l := range logs {
		if len(l.Topics) == 0 {
			continue
		}
		txData, err := r.fetchTxAndReceipt(l.TxHash)
		if err != nil {
			r.logger.Error("Error creating transaction data", "txHash", l.TxHash, "error", err)
			continue
		}
		switch l.Topics[0] {
		case ethadapter.NetworkContractAddressAddedID:
			processed.AddEvent(common.NetworkContractAddressAddedTx, txData)
		case ethadapter.AdditionalContractAddressAddedID:
			processed.AddEvent(common.AdditionalContractAddressAddedTx, txData)
		default:
			// there are known events that we don't care about here
			r.logger.Trace("Unknown log topic", "topic", l.Topics[0], "txHash", l.TxHash)
		}
	}
	return nil
}

func (r *DataService) processRollupLogs(block *types.Header, contractAddr gethcommon.Address, processed *common.ProcessedL1Data) error {
	rollupLogs, err := r.getContractLogs(block, contractAddr)
	if err != nil {
		return err
	}

	for _, l := range rollupLogs {
		if len(l.Topics) == 0 {
			continue
		}
		txData, err := r.fetchTxAndReceipt(l.TxHash)
		if err != nil {
			r.logger.Error("Error creating transaction data", "txHash", l.TxHash, "error", err)
			continue
		}
		switch l.Topics[0] {
		case ethadapter.RollupAddedID:
			err = r.processRollupLog(l, txData, processed)
		}
		if err != nil {
			r.logger.Error("Error processing log", "txHash", l.TxHash, "error", err)
			return fmt.Errorf("error processing log: %w", err)
		}
	}
	return nil
}

func (r *DataService) processRollupLog(l types.Log, txData *common.L1TxData, processed *common.ProcessedL1Data) error {
	abi, err := DataAvailabilityRegistry.DataAvailabilityRegistryMetaData.GetAbi()
	if err != nil {
		r.logger.Error("Error getting DataAvailabilityRegistry ABI", log.ErrKey, err)
		return err
	}
	var event DataAvailabilityRegistry.DataAvailabilityRegistryRollupAdded
	err = abi.UnpackIntoInterface(&event, "RollupAdded", l.Data)
	if err != nil {
		r.logger.Error("Error unpacking RollupAdded event", log.ErrKey, err)
		return err
	}
	blobs, err := r.blobResolver.FetchBlobs(context.Background(), processed.BlockHeader, []gethcommon.Hash{event.RollupHash})
	if err != nil {
		r.logger.Error(fmt.Sprintf("error while fetching blobs. Cause: %s", err))
		return err
	}
	txData.BlobsWithSignature = []common.BlobAndSignature{
		{
			Blob:      blobs[0],
			Signature: event.Signature,
		},
	}
	processed.AddEvent(common.RollupTx, txData)
	return nil
}

// fetchTxAndReceipt creates a new L1TxData instance for a transaction
func (r *DataService) fetchTxAndReceipt(txHash gethcommon.Hash) (*common.L1TxData, error) {
	tx, _, err := r.ethClient.TransactionByHash(txHash)
	if err != nil {
		return nil, fmt.Errorf("error fetching transaction: %w", err)
	}

	receipt, err := r.ethClient.TransactionReceipt(txHash)
	if err != nil {
		return nil, fmt.Errorf("error fetching receipt: %w", err)
	}

	return &common.L1TxData{
		Transaction:        tx,
		Receipt:            receipt,
		CrossChainMessages: common.CrossChainMessages{},
		ValueTransfers:     common.ValueTransferEvents{},
	}, nil
}

// processCrossChainLogs handles cross-chain message logs
func (r *DataService) processCrossChainLogs(l types.Log, txData *common.L1TxData, processed *common.ProcessedL1Data) error {
	if !r.ethClient.SupportsEventLogs() {
		return nil
	}
	messages, err := crosschain.ConvertLogsToMessages([]types.Log{l}, ethadapter.CrossChainEventName, ethadapter.MessageBusABI)
	if err != nil {
		return err
	}

	txData.CrossChainMessages = messages
	processed.AddEvent(common.CrossChainMessageTx, txData)
	return nil
}

// processSequencerLogs handles sequencer logs
func (r *DataService) processSequencerLogs(l types.Log, txData *common.L1TxData, processed *common.ProcessedL1Data, txType common.L1TenEventType) error {
	if !r.ethClient.SupportsEventLogs() {
		// todo - this is a hack that must be addressed for our sanity
		enclaveID, _ := getEnclaveIdFromLog(l)
		txData.SequencerEnclaveID = enclaveID
		processed.AddEvent(txType, txData)
		return nil
	}

	enclaveID, err := getEnclaveIdFromLog(l)
	if err != nil {
		return err
	}
	txData.SequencerEnclaveID = enclaveID
	processed.AddEvent(txType, txData)
	return nil
}

// processEnclaveRegistrationTx handles decoded transaction types
func (r *DataService) processEnclaveRegistrationTx(txData *common.L1TxData, processed *common.ProcessedL1Data) error {
	networkLib := r.contractRegistry.EnclaveRegistryLib()
	decodedTx, err := networkLib.DecodeTx(txData.Transaction)
	if err != nil {
		r.logger.Error("Error decoding transaction", "txHash", txData.Transaction.Hash, "error", err)
		return err
	}
	if decodedTx != nil {
		switch decodedTx.(type) {
		case *common.L1InitializeSecretTx:
			processed.AddEvent(common.InitialiseSecretTx, txData)
		case *common.L1PermissionSeqTx: // FIXME I think this can be deleted?
			return nil // no-op as it was processed in the previous processSequencerLogs call
		default:
			// this should never happen since the specific events should always decode into one of these types
			r.logger.Error("Unknown tx type", "txHash", txData.Transaction.Hash().Hex())
		}
	}
	return nil
}

// stream blocks from L1 as they arrive and forward them to subscribers, no guarantee of perfect ordering or that there won't be gaps.
// If streaming is interrupted it will carry on from latest, it won't try to replay missed blocks.
// Initially, and when it detects a gap, it will eagerly fetch historical blocks to catch up.
func (r *DataService) streamLiveBlocks() {
	// initial catch-up of all historical blocks before starting live streaming
	r.eagerFetchBlocksUpTo(_headBlock)

	liveStream, streamSub := r.resetLiveStream()
	for r.running.Load() {
		select {
		case blockHeader := <-liveStream:
			r.logger.Info("received block from l1 stream", log.BlockHashKey, blockHeader.Hash(), log.BlockHeightKey, blockHeader.Number)
			if blockHeader.Number.Uint64() > r.eagerFetchHeight.Load()+_minGapForEagerFetch {
				// catch up historical blocks if we are behind
				r.eagerFetchBlocksUpTo(blockHeader.Number)
			}

			// check that the block's ancestors are available in the host db
			err := r.fetchAncestors(blockHeader.ParentHash)
			if err != nil {
				r.logger.Error("error fetching ancestors", log.ErrKey, err)
				continue
			}

			r.storeBlockHeader(blockHeader)
			r.updateEagerFetchHeight(blockHeader.Number.Uint64())

			// only notify subscribers (enclave guardians) on the streamed block, because they have their own catch-up mechanism
			r.head = blockHeader.Hash()
			for _, handler := range r.blockSubscribers.Subscribers() {
				go handler.HandleBlock(blockHeader)
			}

		case <-time.After(_timeoutNoBlocks):
			r.logger.Warn("no new blocks received since timeout. Reconnecting..", "timeout", _timeoutNoBlocks)
			if streamSub != nil {
				streamSub.Unsubscribe()
			}
			if liveStream != nil {
				close(liveStream)
			}
			liveStream, streamSub = r.resetLiveStream()

		case <-r.hostInterrupt.Done():
			r.logger.Info("block streaming stopped by stop signal")
			if streamSub != nil {
				streamSub.Unsubscribe()
			}
			if liveStream != nil {
				close(liveStream)
			}
			return
		}
	}
}

func (r *DataService) fetchAncestors(blockHash gethcommon.Hash) error {
	var blocksToStore []*types.Header // accumulate ancestor blocks to store in order if they are missing
	currentHash := blockHash

	// walk back ancestors until we find one in the host db (or reach l1StartHash)
	for currentHash != r.l1StartHash {
		_, blockErr := r.blockResolver.ReadBlock(&currentHash)
		if blockErr == nil {
			break
		}

		r.logger.Info("`newHeads` block ancestor not found in host db. Fetching from l1", log.BlockHashKey, currentHash)

		b, err := r.ethClient.HeaderByHash(currentHash)
		if err != nil {
			return fmt.Errorf("error fetching block: %w", err)
		}

		blocksToStore = append(blocksToStore, b)
		currentHash = b.ParentHash
	}

	for i := len(blocksToStore) - 1; i >= 0; i-- {
		r.storeBlockHeader(blocksToStore[i])
		r.updateEagerFetchHeight(blocksToStore[i].Number.Uint64())
	}

	return nil
}

// eagerFetchBlocksUpTo will fetch historical blocks from L1 starting from eagerFetchHeight until there are no more blocks
// pass _headBlock to fetch up to the latest head,
// otherwise this will fetch up to `upToHeight - 1` (to allow for live block to be processed and not re-fetched)
func (r *DataService) eagerFetchBlocksUpTo(upToHeight *big.Int) {
	height := r.eagerFetchHeight.Load()
	r.logger.Info("Starting eager L1 block fetch to catch up", "startingHeight", height)

	for r.running.Load() && (upToHeight.Cmp(_headBlock) == 0 || big.NewInt(int64(height)).Cmp(upToHeight) < 0) {
		block, err := r.ethClient.HeaderByNumber(big.NewInt(int64(height)))
		if err != nil {
			if errors.Is(err, ethereum.NotFound) {
				// no more blocks to fetch
				r.logger.Info("Eager fetch caught up to latest L1 block", "height", height)
				return
			}
			r.logger.Warn("Eager fetch block by number failed", "height", height, log.ErrKey, err)
			// abandon for now, retry will be triggered by L1 block arrival
			return
		}

		r.storeBlockHeader(block)
		height += 1
		r.updateEagerFetchHeight(height)
	}
}

func (r *DataService) storeBlockHeader(blockHeader *types.Header) {
	err := r.blockResolver.AddBlock(blockHeader)
	if err != nil {
		r.logger.Error("Could not add block to host db.", log.ErrKey, err)
		// todo - handle unexpected errors here
	}
}

func (r *DataService) resetLiveStream() (chan *types.Header, ethereum.Subscription) {
	r.logger.Info("reconnecting to L1 new Heads")
	err := retry.Do(func() error {
		if !r.running.Load() {
			// break out of the loop if repository has stopped
			return retry.FailFast(errors.New("repository is stopped"))
		}
		err := r.ethClient.ReconnectIfClosed()
		if err != nil {
			r.logger.Warn("failed to reconnect to L1", log.ErrKey, err)
			return err
		}
		return nil
	}, retry.NewBackoffAndRetryForeverStrategy([]time.Duration{100 * time.Millisecond, 1 * time.Second, 5 * time.Second}, 10*time.Second))
	if err != nil {
		// this should only happen if repository has been stopped, because we retry reconnecting forever
		r.logger.Warn("unable to reconnect to L1", log.ErrKey, err)
		return nil, nil
	}

	ch, s := r.ethClient.BlockListener()
	r.logger.Info("successfully reconnected to L1 new Heads")
	return ch, s
}

func (r *DataService) FetchBlockByHeight(height *big.Int) (*types.Header, error) {
	return r.ethClient.HeaderByNumber(height)
}

// getEnclaveIdFromLog gets the enclave ID from the log topic
func getEnclaveIdFromLog(log types.Log) (gethcommon.Address, error) {
	// the enclaveID field is not indexed, we read it from the data field
	if len(log.Data) < 32 {
		return gethcommon.Address{}, errors.New("log data too short, expected enclaveID address")
	}
	return gethcommon.BytesToAddress(log.Data[:32]), nil
}

func increment(i *big.Int) *big.Int {
	return big.NewInt(0).Add(i, one)
}
