package host

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

// EthereumAPI implements a subset of the Ethereum JSON RPC operations.
type EthereumAPI struct {
	host *Node
}

func NewEthereumAPI(host *Node) *EthereumAPI {
	return &EthereumAPI{
		host: host,
	}
}

// ChainId returns the Obscuro chain ID.
func (api *EthereumAPI) ChainId() (*hexutil.Big, error) { //nolint:stylecheck,revive
	return (*hexutil.Big)(&api.host.config.ChainID), nil
}

// BlockNumber returns the height of the current head rollup.
func (api *EthereumAPI) BlockNumber() hexutil.Uint64 {
	return hexutil.Uint64(api.host.nodeDB.GetCurrentRollupHead().Number)
}

// GetBalance returns the address's balance on the Obscuro network.
// TODO - Establish what value we should return here (balance in a specific ERC-20 contract?).
// TODO - Encrypt the response with a viewing key on the enclave side once we're returning a non-dummy value.
func (api *EthereumAPI) GetBalance(context.Context, common.Address, rpc.BlockNumberOrHash) (*hexutil.Big, error) {
	return (*hexutil.Big)(big.NewInt(0)), nil
}

// GetBlockByNumber is a placeholder for an RPC method required by MetaMask.
func (api *EthereumAPI) GetBlockByNumber(context.Context, rpc.BlockNumber, bool) (map[string]interface{}, error) {
	return nil, nil //nolint:nilnil
}

// GasPrice is a placeholder for an RPC method required by MetaMask.
func (api *EthereumAPI) GasPrice(context.Context) (*hexutil.Big, error) {
	return (*hexutil.Big)(big.NewInt(0)), nil
}
