package ethereummock

import (
	"bytes"
	"encoding/gob"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/ethclient/txencoder"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

var (
	depositTxAddr       = common.HexToAddress("0x01")
	rollupTxAddr        = common.HexToAddress("0x02")
	storeSecretTxAddr   = common.HexToAddress("0x03")
	requestSecretTxAddr = common.HexToAddress("0x04")
)

type mockEncoder struct {
}

func (m *mockEncoder) CreateRollup(tx *obscurocommon.L1RollupTx, nonce uint64) types.TxData {
	return encodeTx(tx, nonce, rollupTxAddr)
}

func (m *mockEncoder) CreateRequestSecret(tx *obscurocommon.L1RequestSecretTx, nonce uint64) types.TxData {
	return encodeTx(tx, nonce, requestSecretTxAddr)
}

func (m *mockEncoder) CreateStoreSecret(tx *obscurocommon.L1StoreSecretTx, nonce uint64) types.TxData {
	return encodeTx(tx, nonce, storeSecretTxAddr)
}

func (m *mockEncoder) CreateDepositTx(tx *obscurocommon.L1DepositTx, nonce uint64) types.TxData {
	return encodeTx(tx, nonce, depositTxAddr)
}

func NewMockTxEncoder() txencoder.TxEncoder {
	return &mockEncoder{}
}

func encodeTx(tx obscurocommon.L1Transaction, nonce uint64, opType common.Address) types.TxData {
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
