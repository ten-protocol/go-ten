package components

import (
	"bytes"
	"errors"
	"fmt"
	"sync"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
)

type batchRegistryImpl struct {
	storage db.Storage
	logger  gethlog.Logger

	// Channel on which batches will be pushed. It is held by another caller outside the
	// batch registry.
	batchSubscription *chan *core.Batch
	// Channel for pushing batch height numbers which are needed in order
	// to figure out what events to send to subscribers.
	eventSubscription *chan uint64

	subscriptionMutex sync.Mutex
}

func NewBatchRegistry(storage db.Storage, logger gethlog.Logger) BatchRegistry {
	return &batchRegistryImpl{
		storage: storage,
		logger:  logger,
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

// StoreBatch - stores a batch and if it is the new l2 head, then registry will update
// stored head pointers
func (br *batchRegistryImpl) StoreBatch(batch *core.Batch, receipts types.Receipts) error {
	// Check if this batch is already stored.
	if _, err := br.GetBatch(*batch.Hash()); err == nil {
		br.logger.Warn("Attempted to store batch twice! This indicates issues with the batch processing loop")
		return nil
	}

	dbTransaction := br.storage.NewTransaction()

	isHeadBatch, err := br.updateHeadPointers(batch, receipts, dbTransaction)
	if err != nil {
		return fmt.Errorf("failed updating head pointers. Cause: %w", err)
	}

	if err = dbTransaction.StoreBatch(batch, receipts); err != nil {
		return fmt.Errorf("failed to store batch. Cause: %w", err)
	}

	if err = dbTransaction.Commit(); err != nil {
		return fmt.Errorf("failed to commit batch to db. Cause: %w", err)
	}

	br.notifySubscriber(batch, isHeadBatch)

	return nil
}

func (br *batchRegistryImpl) notifySubscriber(batch *core.Batch, isHeadBatch bool) {
	br.subscriptionMutex.Lock()
	subscriptionChan := br.batchSubscription
	eventChan := br.eventSubscription
	br.subscriptionMutex.Unlock()

	if subscriptionChan != nil {
		*subscriptionChan <- batch
	}

	if br.eventSubscription != nil && isHeadBatch {
		*eventChan <- batch.NumberU64()
	}
}

func (br *batchRegistryImpl) updateHeadPointers(batch *core.Batch, receipts types.Receipts, storageUpdater db.StorageUpdater) (bool, error) {
	if err := br.updateBlockPointers(batch, receipts, storageUpdater); err != nil {
		return false, err
	}

	return br.updateBatchPointers(batch, storageUpdater)
}

func (br *batchRegistryImpl) updateBatchPointers(batch *core.Batch, storageUpdater db.StorageUpdater) (bool, error) {
	if head, err := br.storage.FetchHeadBatch(); err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return false, err
	} else if head != nil && batch.NumberU64() < head.NumberU64() {
		return false, nil
	}

	return true, storageUpdater.SetHeadBatchPointer(batch)
}

func (br *batchRegistryImpl) updateBlockPointers(batch *core.Batch, receipts types.Receipts, storageUpdater db.StorageUpdater) error {
	head, err := br.GetHeadBatchFor(batch.Header.L1Proof)

	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return fmt.Errorf("unexpected error while getting head batch for block. Cause: %w", err)
	} else if head != nil && batch.NumberU64() < head.NumberU64() {
		return fmt.Errorf("inappropriate update from previous head with height %d to new head with height %d for same l1 block", head.NumberU64(), batch.NumberU64())
	}

	return storageUpdater.UpdateHeadBatch(batch.Header.L1Proof, batch, receipts)
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

func (br *batchRegistryImpl) Subscribe(lastKnownHead *common.L2BatchHash) (chan *core.Batch, error) {
	br.subscriptionMutex.Lock()
	defer br.subscriptionMutex.Unlock()
	missingBatches, err := br.getMissingBatches(lastKnownHead)
	if err != nil {
		return nil, err
	}

	subChannel := make(chan *core.Batch, len(missingBatches))
	for i := len(missingBatches) - 1; i >= 0; i-- {
		batch := missingBatches[i]
		subChannel <- batch
	}

	br.batchSubscription = &subChannel
	return *br.batchSubscription, nil
}

func (br *batchRegistryImpl) Unsubscribe() {
	br.subscriptionMutex.Lock()
	defer br.subscriptionMutex.Unlock()
	if br.batchSubscription != nil {
		close(*br.batchSubscription)
		br.batchSubscription = nil
	}
}

func (br *batchRegistryImpl) getMissingBatches(fromHash *common.L2BatchHash) ([]*core.Batch, error) {
	if fromHash == nil {
		return nil, nil
	}

	from, err := br.GetBatch(*fromHash)
	if err != nil {
		br.logger.Error("Error while attempting to stream from batch", log.ErrKey, err)
		return nil, err
	}

	to, err := br.GetHeadBatch()
	if err != nil {
		br.logger.Error("Unable to get head batch while attempting to stream", log.ErrKey, err)
		return nil, err
	}

	missingBatches := make([]*core.Batch, 0)
	for !bytes.Equal(to.Hash().Bytes(), from.Hash().Bytes()) {
		if to.NumberU64() == 0 {
			br.logger.Error("Reached genesis when seeking missing batches to stream", log.ErrKey, err)
			return nil, err
		}

		if from.NumberU64() == to.NumberU64() {
			from, err = br.GetBatch(from.Header.ParentHash)
			if err != nil {
				br.logger.Error("Unable to get batch in chain while attempting to stream", log.ErrKey, err)
				return nil, err
			}
		}

		missingBatches = append(missingBatches, to)
		to, err = br.GetBatch(to.Header.ParentHash)
		if err != nil {
			br.logger.Error("Unable to get batch in chain while attempting to stream", log.ErrKey, err)
			return nil, err
		}
	}

	return missingBatches, nil
}

func (br *batchRegistryImpl) FindAncestralBatchFor(block *common.L1Block) (*core.Batch, error) {
	currentBlock := block
	var ancestorBatch *core.Batch = nil
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
			return nil, fmt.Errorf("unable to find parent for block %s in ancestral chain. Cause: %w", parentBlockHash.Hex(), err)
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

func (br *batchRegistryImpl) BatchesAfter(batchHash gethcommon.Hash) ([]*core.Batch, error) {
	batches := make([]*core.Batch, 0)

	var batch *core.Batch
	var err error
	if batchHash == gethcommon.BigToHash(gethcommon.Big0) {
		if batch, err = br.storage.FetchBatchByHeight(0); err != nil {
			return nil, err
		}
		batches = append(batches, batch)
	} else {
		if batch, err = br.storage.FetchBatch(batchHash); err != nil {
			return nil, err
		}
	}

	headBatch, err := br.storage.FetchHeadBatch()
	if err != nil {
		return nil, err
	}

	if headBatch.NumberU64() < batch.NumberU64() {
		return nil, errors.New("head batch height is in the past compared to requested batch")
	}

	for batch.Number().Cmp(headBatch.Number()) != 0 {
		batch, _ = br.storage.FetchBatchByHeight(batch.NumberU64() + 1)
		batches = append(batches, batch)
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
	blockchainState, err := br.storage.CreateStateDB(*batch.Hash())
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
