package batchmanager

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/host/db"
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
	// If this is the genesis block, we don't need to request the parent.
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
	// Fetching the head batch upfront avoids a potential infinite loop if batches are produced very fast.
	headBatchHeader, err := b.db.GetHeadBatchHeader()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve head batch header. Cause: %w", err)
	}

	// We gather the batches from the canonical chain.
	for {
		currentBatchNumber = big.NewInt(0).Add(currentBatchNumber, big.NewInt(1))

		batchHash, err := b.db.GetBatchHash(currentBatchNumber)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve batch hash for batch number %d. Cause: %w", currentBatchNumber, err)
		}
		batch, err := b.db.GetBatch(*batchHash)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve batch for batch hash %s. Cause: %w", batchHash, err)
		}

		batchesToSend = append(batchesToSend, batch)

		if currentBatchNumber.Cmp(headBatchHeader.Number) >= 0 {
			break
		}
		// We only send at most 10 catch-up batches at once, to avoid pressure on the messaging system.
		if len(batchesToSend) >= 10 {
			break
		}
	}

	// todo - joel - this is logging code
	var batchesBeingSent []string //nolint:prealloc
	for _, batch := range batchesToSend {
		batchesBeingSent = append(batchesBeingSent, strconv.FormatInt(batch.Header.Number.Int64(), 10))
	}
	println(fmt.Sprintf("jjj sending catch-up batches to node %s; batches are %s", batchRequest.Requester, strings.Join(batchesBeingSent, ", ")))

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
