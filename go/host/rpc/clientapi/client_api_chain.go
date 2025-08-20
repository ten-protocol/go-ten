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

const (
	// DefaultMinGasPrice is the default minimum gas price for TEN network (10 Gwei)
	// This matches the value from go/config/defaults/0-base-config.yaml
	DefaultMinGasPrice = 10000000 // 10 Gwei in wei
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
	return (*hexutil.Big)(big.NewInt(api.host.Config().TenChainID)), nil
}

// BatchNumber returns the height of the current head batch.
func (api *ChainAPI) BatchNumber() hexutil.Uint64 {
	enclaveHeadBatch, err := api.host.ConfirmedHeadBatch()
	if err != nil {
		api.logger.Error("could not retrieve head batch header", log.ErrKey, err)
		return 0
	}
	return hexutil.Uint64(enclaveHeadBatch.Number.Uint64())
}

// GetBatchByNumber returns the header of the batch with the given height.
func (api *ChainAPI) GetBatchByNumber(ctx context.Context, number rpc.BlockNumber, _ bool) (*common.BatchHeader, error) {
	batchHash, err := api.batchNumberToBatchHash(number)
	if err != nil {
		return nil, fmt.Errorf("could not find batch with height %d. Cause: %w", number, err)
	}
	return api.GetBatchByHash(ctx, *batchHash, true)
}

// GetBatchByHash returns the header of the batch with the given hash.
func (api *ChainAPI) GetBatchByHash(_ context.Context, hash gethcommon.Hash, _ bool) (*common.BatchHeader, error) {
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

// FeeHistory returns historical gas fee data for the requested block range.
// This implements the eth_feeHistory RPC method as specified in EIP-1559.
// Note from Stefan - Method is AI generated according to spec, looks correct.
func (api *ChainAPI) FeeHistory(ctx context.Context, blockCountStr string, newestBlock rpc.BlockNumber, rewardPercentiles []float64) (*FeeHistoryResult, error) {
	// Parse and validate blockCount
	blockCount, err := hexutil.DecodeUint64(blockCountStr)
	if err != nil {
		return nil, fmt.Errorf("invalid blockCount: %w", err)
	}

	// Validate blockCount is within acceptable range (1 to 1024)
	if blockCount == 0 || blockCount > 1024 {
		return nil, fmt.Errorf("blockCount must be between 1 and 1024, got %d", blockCount)
	}

	// Validate reward percentiles are in ascending order and within [0, 100]
	for i, percentile := range rewardPercentiles {
		if percentile < 0 || percentile > 100 {
			return nil, fmt.Errorf("reward percentile must be between 0 and 100, got %f", percentile)
		}
		if i > 0 && percentile < rewardPercentiles[i-1] {
			return nil, fmt.Errorf("reward percentiles must be in ascending order")
		}
	}

	// Get the newest block header to determine the range
	var newestHeader *common.BatchHeader
	if newestBlock == rpc.LatestBlockNumber || newestBlock == rpc.PendingBlockNumber {
		newestHeader, err = api.host.Storage().FetchHeadBatchHeader()
		if err != nil {
			api.logger.Error("Unable to retrieve head batch header for fee history.", log.ErrKey, err)
			return nil, fmt.Errorf("unable to retrieve latest block")
		}
	} else {
		blockNum := uint64(newestBlock.Int64())
		newestHeader, err = api.host.Storage().FetchBatchHeaderByHeight(big.NewInt(int64(blockNum)))
		if err != nil {
			api.logger.Error("Unable to retrieve batch header for fee history.", log.ErrKey, err, "blockNumber", blockNum)
			return nil, fmt.Errorf("unable to retrieve block %d", blockNum)
		}
	}

	// Calculate the range of blocks to fetch
	newestBlockNumber := newestHeader.Number.Uint64()
	var oldestBlockNumber uint64
	if newestBlockNumber >= blockCount-1 {
		oldestBlockNumber = newestBlockNumber - blockCount + 1
	} else {
		// If we don't have enough blocks, start from genesis
		oldestBlockNumber = 0
		blockCount = newestBlockNumber + 1 // Adjust blockCount to available blocks
	}

	// Fetch the block headers for the range
	headers := make([]*common.BatchHeader, 0, blockCount)
	for blockNum := oldestBlockNumber; blockNum <= newestBlockNumber; blockNum++ {
		header, err := api.host.Storage().FetchBatchHeaderByHeight(big.NewInt(int64(blockNum)))
		if err != nil {
			api.logger.Error("Unable to retrieve batch header.", log.ErrKey, err, "blockNumber", blockNum)
			continue // Skip missing blocks
		}
		headers = append(headers, header)
	}

	if len(headers) == 0 {
		return nil, fmt.Errorf("no blocks found in range")
	}

	// Calculate base fees (includes one extra for the next block)
	baseFees := make([]*hexutil.Big, 0, len(headers)+1)
	gasUsedRatios := make([]float64, 0, len(headers))
	rewards := make([][]*hexutil.Big, 0, len(headers))

	for _, header := range headers {
		// Add base fee for this block
		baseFees = append(baseFees, (*hexutil.Big)(header.BaseFee))

		// Calculate gas used ratio
		if header.GasLimit > 0 {
			ratio := float64(header.GasUsed) / float64(header.GasLimit)
			gasUsedRatios = append(gasUsedRatios, ratio)
		} else {
			gasUsedRatios = append(gasUsedRatios, 0.0)
		}

		// Calculate reward percentiles if requested
		if len(rewardPercentiles) > 0 {
			blockRewards, err := api.calculateRewardPercentiles(ctx, header, rewardPercentiles)
			if err != nil {
				api.logger.Error("Unable to calculate reward percentiles.", log.ErrKey, err, "blockNumber", header.Number)
				// Use zero rewards for this block
				blockRewards = make([]*hexutil.Big, len(rewardPercentiles))
				for i := range blockRewards {
					blockRewards[i] = (*hexutil.Big)(big.NewInt(0))
				}
			}
			rewards = append(rewards, blockRewards)
		}
	}

	// Add the base fee for the next block (this is calculated from the latest block)
	// For simplicity, we'll use the same base fee as the latest block
	// In a real implementation, this should be calculated based on the latest block's gas usage
	if len(headers) > 0 {
		lastHeader := headers[len(headers)-1]
		nextBaseFee := api.calculateNextBaseFee(lastHeader)
		baseFees = append(baseFees, (*hexutil.Big)(nextBaseFee))
	}

	return &FeeHistoryResult{
		OldestBlock:  (*hexutil.Big)(big.NewInt(int64(oldestBlockNumber))),
		BaseFee:      baseFees,
		GasUsedRatio: gasUsedRatios,
		Reward:       rewards,
	}, nil
}

// calculateRewardPercentiles calculates the reward percentiles for a given block.
// This is a simplified implementation that returns zero rewards since TEN doesn't
// have traditional miner rewards in the same way as Ethereum.
func (api *ChainAPI) calculateRewardPercentiles(ctx context.Context, header *common.BatchHeader, percentiles []float64) ([]*hexutil.Big, error) {
	// For TEN, since we don't have traditional priority fees like Ethereum,
	// we return zero rewards for all percentiles
	rewards := make([]*hexutil.Big, len(percentiles))
	for i := range rewards {
		rewards[i] = (*hexutil.Big)(big.NewInt(0))
	}
	return rewards, nil
}

// calculateNextBaseFee calculates what the base fee should be for the next block
// Simple EIP-1559 estimation: if gas usage > 50% of limit, increase by ~12.5%, otherwise decrease
func (api *ChainAPI) calculateNextBaseFee(header *common.BatchHeader) *big.Int {
	currentBaseFee := header.BaseFee
	gasUsed := header.GasUsed
	gasLimit := header.GasLimit

	// Gas target is 50% of the limit (EIP-1559)
	gasTarget := gasLimit / 2

	var nextBaseFee *big.Int

	if gasUsed > gasTarget {
		// Above target: increase base fee by up to 12.5%
		// Simple approximation: increase by (gasUsed - gasTarget) / gasTarget * 12.5%
		excess := gasUsed - gasTarget
		increase := new(big.Int).Mul(currentBaseFee, big.NewInt(int64(excess)))
		increase.Div(increase, big.NewInt(int64(gasTarget)))
		increase.Div(increase, big.NewInt(8)) // Divide by 8 for 12.5% max

		nextBaseFee = new(big.Int).Add(currentBaseFee, increase)
	} else if gasUsed < gasTarget {
		// Below target: decrease base fee by up to 12.5%
		deficit := gasTarget - gasUsed
		decrease := new(big.Int).Mul(currentBaseFee, big.NewInt(int64(deficit)))
		decrease.Div(decrease, big.NewInt(int64(gasTarget)))
		decrease.Div(decrease, big.NewInt(8)) // Divide by 8 for 12.5% max

		nextBaseFee = new(big.Int).Sub(currentBaseFee, decrease)
	} else {
		// Exactly at target: no change
		nextBaseFee = new(big.Int).Set(currentBaseFee)
	}

	// Floor at minimum gas price
	minGasPrice := big.NewInt(DefaultMinGasPrice)
	if nextBaseFee.Cmp(minGasPrice) < 0 {
		nextBaseFee = minGasPrice
	}

	return nextBaseFee
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
