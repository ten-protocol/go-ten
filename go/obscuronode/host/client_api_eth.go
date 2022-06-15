package host

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"

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

// GetBalance returns the address's balance on the Obscuro network, encrypted with the viewing key for the address and
// encoded as hex.
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

// Call returns the result of executing the smart contract as a user, encrypted with the viewing key for the address
// and encoded as hex.
// `data` is generally generated from the ABI of a smart contract.
func (api *EthereumAPI) Call(_ context.Context, args TransactionArgs, _ rpc.BlockNumberOrHash, _ *StateOverride) (string, error) {
	encryptedResponse, err := api.host.EnclaveClient.ExecuteOffChainTransaction(*args.From, *args.To, *args.Data)
	return common.Bytes2Hex(encryptedResponse), err
}

// TransactionArgs is a copy of the same class in Geth's `internal/ethapi` package.
type TransactionArgs struct {
	From                 *common.Address `json:"from"`
	To                   *common.Address `json:"to"`
	Gas                  *hexutil.Uint64 `json:"gas"`
	GasPrice             *hexutil.Big    `json:"gasPrice"`
	MaxFeePerGas         *hexutil.Big    `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *hexutil.Big    `json:"maxPriorityFeePerGas"`
	Value                *hexutil.Big    `json:"value"`
	Nonce                *hexutil.Uint64 `json:"nonce"`

	Data  *hexutil.Bytes `json:"data"`
	Input *hexutil.Bytes `json:"input"`

	// Introduced by AccessListTxType transaction.
	AccessList *types.AccessList `json:"accessList,omitempty"`
	ChainID    *hexutil.Big      `json:"chainId,omitempty"`
}

// OverrideAccount is a copy of the same class in Geth's `internal/ethapi` package.
type OverrideAccount struct {
	Nonce     *hexutil.Uint64              `json:"nonce"`
	Code      *hexutil.Bytes               `json:"code"`
	Balance   **hexutil.Big                `json:"balance"`
	State     *map[common.Hash]common.Hash `json:"state"`
	StateDiff *map[common.Hash]common.Hash `json:"stateDiff"`
}

// StateOverride is a copy of the same class in Geth's `internal/ethapi` package.
type StateOverride map[common.Address]OverrideAccount
