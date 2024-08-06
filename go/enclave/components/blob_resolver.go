package components

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ten-protocol/go-ten/go/ethadapter"
)

type beaconBlobResolver struct {
	beaconClient *ethadapter.L1BeaconClient
}

func NewBeaconBlobResolver(beaconClient *ethadapter.L1BeaconClient) BlobResolver {
	return &beaconBlobResolver{beaconClient: beaconClient}
}

func (r *beaconBlobResolver) FetchBlobs(ctx context.Context, b *types.Header, hashes []ethadapter.IndexedBlobHash) ([]*ethadapter.Blob, error) {
	blobs, err := r.beaconClient.FetchBlobs(ctx, b, hashes)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch blobs from beacon client: %w", err)
	}
	return blobs, nil
}
