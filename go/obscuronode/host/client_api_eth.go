package host

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

// EthereumAPI implements a subset of the Ethereum JSON RPC operations. All the method signatures are copied from the
// corresponding Geth implementations.
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
	return hexutil.Uint64(api.host.nodeDB.GetCurrentRollupHead().Number.Uint64())
}

// GetBalance returns the address's balance on the Obscuro network, encrypted with the viewing key corresponding to the
// `address` field and encoded as hex.
func (api *EthereumAPI) GetBalance(_ context.Context, encryptedParams []byte) (string, error) {
	encryptedBalance, err := api.host.EnclaveClient.GetBalance(encryptedParams)
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

// Call returns the result of executing the smart contract as a user, encrypted with the viewing key corresponding to
// the `from` field and encoded as hex.
func (api *EthereumAPI) Call(_ context.Context, encryptedParams []byte) (string, error) {
	encryptedResponse, err := api.host.EnclaveClient.ExecuteOffChainTransaction(encryptedParams)
	return common.Bytes2Hex(encryptedResponse), err
}

// GetTransactionReceipt returns the transaction receipt for the given transaction hash, encrypted with the viewing key
// corresponding to the original transaction submitter and encoded as hex.
func (api *EthereumAPI) GetTransactionReceipt(_ context.Context, encryptedParams []byte) (map[string]interface{}, error) {
	return nil, nil // todo - joel - return something sensible.
}
