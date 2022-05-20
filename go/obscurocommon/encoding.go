package obscurocommon

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/log"
)

func EncodeBlock(b *types.Block) EncodedBlock {
	encoded, err := rlp.EncodeToBytes(b)
	if err != nil {
		log.Panic("could not encode block to bytes. Cause: %s", err)
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
		log.Panic("could not decode block from bytes. Cause: %s", err)
	}

	return b
}
