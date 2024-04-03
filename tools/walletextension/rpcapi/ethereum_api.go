package rpcapi

import (
	"context"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

type EthereumAPI struct {
	we *Services
}

func NewEthereumAPI(we *Services,
) *EthereumAPI {
	return &EthereumAPI{we}
}

func (api *EthereumAPI) GasPrice(ctx context.Context) (*hexutil.Big, error) {
	return UnauthenticatedTenRPCCall[hexutil.Big](ctx, api.we, &CacheCfg{CacheType: LatestBatch}, "eth_gasPrice")
}

func (api *EthereumAPI) MaxPriorityFeePerGas(ctx context.Context) (*hexutil.Big, error) {
	return UnauthenticatedTenRPCCall[hexutil.Big](ctx, api.we, &CacheCfg{CacheType: LatestBatch}, "eth_maxPriorityFeePerGas")
}

type FeeHistoryResult struct {
	OldestBlock  *hexutil.Big     `json:"oldestBlock"`
	Reward       [][]*hexutil.Big `json:"reward,omitempty"`
	BaseFee      []*hexutil.Big   `json:"baseFeePerGas,omitempty"`
	GasUsedRatio []float64        `json:"gasUsedRatio"`
}

func (api *EthereumAPI) FeeHistory(ctx context.Context, blockCount math.HexOrDecimal64, lastBlock rpc.BlockNumber, rewardPercentiles []float64) (*FeeHistoryResult, error) {
	return UnauthenticatedTenRPCCall[FeeHistoryResult](
		ctx,
		api.we,
		&CacheCfg{CacheTypeDynamic: func() CacheStrategy {
			return cacheTTLBlockNumber(lastBlock)
		}},
		"eth_feeHistory",
		blockCount,
		lastBlock,
		rewardPercentiles,
	)
}

func (api *EthereumAPI) Syncing() (interface{}, error) {
	return nil, rpcNotImplemented
}
