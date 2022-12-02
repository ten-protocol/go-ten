package batchmanager

import (
	"errors"
	"fmt"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/host/db"
)

// ErrBatchesMissing indicates that when processing new batches, one or more batches were missing from the database.
var ErrBatchesMissing = errors.New("one or more batches in L2 chain were missing")

// BatchManager handles the creation and processing of batches for the host.
type BatchManager struct {
	db *db.DB
}

func NewBatchManager(db *db.DB) *BatchManager {
	return &BatchManager{
		db: db,
	}
}

// StoreBatches stores the provided batches. If there are missing batches in the chain, it returns a
// `ErrBatchesMissing`.
func (b *BatchManager) StoreBatches(batches []*common.ExtBatch) error {
	for _, batch := range batches {
		// If we have stored the batch's parent, or this batch is the genesis batch, we store the batch.
		_, err := b.db.GetBatch(batch.Header.ParentHash)
		if err == nil || batch.Header.Number.Uint64() == common.L2GenesisHeight {
			err = b.db.AddBatchHeader(batch)
			if err != nil {
				return fmt.Errorf("could not store batch header. Cause: %w", err)
			}
			continue
		}

		// If we could not find the parent, we return an `ErrBatchesMissing`.
		if errors.Is(err, errutil.ErrNotFound) {
			return ErrBatchesMissing
		}

		return fmt.Errorf("could not retrieve batch header. Cause: %w", err)
	}

	return nil
}

// CreateBatchRequest creates a request for missing batches, which contains our address and our view of the canonical
// L2 head.
func (b *BatchManager) CreateBatchRequest(nodeP2PAddress string) (*common.BatchRequest, error) {
	var headBatchHash *gethcommon.Hash

	// We retrieve our view of the canonical L2 head.
	currentHeadBatch, err := b.db.GetHeadBatchHeader()
	if err != nil {
		if !errors.Is(err, errutil.ErrNotFound) {
			return nil, fmt.Errorf("could not retrieve head batch. Cause: %w", err)
		}
		headBatchHash = &gethcommon.Hash{}
	} else {
		hash := currentHeadBatch.Hash()
		headBatchHash = &hash
	}

	return &common.BatchRequest{
		Requester:        nodeP2PAddress,
		CurrentHeadBatch: headBatchHash,
	}, nil
}

// GetBatches retrieves the batches from the host's database matching the batch request.
func (b *BatchManager) GetBatches(batchRequest *common.BatchRequest) ([]*common.ExtBatch, error) {
	// We handle the case where the requester has no batches stored at all.
	requesterHeadBatch := batchRequest.CurrentHeadBatch
	if (*batchRequest.CurrentHeadBatch == gethcommon.Hash{}) {
		var err error
		requesterHeadBatch, err = b.db.GetBatchHash(big.NewInt(0))
		if err != nil {
			return nil, fmt.Errorf("could not retrieve zero'th batch hash. Cause: %w", err)
		}
	}

	// We determine the latest canonical ancestor to start sending batches from.
	canonicalAncestor, err := b.latestCanonicalAncestor(requesterHeadBatch)
	if err != nil {
		return nil, fmt.Errorf("could not determine latest canonical ancestor. Cause: %w", err)
	}
	batchesToSend := []*common.ExtBatch{canonicalAncestor}
	currentBatchNumber := canonicalAncestor.Header.Number

	// We gather the batches from the canonical chain.
	for {
		currentBatchNumber = big.NewInt(0).Add(currentBatchNumber, big.NewInt(1))

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
		batchesToSend = append(batchesToSend, batch)
	}

	return batchesToSend, nil
}

// Determines the latest canonical ancestor between the provided batch hash and the sequencer's canonical chain.
func (b *BatchManager) latestCanonicalAncestor(batchHash *gethcommon.Hash) (*common.ExtBatch, error) {
	batch, err := b.db.GetBatch(*batchHash)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve batch. Cause: %w", err)
	}

	canonicalBatchHashAtSameHeight, err := b.db.GetBatchHash(batch.Header.Number)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve canonical batch hash. Cause: %w", err)
	}

	// If the batch's hash does not match the canonical batch's hash at the same height, we need to keep walking back.
	if batch.Hash() != *canonicalBatchHashAtSameHeight {
		return b.latestCanonicalAncestor(&batch.Header.ParentHash)
	}
	return batch, nil
}
