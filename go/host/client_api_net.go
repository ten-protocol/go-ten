package host

import "fmt"

// NetworkAPI implements a subset of the Ethereum network JSON RPC operations.
type NetworkAPI struct {
	host *Node
}

func NewNetworkAPI(host *Node) *NetworkAPI {
	return &NetworkAPI{
		host: host,
	}
}

// Version returns the protocol version of the Obscuro network.
func (api *NetworkAPI) Version() string {
	return fmt.Sprintf("%d", api.host.config.ObscuroChainID)
}
