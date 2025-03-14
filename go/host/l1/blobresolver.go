package l1

import (
	"context"
	"fmt"
	"time"

	"github.com/ten-protocol/go-ten/go/common/retry"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"

	"github.com/ten-protocol/go-ten/go/ethadapter"
)

var _maxWaitForBlobs = 2 * time.Minute

// BlobResolver is an interface for fetching blobs
type BlobResolver interface {
	// FetchBlobs Fetches the blob data using beacon chain APIs
	FetchBlobs(ctx context.Context, b *types.Header, hashes []gethcommon.Hash) ([]*kzg4844.Blob, error)
	// StoreBlobs is used to store blobs for the in-memory testing nodes
	StoreBlobs(slot uint64, blobs []*kzg4844.Blob) error
}

type beaconBlobResolver struct {
	beaconClient *ethadapter.L1BeaconClient
}

func NewBlobResolver(beaconClient *ethadapter.L1BeaconClient) BlobResolver {
	return &beaconBlobResolver{beaconClient: beaconClient}
}

func (r *beaconBlobResolver) FetchBlobs(ctx context.Context, b *types.Header, hashes []gethcommon.Hash) ([]*kzg4844.Blob, error) {
	var blobs []*kzg4844.Blob
	err := retry.Do(func() error {
		var fetchErr error
		blobs, fetchErr = r.beaconClient.FetchBlobs(ctx, b, hashes)
		return fetchErr
	}, retry.NewTimeoutStrategy(_maxWaitForBlobs, time.Second))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch blobs after retries: %w", err)
	}
	return blobs, nil
}

func (r *beaconBlobResolver) StoreBlobs(_ uint64, _ []*kzg4844.Blob) error {
	panic("provided by the ethereum consensus layer")
}
