package l1

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

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
	// todo (@matt) make this configurable?
	_timeoutNoBlocks = 30 * time.Second
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
func (r *DataService) FetchNextBlock(prevBlockHash gethcommon.Hash) (*types.Block, bool, error) {
	if prevBlockHash == r.head {
		// prevBlock is the latest known head
		return nil, false, ErrNoNextBlock
	}

	if prevBlockHash == (gethcommon.Hash{}) {
		// prevBlock is empty, so we are starting from genesis
		blk, err := r.ethClient.BlockByNumber(big.NewInt(0))
		if err != nil {
			return nil, false, fmt.Errorf("could not find genesis block - %w", err)
		}
		return blk, false, nil
	}

	// the latestCanonAncestor will usually return the prevBlock itself but this step is necessary to walk back if there was a fork
	lca, err := r.latestCanonAncestor(prevBlockHash)
	if err != nil {
		return nil, false, err
	}
	// and send the canonical block at the height after that
	// (which may be a fork, or it may just be the next on the same branch if we are catching-up)
	blk, err := r.ethClient.BlockByNumber(increment(lca.Number()))
	if err != nil {
		if errors.Is(err, ethereum.NotFound) {
			return nil, false, ErrNoNextBlock
		}
		return nil, false, fmt.Errorf("could not find block after latest canon ancestor, height=%s - %w", increment(lca.Number()), err)
	}

	return blk, blk.Hash() == r.head, nil
}

func (r *DataService) latestCanonAncestor(blkHash gethcommon.Hash) (*types.Block, error) {
	blk, err := r.ethClient.BlockByHash(blkHash)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch L1 block with hash=%s - %w", blkHash, err)
	}
	canonAtSameHeight, err := r.ethClient.BlockByNumber(blk.Number())
	if err != nil {
		return nil, fmt.Errorf("unable to fetch L1 block at height=%d - %w", blk.Number(), err)
	}
	if blk.Hash() != canonAtSameHeight.Hash() {
		return r.latestCanonAncestor(blk.ParentHash())
	}
	return blk, nil
}

// GetTenRelevantTransactions processes logs in their natural order without grouping by transaction hash.
func (r *DataService) GetTenRelevantTransactions(block *common.L1Block) (*common.ProcessedL1Data, error) {
	processed := &common.ProcessedL1Data{
		BlockHeader: block.Header(),
		Events:      []common.L1Event{},
	}

	logs, err := r.fetchMessageBusMgmtContractLogs(block)
	if err != nil {
		return nil, err
	}

	for _, l := range logs {
		if len(l.Topics) == 0 {
			r.logger.Warn("Log has no topics", "txHash", l.TxHash)
			continue
		}

		txData, err := r.fetchTxAndReceipt(l.TxHash)
		if err != nil {
			r.logger.Error("Error creating transaction data", "txHash", l.TxHash, "error", err)
			continue
		}

		// first topic is always the event signature
		switch l.Topics[0] {
		case crosschain.CrossChainEventID:
			r.processCrossChainLogs(l, txData, processed)
		case crosschain.ValueTransferEventID:
			r.processValueTransferLogs(l, txData, processed)
		case crosschain.SequencerEnclaveGrantedEventID:
			r.processSequencerLogs(l, txData, processed, common.SequencerAddedTx)
			r.processManagementContractTx(txData, processed) // we need to decode the InitialiseSecretTx
		case crosschain.SequencerEnclaveRevokedEventID:
			r.processSequencerLogs(l, txData, processed, common.SequencerRevokedTx)
		case crosschain.ImportantContractAddressUpdatedID:
			r.processManagementContractTx(txData, processed)
		case crosschain.RollupAddedID:
			r.processManagementContractTx(txData, processed)
		case crosschain.NetworkSecretRequestedID:
			processed.AddEvent(common.SecretRequestTx, txData)
		case crosschain.NetworkSecretRespondedID:
			processed.AddEvent(common.SecretResponseTx, txData)
		default:
			r.logger.Warn("Unknown log topic", "topic", l.Topics[0], "txHash", l.TxHash)
		}
	}

	return processed, nil
}

// fetchMessageBusMgmtContractLogs retrieves all logs from management contract and message bus addresses
func (r *DataService) fetchMessageBusMgmtContractLogs(block *common.L1Block) ([]types.Log, error) {
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
func (r *DataService) processCrossChainLogs(l types.Log, txData *common.L1TxData, processed *common.ProcessedL1Data) {
	if messages, err := crosschain.ConvertLogsToMessages([]types.Log{l}, crosschain.CrossChainEventName, crosschain.MessageBusABI); err == nil {
		txData.CrossChainMessages = messages
		processed.AddEvent(common.CrossChainMessageTx, txData)
	}
}

// processValueTransferLogs handles value transfer logs
func (r *DataService) processValueTransferLogs(l types.Log, txData *common.L1TxData, processed *common.ProcessedL1Data) {
	if transfers, err := crosschain.ConvertLogsToValueTransfers([]types.Log{l}, crosschain.ValueTransferEventName, crosschain.MessageBusABI); err == nil {
		txData.ValueTransfers = transfers
		processed.AddEvent(common.CrossChainValueTranserTx, txData)
	}
}

// processSequencerLogs handles sequencer logs
func (r *DataService) processSequencerLogs(l types.Log, txData *common.L1TxData, processed *common.ProcessedL1Data, txType common.L1TenEventType) {
	if enclaveID, err := getEnclaveIdFromLog(l); err == nil {
		txData.SequencerEnclaveID = enclaveID
		processed.AddEvent(txType, txData)
	}
}

// processManagementContractTx handles decoded transaction types
func (r *DataService) processManagementContractTx(txData *common.L1TxData, processed *common.ProcessedL1Data) {
	b := processed.BlockHeader
	if decodedTx := r.mgmtContractLib.DecodeTx(txData.Transaction); decodedTx != nil {
		switch t := decodedTx.(type) {
		case *common.L1InitializeSecretTx:
			processed.AddEvent(common.InitialiseSecretTx, txData)
		case *common.L1SetImportantContractsTx:
			processed.AddEvent(common.SetImportantContractsTx, txData)
		case *common.L1RollupHashes:
			if blobs, err := r.blobResolver.FetchBlobs(context.Background(), b, t.BlobHashes); err == nil {
				txData.Blobs = blobs
				processed.AddEvent(common.RollupTx, txData)
			}
		default:
			// this should never happen since the specific events should always decode into one of these types
			r.logger.Error("Unknown tx type", "txHash", txData.Transaction.Hash().Hex())
		}
	}
}

// stream blocks from L1 as they arrive and forward them to subscribers, no guarantee of perfect ordering or that there won't be gaps.
// If streaming is interrupted it will carry on from latest, it won't try to replay missed blocks.
func (r *DataService) streamLiveBlocks() {
	liveStream, streamSub := r.resetLiveStream()
	for r.running.Load() {
		select {
		case header := <-liveStream:
			r.head = header.Hash()
			block, err := r.ethClient.BlockByHash(header.Hash())
			if err != nil {
				r.logger.Error("Error fetching new block", log.BlockHashKey, header.Hash(),
					log.BlockHeightKey, header.Number, log.ErrKey, err)
				continue
			}
			for _, handler := range r.blockSubscribers.Subscribers() {
				go handler.HandleBlock(block)
			}
		case <-time.After(_timeoutNoBlocks):
			r.logger.Warn("no new blocks received since timeout", "timeout", _timeoutNoBlocks)
			// reset stream to ensure it has not died
			liveStream, streamSub = r.resetLiveStream()
		}
	}

	if streamSub != nil {
		streamSub.Unsubscribe()
	}
}

func (r *DataService) resetLiveStream() (chan *types.Header, ethereum.Subscription) {
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
	return r.ethClient.BlockListener()
}

func (r *DataService) FetchBlockByHeight(height *big.Int) (*types.Block, error) {
	return r.ethClient.BlockByNumber(height)
}

// getEnclaveIdFromLog gets the enclave ID from the log topic
func getEnclaveIdFromLog(log types.Log) (gethcommon.Address, error) {
	if len(log.Topics) != 1 {
		return gethcommon.Address{}, fmt.Errorf("invalid number of topics in log: %d", len(log.Topics))
	}

	return gethcommon.BytesToAddress(log.Topics[0].Bytes()), nil
}

func increment(i *big.Int) *big.Int {
	return i.Add(i, one)
}
