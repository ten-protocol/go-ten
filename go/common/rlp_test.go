package common

import (
	"testing"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/google/uuid"
)

func TestSerialiseBlock(t *testing.T) {
	testTx := &L1Tx{
		ID:     uuid.New(),
		TxType: DepositTx,
		Amount: 100,
		// Dest:   wallet_mock.New().Address,
	}
	block := Block{
		Header:       GenesisBlock.Header,
		Transactions: Transactions{testTx},
		hash:         GenesisBlock.hash,
		height:       GenesisBlock.height,
		size:         GenesisBlock.size,
	}
	bytes, err := rlp.EncodeToBytes(block)
	if err != nil {
		panic(err)
	}
	block1 := Block{Transactions: Transactions{testTx}}
	err2 := rlp.DecodeBytes(bytes, &block1)
	if err2 != nil {
		panic(err2)
	}
	if block1.Hash() != block.Hash() {
		t.Errorf("block deserialized incorrectly\n")
	}
	if block1.Transactions[0].ID != block.Transactions[0].ID {
		t.Errorf("block deserialized incorrectly\n")
	}
}

func TestPlay(t *testing.T) {
	testTx := &L1Tx{
		ID:     uuid.New(),
		TxType: DepositTx,
		Amount: 100,
		// Dest:   wallet_mock.New().Address,
	}
	block := Block{
		Header:       GenesisBlock.Header,
		Transactions: Transactions{testTx},
		hash:         GenesisBlock.hash,
		height:       GenesisBlock.height,
		size:         GenesisBlock.size,
	}
	bytes, err := rlp.EncodeToBytes(block)
	if err != nil {
		panic(err)
	}
	block1 := Block{Transactions: Transactions{testTx}}
	err2 := rlp.DecodeBytes(bytes, &block1)
	if err2 != nil {
		panic(err2)
	}
	if block1.Hash() != block.Hash() {
		t.Errorf("block deserialized incorrectly\n")
	}
	if block1.Transactions[0].ID != block.Transactions[0].ID {
		t.Errorf("block deserialized incorrectly\n")
	}
}
