package ethereummock

import (
	"context"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/go/host/l1"
)

type BlobResolverInMem struct {
	// map of slots to versioned hashes to match the beacon APIs
	slotToVersionedHashes sync.Map
	// map of versioned hash to blob for efficient lookup
	versionedHashToBlob   sync.Map
	mu                  sync.RWMutex
	genesisTime         uint64
	secondsPerSlot      uint64
}

func NewBlobResolver(genesisTime uint64, secondsPerSlot uint64) l1.BlobResolver {
	return &BlobResolverInMem{
		slotToVersionedHashes: sync.Map{},
		versionedHashToBlob:   sync.Map{},
		mu:                    sync.RWMutex{},
		genesisTime:           genesisTime,
		secondsPerSlot:        secondsPerSlot,
	}
}

func (b *BlobResolverInMem) StoreBlobs(slot uint64, blobs []*kzg4844.Blob) error {
	for _, blob := range blobs {
		commitment, err := kzg4844.BlobToCommitment(blob)
		if err != nil {
			return fmt.Errorf("failed to compute commitment: %w", err)
		}

		versionedHash := ethadapter.KZGToVersionedHash(commitment)
		
		hashes, _ := b.slotToVersionedHashes.LoadOrStore(slot, &sync.Map{})
		hashes.(*sync.Map).Store(versionedHash, struct{}{})
		
		b.versionedHashToBlob.Store(versionedHash, blob)
	}
	return nil
}

func (b *BlobResolverInMem) FetchBlobs(_ context.Context, block *types.Header, hashes []gethcommon.Hash) ([]*kzg4844.Blob, error) {
	slot, _ := ethadapter.CalculateSlot(block.Time, MockGenesisBlock.Time(), b.secondsPerSlot)
	
	storedHashes, exists := b.slotToVersionedHashes.Load(slot)
	if !exists {
		return nil, fmt.Errorf("no blobs found for slot %d: %w", slot, ethereum.NotFound)
	}

	if len(hashes) == 0 {
		var allBlobs []*kzg4844.Blob
		storedHashes.(*sync.Map).Range(func(key, _ interface{}) bool {
			if blob, exists := b.versionedHashToBlob.Load(key); exists {
				allBlobs = append(allBlobs, blob.(*kzg4844.Blob))
			}
			return true
		})
		return allBlobs, nil
	}

	var blobs []*kzg4844.Blob
	for _, vh := range hashes {
		if _, found := storedHashes.(*sync.Map).Load(vh); found {
			if blob, exists := b.versionedHashToBlob.Load(vh); exists {
				blobs = append(blobs, blob.(*kzg4844.Blob))
			} else {
				return nil, fmt.Errorf("blob for hash %s not found", vh.Hex())
			}
		} else {
			return nil, fmt.Errorf("versioned hash %s not found in slot %d", vh.Hex(), slot)
		}
	}

	return blobs, nil
}
