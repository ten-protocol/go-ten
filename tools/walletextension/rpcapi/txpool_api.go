package rpcapi

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	rpc2 "github.com/ten-protocol/go-ten/go/enclave/rpc"
	"github.com/ten-protocol/go-ten/tools/walletextension/services"
)

type TxPoolAPI struct {
	we *services.Services
}

func NewTxPoolAPI(we *services.Services) *TxPoolAPI {
	return &TxPoolAPI{we}
}

func (s *TxPoolAPI) Content() map[string]map[string]map[string]*rpc2.RpcTransaction {
	// not implemented
	return nil
}

func (s *TxPoolAPI) ContentFrom(_ common.Address) map[string]map[string]*rpc2.RpcTransaction {
	// not implemented
	return nil
}

func (s *TxPoolAPI) Status() map[string]hexutil.Uint {
	// not implemented
	return nil
}

func (s *TxPoolAPI) Inspect() map[string]map[string]map[string]string {
	// not implemented
	return nil
}
