package batchmanager

import (
	"errors"
	"fmt"
	"math/big"
	"sort"

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

// IsMissingBatches retruns a bool indicating whether any historical batches are missing, given the state of the host's
// database and the batches provided. If batches are missing, it creates a corresponding batch request.
func (b *BatchManager) IsMissingBatches(batches []*common.ExtBatch) (bool, *common.BatchRequest, error) {
	// We sort the batches, then check for duplicates or gaps. If we don't identify gaps first, there's a risk that
	// we won't request sufficient missing batches (e.g. we have `[0,1]` in our DB, and receive `[3,4,6]`; it is
	// important that we don't "see" the `3` and fail to request the `5`).
	b.sortBatchesByNumber(batches)
	hasGapsOrDupes, err := b.checkForGapsAndDupes(batches)
	if hasGapsOrDupes {
		return false, nil, err
	}

	earliestReceivedBatch := batches[0]
	latestReceivedBatch := batches[len(batches)-1]

	var earliestMissingBatch *big.Int
	parentBatchNumber := big.NewInt(0).Sub(earliestReceivedBatch.Header.Number, big.NewInt(1))
	for {
		// If we have reached the head of the chain, break.
		if parentBatchNumber.Int64() < int64(common.L2GenesisHeight) {
			break
		}

		_, err = b.db.GetBatchHash(parentBatchNumber)
		if err != nil {
			// If the batch is not found, we update the variable tracking the earliest missing batch.
			if errors.Is(err, errutil.ErrNotFound) {
				earliestMissingBatch = parentBatchNumber
				parentBatchNumber = big.NewInt(0).Sub(parentBatchNumber, big.NewInt(1))
				continue
			}
			return false, nil, fmt.Errorf("could not get batch hash by number. Cause: %w", err)
		}

		// If there was no error, we have reached a stored batch.
		break
	}

	if earliestMissingBatch == nil {
		// There are no missing batches to request.
		return false, nil, nil
	}

	return true, &common.BatchRequest{From: earliestMissingBatch, To: latestReceivedBatch.Header.Number}, nil
}

// GetBatches retrieves the batches matching the batch request from the host's database.
func (b *BatchManager) GetBatches(batchRequest *common.BatchRequest) ([]*common.ExtBatch, error) {
	var batches []*common.ExtBatch
	currentBatchToRetrieve := batchRequest.From
	for currentBatchToRetrieve.Cmp(batchRequest.To) != 1 {
		batchHash, err := b.db.GetBatchHash(currentBatchToRetrieve)
		if err != nil {
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

// StoreBatches stores the provided batches in the host's database.
func (b *BatchManager) StoreBatches(batches []*common.ExtBatch) error {
	for _, batch := range batches {
		err := b.db.AddBatchHeader(batch)
		if err != nil {
			return fmt.Errorf("could not store batch header. Cause: %w", err)
		}
	}
	return nil
}

// Sorts a list of batches by batch number.
func (b *BatchManager) sortBatchesByNumber(batches []*common.ExtBatch) {
	sort.Slice(batches, func(i, j int) bool {
		return batches[i].Header.Number.Cmp(batches[i].Header.Number) < 0
	})
}

// Indicates whether a list of batches sorted by number has any gaps or duplicates.
func (b *BatchManager) checkForGapsAndDupes(batches []*common.ExtBatch) (bool, error) {
	for idx := 0; idx < len(batches)-1; idx++ {
		i := batches[idx]
		j := batches[idx+1]

		numberGap := big.NewInt(0).Sub(j.Header.Number, i.Header.Number)
		gapIsZero := numberGap.Cmp(big.NewInt(0)) == 0
		gapIsMoreThanOne := numberGap.Cmp(big.NewInt(1)) != 0

		if gapIsZero {
			return true, fmt.Errorf("duplicates in set of batches to process")
		}
		if gapIsMoreThanOne {
			return true, fmt.Errorf("gaps in chain of set of batches to process")
		}
	}
	return false, nil
}
