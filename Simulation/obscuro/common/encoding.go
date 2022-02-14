package common

import (
	"github.com/ethereum/go-ethereum/rlp"
	"simulation/common"
)

func (r Rollup) Encode() (common.EncodedRollup, error) {
	return rlp.EncodeToBytes(r)
}

func Decode(encoded common.EncodedRollup) (*Rollup, error) {
	r := Rollup{}
	err := rlp.DecodeBytes(encoded, &r)
	return &r, err
}

func EncodeRollup(r *Rollup) common.EncodedRollup {
	encoded, err := r.Encode()
	if err != nil {
		panic(err)
	}
	return encoded
}

func DecodeRollup(rollup common.EncodedRollup) *Rollup {
	r, err := Decode(rollup)
	if err != nil {
		panic(err)
	}
	return r
}
