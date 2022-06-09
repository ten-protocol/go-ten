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

// GetBalance returns the address's balance on the Obscuro network, encrypted with the viewing key for the address.
func (api *EthereumAPI) GetBalance(_ context.Context, address common.Address, _ rpc.BlockNumberOrHash) (string, error) {
	encryptedBalance, err := api.host.EnclaveClient.GetBalance(address)
	if err != nil {
		return "", err
	}
	return common.Bytes2Hex(encryptedBalance), nil
}

// GetBlockByNumber is a placeholder for an RPC method required by MetaMask.
func (api *EthereumAPI) GetBlockByNumber(context.Context, rpc.BlockNumber, bool) (map[string]interface{}, error) {
	return nil, nil //nolint:nilnil
}

// GasPrice is a placeholder for an RPC method required by MetaMask.
func (api *EthereumAPI) GasPrice(context.Context) (*hexutil.Big, error) {
	return (*hexutil.Big)(big.NewInt(0)), nil
}
