package clientapi

import (
	"fmt"

	"github.com/obscuronet/go-obscuro/go/host"
)

// NetworkAPI implements a subset of the Ethereum network JSON RPC operations.
type NetworkAPI struct {
	host host.Host
}

func NewNetworkAPI(host host.Host) *NetworkAPI {
	return &NetworkAPI{
		host: host,
	}
}

// Version returns the protocol version of the Obscuro network.
func (api *NetworkAPI) Version() string {
	return fmt.Sprintf("%d", api.host.Config().ObscuroChainID)
}
