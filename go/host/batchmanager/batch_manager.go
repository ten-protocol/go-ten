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
func (b *BatchManager) StoreBatches(batches []*common.ExtBatch) error {
	for _, batch := range batches {
		parentBatchNumber := big.NewInt(0).Sub(batch.Header.Number, big.NewInt(1))
		_, err := b.db.GetBatchHash(parentBatchNumber)

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
			earliestMissingBatch, err := b.findEarliestMissingBatch(parentBatchNumber)
			if err != nil {
				return fmt.Errorf("could not calculate earliest missing batch. Cause: %w", err)
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

	currentBatchToRetrieve := batchRequest.EarliestMissingBatch
	for {
		batchHash, err := b.db.GetBatchHash(currentBatchToRetrieve)
		if err != nil {
			// We have reached the latest batch. Our work is complete.
			if errors.Is(err, errutil.ErrNotFound) {
				break
			}
			return nil, fmt.Errorf("could not retrieve batch hash for batch number %d. Cause: %w", currentBatchToRetrieve, err)
		}
		batch, err := b.db.GetBatch(*batchHash)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve batch for batch hash %s. Cause: %w", batchHash, err)
		}
		batches = append(batches, batch)
		currentBatchToRetrieve = big.NewInt(0).Add(currentBatchToRetrieve, big.NewInt(1))
	}

	return batches, nil
}

// Starting from the provided number, we walk the chain batch until we find a stored batch.
func (b *BatchManager) findEarliestMissingBatch(startBatchNumber *big.Int) (*big.Int, error) {
	earliestMissingBatch := startBatchNumber

	for {
		// If we have reached the head of the chain, break.
		if earliestMissingBatch.Int64() <= int64(common.L2GenesisHeight) {
			return earliestMissingBatch, nil
		}

		// We check whether the batch is stored.
		_, err := b.db.GetBatchHash(earliestMissingBatch)
		// If there was no error, we have reached a stored batch.
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
