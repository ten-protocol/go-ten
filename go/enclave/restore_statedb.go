package enclave

import (
	"context"
	"errors"
	"fmt"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/enclave/components"
	"github.com/ten-protocol/go-ten/go/enclave/storage"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

// this function looks at the batch chain and makes sure the resulting stateDB snapshots are available, replaying them if needed
// (if there had been a clean shutdown and all stateDB data was persisted this should do nothing)
func syncExecutedBatchesWithEVMStateDB(ctx context.Context, storage storage.Storage, registry components.BatchRegistry, batchExecutor components.BatchExecutor, logger gethlog.Logger) error {
	if registry.HeadBatchSeq() == nil {
		// not initialised yet
		return nil
	}
	batch, err := storage.FetchBatchBySeqNo(ctx, registry.HeadBatchSeq().Uint64())
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			// there is no head batch, this is probably a new node - there is no state to rebuild
			logger.Info("no head batch found in DB after restart", log.ErrKey, err)
			return nil
		}
		return fmt.Errorf("unexpected error fetching head batch to resync- %w", err)
	}
	if !stateDBAvailableForBatch(ctx, registry, batch.Hash()) {
		logger.Info("state not available for latest batch after restart - rebuilding stateDB cache from batches")
		err = markUnexecutedBatches(ctx, storage, registry, batchExecutor, logger)
		if err != nil {
			return fmt.Errorf("unable to replay batches to restore valid state - %w", err)
		}
	}
	return nil
}

// The enclave caches a stateDB instance against each batch hash, this is the input state when producing the following
// batch in the chain and is used to query state at a certain height.
//
// This method checks if the stateDB data is available for a given batch hash (so it can be restored if not)
func stateDBAvailableForBatch(ctx context.Context, registry components.BatchRegistry, hash common.L2BatchHash) bool {
	_, _, err := registry.GetBatchState(ctx, gethrpc.BlockNumberOrHash{BlockHash: &hash})
	// Note: logging here would be too verbose, but we track this in the loop
	return err == nil
}

// markUnexecutedBatches marks the batches for which the statedb is missing as un-executed
func markUnexecutedBatches(ctx context.Context, storage storage.Storage, registry components.BatchRegistry, _ components.BatchExecutor, logger gethlog.Logger) error {
	logger.Info("markUnexecutedBatches - start", "headBatchSeq", registry.HeadBatchSeq())

	// `currentBatch` variable will eventually be the latest batch for which we are able to produce a StateDB
	// - we will then set that as the head of the L2 so that this node can rebuild its missing state
	headSeq := registry.HeadBatchSeq()
	logger.Info("markUnexecutedBatches - fetching head batch", "seqNo", headSeq, "seqNoIsNil", headSeq == nil)

	currentBatch, err := storage.FetchBatchBySeqNo(ctx, headSeq.Uint64())
	if err != nil {
		logger.Error("markUnexecutedBatches - failed to fetch head batch", log.ErrKey, err)
		return fmt.Errorf("no head batch found in DB but expected to replay batches - %w", err)
	}
	logger.Info("markUnexecutedBatches - fetched head batch", log.BatchSeqNoKey, currentBatch.SeqNo(), log.BatchHashKey, currentBatch.Hash(), "batchIsNil", currentBatch == nil, "headerIsNil", currentBatch.Header == nil)

	if currentBatch == nil {
		logger.Crit("markUnexecutedBatches - currentBatch is nil!")
	}
	if currentBatch.Header == nil {
		logger.Crit("markUnexecutedBatches - currentBatch.Header is nil!")
	}

	// loop backwards building a slice of all batches that don't have cached stateDB data available
	logger.Info("markUnexecutedBatches - entering loop")
	loopIteration := 0
	for !stateDBAvailableForBatch(ctx, registry, currentBatch.Hash()) {
		loopIteration++
		logger.Info("markUnexecutedBatches - loop iteration", "iteration", loopIteration, log.BatchSeqNoKey, currentBatch.SeqNo(), log.BatchHashKey, currentBatch.Hash())

		logger.Info("markUnexecutedBatches - marking batch as unexecuted", log.BatchSeqNoKey, currentBatch.SeqNo())
		err = storage.MarkBatchAsUnexecuted(ctx, currentBatch.SeqNo())
		if err != nil {
			logger.Error("markUnexecutedBatches - failed to mark batch as unexecuted", log.ErrKey, err)
			return fmt.Errorf("unable to mark batch as unexecuted - %w", err)
		}
		logger.Info("markUnexecutedBatches - marked batch as unexecuted successfully")

		if currentBatch.NumberU64() == common.L2GenesisHeight {
			// no more parents to check, replaying from genesis
			logger.Info("markUnexecutedBatches - reached genesis, breaking loop")
			break
		}

		/*		// optional - try to execute the batch
				canExecute, err := registry.CanExecute(ctx, currentBatch.Header)
				if err != nil {
					return fmt.Errorf("could not determine the execution prerequisites for batchHeader %s. Cause: %w", currentBatch.Hash(), err)
				}
				logger.Trace("Can execute stored batch", log.BatchSeqNoKey, currentBatch.SeqNo(), "can", canExecute)

				if canExecute {
					err = registry.ExecuteBatch(ctx, batchExecutor, currentBatch.Header)
					if err != nil {
						return fmt.Errorf("could not execute batch %s. Cause: %w", currentBatch.Hash(), err)
					}
				}
		*/

		logger.Info("markUnexecutedBatches - fetching parent batch", "currentBatchHeader", currentBatch.Header != nil, "parentHash", currentBatch.Header.ParentHash)
		currentBatch, err = storage.FetchBatch(ctx, currentBatch.Header.ParentHash)
		if err != nil {
			logger.Error("markUnexecutedBatches - failed to fetch parent batch", log.ErrKey, err)
			return fmt.Errorf("unable to fetch previous batch while rolling back to stable state - %w", err)
		}
		logger.Info("markUnexecutedBatches - fetched parent batch", log.BatchSeqNoKey, currentBatch.SeqNo(), log.BatchHashKey, currentBatch.Hash(), "batchIsNil", currentBatch == nil, "headerIsNil", currentBatch.Header == nil)

		if currentBatch == nil {
			logger.Crit("markUnexecutedBatches - parent batch is nil!")
		}
		if currentBatch.Header == nil {
			logger.Crit("markUnexecutedBatches - parent batch Header is nil!")
		}
	}
	logger.Info("markUnexecutedBatches - exited loop successfully", "totalIterations", loopIteration)
	return nil
}
