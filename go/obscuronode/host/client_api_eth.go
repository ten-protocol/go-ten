package host

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// EthereumAPI implements a subset of the Ethereum JSON RPC operations.
type EthereumAPI struct{}

func NewEthereumAPI() *ObscuroAPI {
	return &ObscuroAPI{}
}

// ChainId returns the Obscuro chain ID.
func (api *ObscuroAPI) ChainId() (*hexutil.Big, error) { //nolint:stylecheck,revive
	return (*hexutil.Big)(&api.host.config.ChainID), nil
}
