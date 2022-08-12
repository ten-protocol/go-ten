package clientapi

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/host"
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

// GetID returns the ID of the host.
func (api *TestAPI) GetID() gethcommon.Address {
	return api.host.Config().ID
}

// GetCurrentBlockHead returns the current head block's header.
func (api *TestAPI) GetCurrentBlockHead() *types.Header {
	return api.host.DB().GetCurrentBlockHead()
}

// GetRollupHeader returns the header of the rollup with the given hash.
func (api *TestAPI) GetRollupHeader(hash gethcommon.Hash) *common.Header {
	headerWithHashes := api.host.DB().GetRollupHeader(hash)
	if headerWithHashes == nil {
		return nil
	}
	return headerWithHashes.Header
}

// StopHost gracefully stops the host.
// TODO - Investigate how to authenticate this and other sensitive methods in production (Geth uses JWT).
func (api *TestAPI) StopHost() {
	go api.host.Stop()
}
