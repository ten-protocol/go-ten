package batchmanager

import (
	"errors"
	"fmt"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/host/db"
	"math/big"
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
	CurrentHeadBatch *gethcommon.Hash // Our view of the current head batch, or nil if no batches are stored.
}

func (b BatchesMissingError) Error() string {
	return fmt.Sprintf("missing batches; earliest missing batch is %s", b.CurrentHeadBatch.Hex())
}

// StoreBatches stores the provided batches. If there are missing batches in the chain, it returns a
// `BatchesMissingError`.
func (b *BatchManager) StoreBatches(batches []*common.ExtBatch, nodeId uint64) error { //nolint:gocognit
	//if nodeId == 2 {
	print(fmt.Sprintf("jjj received the following batches on node %d: ", nodeId))
	for _, batch := range batches {
		print(fmt.Sprintf("%d, ", batch.Header.Number))
	}
	println()
	//}

	for _, batch := range batches {
		// We do not have the corresponding L1 block stored yet, so we discard the batch. We'll request the
		// batch later as part of catch-up, once we have the L1 block stored.
		_, err := b.db.GetBlockHeader(batch.Header.L1Proof)
		if err != nil {
			if errors.Is(err, errutil.ErrNotFound) {
				if nodeId == 2 {
					println(fmt.Sprintf("jjj skipping batch %d on node %d because don't have block yet", batch.Header.Number.Uint64(), nodeId))
				}
				// todo - joel - do this, or re-request the batch?
				return nil
				//return &BatchesMissingError{batch.Header.Number}
			}
			return fmt.Errorf("could not retrieve L1 block for batch. Cause: %w", err)
		}

		if nodeId == 2 {
			println(fmt.Sprintf("jjj working on batch %d on node %d because we have block. Hash: %s; parent hash: %s",
				batch.Header.Number.Uint64(), nodeId, batch.Hash(), batch.Header.ParentHash))
		}

		_, err = b.db.GetBatch(batch.Header.ParentHash)

		// We have stored the batch's parent, or this batch is the genesis batch, so we store the batch.
		if err == nil || batch.Header.Number.Uint64() == common.L2GenesisHeight {
			if nodeId == 2 {
				println(fmt.Sprintf("jjj storing batch %d on node %d", batch.Header.Number.Uint64(), nodeId))
			}
			err = b.db.AddBatchHeader(batch)
			if err != nil {
				return fmt.Errorf("could not store batch header. Cause: %w", err)
			}
			continue
		}

		// If we could not find the parent, we have at least one missing batch.
		// todo - joel - update descriptions around here
		if errors.Is(err, errutil.ErrNotFound) {
			if nodeId == 2 {
				println(fmt.Sprintf("jjj could not store batch %d on node %d; requesting", batch.Header.Number.Uint64(), nodeId))
			}
			headBatchHeader, err := b.db.GetHeadBatchHeader()
			if err != nil {
				if !errors.Is(err, errutil.ErrNotFound) {
					return fmt.Errorf("could not retrieve head batch header")
				}
				return &BatchesMissingError{&gethcommon.Hash{}}
			}

			headBatchHash := headBatchHeader.Hash()
			return &BatchesMissingError{&headBatchHash}
		}

		return fmt.Errorf("could not retrieve batch header. Cause: %w", err)
	}

	return nil
}

// GetBatches retrieves the batches matching the batch request from the host's database.
func (b *BatchManager) GetBatches(batchRequest *common.BatchRequest) ([]*common.ExtBatch, error) {
	var currentBatch *common.ExtBatch

	// todo - joel - big tidy up
	if (*batchRequest.CurrentHeadBatch == gethcommon.Hash{}) {
		var err error
		batchHash, err := b.db.GetBatchHash(big.NewInt(0))
		if err != nil {
			panic("todo - joel")
		}
		currentBatch, err = b.db.GetBatch(*batchHash)
		if err != nil {
			panic("todo - joel")
		}
	} else {

		requesterHead := batchRequest.CurrentHeadBatch
		latestCanonicalAncestor, err := b.latestCanonAncestor(requesterHead)
		if err != nil {
			panic("todo - joel")
		}

		currentBatch = latestCanonicalAncestor
	}

	batches := []*common.ExtBatch{currentBatch}

	for {
		currentBatchNumber := big.NewInt(0).Add(currentBatch.Header.Number, big.NewInt(1))

		batchHash, err := b.db.GetBatchHash(currentBatchNumber)
		if err != nil {
			// We have reached the latest batch.
			if errors.Is(err, errutil.ErrNotFound) {
				break
			}
			return nil, fmt.Errorf("could not retrieve batch hash for batch number %d. Cause: %w", currentBatchNumber, err)
		}

		batch, err := b.db.GetBatch(*batchHash)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve batch for batch hash %s. Cause: %w", batchHash, err)
		}
		batches = append(batches, batch)

		currentBatch = batch
	}

	print("jjj sending following batches: ")
	for _, batch := range batches {
		print(fmt.Sprintf("%d, ", batch.Header.Number))
	}
	print("\n")

	// todo - right batches are sent, but never appear to be worked on

	return batches, nil
}

// todo - joel - describe
func (b *BatchManager) latestCanonAncestor(batchHash *gethcommon.Hash) (*common.ExtBatch, error) {
	batch, err := b.db.GetBatch(*batchHash)
	if err != nil {
		return nil, err
	}

	canonAtSameHeight, err := b.db.GetBatchHash(batch.Header.Number)
	if err != nil {
		return nil, err
	}

	if batch.Hash() != *canonAtSameHeight {
		return b.latestCanonAncestor(&batch.Header.ParentHash)
	}

	return batch, nil
}
