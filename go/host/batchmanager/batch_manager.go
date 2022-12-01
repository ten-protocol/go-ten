package batchmanager

import (
	"errors"
	"fmt"

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

// BatchMissingError indicates that when processing new batches, a batch was missing from the database.
type BatchMissingError struct {
	MissingBatch *common.L2RootHash
}

func (b BatchMissingError) Error() string {
	return fmt.Sprintf("missing batch %d", b.MissingBatch)
}

// StoreBatch stores the provided batch. If there is a batch missing in the chain, it returns a `BatchMissingError`.
// There is no way to identify more than one missing batch in the chain - we cannot go by the batch numbers we have
// stored, since these batches may have been stored as part of a discarded fork.
func (b *BatchManager) StoreBatch(batch *common.ExtBatch) error {
	// We store the batch.
	err := b.db.AddBatchHeader(batch)
	if err != nil {
		return fmt.Errorf("could not store batch header. Cause: %w", err)
	}

	// If this is the genesis batch, there can be no missing parent.
	if batch.Header.Number.Uint64() == common.L2GenesisHeight {
		return nil
	}

	// We check if we are missing the parent, in which case we request it.
	if _, err = b.db.GetBatch(batch.Header.ParentHash); err != nil {
		// If we could not find the parent, we return a `BatchMissingError`.
		if errors.Is(err, errutil.ErrNotFound) {
			return &BatchMissingError{&batch.Header.ParentHash}
		}
		return fmt.Errorf("could not retrieve batch's parent. Cause: %w", err)
	}

	return nil
}

// GetBatch retrieves the batch matching the batch request from the host's database.
func (b *BatchManager) GetBatch(batchRequest *common.BatchRequest) (*common.ExtBatch, error) {
	batch, err := b.db.GetBatch(*batchRequest.MissingBatch)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve batch for batch hash %s. Cause: %w", batchRequest.MissingBatch, err)
	}

	return batch, nil
}
