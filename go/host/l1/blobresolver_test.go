package l1

import (
	"context"
	"math/big"
	"net/http"
	"testing"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/ten-protocol/go-ten/go/ethadapter"
)

const (
	vHash1 = "0x012b7a6a22399aa9eecd8eda6ec658679e81be21af6ff296116aee205e2218f2"
	vHash2 = "0x012374e04a848591844b75bc2f500318cf640552379b5e3a1a77bb828620690e"
)

func TestBlobResolver(t *testing.T) {
	beaconClient := ethadapter.NewBeaconHTTPClient(new(http.Client), "https://docs-demo.quiknode.pro/")
	fallback := ethadapter.NewArchivalHTTPClient(new(http.Client), "https://api.ethernow.xyz")
	blobResolver := NewBlobResolver(ethadapter.NewL1BeaconClient(beaconClient, fallback))

	// this will convert to slot 5 which will return 404 from the quicknode api, causing the fallback to be used
	b := &types.Header{
		Time: 1606824090,
	}

	blobs, err := blobResolver.FetchBlobs(context.Background(), b, []gethcommon.Hash{gethcommon.HexToHash(vHash1), gethcommon.HexToHash(vHash2)})
	require.NoError(t, err)
	require.Len(t, blobs, 2)
}

// TestSepoliaBlobResolver checks the public node sepolia beacon APIs work as expected
func TestSepoliaBlobResolver(t *testing.T) {
	// l1_beacon_url for sepolia
	beaconClient := ethadapter.NewBeaconHTTPClient(new(http.Client), "https://ethereum-sepolia-beacon-api.publicnode.com")
	// l1_blob_archive_url for sepolia
	fallback := ethadapter.NewBeaconHTTPClient(new(http.Client), "https://eth-beacon-chain-sepolia.drpc.org/rest/")
	blobResolver := NewBlobResolver(ethadapter.NewL1BeaconClient(beaconClient, fallback))

	// this is a moving point in time so we can't compare hashes or be certain there will be blobs in the block
	historicalBlock := &types.Header{
		Time:   uint64(time.Now().Add(-30 * 24 * time.Hour).Unix()), // 30 days ago
		Number: big.NewInt(1234567),
	}

	_, err := blobResolver.FetchBlobs(context.Background(), historicalBlock, []gethcommon.Hash{})
	require.NoError(t, err)
}
