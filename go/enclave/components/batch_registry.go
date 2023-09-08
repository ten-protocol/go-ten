package components

import (
	"errors"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/go-obscuro/go/enclave/storage"

	"github.com/ethereum/go-ethereum/core/state"
	gethlog "github.com/ethereum/go-ethereum/log"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/measure"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/limiters"
)

type batchRegistry struct {
	storage storage.Storage
	logger  gethlog.Logger

	batchesCallback func(*core.Batch, types.Receipts)
	callbackMutex   sync.RWMutex
}

func NewBatchRegistry(storage storage.Storage, logger gethlog.Logger) BatchRegistry {
	return &batchRegistry{
		storage: storage,
		logger:  logger,
	}
}

func (br *batchRegistry) SubscribeForExecutedBatches(callback func(*core.Batch, types.Receipts)) {
	br.callbackMutex.Lock()
	defer br.callbackMutex.Unlock()
	br.batchesCallback = callback
}

func (br *batchRegistry) UnsubscribeFromBatches() {
	br.callbackMutex.Lock()
	defer br.callbackMutex.Unlock()

	br.batchesCallback = nil
}

func (br *batchRegistry) OnBatchExecuted(batch *core.Batch, receipts types.Receipts) {
	br.callbackMutex.RLock()
	defer br.callbackMutex.RUnlock()

	defer br.logger.Debug("Sending batch and events", log.BatchHashKey, batch.Hash(), log.DurationKey, measure.NewStopwatch())

	if br.batchesCallback != nil {
		br.batchesCallback(batch, receipts)
	}
}

func (br *batchRegistry) HasGenesisBatch() (bool, error) {
	genesisBatchStored := true
	_, err := br.storage.FetchHeadBatch()
	if err != nil {
		if !errors.Is(err, errutil.ErrNotFound) {
			return false, fmt.Errorf("could not retrieve current head batch. Cause: %w", err)
		}
		genesisBatchStored = false
	}

	return genesisBatchStored, nil
}

func (br *batchRegistry) BatchesAfter(batchSeqNo uint64, rollupLimiter limiters.RollupLimiter) ([]*core.Batch, error) {
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
		br.logger.Info("Added batch to rollup", log.BatchHashKey, batch.Hash(), log.BatchSeqNoKey, batch.SeqNo())
	}

	// Sanity check that the rollup includes consecutive batches (according to the seqNo)
	current := batches[0].SeqNo().Uint64()
	for i, b := range batches {
		if current+uint64(i) != b.SeqNo().Uint64() {
			return nil, fmt.Errorf("created invalid rollup with batches out of sequence")
		}
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
	blockchainState, err := br.storage.CreateStateDB(batch.Hash())
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
