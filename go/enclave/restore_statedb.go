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
func syncExecutedBatchesWithEVMStateDB(ctx context.Context, storage storage.Storage, registry components.BatchRegistry, logger gethlog.Logger) error {
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
		err = markUnexecutedBatches(ctx, storage, registry, logger)
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
	_, err := registry.GetBatchState(ctx, gethrpc.BlockNumberOrHash{BlockHash: &hash})
	return err == nil
}

// markUnexecutedBatches marks the batches for which the statedb is missing as un-executed
func markUnexecutedBatches(ctx context.Context, storage storage.Storage, registry components.BatchRegistry, logger gethlog.Logger) error {
	// `currentBatch` variable will eventually be the latest batch for which we are able to produce a StateDB
	// - we will then set that as the head of the L2 so that this node can rebuild its missing state
	currentBatch, err := storage.FetchBatchBySeqNo(ctx, registry.HeadBatchSeq().Uint64())
	if err != nil {
		return fmt.Errorf("no head batch found in DB but expected to replay batches - %w", err)
	}
	// loop backwards building a slice of all batches that don't have cached stateDB data available
	for !stateDBAvailableForBatch(ctx, registry, currentBatch.Hash()) {
		err = storage.MarkBatchAsUnexecuted(ctx, currentBatch.SeqNo())
		if err != nil {
			return fmt.Errorf("unable to mark batch as unexecuted - %w", err)
		}
		if currentBatch.NumberU64() == common.L2GenesisHeight {
			// no more parents to check, replaying from genesis
			break
		}
		currentBatch, err = storage.FetchBatch(ctx, currentBatch.Header.ParentHash)
		if err != nil {
			return fmt.Errorf("unable to fetch previous batch while rolling back to stable state - %w", err)
		}
	}
	return nil
}
