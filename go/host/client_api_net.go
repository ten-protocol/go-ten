package host

import "fmt"

const obscuroNetworkVersion = 1

// NetworkAPI implements a subset of the Ethereum network JSON RPC operations.
type NetworkAPI struct{}

func NewNetworkAPI() *NetworkAPI {
	return &NetworkAPI{}
}

// Version returns the protocol version of the Obscuro network.
func (api *NetworkAPI) Version() string {
	return fmt.Sprintf("%d", obscuroNetworkVersion)
}
