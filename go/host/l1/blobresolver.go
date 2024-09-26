package l1

import (
	"context"
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"

	"github.com/ten-protocol/go-ten/go/ethadapter"
)

type BlobResolver interface {
	// FetchBlobs Fetches the blob data using beacon chain APIs
	FetchBlobs(ctx context.Context, b *types.Header, hashes []gethcommon.Hash) ([]*kzg4844.Blob, error)
	StoreBlobs(slot uint64, blobs []*kzg4844.Blob) error
}

type beaconBlobResolver struct {
	beaconClient *ethadapter.L1BeaconClient
}

func NewBlobResolver(beaconClient *ethadapter.L1BeaconClient) BlobResolver {
	return &beaconBlobResolver{beaconClient: beaconClient}
}

func (r *beaconBlobResolver) FetchBlobs(ctx context.Context, b *types.Header, hashes []gethcommon.Hash) ([]*kzg4844.Blob, error) {
	blobs, err := r.beaconClient.FetchBlobs(ctx, b, hashes)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch blobs from beacon client: %w", err)
	}
	return blobs, nil
}

func (r *beaconBlobResolver) StoreBlobs(slot uint64, blobs_ []*kzg4844.Blob) error {

}
