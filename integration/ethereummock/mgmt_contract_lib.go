package ethereummock

import (
	"bytes"
	"encoding/gob"

	"github.com/obscuronet/go-obscuro/go/ethadapter"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
)

var (
	depositTxAddr          = gethcommon.HexToAddress("0x01")
	rollupTxAddr           = gethcommon.HexToAddress("0x02")
	storeSecretTxAddr      = gethcommon.HexToAddress("0x03")
	requestSecretTxAddr    = gethcommon.HexToAddress("0x04")
	initializeSecretTxAddr = gethcommon.HexToAddress("0x05")
)

// mockContractLib is an implementation of the mgmtcontractlib.MgmtContractLib
// it creates ethereum mocked transactions from common.L1Transaction
// and converts ethereum mocked transactions to common.L1Transaction
type mockContractLib struct{}

func NewMgmtContractLibMock() mgmtcontractlib.MgmtContractLib {
	return &mockContractLib{}
}

func (m *mockContractLib) DecodeTx(tx *types.Transaction) ethadapter.L1Transaction {
	return decodeTx(tx)
}

func (m *mockContractLib) CreateRollup(tx *ethadapter.L1RollupTx, nonce uint64) types.TxData {
	return encodeTx(tx, nonce, rollupTxAddr)
}

func (m *mockContractLib) CreateRequestSecret(tx *ethadapter.L1RequestSecretTx, nonce uint64) types.TxData {
	return encodeTx(tx, nonce, requestSecretTxAddr)
}

func (m *mockContractLib) CreateRespondSecret(tx *ethadapter.L1RespondSecretTx, nonce uint64, _ bool) types.TxData {
	return encodeTx(tx, nonce, storeSecretTxAddr)
}

func (m *mockContractLib) CreateInitializeSecret(tx *ethadapter.L1InitializeSecretTx, nonce uint64) types.TxData {
	return encodeTx(tx, nonce, initializeSecretTxAddr)
}

func (m *mockContractLib) GetHostAddresses() (ethereum.CallMsg, error) {
	return ethereum.CallMsg{}, nil
}

func (m *mockContractLib) DecodeCallResponse([]byte) ([][]string, error) {
	return [][]string{{""}}, nil
}

func decodeTx(tx *types.Transaction) ethadapter.L1Transaction {
	if len(tx.Data()) == 0 {
		panic("Data cannot be 0 in the mock implementation")
	}

	// prepare byte buffer
	buf := bytes.NewBuffer(tx.Data())
	dec := gob.NewDecoder(buf)

	// in the mock implementation we use the To address field to specify the L1 operation (rollup/storesecret/requestsecret)
	// the mock implementation does not process contracts
	// so this is a way that we can differentiate different contract calls
	var t ethadapter.L1Transaction
	switch tx.To().Hex() {
	case rollupTxAddr.Hex():
		t = &ethadapter.L1RollupTx{}
	case storeSecretTxAddr.Hex():
		t = &ethadapter.L1RespondSecretTx{}
	case depositTxAddr.Hex():
		t = &ethadapter.L1DepositTx{}
	case requestSecretTxAddr.Hex():
		t = &ethadapter.L1RequestSecretTx{}
	case initializeSecretTxAddr.Hex():
		t = &ethadapter.L1InitializeSecretTx{}
	default:
		panic("unexpected type")
	}

	// decode to interface implementation
	if err := dec.Decode(t); err != nil {
		panic(err)
	}
	return t
}

func encodeTx(tx ethadapter.L1Transaction, nonce uint64, opType gethcommon.Address) types.TxData {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(tx); err != nil {
		panic(err)
	}

	// the mock implementation does not process contract calls
	// this uses the To address to distinguish between different contract calls / different l1 transactions
	return &types.LegacyTx{
		Nonce: nonce,
		Data:  buf.Bytes(),
		To:    &opType,
	}
}
