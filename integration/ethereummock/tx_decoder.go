package ethereummock

import (
	"bytes"
	"encoding/gob"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/ethclient/txdecoder"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

type mockTxDecoder struct{}

func (m *mockTxDecoder) DecodeTx(tx *types.Transaction) obscurocommon.L1Transaction {
	if len(tx.Data()) == 0 {
		panic("Data cannot be 0 in the mock implementation")
	}

	// prepare byte buffer
	buf := bytes.NewBuffer(tx.Data())
	dec := gob.NewDecoder(buf)

	// in the mock implementation we use the To address field to specify what is the L1 operation
	// the mock implementation does not process contracts
	// so this is a way that we can differentiate different contract calls
	var t obscurocommon.L1Transaction
	switch tx.To().Hex() {
	case rollupTxAddr.Hex():
		t = &obscurocommon.L1RollupTx{}
	case storeSecretTxAddr.Hex():
		t = &obscurocommon.L1StoreSecretTx{}
	case depositTxAddr.Hex():
		t = &obscurocommon.L1DepositTx{}
	case requestSecretTxAddr.Hex():
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

func NewMockTxDecoder() txdecoder.TxDecoder {
	return &mockTxDecoder{}
}
