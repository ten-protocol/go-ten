package l1

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ten-protocol/go-ten/contracts/generated/ManagementContract"

	"github.com/ten-protocol/go-ten/go/host/storage"

	"github.com/ten-protocol/go-ten/go/common/gethutil"

	"github.com/ten-protocol/go-ten/go/enclave/crosschain"
	"github.com/ten-protocol/go-ten/go/ethadapter/mgmtcontractlib"

	"github.com/ten-protocol/go-ten/go/common/subscription"

	"github.com/ten-protocol/go-ten/go/common/host"

	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/retry"
	"github.com/ten-protocol/go-ten/go/ethadapter"
)

var (
	l1BlockTime      = 12 * time.Second
	_timeoutNoBlocks = 2 * l1BlockTime // after this timeout we assume the subscription to the L1 node is not working
	one              = big.NewInt(1)
	ErrNoNextBlock   = errors.New("no next block")
)

type ContractType int

const (
	MgmtContract ContractType = iota
	MsgBus
)

// DataService is a host service for subscribing to new blocks and looking up L1 data
type DataService struct {
	blockSubscribers *subscription.Manager[host.L1BlockHandler]
	// this eth client should only be used by the repository, the repository may "reconnect" it at any time and don't want to interfere with other processes
	ethClient       ethadapter.EthClient
	logger          gethlog.Logger
	mgmtContractLib mgmtcontractlib.MgmtContractLib
	blobResolver    BlobResolver
	blockResolver   storage.BlockResolver

	running           atomic.Bool
	head              gethcommon.Hash
	contractAddresses map[ContractType][]gethcommon.Address
}

func NewL1DataService(
	ethClient ethadapter.EthClient,
	logger gethlog.Logger,
	mgmtContractLib mgmtcontractlib.MgmtContractLib,
	blobResolver BlobResolver,
	contractAddresses map[ContractType][]gethcommon.Address,
) *DataService {
	return &DataService{
		blockSubscribers:  subscription.NewManager[host.L1BlockHandler](),
		ethClient:         ethClient,
		running:           atomic.Bool{},
		logger:            logger,
		mgmtContractLib:   mgmtContractLib,
		blobResolver:      blobResolver,
		contractAddresses: contractAddresses,
	}
}

func (r *DataService) SetBlockResolver(br storage.BlockResolver) {
	r.blockResolver = br
}

func (r *DataService) Start() error {
	r.running.Store(true)

	// Repository constantly streams new blocks and forwards them to subscribers
	go r.streamLiveBlocks()
	return nil
}

func (r *DataService) Stop() error {
	r.running.Store(false)
	return nil
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

	logs, err := r.fetchMessageBusMgmtContractLogs(block)
	if err != nil {
		return nil, err
	}

	for _, l := range logs {
		if len(l.Topics) == 0 {
			r.logger.Error("Log has no topics. Should not happen", "txHash", l.TxHash)
			return nil, errors.New("log has no topics")
		}

		txData, err := r.fetchTxAndReceipt(l.TxHash)
		if err != nil {
			r.logger.Error("Error creating transaction data. Should not happen", "txHash", l.TxHash, "error", err)
			return nil, fmt.Errorf("error creating transaction data. Should not happen: %w", err)
		}

		// first topic is always the event signature
		switch l.Topics[0] {
		case crosschain.CrossChainEventID:
			err = r.processCrossChainLogs(l, txData, processed)
		case crosschain.ValueTransferEventID:
			err = r.processValueTransferLogs(l, txData, processed)
		case crosschain.SequencerEnclaveGrantedEventID:
			err = r.processSequencerLogs(l, txData, processed, common.SequencerAddedTx)
			err = r.processManagementContractTx(txData, processed) // we need to decode the InitialiseSecretTx
		case crosschain.SequencerEnclaveRevokedEventID:
			err = r.processSequencerLogs(l, txData, processed, common.SequencerRevokedTx)
		case crosschain.ImportantContractAddressUpdatedID:
			err = r.processManagementContractTx(txData, processed)
		case crosschain.RollupAddedID:
			err = r.processRollupLogs(l, txData, processed)
		case crosschain.NetworkSecretRequestedID:
			processed.AddEvent(common.SecretRequestTx, txData)
		case crosschain.NetworkSecretRespondedID:
			processed.AddEvent(common.SecretResponseTx, txData)
		default:
			// there are known events that we don't care about here
			r.logger.Debug("Unknown log topic", "topic", l.Topics[0], "txHash", l.TxHash)
		}

		if err != nil {
			r.logger.Error("Error processing log", "txHash", l.TxHash, "error", err)
			return nil, fmt.Errorf("error processing log: %w", err)
		}

	}
	return processed, nil
}

// fetchMessageBusMgmtContractLogs retrieves all logs from management contract and message bus addresses
func (r *DataService) fetchMessageBusMgmtContractLogs(block *types.Header) ([]types.Log, error) {
	blkHash := block.Hash()
	var allAddresses []gethcommon.Address
	allAddresses = append(allAddresses, r.contractAddresses[MgmtContract]...)
	allAddresses = append(allAddresses, r.contractAddresses[MsgBus]...)

	logs, err := r.ethClient.GetLogs(ethereum.FilterQuery{BlockHash: &blkHash, Addresses: allAddresses})
	if err != nil {
		return nil, fmt.Errorf("unable to fetch logs for L1 block - %w", err)
	}
	return logs, nil
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
	messages, err := crosschain.ConvertLogsToMessages([]types.Log{l}, crosschain.CrossChainEventName, crosschain.MessageBusABI)
	if err != nil {
		return err
	}

	txData.CrossChainMessages = messages
	processed.AddEvent(common.CrossChainMessageTx, txData)
	return nil
}

// processValueTransferLogs handles value transfer logs
func (r *DataService) processValueTransferLogs(l types.Log, txData *common.L1TxData, processed *common.ProcessedL1Data) error {
	transfers, err := crosschain.ConvertLogsToValueTransfers([]types.Log{l}, crosschain.ValueTransferEventName, crosschain.MessageBusABI)
	if err != nil {
		return err
	}
	txData.ValueTransfers = transfers
	processed.AddEvent(common.CrossChainValueTranserTx, txData)
	return nil
}

// processSequencerLogs handles sequencer logs
func (r *DataService) processSequencerLogs(l types.Log, txData *common.L1TxData, processed *common.ProcessedL1Data, txType common.L1TenEventType) error {
	enclaveID, err := getEnclaveIdFromLog(l)
	if err != nil {
		return err
	}
	txData.SequencerEnclaveID = enclaveID
	processed.AddEvent(txType, txData)
	return nil
}

func (r *DataService) processRollupLogs(l types.Log, txData *common.L1TxData, processed *common.ProcessedL1Data) error {
	abi, err := ManagementContract.ManagementContractMetaData.GetAbi()
	if err != nil {
		r.logger.Error("Error getting ManagementContract ABI", log.ErrKey, err)
		return err
	}
	var event ManagementContract.ManagementContractRollupAdded
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

// processManagementContractTx handles decoded transaction types
func (r *DataService) processManagementContractTx(txData *common.L1TxData, processed *common.ProcessedL1Data) error {
	decodedTx, _ := r.mgmtContractLib.DecodeTx(txData.Transaction)
	if decodedTx != nil {
		switch decodedTx.(type) {
		case *common.L1InitializeSecretTx:
			processed.AddEvent(common.InitialiseSecretTx, txData)
		case *common.L1SetImportantContractsTx:
			processed.AddEvent(common.SetImportantContractsTx, txData)
		case *common.L1PermissionSeqTx:
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
func (r *DataService) streamLiveBlocks() {
	liveStream, streamSub := r.resetLiveStream()
	for r.running.Load() {
		select {
		case blockHeader := <-liveStream:
			r.logger.Info(fmt.Sprintf("received block from l1 stream: %v", blockHeader))
			err := r.blockResolver.AddBlock(blockHeader)
			if err != nil {
				r.logger.Error("Could not add block to host db.", log.ErrKey, err)
				// todo - handle unexpected errors here
			}

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
			// reset stream to ensure it has not died
			liveStream, streamSub = r.resetLiveStream()
		}
	}
	r.logger.Info("block streaming stopped")
	if streamSub != nil {
		streamSub.Unsubscribe()
	}
	if liveStream != nil {
		close(liveStream)
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
