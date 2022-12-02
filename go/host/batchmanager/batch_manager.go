package batchmanager

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/host/db"
)

// BatchManager handles the creation and processing of batches for the host.
type BatchManager struct {
	db *db.DB
}

func NewBatchManager(db *db.DB) *BatchManager {
	return &BatchManager{
		db: db,
	}
}

// BatchesMissingError indicates that when processing new batches, one or more batches were missing from the database.
type BatchesMissingError struct {
	EarliestMissingBatch *big.Int
}

func (b BatchesMissingError) Error() string {
	return fmt.Sprintf("missing batches; earliest missing batch is %d", b.EarliestMissingBatch)
}

// StoreBatches stores the provided batches. If there are missing batches in the chain, it returns a
// `BatchesMissingError`.
func (b *BatchManager) StoreBatches(batches []*common.ExtBatch, nodeId uint64) error {
	if nodeId == 2 {
		print(fmt.Sprintf("jjj received the following batches on node %d: ", nodeId))
		for _, batch := range batches {
			print(fmt.Sprintf("%d", batch.Header.Number))
		}
		println()
	}

	for _, batch := range batches {
		_, err := b.db.GetBlockHeader(batch.Header.L1Proof)
		if err != nil {
			if errors.Is(err, errutil.ErrNotFound) {
				if nodeId == 2 {
					println(fmt.Sprintf("jjj skipping batch %d on node %d because don't have block yet", batch.Header.Number.Uint64(), nodeId))
				}
				// We do not have the corresponding L1 block stored yet, so we discard the batch. We'll request the
				// batch later as part of catch-up, once we have the L1 block stored.
				return &BatchesMissingError{batch.Header.Number} // todo - joel - think about this
			}
			return fmt.Errorf("could not retrieve L1 block for batch. Cause: %w", err)
		}

		if nodeId == 2 {
			println(fmt.Sprintf("jjj working on batch %d on node %d because we have block", batch.Header.Number.Uint64(), nodeId))
		}

		_, err = b.db.GetBatch(batch.Header.ParentHash)

		// We have stored the batch's parent, or this batch is the genesis batch, so we store the batch.
		if err == nil || batch.Header.Number.Uint64() == common.L2GenesisHeight {
			err = b.db.AddBatchHeader(batch)
			if err != nil {
				return fmt.Errorf("could not store batch header. Cause: %w", err)
			}
			continue
		}

		// If we could not find the parent, we have at least one missing batch.
		if errors.Is(err, errutil.ErrNotFound) {
			parentBatchNumber := big.NewInt(0).Sub(batch.Header.Number, big.NewInt(1))
			// This is not foolproof. We may find that we have a batch stored for a given number, but unbeknownst to
			// us, it is for a different fork. This means that we may have to go through several rounds of requests,
			// getting only one additional link in the chain each time.
			earliestMissingBatch, err := b.findEarliestMissingBatch(parentBatchNumber)
			if err != nil {
				return fmt.Errorf("could not calculate earliest missing batch. Cause: %w", err)
			}
			if nodeId == 2 {
				println(fmt.Sprintf("jjj requesting batches from %d on node %d", batch.Header.Number.Uint64(), nodeId))
			}
			return &BatchesMissingError{earliestMissingBatch}
		}

		return fmt.Errorf("could not retrieve batch header. Cause: %w", err)
	}

	return nil
}

// GetBatches retrieves the batches matching the batch request from the host's database.
func (b *BatchManager) GetBatches(batchRequest *common.BatchRequest) ([]*common.ExtBatch, error) {
	var batches []*common.ExtBatch

	currentBatch := batchRequest.EarliestMissingBatch
	for {
		batchHash, err := b.db.GetBatchHash(currentBatch)
		if err != nil {
			// We have reached the latest batch.
			if errors.Is(err, errutil.ErrNotFound) {
				break
			}
			return nil, fmt.Errorf("could not retrieve batch hash for batch number %d. Cause: %w", currentBatch, err)
		}

		batch, err := b.db.GetBatch(*batchHash)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve batch for batch hash %s. Cause: %w", batchHash, err)
		}
		batches = append(batches, batch)

		currentBatch = big.NewInt(0).Add(currentBatch, big.NewInt(1))
	}

	return batches, nil
}

// Starting from the provided number, we walk the chain batch until we find a batch number against which we have stored
// a batch.
func (b *BatchManager) findEarliestMissingBatch(startBatchNumber *big.Int) (*big.Int, error) {
	earliestMissingBatch := startBatchNumber

	for {
		// If we have reached the head of the chain, break.
		if earliestMissingBatch.Int64() <= int64(common.L2GenesisHeight) {
			return earliestMissingBatch, nil
		}

		// We check whether the batch is stored. If there was no error, we have reached a stored batch.
		_, err := b.db.GetBatchHash(earliestMissingBatch)
		if err == nil {
			return earliestMissingBatch, nil
		}

		// If the batch is not found, we update the variable tracking the earliest missing batch.
		if errors.Is(err, errutil.ErrNotFound) {
			earliestMissingBatch = big.NewInt(0).Sub(earliestMissingBatch, big.NewInt(1))
			continue
		}
		return nil, fmt.Errorf("could not get batch hash by number. Cause: %w", err)
	}
}
