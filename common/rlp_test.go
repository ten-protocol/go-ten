package common

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/google/uuid"
	"github.com/otherview/obscuro-playground/obscuro/common"
	wallet_mock "github.com/otherview/obscuro-playground/wallet-mock"
)

func TestSerialiseL2Tx(t *testing.T) {
	tx := common.L2Tx{
		Id:     uuid.New(),
		TxType: common.TransferTx,
		Amount: 100,
		From:   wallet_mock.New().Address,
		To:     wallet_mock.New().Address,
	}
	bytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		panic(err)
	}
	tx1 := common.L2Tx{}
	err2 := rlp.DecodeBytes(bytes, &tx1)
	if err2 != nil {
		panic(err2)
	}
	if tx1.Id != tx.Id {
		t.Errorf("tx deserialized incorrectly\n")
	}
}

func TestSerialiseRollup(t *testing.T) {
	tx := common.L2Tx{
		Id:     uuid.New(),
		TxType: common.TransferTx,
		Amount: 100,
		From:   wallet_mock.New().Address,
		To:     wallet_mock.New().Address,
	}
	rollup := common.Rollup{
		Height:       1,
		RootHash:     uuid.New(),
		Agg:          1,
		ParentHash:   uuid.New(),
		CreationTime: time.Now(),
		L1Proof:      uuid.New(),
		Nonce:        100,
		State:        "",
		Withdrawals:  nil,
		Transactions: []common.L2Tx{tx},
	}
	_, read, err := rlp.EncodeToReader(&rollup)
	if err != nil {
		panic(err)
	}
	r1 := common.Rollup{}

	err2 := rlp.Decode(read, &r1)

	if err2 != nil {
		panic(err2)
	}
	if r1.Hash() != rollup.Hash() {
		t.Errorf("rollup deserialized incorrectly\n")
	}
	if r1.Transactions[0].Id != rollup.Transactions[0].Id {
		t.Errorf("rollup deserialized incorrectly\n")
	}
}

func TestSerialiseBlock(t *testing.T) {
	tx := L1Tx{
		Id:     uuid.New(),
		TxType: DepositTx,
		Amount: 100,
		Dest:   wallet_mock.New().Address,
	}
	block := Block{
		Height:       1,
		RootHash:     uuid.New(),
		Miner:        1,
		ParentHash:   uuid.New(),
		ReceiveTime:  time.Now(),
		Nonce:        100,
		Transactions: []L1Tx{tx},
	}
	bytes, err := rlp.EncodeToBytes(block)
	if err != nil {
		panic(err)
	}
	b1 := Block{Transactions: make([]L1Tx, 1)}
	err2 := rlp.DecodeBytes(bytes, &b1)
	if err2 != nil {
		panic(err2)
	}
	if b1.Hash() != block.Hash() {
		t.Errorf("block deserialized incorrectly\n")
	}
	if b1.Transactions[0].Id != block.Transactions[0].Id {
		t.Errorf("block deserialized incorrectly\n")
	}
}

func TestPlay(t *testing.T) {
	tx := L1Tx{
		Id:     uuid.New(),
		TxType: DepositTx,
		Amount: 100,
		Dest:   wallet_mock.New().Address,
	}
	block := Block{
		Height:       1,
		RootHash:     uuid.New(),
		Miner:        1,
		ParentHash:   uuid.New(),
		ReceiveTime:  time.Now(),
		Nonce:        100,
		Transactions: []L1Tx{tx},
	}
	bytes, err := rlp.EncodeToBytes(block)
	if err != nil {
		panic(err)
	}
	b1 := Block{Transactions: make([]L1Tx, 1)}
	err2 := rlp.DecodeBytes(bytes, &b1)
	if err2 != nil {
		panic(err2)
	}
	if b1.Hash() != block.Hash() {
		t.Errorf("block deserialized incorrectly\n")
	}
	if b1.Transactions[0].Id != block.Transactions[0].Id {
		t.Errorf("block deserialized incorrectly\n")
	}
}
