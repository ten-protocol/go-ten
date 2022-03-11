package common

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func EncodeBlockErr(b *types.Block) (EncodedBlock, error) {
	return rlp.EncodeToBytes(b)
}

func EncodeBlock(b *types.Block) EncodedBlock {
	encoded, err := EncodeBlockErr(b)
	if err != nil {
		panic(err)
	}

	return encoded
}

func (eb EncodedBlock) Decode() (*types.Block, error) {
	bl := types.Block{}
	err := rlp.DecodeBytes(eb, &bl)

	return &bl, err
}

func (eb EncodedBlock) DecodeBlock() *types.Block {
	b, err := eb.Decode()
	if err != nil {
		panic(err)
	}

	return b
}

func EncodeTx(tx *L1Tx) (EncodedL1Tx, error) {
	return rlp.EncodeToBytes(tx)
}

func (tx EncodedL1Tx) Decode() (L1Tx, error) {
	tx1 := L1Tx{}
	err := rlp.DecodeBytes(tx, &tx1)

	return tx1, err
}
