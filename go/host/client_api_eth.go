package host

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/obscuro-playground/go/common"

	gethcommon "github.com/ethereum/go-ethereum/common"
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
	return (*hexutil.Big)(big.NewInt(api.host.config.ObscuroChainID)), nil
}

// BlockNumber returns the height of the current head rollup.
func (api *EthereumAPI) BlockNumber() hexutil.Uint64 {
	return hexutil.Uint64(api.host.nodeDB.GetCurrentRollupHead().Number.Uint64())
}

// GetBalance returns the address's balance on the Obscuro network, encrypted with the viewing key corresponding to the
// `address` field and encoded as hex.
func (api *EthereumAPI) GetBalance(_ context.Context, encryptedParams common.EncryptedParamsGetBalance) (string, error) {
	encryptedBalance, err := api.host.EnclaveClient.GetBalance(encryptedParams)
	if err != nil {
		return "", err
	}
	return gethcommon.Bytes2Hex(encryptedBalance), nil
}

// GetBlockByNumber returns the rollup with the given height as a block. No transactions are included.
func (api *EthereumAPI) GetBlockByNumber(_ context.Context, number rpc.BlockNumber, _ bool) (map[string]interface{}, error) {
	extRollup := api.host.EnclaveClient.GetRollupByHeight(number.Int64())
	return extRollupToBlock(extRollup), nil
}

// GetBlockByHash returns the rollup with the given hash as a block. No transactions are included.
func (api *EthereumAPI) GetBlockByHash(_ context.Context, hash gethcommon.Hash, _ bool) (map[string]interface{}, error) {
	extRollup := api.host.EnclaveClient.GetRollup(hash)
	return extRollupToBlock(extRollup), nil
}

// GasPrice is a placeholder for an RPC method required by MetaMask/Remix.
func (api *EthereumAPI) GasPrice(context.Context) (*hexutil.Big, error) {
	return (*hexutil.Big)(big.NewInt(0)), nil
}

// Call returns the result of executing the smart contract as a user, encrypted with the viewing key corresponding to
// the `from` field and encoded as hex.
func (api *EthereumAPI) Call(_ context.Context, encryptedParams common.EncryptedParamsCall) (string, error) {
	encryptedResponse, err := api.host.EnclaveClient.ExecuteOffChainTransaction(encryptedParams)
	if err != nil {
		return "", err
	}
	return gethcommon.Bytes2Hex(encryptedResponse), nil
}

// GetTransactionReceipt returns the transaction receipt for the given transaction hash, encrypted with the viewing key
// corresponding to the original transaction submitter and encoded as hex.
func (api *EthereumAPI) GetTransactionReceipt(_ context.Context, encryptedParams common.EncryptedParamsGetTxReceipt) (string, error) {
	encryptedResponse, err := api.host.EnclaveClient.GetTransactionReceipt(encryptedParams)
	if err != nil {
		return "", err
	}
	return gethcommon.Bytes2Hex(encryptedResponse), nil
}

// EstimateGas is a placeholder for an RPC method required by MetaMask/Remix.
func (api *EthereumAPI) EstimateGas(_ context.Context, _ interface{}, _ *rpc.BlockNumberOrHash) (hexutil.Uint64, error) {
	// TODO - Return a non-dummy gas estimate.
	return 0, nil
}

// SendRawTransaction sends the encrypted transaction.
func (api *EthereumAPI) SendRawTransaction(_ context.Context, encryptedParams common.EncryptedParamsSendRawTx) (string, error) {
	encryptedResponse, err := api.host.SubmitAndBroadcastTx(encryptedParams)
	if err != nil {
		return "", err
	}
	return gethcommon.Bytes2Hex(encryptedResponse), nil
}

// GetCode returns the code stored at the given address in the state for the given rollup height or rollup hash.
func (api *EthereumAPI) GetCode(_ context.Context, address gethcommon.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	rollupHeight, ok := blockNrOrHash.Number()
	if ok {
		rollup := api.host.EnclaveClient.GetRollupByHeight(rollupHeight.Int64())
		if rollup == nil {
			return nil, errors.New("invalid arguments; rollup with given height does not exist")
		}
		rollupHash := rollup.Header.Hash()
		return api.host.EnclaveClient.GetCode(address, &rollupHash)
	}

	rollupHash, ok := blockNrOrHash.Hash()
	if ok {
		return api.host.EnclaveClient.GetCode(address, &rollupHash)
	}

	return nil, errors.New("invalid arguments; neither rollup height nor rollup hash specified")
}

// TODO - Temporary. Will be replaced by encrypted implementation.
func (api *EthereumAPI) GetTransactionCount(_ context.Context, address gethcommon.Address, _ rpc.BlockNumberOrHash) (*hexutil.Uint64, error) {
	nonce := api.host.EnclaveClient.Nonce(address)
	return (*hexutil.Uint64)(&nonce), nil
}

// GetTransactionByHash returns the transaction with the given hash.
func (api *EthereumAPI) GetTransactionByHash(_ context.Context, hash gethcommon.Hash) (*RPCTransaction, error) {
	// Unlike in the Geth impl, we do not try and retrieve unconfirmed transactions from the mempool.
	tx, blockHash, blockNumber, index, err := api.host.EnclaveClient.GetTransaction(hash)
	if err != nil {
		return nil, err
	}
	// Unlike in the Geth impl, we hardcode the use of a London signer.
	signer := types.NewLondonSigner(tx.ChainId())
	return newRPCTransaction(tx, blockHash, blockNumber, index, gethcommon.Big0, signer), nil
}

// Maps an external rollup to a block.
func extRollupToBlock(extRollup *common.ExtRollup) map[string]interface{} {
	return map[string]interface{}{
		"number":           (*hexutil.Big)(extRollup.Header.Number),
		"hash":             extRollup.Header.Hash(),
		"parenthash":       extRollup.Header.ParentHash,
		"nonce":            extRollup.Header.Nonce,
		"logsbloom":        extRollup.Header.Bloom,
		"stateroot":        extRollup.Header.Root,
		"receiptsroot":     extRollup.Header.ReceiptHash,
		"miner":            extRollup.Header.Agg,
		"extradata":        hexutil.Bytes(extRollup.Header.Extra),
		"transactionsroot": extRollup.Header.TxHash,
		"transactions":     extRollup.TxHashes,
	}
}

// RPCTransaction is lifted from Geth's internal `ethapi` package.
type RPCTransaction struct {
	BlockHash        *gethcommon.Hash    `json:"blockHash"`
	BlockNumber      *hexutil.Big        `json:"blockNumber"`
	From             gethcommon.Address  `json:"from"`
	Gas              hexutil.Uint64      `json:"gas"`
	GasPrice         *hexutil.Big        `json:"gasPrice"`
	GasFeeCap        *hexutil.Big        `json:"maxFeePerGas,omitempty"`
	GasTipCap        *hexutil.Big        `json:"maxPriorityFeePerGas,omitempty"`
	Hash             gethcommon.Hash     `json:"hash"`
	Input            hexutil.Bytes       `json:"input"`
	Nonce            hexutil.Uint64      `json:"nonce"`
	To               *gethcommon.Address `json:"to"`
	TransactionIndex *hexutil.Uint64     `json:"transactionIndex"`
	Value            *hexutil.Big        `json:"value"`
	Type             hexutil.Uint64      `json:"type"`
	Accesses         *types.AccessList   `json:"accessList,omitempty"`
	ChainID          *hexutil.Big        `json:"chainId,omitempty"`
	V                *hexutil.Big        `json:"v"`
	R                *hexutil.Big        `json:"r"`
	S                *hexutil.Big        `json:"s"`
}

// Lifted from Geth's internal `ethapi` package.
func newRPCTransaction(tx *types.Transaction, blockHash gethcommon.Hash, blockNumber uint64, index uint64, baseFee *big.Int, signer types.Signer) *RPCTransaction {
	from, _ := types.Sender(signer, tx)
	v, r, s := tx.RawSignatureValues()
	result := &RPCTransaction{
		Type:     hexutil.Uint64(tx.Type()),
		From:     from,
		Gas:      hexutil.Uint64(tx.Gas()),
		GasPrice: (*hexutil.Big)(tx.GasPrice()),
		Hash:     tx.Hash(),
		Input:    hexutil.Bytes(tx.Data()),
		Nonce:    hexutil.Uint64(tx.Nonce()),
		To:       tx.To(),
		Value:    (*hexutil.Big)(tx.Value()),
		V:        (*hexutil.Big)(v),
		R:        (*hexutil.Big)(r),
		S:        (*hexutil.Big)(s),
	}
	if blockHash != (gethcommon.Hash{}) {
		result.BlockHash = &blockHash
		result.BlockNumber = (*hexutil.Big)(new(big.Int).SetUint64(blockNumber))
		result.TransactionIndex = (*hexutil.Uint64)(&index)
	}
	switch tx.Type() {
	case types.AccessListTxType:
		al := tx.AccessList()
		result.Accesses = &al
		result.ChainID = (*hexutil.Big)(tx.ChainId())
	case types.DynamicFeeTxType:
		al := tx.AccessList()
		result.Accesses = &al
		result.ChainID = (*hexutil.Big)(tx.ChainId())
		result.GasFeeCap = (*hexutil.Big)(tx.GasFeeCap())
		result.GasTipCap = (*hexutil.Big)(tx.GasTipCap())
		// if the transaction has been mined, compute the effective gas price
		if baseFee != nil && blockHash != (gethcommon.Hash{}) {
			// price = min(tip, gasFeeCap - baseFee) + baseFee
			price := math.BigMin(new(big.Int).Add(tx.GasTipCap(), baseFee), tx.GasFeeCap())
			result.GasPrice = (*hexutil.Big)(price)
		} else {
			result.GasPrice = (*hexutil.Big)(tx.GasFeeCap())
		}
	}
	return result
}
