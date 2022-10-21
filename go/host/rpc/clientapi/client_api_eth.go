package clientapi

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/host"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// EthereumAPI implements a subset of the Ethereum JSON RPC operations. All the method signatures are copied from the
// corresponding Geth implementations.
type EthereumAPI struct {
	host host.Host
}

func NewEthereumAPI(host host.Host) *EthereumAPI {
	return &EthereumAPI{
		host: host,
	}
}

// ChainId returns the Obscuro chain ID.
func (api *EthereumAPI) ChainId() (*hexutil.Big, error) { //nolint:stylecheck,revive
	return (*hexutil.Big)(big.NewInt(api.host.Config().ObscuroChainID)), nil
}

// BlockNumber returns the height of the current head rollup.
func (api *EthereumAPI) BlockNumber() hexutil.Uint64 {
	head := api.host.DB().GetCurrentRollupHead()
	if head == nil {
		return 0
	}
	number := head.Header.Number.Uint64()
	return hexutil.Uint64(number)
}

// GetBalance returns the address's balance on the Obscuro network, encrypted with the viewing key corresponding to the
// `address` field and encoded as hex.
func (api *EthereumAPI) GetBalance(_ context.Context, encryptedParams common.EncryptedParamsGetBalance) (string, error) {
	encryptedBalance, err := api.host.EnclaveClient().GetBalance(encryptedParams)
	if err != nil {
		return "", err
	}
	return gethcommon.Bytes2Hex(encryptedBalance), nil
}

// GetBlockByNumber returns the rollup with the given height as a block. No transactions are included.
func (api *EthereumAPI) GetBlockByNumber(ctx context.Context, number rpc.BlockNumber, _ bool) (map[string]interface{}, error) {
	rollupHash, err := api.blockNumberToHash(number)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch block number: %w", err)
	}
	return api.GetBlockByHash(ctx, *rollupHash, true)
}

// GetBlockByHash returns the rollup with the given hash as a block. No transactions are included.
func (api *EthereumAPI) GetBlockByHash(_ context.Context, hash gethcommon.Hash, _ bool) (map[string]interface{}, error) {
	rollupHeaderWithHashes := api.host.DB().GetRollupHeader(hash)
	if rollupHeaderWithHashes == nil {
		return nil, nil //nolint:nilnil
	}
	return headerWithHashesToBlock(rollupHeaderWithHashes), nil
}

// GasPrice is a placeholder for an RPC method required by MetaMask/Remix.
func (api *EthereumAPI) GasPrice(context.Context) (*hexutil.Big, error) {
	return (*hexutil.Big)(big.NewInt(1)), nil
}

// Call returns the result of executing the smart contract as a user, encrypted with the viewing key corresponding to
// the `from` field and encoded as hex.
func (api *EthereumAPI) Call(_ context.Context, encryptedParams common.EncryptedParamsCall) (string, error) {
	encryptedResponse, err := api.host.EnclaveClient().ExecuteOffChainTransaction(encryptedParams)
	if err != nil {
		return "", err
	}
	return gethcommon.Bytes2Hex(encryptedResponse), nil
}

// GetTransactionReceipt returns the transaction receipt for the given transaction hash, encrypted with the viewing key
// corresponding to the original transaction submitter and encoded as hex, or nil if no matching transaction exists.
func (api *EthereumAPI) GetTransactionReceipt(_ context.Context, encryptedParams common.EncryptedParamsGetTxReceipt) (*string, error) {
	encryptedResponse, err := api.host.EnclaveClient().GetTransactionReceipt(encryptedParams)
	if err != nil {
		return nil, err
	}
	if encryptedResponse == nil {
		return nil, nil //nolint:nilnil
	}
	encryptedResponseHex := gethcommon.Bytes2Hex(encryptedResponse)
	return &encryptedResponseHex, nil
}

// EstimateGas requests the enclave the gas estimation based on the callMsg supplied params (encrypted)
func (api *EthereumAPI) EstimateGas(_ context.Context, encryptedParams common.EncryptedParamsEstimateGas) (*string, error) {
	encryptedResponse, err := api.host.EnclaveClient().EstimateGas(encryptedParams)
	if err != nil {
		return nil, err
	}

	encryptedResponseHex := gethcommon.Bytes2Hex(encryptedResponse)
	return &encryptedResponseHex, nil
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
	// requested a number
	if rollupNumber, ok := blockNrOrHash.Number(); ok {
		rollupHash, err := api.blockNumberToHash(rollupNumber)
		if err != nil {
			return nil, fmt.Errorf("unable to fetch block number: %w", err)
		}

		return api.host.EnclaveClient().GetCode(address, rollupHash)
	}

	// requested a hash
	if rollupHash, ok := blockNrOrHash.Hash(); ok {
		return api.host.EnclaveClient().GetCode(address, &rollupHash)
	}

	return nil, errors.New("invalid arguments; neither rollup height nor rollup hash specified")
}

func (api *EthereumAPI) GetTransactionCount(_ context.Context, encryptedParams common.EncryptedParamsGetTxCount) (string, error) {
	encryptedResponse, err := api.host.EnclaveClient().GetTransactionCount(encryptedParams)
	if err != nil {
		return "", err
	}
	if encryptedResponse == nil {
		return "", err
	}
	return gethcommon.Bytes2Hex(encryptedResponse), nil
}

// GetTransactionByHash returns the transaction with the given hash, encrypted with the viewing key corresponding to the
// `from` field and encoded as hex, or nil if no matching transaction exists.
func (api *EthereumAPI) GetTransactionByHash(_ context.Context, encryptedParams common.EncryptedParamsGetTxByHash) (*string, error) {
	encryptedResponse, err := api.host.EnclaveClient().GetTransaction(encryptedParams)
	if err != nil {
		return nil, err
	}
	if encryptedResponse == nil {
		return nil, err
	}
	encryptedResponseHex := gethcommon.Bytes2Hex(encryptedResponse)
	return &encryptedResponseHex, nil
}

// FeeHistory is a placeholder for an RPC method required by MetaMask/Remix.
func (api *EthereumAPI) FeeHistory(context.Context, rpc.DecimalOrHex, rpc.BlockNumber, []float64) (*FeeHistoryResult, error) {
	// TODO - Return a non-dummy fee history.
	return &FeeHistoryResult{
		OldestBlock:  (*hexutil.Big)(big.NewInt(0)),
		Reward:       [][]*hexutil.Big{},
		BaseFee:      []*hexutil.Big{},
		GasUsedRatio: []float64{},
	}, nil
}

// Maps an external rollup to a block.
func headerWithHashesToBlock(headerWithHashes *common.HeaderWithTxHashes) map[string]interface{} {
	header := headerWithHashes.Header
	return map[string]interface{}{
		"number":           (*hexutil.Big)(header.Number),
		"hash":             header.Hash(),
		"parentHash":       header.ParentHash,
		"nonce":            header.Nonce,
		"logsBloom":        header.Bloom,
		"stateRoot":        header.Root,
		"receiptsRoot":     header.ReceiptHash,
		"miner":            header.Agg,
		"extraData":        hexutil.Bytes(header.Extra),
		"transactionsRoot": header.TxHash,
		"transactions":     headerWithHashes.TxHashes,

		"sha3Uncles":    header.UncleHash,
		"difficulty":    header.Difficulty,
		"gasLimit":      header.GasLimit,
		"gasUsed":       header.GasUsed,
		"timestamp":     header.Time,
		"mixHash":       header.MixDigest,
		"baseFeePerGas": header.BaseFee,
	}
}

// FeeHistoryResult is the structure returned by Geth `eth_feeHistory` API.
type FeeHistoryResult struct {
	OldestBlock  *hexutil.Big     `json:"oldestBlock"`
	Reward       [][]*hexutil.Big `json:"reward,omitempty"`
	BaseFee      []*hexutil.Big   `json:"baseFeePerGas,omitempty"`
	GasUsedRatio []float64        `json:"gasUsedRatio"`
}

func (api *EthereumAPI) blockNumberToHash(blockNumber rpc.BlockNumber) (*gethcommon.Hash, error) {
	// Predefined constants to support Geth's API
	switch blockNumber {
	case rpc.LatestBlockNumber:
		hash := api.host.DB().GetCurrentRollupHead().Header.Hash()
		return &hash, nil
	case rpc.EarliestBlockNumber:
		hash := api.host.DB().GetRollupHeader(gethcommon.BigToHash(big.NewInt(0))).Header.Hash()
		return &hash, nil
	case rpc.PendingBlockNumber:
		// todo Dependent on the current pending rollup - leaving it for a different iteration as it will need more thought
		return nil, nil //nolint
	}

	rollupHash := api.host.DB().GetRollupHash(big.NewInt(blockNumber.Int64()))
	if rollupHash == nil {
		return nil, fmt.Errorf("unable to fetch rollup at height: %d", blockNumber.Int64())
	}
	return rollupHash, nil
}
