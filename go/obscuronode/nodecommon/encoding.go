package nodecommon

import (
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

func EncodeRollup(r *Rollup) obscurocommon.EncodedRollup {
	encoded, err := rlp.EncodeToBytes(r)
	if err != nil {
		log.Error(fmt.Sprintf("could not encode rollup. Cause: %s", err))
		panic(err)
	}

	return encoded
}

func DecodeRollup(encoded obscurocommon.EncodedRollup) (*Rollup, error) {
	r := new(Rollup)
	err := rlp.DecodeBytes(encoded, r)

	return r, err
}

func DecodeRollupOrPanic(rollup obscurocommon.EncodedRollup) *Rollup {
	r, err := DecodeRollup(rollup)
	if err != nil {
		log.Error(fmt.Sprintf("could not decode rollup. Cause: %s", err))
		panic(err)
	}

	return r
}
