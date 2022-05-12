package ethereummock

import (
	"bytes"
	"encoding/gob"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/ethclient/mgmtcontractlib"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

const (
	depositTxID uint64 = iota
	rollupTxID
	storeSecretTxID
	requestSecretTxID
)

// MockTxHandler implements mgmtcontractlib.TxHandler for the ethereummock package
// PackTx encodes the obscurocommon.L1Transaction in the Data field of the types.LegacyTx
// and specifies the obscurocommon.L1Transaction type in the Gas Field (since it is not used for anything else)
//
// UnPackTx does the reverse steps - understands obscurocommon.L1Transaction to return based on the types.Transaction Gas
// and decodes the correct object
type MockTxHandler struct{}

func (m *MockTxHandler) PackTx(tx obscurocommon.L1Transaction, _ common.Address, nonce uint64) (types.TxData, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(tx); err != nil {
		panic(err)
	}

	gasType := uint64(0)

	switch tx.(type) {
	case *obscurocommon.L1RollupTx:
		gasType = rollupTxID
	case *obscurocommon.L1DepositTx:
		gasType = depositTxID
	case *obscurocommon.L1RequestSecretTx:
		gasType = requestSecretTxID
	case *obscurocommon.L1StoreSecretTx:
		gasType = storeSecretTxID
	}

	return &types.LegacyTx{
		Gas:   gasType,
		Nonce: nonce,
		Data:  buf.Bytes(),
	}, nil
}

func (m *MockTxHandler) UnPackTx(tx *types.Transaction) obscurocommon.L1Transaction {
	if len(tx.Data()) == 0 {
		panic("Data cannot be 0 in the mock implementation")
	}

	// prepare byte buffer
	buf := bytes.NewBuffer(tx.Data())
	dec := gob.NewDecoder(buf)

	// check type
	var t obscurocommon.L1Transaction
	switch tx.Gas() {
	case rollupTxID:
		t = &obscurocommon.L1RollupTx{}
	case storeSecretTxID:
		t = &obscurocommon.L1StoreSecretTx{}
	case depositTxID:
		t = &obscurocommon.L1DepositTx{}
	case requestSecretTxID:
		t = &obscurocommon.L1RequestSecretTx{}
	default:
		panic("unexpected type")
	}

	// decode to interface implementation
	if err := dec.Decode(t); err != nil {
		panic(err)
	}
	return t
}

func NewMockTxHandler() mgmtcontractlib.TxHandler {
	return &MockTxHandler{}
}
