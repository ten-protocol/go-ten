package clientapi

import (
	"fmt"
)

// NetworkAPI implements a subset of the Ethereum network JSON RPC operations.
type NetworkAPI struct {
	chainID int64
}

func NewNetworkAPI(chainID int64) *NetworkAPI {
	return &NetworkAPI{
		chainID: chainID,
	}
}

// Version returns the protocol version of the Obscuro network.
func (api *NetworkAPI) Version() string {
	return fmt.Sprintf("%d", api.chainID)
}
