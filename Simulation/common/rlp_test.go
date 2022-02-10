package common

import (
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/google/uuid"
	"simulation/obscuro"
	wallet_mock "simulation/wallet-mock"
	"testing"
	"time"
)

func TestSerialiseL2Tx(t *testing.T) {
	tx := obscuro.L2Tx{
		Id:     uuid.New(),
		TxType: obscuro.TransferTx,
		Amount: 100,
		From:   wallet_mock.New().Address,
		To:     wallet_mock.New().Address,
	}
	bytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		panic(err)
	}
	tx1 := obscuro.L2Tx{}
	err2 := rlp.DecodeBytes(bytes, &tx1)
	if err2 != nil {
		panic(err2)
	}
	if tx1.Id != tx.Id {
		t.Errorf("tx deserialized incorrectly\n")
	}
}

func TestSerialiseRollup(t *testing.T) {
	tx := obscuro.L2Tx{
		Id:     uuid.New(),
		TxType: obscuro.TransferTx,
		Amount: 100,
		From:   wallet_mock.New().Address,
		To:     wallet_mock.New().Address,
	}
	rollup := obscuro.Rollup{
		Height:       1,
		RootHash:     uuid.New(),
		Agg:          1,
		ParentHash:   uuid.New(),
		CreationTime: time.Now(),
		L1Proof:      uuid.New(),
		Nonce:        100,
		State:        "",
		Withdrawals:  nil,
		Transactions: []obscuro.L2Tx{tx},
	}
	_, read, err := rlp.EncodeToReader(&rollup)
	if err != nil {
		panic(err)
	}
	r1 := obscuro.Rollup{}

	err2 := rlp.Decode(read, &r1)

	if err2 != nil {
		panic(err2)
	}
	if r1.RootHash != rollup.RootHash {
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
	if b1.RootHash != block.RootHash {
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
	if b1.RootHash != block.RootHash {
		t.Errorf("block deserialized incorrectly\n")
	}
	if b1.Transactions[0].Id != block.Transactions[0].Id {
		t.Errorf("block deserialized incorrectly\n")
	}
}
