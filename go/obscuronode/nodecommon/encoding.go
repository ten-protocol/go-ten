package nodecommon

import (
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/log"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

func EncodeRollup(r *EncryptedRollup) obscurocommon.EncodedRollup {
	encoded, err := rlp.EncodeToBytes(r)
	if err != nil {
		log.Panic("could not encode rollup. Cause: %s", err)
	}

	return encoded
}

func DecodeRollup(encoded obscurocommon.EncodedRollup) (*EncryptedRollup, error) {
	r := new(EncryptedRollup)
	err := rlp.DecodeBytes(encoded, r)

	return r, err
}

func DecodeRollupOrPanic(rollup obscurocommon.EncodedRollup) *EncryptedRollup {
	r, err := DecodeRollup(rollup)
	if err != nil {
		log.Panic("could not decode rollup. Cause: %s", err)
	}

	return r
}

func EncodeAttestation(att *obscurocommon.AttestationReport) obscurocommon.EncodedAttestationReport {
	encoded, err := rlp.EncodeToBytes(att)
	if err != nil {
		panic(err)
	}

	return encoded
}

func DecodeAttestation(encoded obscurocommon.EncodedAttestationReport) (*obscurocommon.AttestationReport, error) {
	att := new(obscurocommon.AttestationReport)
	err := rlp.DecodeBytes(encoded, att)

	return att, err
}
