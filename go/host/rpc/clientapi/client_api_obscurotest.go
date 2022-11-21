package clientapi

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
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

// BlockNumber returns the height of the current head block.
func (api *TestAPI) BlockNumber() (hexutil.Uint64, error) {
	head, err := api.host.DB().GetHeadBlockHeader()
	if err != nil {
		return 0, err
	}

	number := head.Number.Uint64()
	return hexutil.Uint64(number), nil
}

// StopHost gracefully stops the host.
// TODO - Investigate how to authenticate this and other sensitive methods in production (Geth uses JWT).
func (api *TestAPI) StopHost() {
	api.host.Stop()
}
