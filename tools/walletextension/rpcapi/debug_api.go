package rpcapi

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

type DebugAPI struct {
	we *Services
}

func NewDebugAPI(we *Services) *DebugAPI {
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

func (s *DebugAPI) GetRawTransaction(ctx context.Context, hash common.Hash) (hexutil.Bytes, error) {
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

// EventLogRelevancy - specific to TEN - todo
func (api *DebugAPI) EventLogRelevancy(_ context.Context, _ common.Hash) (interface{}, error) {
	// todo
	return nil, rpcNotImplemented
}
