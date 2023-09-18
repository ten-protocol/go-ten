package clientapi

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/host"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/responses"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
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

// GetBlockByNumber returns the header of the batch with the given height.
func (api *EthereumAPI) GetBlockByNumber(ctx context.Context, number rpc.BlockNumber, _ bool) (common.BatchHeader, error) {
	batchHash, err := api.batchNumberToBatchHash(number)
	if err != nil {
		return common.BatchHeader{}, fmt.Errorf("could not find batch with height %d. Cause: %w", number, err)
	}
	return api.GetBlockByHash(ctx, *batchHash, true)
}

// GetBlockByHash returns the header of the batch with the given hash.
func (api *EthereumAPI) GetBlockByHash(_ context.Context, hash gethcommon.Hash, _ bool) (common.BatchHeader, error) {
	batchHeader, err := api.host.DB().GetBatchHeader(hash)
	if err != nil {
		return common.BatchHeader{}, err
	}
	return *batchHeader, nil
}

// GasPrice is a placeholder for an RPC method required by MetaMask/Remix.
func (api *EthereumAPI) GasPrice(context.Context) (*hexutil.Big, error) {
	header, err := api.host.DB().GetHeadBatchHeader()
	if err != nil {
		return nil, err
	}

	if header.BaseFee == nil || header.BaseFee.Cmp(gethcommon.Big0) == 0 {
		return (*hexutil.Big)(big.NewInt(1)), nil
	}

	return (*hexutil.Big)(big.NewInt(0).Set(header.BaseFee)), nil
}

// GetBalance returns the address's balance on the Obscuro network, encrypted with the viewing key corresponding to the
// `address` field and encoded as hex.
func (api *EthereumAPI) GetBalance(_ context.Context, encryptedParams common.EncryptedParamsGetBalance) (responses.EnclaveResponse, error) {
	enclaveResponse, sysError := api.host.EnclaveClient().GetBalance(encryptedParams)
	if sysError != nil {
		return api.handleSysError("GetBalance", sysError)
	}
	return *enclaveResponse, nil
}

// Call returns the result of executing the smart contract as a user, encrypted with the viewing key corresponding to
// the `from` field and encoded as hex.
func (api *EthereumAPI) Call(_ context.Context, encryptedParams common.EncryptedParamsCall) (responses.EnclaveResponse, error) {
	enclaveResponse, sysError := api.host.EnclaveClient().ObsCall(encryptedParams)
	if sysError != nil {
		return api.handleSysError("Call", sysError)
	}
	return *enclaveResponse, nil
}

// GetTransactionReceipt returns the transaction receipt for the given transaction hash, encrypted with the viewing key
// corresponding to the original transaction submitter and encoded as hex, or nil if no matching transaction exists.
func (api *EthereumAPI) GetTransactionReceipt(_ context.Context, encryptedParams common.EncryptedParamsGetTxReceipt) (responses.EnclaveResponse, error) {
	enclaveResponse, sysError := api.host.EnclaveClient().GetTransactionReceipt(encryptedParams)
	if sysError != nil {
		return api.handleSysError("GetTransactionReceipt", sysError)
	}
	return *enclaveResponse, nil
}

// EstimateGas requests the enclave the gas estimation based on the callMsg supplied params (encrypted)
func (api *EthereumAPI) EstimateGas(_ context.Context, encryptedParams common.EncryptedParamsEstimateGas) (responses.EnclaveResponse, error) {
	enclaveResponse, sysError := api.host.EnclaveClient().EstimateGas(encryptedParams)
	if sysError != nil {
		return api.handleSysError("EstimateGas", sysError)
	}
	return *enclaveResponse, nil
}

// SendRawTransaction sends the encrypted transaction.
func (api *EthereumAPI) SendRawTransaction(_ context.Context, encryptedParams common.EncryptedParamsSendRawTx) (responses.EnclaveResponse, error) {
	enclaveResponse, sysError := api.host.SubmitAndBroadcastTx(encryptedParams)
	if sysError != nil {
		return api.handleSysError("SendRawTransaction", sysError)
	}
	return *enclaveResponse, nil
}

// GetCode returns the code stored at the given address in the state for the given batch height or batch hash.
// todo (#1620) - instead of converting the block number of hash client-side, do it on the enclave
func (api *EthereumAPI) GetCode(_ context.Context, address gethcommon.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	var batchHash *gethcommon.Hash

	// requested a number
	if batchNumber, ok := blockNrOrHash.Number(); ok {
		hash, err := api.batchNumberToBatchHash(batchNumber)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve batch with height %d. Cause: %w", batchNumber, err)
		}
		batchHash = hash
	}

	// requested a hash
	if hash, ok := blockNrOrHash.Hash(); ok {
		batchHash = &hash
	}

	if batchHash == nil {
		return nil, errors.New("invalid arguments; neither batch height nor batch hash specified")
	}

	code, sysError := api.host.EnclaveClient().GetCode(address, batchHash)
	if sysError != nil {
		api.logger.Error(fmt.Sprintf("Enclave System Error. Function %s", "GetCode"), log.ErrKey, sysError)
		return nil, fmt.Errorf(responses.InternalErrMsg)
	}

	return code, nil
}

func (api *EthereumAPI) GetTransactionCount(_ context.Context, encryptedParams common.EncryptedParamsGetTxCount) (responses.EnclaveResponse, error) {
	enclaveResponse, sysError := api.host.EnclaveClient().GetTransactionCount(encryptedParams)
	if sysError != nil {
		return api.handleSysError("GetTransactionCount", sysError)
	}
	return *enclaveResponse, nil
}

// GetTransactionByHash returns the transaction with the given hash, encrypted with the viewing key corresponding to the
// `from` field and encoded as hex, or nil if no matching transaction exists.
func (api *EthereumAPI) GetTransactionByHash(_ context.Context, encryptedParams common.EncryptedParamsGetTxByHash) (responses.EnclaveResponse, error) {
	enclaveResponse, sysError := api.host.EnclaveClient().GetTransaction(encryptedParams)
	if sysError != nil {
		return api.handleSysError("GetTransactionByHash", sysError)
	}
	return *enclaveResponse, nil
}

// GetStorageAt is a reused method for listing the users transactions
func (api *EthereumAPI) GetStorageAt(_ context.Context, encryptedParams common.EncryptedParamsGetStorageAt) (*responses.Receipts, error) {
	return api.host.EnclaveClient().GetCustomQuery(encryptedParams)
}

// FeeHistory is a placeholder for an RPC method required by MetaMask/Remix.
// rpc.DecimalOrHex -> []byte
func (api *EthereumAPI) FeeHistory(context.Context, []byte, rpc.BlockNumber, []float64) (*FeeHistoryResult, error) {
	// todo (#1621) - return a non-dummy fee history

	return &FeeHistoryResult{
		OldestBlock:  (*hexutil.Big)(big.NewInt(0)),
		Reward:       [][]*hexutil.Big{},
		BaseFee:      []*hexutil.Big{},
		GasUsedRatio: []float64{},
	}, nil
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
		return nil, errutil.ErrNoImpl
	}

	batchNumberBig := big.NewInt(batchNumber.Int64())
	batchHash, err := api.host.DB().GetBatchHash(batchNumberBig)
	if err != nil {
		return nil, err
	}
	return batchHash, nil
}

func (api *EthereumAPI) handleSysError(function string, sysError common.SystemError) (responses.EnclaveResponse, error) {
	api.logger.Error(fmt.Sprintf("Enclave System Error. Function %s", function), log.ErrKey, sysError)
	return responses.EnclaveResponse{
		Err: &responses.InternalErrMsg,
	}, nil
}
