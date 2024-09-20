package ethereummock

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/go/host/l1"
	"sync"
)

type BlobResolverInMem struct {
	// map of slots to versioned hashes to match the beacon APIs
	slotToVersionedHashes map[uint64][]gethcommon.Hash
	// map of versioned hash to blob for efficient lookup
	versionedHashToBlob map[gethcommon.Hash]*kzg4844.Blob
	mu                  sync.RWMutex
	genesisTime         uint64
	secondsPerSlot      uint64
	port                int
}

func NewBlobResolver(genesisTime uint64, secondsPerSlot uint64) l1.BlobResolver {
	return &BlobResolverInMem{
		slotToVersionedHashes: make(map[uint64][]gethcommon.Hash),
		versionedHashToBlob:   make(map[gethcommon.Hash]*kzg4844.Blob),
		mu:                    sync.RWMutex{},
		genesisTime:           genesisTime,
		secondsPerSlot:        secondsPerSlot,
	}
}

func (b *BlobResolverInMem) StoreBlobs(slot uint64, blobs []*kzg4844.Blob) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, blob := range blobs {
		commitment, err := kzg4844.BlobToCommitment(blob)
		if err != nil {
			return fmt.Errorf("failed to compute commitment: %w", err)
		}

		versionedHash := ethadapter.KZGToVersionedHash(commitment)
		b.slotToVersionedHashes[slot] = append(b.slotToVersionedHashes[slot], versionedHash)
		b.versionedHashToBlob[versionedHash] = blob
	}
	return nil
}

func (b *BlobResolverInMem) FetchBlobs(_ context.Context, block *types.Header, hashes []gethcommon.Hash) ([]*kzg4844.Blob, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	slot, _ := ethadapter.TimeToSlot(block.Time, MockGenesisBlock.Time(), b.secondsPerSlot)
	// Retrieve the list of versioned hashes stored for the slot.
	storedHashes, exists := b.slotToVersionedHashes[slot]
	if !exists {
		return nil, fmt.Errorf("no blobs found for slot %d: %w", slot, ethereum.NotFound)
	}

	// If no specific versionedHashes are provided, return all blobs for the slot.
	if len(hashes) == 0 {
		var allBlobs []*kzg4844.Blob
		for _, vh := range storedHashes {
			blob, exists := b.versionedHashToBlob[vh]
			if !exists {
				return nil, fmt.Errorf("blob for hash %s not found", vh.Hex())
			}
			allBlobs = append(allBlobs, blob)
		}
		return allBlobs, nil
	}

	// Create a map for quick lookup of stored hashes.
	hashSet := make(map[gethcommon.Hash]struct{}, len(storedHashes))
	for _, h := range storedHashes {
		hashSet[h] = struct{}{}
	}

	// Retrieve the blobs that match the provided versioned hashes.
	var blobs []*kzg4844.Blob
	for _, vh := range hashes {
		if _, found := hashSet[vh]; found {
			blob, exists := b.versionedHashToBlob[vh]
			if !exists {
				return nil, fmt.Errorf("blob for hash %s not found", vh.Hex())
			}
			blobs = append(blobs, blob)
		} else {
			return nil, fmt.Errorf("versioned hash %s not found in slot %d", vh.Hex(), slot)
		}
	}

	return blobs, nil
}
