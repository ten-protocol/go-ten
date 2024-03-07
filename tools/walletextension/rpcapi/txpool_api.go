package rpcapi

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	rpc2 "github.com/ten-protocol/go-ten/go/enclave/rpc"
)

// todo
type TxPoolAPI struct {
	we *Services
}

func NewTxPoolAPI(we *Services) *TxPoolAPI {
	return &TxPoolAPI{we}
}

func (s *TxPoolAPI) Content() map[string]map[string]map[string]*rpc2.RpcTransaction {
	content := map[string]map[string]map[string]*rpc2.RpcTransaction{
		"pending": make(map[string]map[string]*rpc2.RpcTransaction),
		"queued":  make(map[string]map[string]*rpc2.RpcTransaction),
	}
	return content
}

func (s *TxPoolAPI) ContentFrom(addr common.Address) map[string]map[string]*rpc2.RpcTransaction {
	content := make(map[string]map[string]*rpc2.RpcTransaction, 2)
	return content
}

func (s *TxPoolAPI) Status() map[string]hexutil.Uint {
	pending, queue := 0, 0
	return map[string]hexutil.Uint{
		"pending": hexutil.Uint(pending),
		"queued":  hexutil.Uint(queue),
	}
}

func (s *TxPoolAPI) Inspect() map[string]map[string]map[string]string {
	content := map[string]map[string]map[string]string{
		"pending": make(map[string]map[string]string),
		"queued":  make(map[string]map[string]string),
	}
	return content
}
