package rpcapi

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

type DebugAPI struct {
	we *Services
}

func NewDebugAPI(we *Services) *DebugAPI {
	return &DebugAPI{we}
}

/*func (api *DebugAPI) GetRawHeader(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	// not implemented
	return nil, nil
}

func (api *DebugAPI) GetRawBlock(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	// not implemented
	return nil, nil
}

func (api *DebugAPI) GetRawReceipts(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) ([]hexutil.Bytes, error) {
	// not implemented
	return nil, nil
}

func (s *DebugAPI) GetRawTransaction(ctx context.Context, hash common.Hash) (hexutil.Bytes, error) {
	// not implemented
	return nil, nil
}

func (api *DebugAPI) PrintBlock(ctx context.Context, number uint64) (string, error) {
	// not implemented
	return "", nil
}

func (api *DebugAPI) ChaindbProperty(property string) (string, error) {
	// not implemented
	return "", nil
}

func (api *DebugAPI) ChaindbCompact() error {
	// not implemented
	return nil
}

func (api *DebugAPI) SetHead(number hexutil.Uint64) {
	// not implemented
}
*/

// EventLogRelevancy - specific to Ten - todo
func (api *DebugAPI) EventLogRelevancy(_ context.Context, _ common.Hash) (interface{}, error) {
	// todo
	return nil, nil
}
