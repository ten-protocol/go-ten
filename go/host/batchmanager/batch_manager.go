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

const (
	// A limit on the number of batches that can be served in a single catch-up request.
	maxBatchesPerRequest = 10
)

// BatchManager handles the creation and processing of batches for the host.
type BatchManager struct {
	db               *db.DB
	p2pPublicAddress string
}

func NewBatchManager(db *db.DB, p2pPublicAddress string) *BatchManager {
	return &BatchManager{
		db:               db,
		p2pPublicAddress: p2pPublicAddress,
	}
}

// IsParentStored indicates whether the batch has already been stored. If not, it returns the batch request to send to
// the sequencer.
func (b *BatchManager) IsParentStored(batch *common.ExtBatch) (bool, *common.BatchRequest, error) {
	// If this is the genesis batch, we don't need to request the parent.
	if batch.Header.Number.Uint64() == common.L2GenesisHeight {
		return true, nil, nil
	}

	_, err := b.db.GetBatch(batch.Header.ParentHash)
	if err != nil {
		// The parent is missing.
		if errors.Is(err, errutil.ErrNotFound) {
			batchRequest, err := b.createBatchRequest()
			return false, batchRequest, err
		}
		return false, nil, fmt.Errorf("could not retrieve batch header. Cause: %w", err)
	}
	return true, nil, nil
}

type batchResolverFunc = func(gethcommon.Hash) (*common.ExtBatch, error)

// GetBatches retrieves the batches from the host's database matching the batch request.
func (b *BatchManager) GetBatches(batchRequest *common.BatchRequest, searchMissingBatch batchResolverFunc) ([]*common.ExtBatch, error) {
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
	firstBatch, err := b.latestCanonicalAncestor(requesterHeadBatch)
	if err != nil {
		return nil, fmt.Errorf("could not determine latest canonical ancestor. Cause: %w", err)
	}
	var batchesToSend []*common.ExtBatch

	// We find the batch we want to send up to - either the head batch, or the max number of requested batches,
	// whichever is lower.
	lastBatch, err := b.db.GetHeadBatchHeader()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve head batch header. Cause: %w", err)
	}
	if lastBatch.Number.Int64()-firstBatch.Header.Number.Int64() > maxBatchesPerRequest {
		lastBatchNumber := big.NewInt(0).Add(firstBatch.Header.Number, big.NewInt(maxBatchesPerRequest))
		batchHash, err := b.db.GetBatchHash(lastBatchNumber)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve batch hash for batchHeight=%d. Cause: %w", lastBatchNumber, err)
		}
		lastBatch, err = b.db.GetBatchHeader(*batchHash)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve batch for batch hash %s. Cause: %w", batchHash, err)
		}
	}

	// We gather the batches by walking backwards from the final batch.
	currentBatch, err := b.db.GetBatch(lastBatch.Hash())
	if err != nil {
		return nil, fmt.Errorf("could not retrieve batch. Cause: %w", err)
	}
	for {
		if currentBatch.Hash().Hex() == firstBatch.Hash().Hex() {
			break
		}
		batchesToSend = append(batchesToSend, currentBatch)
		hashToLookFor := currentBatch.Header.ParentHash
		currentBatch, err = b.db.GetBatch(hashToLookFor)
		if err != nil {
			currentBatch, err = searchMissingBatch(hashToLookFor)
			if err != nil {
				return nil, fmt.Errorf("could not retrieve batch header. Cause: %w", err)
			}
		}
	}
	batchesToSend = append(batchesToSend, firstBatch)

	// We reverse the batches so that the recipient can process them in order.
	for i, j := 0, len(batchesToSend)-1; i < j; i, j = i+1, j-1 {
		batchesToSend[i], batchesToSend[j] = batchesToSend[j], batchesToSend[i]
	}
	return batchesToSend, nil
}

// Creates a request for missing batches, which contains our address and our view of the canonical L2 head.
func (b *BatchManager) createBatchRequest() (*common.BatchRequest, error) {
	var headBatchHash gethcommon.Hash

	// We retrieve our view of the canonical L2 head.
	currentHeadBatch, err := b.db.GetHeadBatchHeader()
	if err != nil {
		if !errors.Is(err, errutil.ErrNotFound) {
			return nil, fmt.Errorf("could not retrieve head batch. Cause: %w", err)
		}
		headBatchHash = gethcommon.Hash{}
	} else {
		headBatchHash = currentHeadBatch.Hash()
	}

	return &common.BatchRequest{
		Requester:        b.p2pPublicAddress,
		CurrentHeadBatch: &headBatchHash,
	}, nil
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
