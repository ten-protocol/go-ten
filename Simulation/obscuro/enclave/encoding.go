package enclave

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

func (tx L2Tx) Encode() (EncodedL2Tx, error) {
	return rlp.EncodeToBytes(tx)
}

func (encoded EncodedL2Tx) Decode() (tx L2Tx, err error) {
	err = rlp.DecodeBytes(encoded, &tx)
	return
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

func DecodeTx(tx EncodedL2Tx) L2Tx {
	r, err := tx.Decode()
	if err != nil {
		panic(err)
	}
	return r
}
