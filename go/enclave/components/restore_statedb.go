package components

import (
	"context"
	"errors"
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/enclave/storage"
)

// this function looks at the batch chain and makes sure the resulting stateDB snapshots are available, replaying them if needed
// (if there had been a clean shutdown and all stateDB data was persisted this should do nothing)
func syncExecutedBatchesWithEVMStateDB(ctx context.Context, storage storage.Storage, logger gethlog.Logger) (*common.BatchHeader, error) {
	headBatch, err := storage.FetchHeadBatchHeader(context.Background())
	if err != nil {
		if !errors.Is(err, errutil.ErrNotFound) {
			return nil, fmt.Errorf("failed to read head batch header: %w", err)
		}
		return nil, nil
	}

	// `headBatch` variable will eventually be the latest batch for which we are able to produce a StateDB
	// - we will then set that as the head of the L2 so that this node can rebuild its missing state
	// loop backwards building a slice of all batches that don't have cached stateDB data available
	for !stateDBAvailableForBatch(storage, headBatch.Root, logger) {
		logger.Info("StateDB not available for batch, rolling back", "batchHash", headBatch.Hash(), "sequencerOrderNo", headBatch.SequencerOrderNo)
		err = storage.MarkBatchAsUnexecuted(ctx, headBatch.SequencerOrderNo)
		if err != nil {
			return nil, fmt.Errorf("unable to mark batch as unexecuted - %w", err)
		}
		if headBatch.Number.Uint64() == common.L2GenesisHeight {
			// no more parents to check, replaying from genesis
			break
		}

		headBatch, err = storage.FetchBatchHeader(ctx, headBatch.ParentHash)
		if err != nil {
			return nil, fmt.Errorf("unable to fetch previous batch while rolling back to stable state - %w", err)
		}
	}
	return headBatch, nil
}

// The enclave caches a stateDB instance against each batch hash, this is the input state when producing the following
// batch in the chain and is used to query state at a certain height.
//
// This method checks if the stateDB data is available for a given batch hash (so it can be restored if not)
func stateDBAvailableForBatch(storage storage.Storage, root gethcommon.Hash, logger gethlog.Logger) bool {
	_, err := storage.OpenTrie(root)
	if err != nil {
		logger.Warn("failed to open state trie for batch", "batch.root", root, "err", err)
	}
	return err == nil
}
