package ethereummock

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/l1client/rollupcontractlib"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

// MockTxHandler implements rollupcontractlib.TxHandler for the ethereummock package
// it never PackTx (because the mock eth takes care of it)
// it always UnPackTx given the expected direct conversion of types.Transaction.Data -> obscurocommon.L1TxData
type MockTxHandler struct{}

func (m MockTxHandler) PackTx(tx *obscurocommon.L1TxData, from common.Address, nonce uint64) (types.TxData, error) {
	panic("implement me")
}

func (m MockTxHandler) UnPackTx(tx *types.Transaction) *obscurocommon.L1TxData {
	t := obscurocommon.TxData(tx)
	return &t
}

func NewMockTxHandler() rollupcontractlib.TxHandler {
	return &MockTxHandler{}
}
