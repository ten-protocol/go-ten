package enclave

import (
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/google/uuid"
	common2 "github.com/obscuronet/obscuro-playground/go/obscuronode/common"
	wallet_mock "github.com/obscuronet/obscuro-playground/integration/walletmock"
)

func TestSerialiseL2Tx(t *testing.T) {
	tx := L2Tx{
		ID:     uuid.New(),
		TxType: TransferTx,
		Amount: 100,
		From:   wallet_mock.New().Address,
		To:     wallet_mock.New().Address,
	}
	bytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		panic(err)
	}
	tx1 := L2Tx{}
	err2 := rlp.DecodeBytes(bytes, &tx1)
	if err2 != nil {
		panic(err2)
	}
	if tx1.ID != tx.ID {
		t.Errorf("tx deserialized incorrectly\n")
	}
}

func TestSerialiseRollup(t *testing.T) {
	tx := L2Tx{
		ID:     uuid.New(),
		TxType: TransferTx,
		Amount: 100,
		From:   wallet_mock.New().Address,
		To:     wallet_mock.New().Address,
	}
	height := atomic.Value{}
	height.Store(1)
	rollup := common2.Rollup{
		Header:       GenesisRollup.Header,
		Height:       height,
		Transactions: EncryptTransactions(Transactions{tx}),
	}
	_, read, err := rlp.EncodeToReader(&rollup)
	if err != nil {
		panic(err)
	}
	r1 := common2.Rollup{}

	err = rlp.Decode(read, &r1)
	if err != nil {
		panic(err)
	}
	if r1.Hash() != rollup.Hash() {
		t.Errorf("rollup deserialized incorrectly\n")
	}
	//if r1.Transactions[0].ID != rollup.Transactions[0].ID {
	//	t.Errorf("rollup deserialized incorrectly\n")
	//}
}
