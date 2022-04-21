package ethereummock

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/l1client/txhandler"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

type MockTxHandler struct {
}

func (m MockTxHandler) PackTx(tx *obscurocommon.L1TxData, from common.Address, nonce uint64) (types.TxData, error) {
	panic("implement me")
}

func (m MockTxHandler) UnPackTx(tx *types.Transaction) *obscurocommon.L1TxData {
	t := obscurocommon.TxData(tx)
	return &t
}

func NewMockTxHandler() txhandler.TxHandler {
	return &MockTxHandler{}
}
