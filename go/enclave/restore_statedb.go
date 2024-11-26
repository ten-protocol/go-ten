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
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"github.com/ten-protocol/go-ten/go/enclave/genesis"
	"github.com/ten-protocol/go-ten/go/enclave/storage"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

// this function looks at the batch chain and makes sure the resulting stateDB snapshots are available, replaying them if needed
// (if there had been a clean shutdown and all stateDB data was persisted this should do nothing)
func restoreStateDBCache(ctx context.Context, storage storage.Storage, registry components.BatchRegistry, producer components.BatchExecutor, gen *genesis.Genesis, logger gethlog.Logger) error {
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
		err = replayBatchesToValidState(ctx, storage, registry, producer, gen, logger)
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

// replayBatchesToValidState is used to repopulate the stateDB cache with data from persisted batches. Two step process:
// 1. step backwards from head batch until we find a batch that is already in stateDB cache, builds list of batches to replay
// 2. iterate that list of batches from the earliest, process the transactions to calculate and cache the stateDB
// todo (#1416) - get unit test coverage around this (and L2 Chain code more widely, see ticket #1416 )
func replayBatchesToValidState(ctx context.Context, storage storage.Storage, registry components.BatchRegistry, batchExecutor components.BatchExecutor, gen *genesis.Genesis, logger gethlog.Logger) error {
	// this slice will be a stack of batches to replay as we walk backwards in search of latest valid state
	// todo - consider capping the size of this batch list using FIFO to avoid memory issues, and then repeating as necessary
	var batchesToReplay []*core.Batch
	// `batchToReplayFrom` variable will eventually be the latest batch for which we are able to produce a StateDB
	// - we will then set that as the head of the L2 so that this node can rebuild its missing state
	batchToReplayFrom, err := storage.FetchBatchBySeqNo(ctx, registry.HeadBatchSeq().Uint64())
	if err != nil {
		return fmt.Errorf("no head batch found in DB but expected to replay batches - %w", err)
	}
	// loop backwards building a slice of all batches that don't have cached stateDB data available
	for !stateDBAvailableForBatch(ctx, registry, batchToReplayFrom.Hash()) {
		batchesToReplay = append(batchesToReplay, batchToReplayFrom)
		if batchToReplayFrom.NumberU64() == 0 {
			// no more parents to check, replaying from genesis
			break
		}
		batchToReplayFrom, err = storage.FetchBatch(ctx, batchToReplayFrom.Header.ParentHash)
		if err != nil {
			return fmt.Errorf("unable to fetch previous batch while rolling back to stable state - %w", err)
		}
	}
	logger.Info("replaying batch data into stateDB cache", "fromBatch", batchesToReplay[len(batchesToReplay)-1].NumberU64(),
		"toBatch", batchesToReplay[0].NumberU64())
	// loop through the slice of batches without stateDB data to cache the state (loop in reverse because slice is newest to oldest)
	for i := len(batchesToReplay) - 1; i >= 0; i-- {
		batch := batchesToReplay[i]

		// if genesis batch then create the genesis state before continuing on with remaining batches
		if batch.NumberU64() == 0 {
			err := gen.CommitGenesisState(storage)
			if err != nil {
				return err
			}
			continue
		}

		// calculate the stateDB after this batch and store it in the cache
		_, err := batchExecutor.ExecuteBatch(ctx, batch)
		if err != nil {
			return err
		}
	}

	return nil
}
