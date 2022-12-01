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
	db   *db.DB
	node int
}

func NewBatchManager(db *db.DB, node int) *BatchManager {
	return &BatchManager{
		db:   db,
		node: node,
	}
}

// BatchesMissingError indicates that when processing new batches, one or more batches were missing from the database.
type BatchesMissingError struct {
	EarliestMissingBatch *big.Int
}

// todo - joel - describe
type L1BlockMissingError struct {
	BatchNumber *big.Int
}

func (b BatchesMissingError) Error() string {
	return fmt.Sprintf("missing batches; earliest missing batch is %d", b.EarliestMissingBatch)
}

// StoreBatches stores the provided batches. If there are missing batches in the chain, it returns a
// `BatchesMissingError`.
func (b *BatchManager) StoreBatches(batches []*common.ExtBatch) error {
	for _, batch := range batches {
		// If we don't have the parent batch and this is not the genesis batch, we request the missing batches.
		if batch.Header.Number.Uint64() != common.L2GenesisHeight { //nolint:nestif
			parentBatchNumber := big.NewInt(0).Sub(batch.Header.Number, big.NewInt(1))
			if _, err := b.db.GetBatchHash(parentBatchNumber); err != nil {
				if errors.Is(err, errutil.ErrNotFound) {
					earliestMissingBatch, err := b.findEarliestMissingBatch(parentBatchNumber)
					if err != nil {
						return fmt.Errorf("could not calculate earliest missing batch. Cause: %w", err)
					}
					return &BatchesMissingError{earliestMissingBatch}
				}
				return fmt.Errorf("could not retrieve parent batch. Cause: %w", err)
			}
		}

		// If we don't have the L1 block, we request the batch to be resent to gain time.
		// TODO - #718 - Find a more efficient solution, rather than forcing the sequencer to resend.
		if _, err := b.db.GetBlockHeader(batch.Header.L1Proof); err != nil {
			if errors.Is(err, errutil.ErrNotFound) {
				if b.node == 1 {
					println(fmt.Sprintf("jjj waiting for block %s for rollup %d", batch.Header.L1Proof.Hex(), batch.Header.Number))
				}
				earliestMissingBatch, err := b.findEarliestMissingBatch(batch.Header.Number)
				if err != nil {
					return fmt.Errorf("could not calculate earliest missing batch. Cause: %w", err)
				}
				return &BatchesMissingError{earliestMissingBatch}
			}
			println("jjj could not query for block", batch.Header.L1Proof.Hex())
			return fmt.Errorf("could not retrieve batch's L1 block. Cause: %w", err)
		}

		if b.node == 1 {
			println("jjj got block", batch.Header.L1Proof.Hex())
		}

		// We store the batch.
		if err := b.db.AddBatchHeader(batch); err != nil {
			return fmt.Errorf("could not store batch header. Cause: %w", err)
		}
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
