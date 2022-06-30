package host

import (
	"context"
	"math/big"
	"strings"

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
	extRollup := api.host.EnclaveClient.GetRollupByHeight(uint64(number))
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

// GetTransactionCount returns the nonce of the wallet with the given address.
func (api *EthereumAPI) GetTransactionCount(_ context.Context, encryptedParams common.EncryptedParamsGetTxCount) (string, error) {
	encryptedResponse, err := api.host.EnclaveClient.GetTransactionCount(encryptedParams)
	if err != nil {
		return "", err
	}
	// todo: get rid of this hack when we stop supporting unencrypted params for testing purposes
	// 		(the unencrypted hex number doesn't like to be bytes to hex encoded)
	if strings.HasPrefix(string(encryptedResponse), "0x") {
		// the response wasn't encrypted
		return string(encryptedResponse), nil
	}
	return gethcommon.Bytes2Hex(encryptedResponse), nil
}

// EstimateGas is a placeholder for an RPC method required by MetaMask/Remix.
func (api *EthereumAPI) EstimateGas(_ context.Context, _ interface{}, _ *rpc.BlockNumberOrHash) (hexutil.Uint64, error) {
	// TODO - Return a non-dummy gas estimate.
	return 0, nil
}

// SendRawTransaction sends the encrypted transaction
func (api *EthereumAPI) SendRawTransaction(_ context.Context, encryptedParams common.EncryptedParamsSendRawTx) (string, error) {
	encryptedResponse, err := api.host.SubmitAndBroadcastTx(encryptedParams)
	if err != nil {
		return "", err
	}
	return gethcommon.Bytes2Hex(encryptedResponse), nil
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
