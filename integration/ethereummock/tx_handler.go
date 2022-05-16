package ethereummock

import (
	"bytes"
	"encoding/gob"

	"github.com/obscuronet/obscuro-playground/go/ethclient/txhandler"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

var (
	depositTxAddr       = common.HexToAddress("0x01")
	rollupTxAddr        = common.HexToAddress("0x02")
	storeSecretTxAddr   = common.HexToAddress("0x03")
	requestSecretTxAddr = common.HexToAddress("0x04")
)

// MockTxHandler implements mgmtcontractlib.TxHandler for the ethereummock package
// The ethereummock does not execute contracts. As such we need a way to emulate transactions that execute contracts.
// The MockTxHandler uses the To Field to make contract executions. Both the PackTx and UnPackTx know how to handle these.
//
// PackTx encodes the obscurocommon.L1Transaction in the Data field of the types.LegacyTx
// and specifies the obscurocommon.L1Transaction type in the To Field (since it is not used for anything else)
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

	// in the mock implementation we use the TO address field to specify what is the L1 operation
	contractAddress := common.Address{}
	switch tx.(type) {
	case *obscurocommon.L1RollupTx:
		contractAddress = rollupTxAddr
	case *obscurocommon.L1DepositTx:
		contractAddress = depositTxAddr
	case *obscurocommon.L1RequestSecretTx:
		contractAddress = requestSecretTxAddr
	case *obscurocommon.L1StoreSecretTx:
		contractAddress = storeSecretTxAddr
	}

	return &types.LegacyTx{
		Nonce: nonce,
		Data:  buf.Bytes(),
		To:    &contractAddress,
	}, nil
}

func (m *MockTxHandler) UnPackTx(tx *types.Transaction) obscurocommon.L1Transaction {
	if len(tx.Data()) == 0 {
		panic("Data cannot be 0 in the mock implementation")
	}

	// prepare byte buffer
	buf := bytes.NewBuffer(tx.Data())
	dec := gob.NewDecoder(buf)

	// in the mock implementation we use the To address field to specify what is the L1 operation
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

func NewMockTxHandler() txhandler.TxHandler {
	return &MockTxHandler{}
}
