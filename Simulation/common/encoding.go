package common

import "github.com/ethereum/go-ethereum/rlp"

func (b Block) Encode() (EncodedBlock, error) {
	return rlp.EncodeToBytes(b)
}

func (b EncodedBlock) Decode() (*Block, error) {
	bl := Block{}
	err := rlp.DecodeBytes(b, &bl)
	return &bl, err
}

func (b Block) EncodeBlock() EncodedBlock {
	encoded, err := b.Encode()
	if err != nil {
		panic(err)
	}
	return encoded
}

func (block EncodedBlock) DecodeBlock() *Block {
	b, err := block.Decode()
	if err != nil {
		panic(err)
	}
	return b
}

func (tx L1Tx) Encode() (EncodedL1Tx, error) {
	return rlp.EncodeToBytes(tx)
}

func (tx EncodedL1Tx) Decode() (L1Tx, error) {
	tx1 := L1Tx{}
	err := rlp.DecodeBytes(tx, &tx1)
	return tx1, err
}
