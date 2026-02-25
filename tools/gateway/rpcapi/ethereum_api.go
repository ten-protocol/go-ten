package rpcapi

import (
	"context"

	tenrpc "github.com/ten-protocol/go-ten/go/rpc"

	"github.com/ten-protocol/go-ten/tools/walletextension/cache"

	"github.com/ten-protocol/go-ten/tools/walletextension/services"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

type EthereumAPI struct {
	we *services.Services
}

func NewEthereumAPI(we *services.Services,
) *EthereumAPI {
	return &EthereumAPI{we}
}

func (api *EthereumAPI) GasPrice(ctx context.Context) (*hexutil.Big, error) {
	return UnauthenticatedTenRPCCall[hexutil.Big](ctx, api.we, &cache.Cfg{Type: cache.LatestBatch}, tenrpc.GasPrice)
}

func (api *EthereumAPI) MaxPriorityFeePerGas(ctx context.Context) (*hexutil.Big, error) {
	return UnauthenticatedTenRPCCall[hexutil.Big](ctx, api.we, &cache.Cfg{Type: cache.LatestBatch}, "ten_maxPriorityFeePerGas")
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
		&cache.Cfg{DynamicType: func() cache.Strategy {
			return cacheBlockNumber(lastBlock)
		}},
		"ten_feeHistory",
		blockCount,
		lastBlock,
		rewardPercentiles,
	)
}

func (api *EthereumAPI) Syncing() (interface{}, error) {
	return nil, ErrRPCNotImplemented
}
