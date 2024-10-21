package ethereummock

import (
	"context"
	"fmt"
	"sync"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/go/host/l1"
)

type BlobResolverInMem struct {
	// map of versioned hash to blob for efficient lookup
	versionedHashToBlob sync.Map
	mu                  sync.RWMutex
}

func NewMockBlobResolver() l1.BlobResolver {
	return &BlobResolverInMem{
		versionedHashToBlob: sync.Map{},
		mu:                  sync.RWMutex{},
	}
}

func (b *BlobResolverInMem) StoreBlobs(_ uint64, blobs []*kzg4844.Blob) error {
	for _, blob := range blobs {
		commitment, err := kzg4844.BlobToCommitment(blob)
		if err != nil {
			return fmt.Errorf("failed to compute commitment: %w", err)
		}

		versionedHash := ethadapter.KZGToVersionedHash(commitment)
		b.versionedHashToBlob.Store(versionedHash, blob)
	}
	return nil
}

func (b *BlobResolverInMem) FetchBlobs(_ context.Context, _ *types.Header, hashes []gethcommon.Hash) ([]*kzg4844.Blob, error) {
	var blobs []*kzg4844.Blob
	var missingHashes []string

	for _, vh := range hashes {
		if blob, exists := b.versionedHashToBlob.Load(vh); exists {
			blobs = append(blobs, blob.(*kzg4844.Blob))
		} else {
			missingHashes = append(missingHashes, vh.Hex())
		}
	}

	if len(blobs) == 0 && len(missingHashes) > 0 {
		return nil, fmt.Errorf("blobs not found for hashes: %v", missingHashes)
	}

	return blobs, nil
}
