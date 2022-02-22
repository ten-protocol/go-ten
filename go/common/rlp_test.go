package common

import (
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/google/uuid"
	"testing"
)


func TestSerialiseBlock(t *testing.T) {
	tx := &L1Tx{
		Id:     uuid.New(),
		TxType: DepositTx,
		Amount: 100,
		//Dest:   wallet_mock.New().Address,
	}
	block := Block{
		Header:       GenesisBlock.Header,
		Transactions: Transactions{tx},
		hash:         GenesisBlock.hash,
		height:       GenesisBlock.height,
		size:         GenesisBlock.size,
	}
	bytes, err := rlp.EncodeToBytes(block)
	if err != nil {
		panic(err)
	}
	b1 := Block{Transactions: Transactions{tx}}
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
	tx := &L1Tx{
		Id:     uuid.New(),
		TxType: DepositTx,
		Amount: 100,
		//Dest:   wallet_mock.New().Address,
	}
	block := Block{
		Header:       GenesisBlock.Header,
		Transactions: Transactions{tx},
		hash:         GenesisBlock.hash,
		height:       GenesisBlock.height,
		size:         GenesisBlock.size,
	}
	bytes, err := rlp.EncodeToBytes(block)
	if err != nil {
		panic(err)
	}
	b1 := Block{Transactions: Transactions{tx}}
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
