package nodecommon

import (
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

func EncodeRollup(r *Rollup) obscurocommon.EncodedRollup {
	encoded, err := rlp.EncodeToBytes(r)
	if err != nil {
		panic(err)
	}

	return encoded
}

func DecodeRollup(encoded obscurocommon.EncodedRollup) (*Rollup, error) {
	r := Rollup{}
	err := rlp.DecodeBytes(encoded, &r)

	return &r, err
}

func DecodeRollupOrPanic(rollup obscurocommon.EncodedRollup) *Rollup {
	r, err := DecodeRollup(rollup)
	if err != nil {
		panic(err)
	}

	return r
}
