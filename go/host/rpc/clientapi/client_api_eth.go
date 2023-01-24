package clientapi

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/obscuronet/go-obscuro/go/common/errutil"

	"github.com/obscuronet/go-obscuro/go/common/host"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common"
)

// EthereumAPI implements a subset of the Ethereum JSON RPC operations. All the method signatures are copied from the
// corresponding Geth implementations.
type EthereumAPI struct {
	host   host.Host
	logger gethlog.Logger
}

func NewEthereumAPI(host host.Host, logger gethlog.Logger) *EthereumAPI {
	return &EthereumAPI{
		host:   host,
		logger: logger,
	}
}

// ChainId returns the Obscuro chain ID.
func (api *EthereumAPI) ChainId() (*hexutil.Big, error) { //nolint:stylecheck,revive
	return (*hexutil.Big)(big.NewInt(api.host.Config().ObscuroChainID)), nil
}

// BlockNumber returns the height of the current head batch.
func (api *EthereumAPI) BlockNumber() hexutil.Uint64 {
	header, err := api.host.DB().GetHeadBatchHeader()
	if err != nil {
		// This error may be nefarious, but unfortunately the Eth API doesn't allow us to return an error.
		api.logger.Error("could not retrieve head batch header", log.ErrKey, err)
		return 0
	}
	return hexutil.Uint64(header.Number.Uint64())
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

// GetBlockByNumber returns the header of the batch with the given height.
func (api *EthereumAPI) GetBlockByNumber(ctx context.Context, number rpc.BlockNumber, _ bool) (map[string]interface{}, error) {
	batchHash, err := api.batchNumberToBatchHash(number)
	if err != nil {
		return nil, fmt.Errorf("could not find batch with height %d. Cause: %w", number, err)
	}
	return api.GetBlockByHash(ctx, *batchHash, true)
}

// GetBlockByHash returns the header of the batch with the given hash.
func (api *EthereumAPI) GetBlockByHash(_ context.Context, hash gethcommon.Hash, _ bool) (map[string]interface{}, error) {
	batchHeader, err := api.host.DB().GetBatchHeader(hash)
	if err != nil {
		return nil, err
	}
	return headerToMap(batchHeader), nil
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

// GetCode returns the code stored at the given address in the state for the given batch height or batch hash.
// TODO - Instead of converting the block number of hash client-side, do it on the enclave.
func (api *EthereumAPI) GetCode(_ context.Context, address gethcommon.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	// requested a number
	if batchNumber, ok := blockNrOrHash.Number(); ok {
		batchHash, err := api.batchNumberToBatchHash(batchNumber)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve batch with height %d. Cause: %w", batchNumber, err)
		}

		return api.host.EnclaveClient().GetCode(address, batchHash)
	}

	// requested a hash
	if batchHash, ok := blockNrOrHash.Hash(); ok {
		return api.host.EnclaveClient().GetCode(address, &batchHash)
	}

	return nil, errors.New("invalid arguments; neither batch height nor batch hash specified")
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

// Converts a batch header to a key/value map.
// TODO - Include all the fields of the rollup header that do not exist in the Geth block headers as well (not just withdrawals).
func headerToMap(header *common.BatchHeader) map[string]interface{} {
	return map[string]interface{}{
		// The fields present in Geth's `types/Header` struct.
		"parentHash":       header.ParentHash,
		"sha3Uncles":       header.UncleHash,
		"miner":            header.Coinbase,
		"stateRoot":        header.Root,
		"transactionsRoot": header.TxHash,
		"receiptsRoot":     header.ReceiptHash,
		"logsBloom":        header.Bloom,
		"difficulty":       header.Difficulty,
		"number":           header.Number,
		"gasLimit":         header.GasLimit,
		"gasUsed":          header.GasUsed,
		"timestamp":        header.Time,
		"extraData":        header.Extra,
		"mixHash":          header.MixDigest,
		"nonce":            header.Nonce,
		"baseFeePerGas":    header.BaseFee,

		// The custom Obscuro fields.
		"agg":                     header.Agg,
		"l1Proof":                 header.L1Proof,
		"crossChainMessages":      header.CrossChainMessages,
		"inboundCrossChainHash":   header.LatestInboudCrossChainHash,
		"inboundCrossChainHeight": header.LatestInboundCrossChainHeight,
	}
}

// FeeHistoryResult is the structure returned by Geth `eth_feeHistory` API.
type FeeHistoryResult struct {
	OldestBlock  *hexutil.Big     `json:"oldestBlock"`
	Reward       [][]*hexutil.Big `json:"reward,omitempty"`
	BaseFee      []*hexutil.Big   `json:"baseFeePerGas,omitempty"`
	GasUsedRatio []float64        `json:"gasUsedRatio"`
}

// Given a batch number, returns the hash of the batch with that number.
func (api *EthereumAPI) batchNumberToBatchHash(batchNumber rpc.BlockNumber) (*gethcommon.Hash, error) {
	// Handling the special cases first. No special handling is required for rpc.EarliestBlockNumber.
	if batchNumber == rpc.LatestBlockNumber {
		batchHeader, err := api.host.DB().GetHeadBatchHeader()
		if err != nil {
			return nil, err
		}
		batchHash := batchHeader.Hash()
		return &batchHash, nil
	}

	if batchNumber == rpc.PendingBlockNumber {
		// todo Dependent on the current pending batch - leaving it for a different iteration as it will need more thought
		return nil, errutil.ErrNoImpl
	}

	batchNumberBig := big.NewInt(batchNumber.Int64())
	batchHash, err := api.host.DB().GetBatchHash(batchNumberBig)
	if err != nil {
		return nil, err
	}
	return batchHash, nil
}
