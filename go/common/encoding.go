package common

import (
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/common/log"
)

func EncodeRollup(r *EncryptedRollup) EncodedRollup {
	encoded, err := rlp.EncodeToBytes(r)
	if err != nil {
		log.Panic("could not encode rollup. Cause: %s", err)
	}

	return encoded
}

func DecodeRollup(encoded EncodedRollup) (*EncryptedRollup, error) {
	r := new(EncryptedRollup)
	err := rlp.DecodeBytes(encoded, r)

	return r, err
}

func DecodeRollupOrPanic(rollup EncodedRollup) *EncryptedRollup {
	r, err := DecodeRollup(rollup)
	if err != nil {
		log.Panic("could not decode rollup. Cause: %s", err)
	}

	return r
}

func EncodeAttestation(att *AttestationReport) EncodedAttestationReport {
	encoded, err := rlp.EncodeToBytes(att)
	if err != nil {
		panic(err)
	}

	return encoded
}

func DecodeAttestation(encoded EncodedAttestationReport) (*AttestationReport, error) {
	att := new(AttestationReport)
	err := rlp.DecodeBytes(encoded, att)

	return att, err
}
