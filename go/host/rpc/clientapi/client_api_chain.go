package clientapi

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/host"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/responses"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
)

// ChainAPI exposes public chain data
type ChainAPI struct {
	host   host.Host
	logger gethlog.Logger
}

func NewChainAPI(host host.Host, logger gethlog.Logger) *ChainAPI {
	return &ChainAPI{
		host:   host,
		logger: logger,
	}
}

// ChainId returns the Obscuro chain ID.
func (api *ChainAPI) ChainId() (*hexutil.Big, error) { //nolint:stylecheck,revive
	return (*hexutil.Big)(big.NewInt(api.host.Config().ObscuroChainID)), nil
}

// BlockNumber returns the height of the current head batch.
func (api *ChainAPI) BlockNumber() hexutil.Uint64 {
	header, err := api.host.Storage().FetchHeadBatchHeader()
	if err != nil {
		// This error may be nefarious, but unfortunately the Eth API doesn't allow us to return an error.
		api.logger.Error("could not retrieve head batch header", log.ErrKey, err)
		return 0
	}
	return hexutil.Uint64(header.Number.Uint64())
}

// GetBlockByNumber returns the header of the batch with the given height.
func (api *ChainAPI) GetBlockByNumber(ctx context.Context, number rpc.BlockNumber, _ bool) (*common.BatchHeader, error) {
	batchHash, err := api.batchNumberToBatchHash(number)
	if err != nil {
		return nil, fmt.Errorf("could not find batch with height %d. Cause: %w", number, err)
	}
	return api.GetBlockByHash(ctx, *batchHash, true)
}

// GetBlockByHash returns the header of the batch with the given hash.
func (api *ChainAPI) GetBlockByHash(_ context.Context, hash gethcommon.Hash, _ bool) (*common.BatchHeader, error) {
	batchHeader, err := api.host.Storage().FetchBatchHeaderByHash(hash)
	if err != nil {
		return nil, err
	}
	return batchHeader, nil
}

// GasPrice is a placeholder for an RPC method required by MetaMask/Remix.
func (api *ChainAPI) GasPrice(context.Context) (*hexutil.Big, error) {
	header, err := api.host.Storage().FetchHeadBatchHeader()
	if err != nil {
		return nil, err
	}

	if header.BaseFee == nil || header.BaseFee.Cmp(gethcommon.Big0) == 0 {
		return (*hexutil.Big)(big.NewInt(params.InitialBaseFee)), nil
	}

	return (*hexutil.Big)(big.NewInt(0).Set(header.BaseFee)), nil
}

// GetCode returns the code stored at the given address in the state for the given batch height or batch hash.
// todo (#1620) - instead of converting the block number of hash client-side, do it on the enclave
func (api *ChainAPI) GetCode(ctx context.Context, address gethcommon.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	code, sysError := api.host.EnclaveClient().GetCode(ctx, address, blockNrOrHash)
	if sysError != nil {
		api.logger.Error(fmt.Sprintf("Enclave System Error. Function %s", "GetCode"), log.ErrKey, sysError)
		return nil, fmt.Errorf(responses.InternalErrMsg)
	}

	return code, nil
}

func (api *ChainAPI) MaxPriorityFeePerGas(_ context.Context) (*hexutil.Big, error) {
	// todo - implement with the gas mechanics
	header, err := api.host.Storage().FetchHeadBatchHeader()
	if err != nil {
		api.logger.Error("Unable to retrieve header for fee history.", log.ErrKey, err)
		return nil, fmt.Errorf("unable to retrieve MaxPriorityFeePerGas")
	}

	// just return the base fee?
	return (*hexutil.Big)(header.BaseFee), err
}

// FeeHistory is a placeholder for an RPC method required by MetaMask/Remix.
// rpc.DecimalOrHex -> []byte
func (api *ChainAPI) FeeHistory(context.Context, string, rpc.BlockNumber, []float64) (*FeeHistoryResult, error) {
	// todo (#1621) - return a non-dummy fee history
	header, err := api.host.Storage().FetchHeadBatchHeader()
	if err != nil {
		api.logger.Error("Unable to retrieve header for fee history.", log.ErrKey, err)
		return nil, fmt.Errorf("unable to retrieve fee history")
	}

	batches := make([]*common.BatchHeader, 0)
	batches = append(batches, header)

	feeHist := &FeeHistoryResult{
		OldestBlock:  (*hexutil.Big)(header.Number),
		Reward:       [][]*hexutil.Big{},
		BaseFee:      []*hexutil.Big{},
		GasUsedRatio: []float64{},
	}

	for _, header := range batches {
		// 0.9 - This number represents how full the block is. As we dont have a dynamic base fee, we tell whomever is requesting that
		// we expect the baseFee to increase, rather than decrease in order to avoid underpriced transactions.
		feeHist.GasUsedRatio = append(feeHist.GasUsedRatio, 0.9)
		feeHist.BaseFee = append(feeHist.BaseFee, (*hexutil.Big)(header.BaseFee))
	}
	return feeHist, nil
}

// FeeHistoryResult is the structure returned by Geth `eth_feeHistory` API.
type FeeHistoryResult struct {
	OldestBlock  *hexutil.Big     `json:"oldestBlock"`
	Reward       [][]*hexutil.Big `json:"reward,omitempty"`
	BaseFee      []*hexutil.Big   `json:"baseFeePerGas,omitempty"`
	GasUsedRatio []float64        `json:"gasUsedRatio"`
}

// Given a batch number, returns the hash of the batch with that number.
func (api *ChainAPI) batchNumberToBatchHash(batchNumber rpc.BlockNumber) (*gethcommon.Hash, error) {
	// Handling the special cases first. No special handling is required for rpc.EarliestBlockNumber.
	// note: our API currently treats all these block statuses the same for obscuro batches
	if batchNumber == rpc.LatestBlockNumber || batchNumber == rpc.PendingBlockNumber ||
		batchNumber == rpc.FinalizedBlockNumber || batchNumber == rpc.SafeBlockNumber {
		batchHeader, err := api.host.Storage().FetchHeadBatchHeader()
		if err != nil {
			return nil, err
		}
		batchHash := batchHeader.Hash()
		return &batchHash, nil
	}
	batchNumberBig := big.NewInt(batchNumber.Int64())
	batchHash, err := api.host.Storage().FetchBatchHashByHeight(batchNumberBig)
	if err != nil {
		return nil, err
	}
	return batchHash, nil
}
