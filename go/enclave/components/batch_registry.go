package components

import (
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
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
)

type batchRegistry struct {
	storage db.Storage
	logger  gethlog.Logger

	subscription      *chan *core.Batch
	subscriptionMutex sync.Mutex
}

func NewBatchRegistry(storage db.Storage, logger gethlog.Logger) BatchRegistry {
	return &batchRegistry{
		storage: storage,
		logger:  logger,
	}
}

func (br *batchRegistry) StoreBatch(batch *core.Batch, receipts types.Receipts) error {
	if err := br.updateHeadPointers(batch, receipts); err != nil {
		return fmt.Errorf("failed updating head pointers. Cause: %w", err)
	}

	if err := br.storage.StoreBatch(batch, receipts); err != nil {
		return fmt.Errorf("failed to store batch. Cause: %w", err)
	}

	// TODO - we probably dont want to spawn a goroutine everytime
	go br.updateSubscribers(batch)

	return nil
}

func (br *batchRegistry) updateSubscribers(batch *core.Batch) {
	br.subscriptionMutex.Lock()
	subscriptionChan := br.subscription
	br.subscriptionMutex.Unlock()

	if subscriptionChan != nil {
		*subscriptionChan <- batch
	}
}

func (br *batchRegistry) updateHeadPointers(batch *core.Batch, receipts types.Receipts) error {
	if err := br.updateBlockPointers(batch, receipts); err != nil {
		return err
	}

	return br.updateBatchPointers(batch)
}

func (br *batchRegistry) updateBatchPointers(batch *core.Batch) error {
	if head, err := br.storage.FetchHeadBatch(); err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return err
	} else if head != nil && batch.NumberU64() < head.NumberU64() {
		return nil
	}

	return br.storage.SetHeadBatchPointer(batch)
}

func (br *batchRegistry) updateBlockPointers(batch *core.Batch, receipts types.Receipts) error {
	head, err := br.GetHeadBatchFor(batch.Header.L1Proof)

	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return fmt.Errorf("unexpected error while getting head batch for block. Cause: %w", err)
	} else if head != nil && batch.NumberU64() < head.NumberU64() {
		return fmt.Errorf("inappropriate update from previous head with height %d to new head with height %d for same l1 block", head.NumberU64(), batch.NumberU64())
	}

	return br.storage.UpdateHeadBatch(batch.Header.L1Proof, batch, receipts)
}

func (br *batchRegistry) GetHeadBatch() (*core.Batch, error) {
	return br.storage.FetchHeadBatch()
}

func (br *batchRegistry) GetHeadBatchFor(blockHash common.L1BlockHash) (*core.Batch, error) {
	return br.storage.FetchHeadBatchForBlock(blockHash)
}

func (br *batchRegistry) GetBatch(batchHash common.L2BatchHash) (*core.Batch, error) {
	return br.storage.FetchBatch(batchHash)
}

func (br *batchRegistry) Subscribe() chan *core.Batch {
	br.subscriptionMutex.Lock()
	defer br.subscriptionMutex.Unlock()
	subChannel := make(chan *core.Batch, 100)
	br.subscription = &subChannel
	return *br.subscription
}

func (br *batchRegistry) Unsubscribe() {
	br.subscriptionMutex.Lock()
	defer br.subscriptionMutex.Unlock()
	if br.subscription != nil {
		close(*br.subscription)
		br.subscription = nil
	}
}

func (br *batchRegistry) FindAncestralBatchFor(block *common.L1Block) (*core.Batch, error) {
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

		currentBlock, err = br.storage.FetchBlock(currentBlock.ParentHash())
		if err != nil {
			return nil, fmt.Errorf("unable to find parent for block in ancestral chain. Cause: %w", err)
		}
	}

	br.logger.Trace("Found ancestral batch")

	return ancestorBatch, nil
}

func (br *batchRegistry) HasGenesisBatch() (bool, error) {
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

func (br *batchRegistry) BatchesAfter(batchHash gethcommon.Hash) ([]*core.Batch, error) {
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

func (br *batchRegistry) GetBatchStateAtHeight(blockNumber *gethrpc.BlockNumber) (*state.StateDB, error) {
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

func (br *batchRegistry) GetBatchAtHeight(height gethrpc.BlockNumber) (*core.Batch, error) {
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
