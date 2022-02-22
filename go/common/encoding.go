package common

import "github.com/ethereum/go-ethereum/rlp"

func (b Block) Encode() (EncodedBlock, error) {
	return rlp.EncodeToBytes(b)
}

func (b Block) EncodeBlock() EncodedBlock {
	encoded, err := b.Encode()
	if err != nil {
		panic(err)
	}

	return encoded
}

func (eb EncodedBlock) Decode() (*Block, error) {
	bl := Block{}
	err := rlp.DecodeBytes(eb, &bl)

	return &bl, err
}

func (eb EncodedBlock) DecodeBlock() *Block {
	b, err := eb.Decode()
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
