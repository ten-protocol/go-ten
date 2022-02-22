package enclave

import (
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/google/uuid"
	common2 "github.com/otherview/obscuro-playground/go/obscuro-node/common"
	wallet_mock "github.com/otherview/obscuro-playground/integration/wallet-mock"
	"sync/atomic"
	"testing"
)

func TestSerialiseL2Tx(t *testing.T) {
	tx := L2Tx{
		Id:     uuid.New(),
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
	if tx1.Id != tx.Id {
		t.Errorf("tx deserialized incorrectly\n")
	}
}

func TestSerialiseRollup(t *testing.T) {
	tx := L2Tx{
		Id:     uuid.New(),
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

	err2 := rlp.Decode(read, &r1)

	if err2 != nil {
		panic(err2)
	}
	if r1.Hash() != rollup.Hash() {
		t.Errorf("rollup deserialized incorrectly\n")
	}
	//if r1.Transactions[0].Id != rollup.Transactions[0].Id {
	//	t.Errorf("rollup deserialized incorrectly\n")
	//}
}