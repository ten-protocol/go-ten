package nodecommon

import (
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

func Decode(encoded obscurocommon.EncodedRollup) (*ExtRollup, error) {
	r := ExtRollup{}
	err := rlp.DecodeBytes(encoded, &r)

	return &r, err
}

func EncodeRollup(r *ExtRollup) obscurocommon.EncodedRollup {
	encoded, err := r.encode()
	if err != nil {
		panic(err)
	}

	return encoded
}

func DecodeRollup(rollup obscurocommon.EncodedRollup) *ExtRollup {
	r, err := Decode(rollup)
	if err != nil {
		panic(err)
	}

	return r
}

func (r ExtRollup) encode() (obscurocommon.EncodedRollup, error) {
	return rlp.EncodeToBytes(r)
}
