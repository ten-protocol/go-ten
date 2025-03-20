package l1

import (
	"context"
	"net/http"
	"testing"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/ten-protocol/go-ten/go/ethadapter"
)

const (
	vHash1 = "0x01c509cb5b108f8edf2333fa35a93acc5e4b24808179b0604f1945a8c7e98a3a"
	vHash2 = "0x013cc291176322a963ca75c09d27e31ab7690d88afdefb87c436815baeeb1078"
)

func TestBlobResolver(t *testing.T) {
	t.Skipf("TODO Test needs updating with new params that work")
	beaconClient := ethadapter.NewBeaconHTTPClient(new(http.Client), "https://docs-demo.quiknode.pro/")
	//fallback := ethadapter.NewArchivalHTTPClient(new(http.Client), "https://eth-beacon-chain.drpc.org/rest/")
	blobResolver := NewBlobResolver(ethadapter.NewL1BeaconClient(beaconClient), nil)

	// this will convert to slot 5 which will return 404 from the quicknode api, causing the fallback to be used
	b := &types.Header{
		Time: 1742476343,
	}

	blobs, err := blobResolver.FetchBlobs(context.Background(), b, []gethcommon.Hash{gethcommon.HexToHash(vHash1), gethcommon.HexToHash(vHash2)})
	require.NoError(t, err)
	require.Len(t, blobs, 2)
}

// TestSepoliaBlobResolver checks the public node sepolia beacon APIs work as expected
func TestSepoliaBlobResolver(t *testing.T) {
	t.Skipf("Test will occasionally not pass due to the time window landing on a block with no blobs")
	// l1_beacon_url for sepolia
	beaconClient := ethadapter.NewBeaconHTTPClient(new(http.Client), "https://ethereum-sepolia-beacon-api.publicnode.com")
	// l1_blob_archive_url for sepolia
	fallback := ethadapter.NewBeaconHTTPClient(new(http.Client), "https://eth-beacon-chain-sepolia.drpc.org/rest/")
	blobResolver := NewBlobResolver(ethadapter.NewL1BeaconClient(beaconClient, fallback), nil)

	// this is a moving point in time so we can't compare hashes or be certain there will be blobs in the block
	// create block with timestamp 30 days ago relative to current time
	historicalBlock := &types.Header{
		Time: uint64(time.Now().Unix()) - (30 * 24 * 60 * 60), // 30 days in seconds
	}

	_, err := blobResolver.FetchBlobs(context.Background(), historicalBlock, []gethcommon.Hash{})
	require.NoError(t, err)
}
