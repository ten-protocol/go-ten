package l1

import (
	"context"
	"net/http"
	"testing"

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
