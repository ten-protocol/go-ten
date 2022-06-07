package host

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// EthereumAPI implements a subset of the Ethereum JSON RPC operations.
type EthereumAPI struct {
	host *Node
}

func NewEthereumAPI(host *Node) *ObscuroAPI {
	return &ObscuroAPI{
		host: host,
	}
}

// ChainId returns the Obscuro chain ID.
func (api *ObscuroAPI) ChainId() (*hexutil.Big, error) {
	return (*hexutil.Big)(&api.host.config.ChainID), nil
}
