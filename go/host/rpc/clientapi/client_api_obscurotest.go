package clientapi

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common/host"
)

// TODO - Disable these methods on a non-test node.

// TestAPI implements JSON RPC operations required for testing.
type TestAPI struct {
	host host.Host
}

func NewTestAPI(host host.Host) *TestAPI {
	return &TestAPI{
		host: host,
	}
}

// GetHeadBlockHeader returns the current head block's header.
func (api *TestAPI) GetHeadBlockHeader() *types.Header {
	return api.host.DB().GetHeadBlockHeader()
}

// StopHost gracefully stops the host.
// TODO - Investigate how to authenticate this and other sensitive methods in production (Geth uses JWT).
func (api *TestAPI) StopHost() {
	api.host.Stop()
}
