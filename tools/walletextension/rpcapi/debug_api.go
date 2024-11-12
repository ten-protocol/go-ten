package rpcapi

import (
	"context"

	"github.com/ten-protocol/go-ten/tools/walletextension/cache"

	"github.com/ten-protocol/go-ten/tools/walletextension/services"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

type DebugAPI struct {
	we *services.Services
}

func NewDebugAPI(we *services.Services) *DebugAPI {
	return &DebugAPI{we}
}

func (api *DebugAPI) GetRawHeader(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	return nil, rpcNotImplemented
}

func (api *DebugAPI) GetRawBlock(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	return nil, rpcNotImplemented
}

func (api *DebugAPI) GetRawReceipts(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) ([]hexutil.Bytes, error) {
	return nil, rpcNotImplemented
}

func (s *DebugAPI) GetRawTransaction(ctx context.Context, hash gethcommon.Hash) (hexutil.Bytes, error) {
	return nil, rpcNotImplemented
}

func (api *DebugAPI) PrintBlock(ctx context.Context, number uint64) (string, error) {
	return "", rpcNotImplemented
}

func (api *DebugAPI) ChaindbProperty(property string) (string, error) {
	return "", rpcNotImplemented
}

func (api *DebugAPI) ChaindbCompact() error {
	return rpcNotImplemented
}

func (api *DebugAPI) SetHead(number hexutil.Uint64) {
	// not implemented
}

// EventLogRelevancy - specific to TEN
// the FilterCriteria must have a single contract address and only the topics on the position 0 ( event signatures)
// the caller must be the contract deployer
// Intended for debug purposes for the owner of a contract. It doesn't reveal any user information.
func (api *DebugAPI) EventLogRelevancy(ctx context.Context, crit common.FilterCriteria) ([]*common.DebugLogVisibility, error) {
	l, err := ExecAuthRPC[[]*common.DebugLogVisibility](
		ctx,
		api.we,
		&AuthExecCfg{
			cacheCfg: &cache.Cfg{
				Type: cache.NoCache,
			},
			tryUntilAuthorised: true,
		},
		"debug_eventLogRelevancy",
		common.SerializableFilterCriteria(crit),
	)
	if err != nil {
		return nil, err
	}
	return *l, nil
}
