package components

import (
	"errors"
	"fmt"
	"sync"

	"github.com/obscuronet/go-obscuro/go/enclave/storage"

	"github.com/ethereum/go-ethereum/params"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/measure"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/limiters"
)

type batchRegistryImpl struct {
	storage       storage.Storage
	logger        gethlog.Logger
	chainConfig   *params.ChainConfig
	batchProducer BatchProducer
	sigValidator  *SignatureValidator
	// Channel on which batches will be pushed. It is held by another caller outside the
	// batch registry.
	batchSubscription *chan *core.Batch
	// Channel for pushing batch height numbers which are needed in order
	// to figure out what events to send to subscribers.
	eventSubscription *chan uint64

	subscriptionMutex sync.Mutex
}

func NewBatchRegistry(storage storage.Storage, batchProducer BatchProducer, sigValidator *SignatureValidator, chainConfig *params.ChainConfig, logger gethlog.Logger) BatchRegistry {
	return &batchRegistryImpl{
		storage:       storage,
		batchProducer: batchProducer,
		sigValidator:  sigValidator,
		chainConfig:   chainConfig,
		logger:        logger,
	}
}

func (br *batchRegistryImpl) CommitBatch(cb *ComputedBatch) error {
	_, err := cb.Commit(true)
	return err
}

func (br *batchRegistryImpl) SubscribeForEvents() chan uint64 {
	evSub := make(chan uint64)
	br.eventSubscription = &evSub
	return *br.eventSubscription
}

func (br *batchRegistryImpl) UnsubscribeFromEvents() {
	br.eventSubscription = nil
}

// StoreBatch - stores a batch and if it is canonical, it sends the events to subscribers
// Deprecated
func (br *batchRegistryImpl) StoreBatch(batch *core.Batch, receipts types.Receipts) error {
	defer br.logger.Info("Registry StoreBatch() exit", log.BatchHashKey, batch.Hash(), log.DurationKey, measure.NewStopwatch())

	// Check if this batch is already stored.
	if _, err := br.storage.FetchBatchHeader(batch.Hash()); err == nil {
		br.logger.Warn("Attempted to store batch twice! This indicates issues with the batch processing loop")
		return nil
	}

	if err := br.storage.StoreBatch(batch); err != nil {
		return fmt.Errorf("failed to store batch. Cause: %w", err)
	}

	if err := br.storage.StoreExecutedBatch(batch, receipts); err != nil {
		return fmt.Errorf("failed to store batch. Cause: %w", err)
	}

	br.notifySubscriber(batch)

	return nil
}

func (br *batchRegistryImpl) ExecuteBatch(batch *core.Batch) (types.Receipts, error) {

	defer br.logger.Info("Executed batch", log.BatchHashKey, batch.Hash(), log.DurationKey, measure.NewStopwatch())

	// Validators recompute the entire batch using the same batch context
	// if they have all necessary prerequisites like having the l1 block processed
	// and the parent hash. This recomputed batch is then checked against the incoming batch.
	// If the sequencer has tampered with something the hash will not add up and validation will
	// produce an error.
	cb, err := br.batchProducer.ComputeBatch(&BatchExecutionContext{
		BlockPtr:     batch.Header.L1Proof,
		ParentPtr:    batch.Header.ParentHash,
		Transactions: batch.Transactions,
		AtTime:       batch.Header.Time,
		ChainConfig:  br.chainConfig,
		SequencerNo:  batch.Header.SequencerOrderNo,
	})
	if err != nil {
		return nil, fmt.Errorf("failed computing batch %s. Cause: %w", batch.Hash(), err)
	}

	if cb.Batch.Hash() != batch.Hash() {
		// todo @stefan - generate a validator challenge here and return it
		br.logger.Error(fmt.Sprintf("Error validating batch. Calculated: %+v\n Incoming: %+v\n", cb.Batch.Header, batch.Header))
		return nil, fmt.Errorf("batch is in invalid state. Incoming hash: %s  Computed hash: %s", batch.Hash(), cb.Batch.Hash())
	}

	if _, err := cb.Commit(true); err != nil {
		return nil, fmt.Errorf("cannot commit stateDB for incoming valid batch %s. Cause: %w", batch.Hash(), err)
	}

	return cb.Receipts, nil
}

func (br *batchRegistryImpl) notifySubscriber(batch *core.Batch) {
	defer br.logger.Info("Registry notified subscribers of batch", log.BatchHashKey, batch.Hash(), log.DurationKey, measure.NewStopwatch())

	br.subscriptionMutex.Lock()
	subscriptionChan := br.batchSubscription
	eventChan := br.eventSubscription
	br.subscriptionMutex.Unlock()

	if subscriptionChan != nil {
		*subscriptionChan <- batch
	}

	if br.eventSubscription != nil {
		*eventChan <- batch.NumberU64()
	}
}

func (br *batchRegistryImpl) GetHeadBatch() (*core.Batch, error) {
	return br.storage.FetchHeadBatch()
}

func (br *batchRegistryImpl) GetHeadBatchFor(blockHash common.L1BlockHash) (*core.Batch, error) {
	return br.storage.FetchHeadBatchForBlock(blockHash)
}

func (br *batchRegistryImpl) GetBatch(batchHash common.L2BatchHash) (*core.Batch, error) {
	return br.storage.FetchBatch(batchHash)
}

func (br *batchRegistryImpl) Subscribe() chan *core.Batch {
	br.subscriptionMutex.Lock()
	defer br.subscriptionMutex.Unlock()
	subChannel := make(chan *core.Batch)

	br.batchSubscription = &subChannel
	return *br.batchSubscription
}

func (br *batchRegistryImpl) Unsubscribe() {
	br.subscriptionMutex.Lock()
	defer br.subscriptionMutex.Unlock()
	if br.batchSubscription != nil {
		close(*br.batchSubscription)
		br.batchSubscription = nil
	}
}

func (br *batchRegistryImpl) FindAncestralBatchFor(block *common.L1Block) (*core.Batch, error) {
	currentBlock := block
	var ancestorBatch *core.Batch
	var err error

	br.logger.Trace("Searching for ancestral batch")
	// todo - this for loop should have more edge cases.
	for ancestorBatch == nil {
		if currentBlock.NumberU64() == common.L1GenesisHeight {
			return nil, fmt.Errorf("reached genesis block")
		}

		ancestorBatch, err = br.GetHeadBatchFor(currentBlock.Hash())
		if err != nil && !errors.Is(err, errutil.ErrNotFound) {
			return nil, fmt.Errorf("unable to get latest ancestral batch. Cause: %w", err)
		}

		parentBlockHash := currentBlock.ParentHash()
		currentBlock, err = br.storage.FetchBlock(parentBlockHash)
		if err != nil {
			return nil, fmt.Errorf("unable to find block %s in ancestral chain. height %d. Cause: %w", parentBlockHash.Hex(), block.Number(), err)
		}
	}

	br.logger.Trace("Found ancestral batch")

	return ancestorBatch, nil
}

func (br *batchRegistryImpl) HasGenesisBatch() (bool, error) {
	genesisBatchStored := true
	_, err := br.GetHeadBatch()
	if err != nil {
		if !errors.Is(err, errutil.ErrNotFound) {
			return false, fmt.Errorf("could not retrieve current head batch. Cause: %w", err)
		}
		genesisBatchStored = false
	}

	return genesisBatchStored, nil
}

func (br *batchRegistryImpl) BatchesAfter(batchSeqNo uint64, rollupLimiter limiters.RollupLimiter) ([]*core.Batch, error) {
	batches := make([]*core.Batch, 0)

	var batch *core.Batch
	var err error
	if batch, err = br.storage.FetchBatchBySeqNo(batchSeqNo); err != nil {
		return nil, err
	}
	batches = append(batches, batch)

	headBatch, err := br.storage.FetchHeadBatch()
	if err != nil {
		return nil, err
	}

	if headBatch.SeqNo().Uint64() < batch.SeqNo().Uint64() {
		return nil, fmt.Errorf("head batch height %d is in the past compared to requested batch %d",
			headBatch.SeqNo().Uint64(),
			batch.SeqNo().Uint64())
	}
	for batch.SeqNo().Cmp(headBatch.SeqNo()) != 0 {
		if didAcceptBatch, err := rollupLimiter.AcceptBatch(batch); err != nil {
			return nil, err
		} else if !didAcceptBatch {
			return batches, nil
		}

		batch, err = br.storage.FetchBatchBySeqNo(batch.SeqNo().Uint64() + 1)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve batch by sequence number less than the head batch. Cause: %w", err)
		}

		batches = append(batches, batch)
		br.logger.Info("Added batch to rollup", log.BatchHashKey, batch.Hash(), "seqNo", batch.SeqNo())
	}

	return batches, nil
}

func (br *batchRegistryImpl) GetBatchStateAtHeight(blockNumber *gethrpc.BlockNumber) (*state.StateDB, error) {
	// We retrieve the batch of interest.
	batch, err := br.GetBatchAtHeight(*blockNumber)
	if err != nil {
		return nil, err
	}

	// We get that of the chain at that height
	blockchainState, err := br.storage.CreateStateDB(batch.Hash())
	if err != nil {
		return nil, fmt.Errorf("could not create stateDB. Cause: %w", err)
	}

	if blockchainState == nil {
		return nil, fmt.Errorf("unable to fetch chain state for batch %s", batch.Hash().Hex())
	}

	return blockchainState, err
}

func (br *batchRegistryImpl) GetBatchAtHeight(height gethrpc.BlockNumber) (*core.Batch, error) {
	var batch *core.Batch
	switch height {
	case gethrpc.EarliestBlockNumber:
		genesisBatch, err := br.storage.FetchBatchByHeight(0)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve genesis rollup. Cause: %w", err)
		}
		batch = genesisBatch
	case gethrpc.PendingBlockNumber:
		// todo - depends on the current pending rollup; leaving it for a different iteration as it will need more thought
		return nil, fmt.Errorf("requested balance for pending block. This is not handled currently")
	case gethrpc.LatestBlockNumber:
		headBatch, err := br.storage.FetchHeadBatch()
		if err != nil {
			return nil, fmt.Errorf("batch with requested height %d was not found. Cause: %w", height, err)
		}
		batch = headBatch
	default:
		maybeBatch, err := br.storage.FetchBatchByHeight(uint64(height))
		if err != nil {
			return nil, fmt.Errorf("batch with requested height %d could not be retrieved. Cause: %w", height, err)
		}
		batch = maybeBatch
	}
	return batch, nil
}
