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

// todo - joel - comment
func (b *BatchManager) RetrieveBatches(batchRequest *common.BatchRequest) ([]*common.ExtBatch, error) {
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

// todo - joel - comment
func (b *BatchManager) SortBatches(batches []*common.ExtBatch) {
	sort.Slice(batches, func(i, j int) bool {
		return batches[i].Header.Number.Cmp(batches[i].Header.Number) < 0
	})
}

// todo - joel - comment
func (b *BatchManager) CheckForGapsAndDupes(batches []*common.ExtBatch) error {
	for idx := 0; idx < len(batches)-1; idx++ {
		i := batches[idx]
		j := batches[idx+1]

		numberGap := big.NewInt(0).Sub(j.Header.Number, i.Header.Number)
		gapIsZero := numberGap.Cmp(big.NewInt(0)) == 0
		gapIsMoreThanOne := numberGap.Cmp(big.NewInt(1)) != 0

		if gapIsZero {
			return fmt.Errorf("duplicates in set of batches to process")
		}
		if gapIsMoreThanOne {
			return fmt.Errorf("gaps in chain of set of batches to process")
		}
	}
	return nil
}

// todo - joel - comment
func (b *BatchManager) CreateBatchRequest(batch *common.ExtBatch) (*common.BatchRequest, error) {
	var earliestMissingBatch *big.Int
	parentBatchNumber := big.NewInt(0).Sub(batch.Header.Number, big.NewInt(1))
	for {
		// If we have reached the head of the chain, break.
		if parentBatchNumber.Int64() < int64(common.L2GenesisHeight) {
			break
		}

		_, err := b.db.GetBatchHash(parentBatchNumber)
		if err != nil {
			// If the batch is not found, we update the variable tracking the earliest missing batch.
			if errors.Is(err, errutil.ErrNotFound) {
				earliestMissingBatch = parentBatchNumber
				parentBatchNumber = big.NewInt(0).Sub(parentBatchNumber, big.NewInt(1))
				continue
			}
			return nil, fmt.Errorf("could not get batch hash by number. Cause: %w", err)
		}

		// If there was no error, we have reach a stored batch.
		break
	}

	if earliestMissingBatch == nil {
		// There are no missing batches to request.
		return nil, nil //nolint:nilnil
	}

	return &common.BatchRequest{From: earliestMissingBatch, To: batch.Header.Number}, nil
}

// todo - joel - comment
func (b *BatchManager) StoreBatches(batches []*common.ExtBatch) error {
	for _, batch := range batches {
		err := b.db.AddBatchHeader(batch)
		if err != nil {
			return fmt.Errorf("could not store batch header. Cause: %w", err)
		}
	}
	return nil
}
