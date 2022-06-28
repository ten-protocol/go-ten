package nodecommon

import (
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/common"
	"github.com/obscuronet/obscuro-playground/go/log"
)

func EncodeRollup(r *EncryptedRollup) common.EncodedRollup {
	encoded, err := rlp.EncodeToBytes(r)
	if err != nil {
		log.Panic("could not encode rollup. Cause: %s", err)
	}

	return encoded
}

func DecodeRollup(encoded common.EncodedRollup) (*EncryptedRollup, error) {
	r := new(EncryptedRollup)
	err := rlp.DecodeBytes(encoded, r)

	return r, err
}

func DecodeRollupOrPanic(rollup common.EncodedRollup) *EncryptedRollup {
	r, err := DecodeRollup(rollup)
	if err != nil {
		log.Panic("could not decode rollup. Cause: %s", err)
	}

	return r
}

func EncodeAttestation(att *common.AttestationReport) common.EncodedAttestationReport {
	encoded, err := rlp.EncodeToBytes(att)
	if err != nil {
		panic(err)
	}

	return encoded
}

func DecodeAttestation(encoded common.EncodedAttestationReport) (*common.AttestationReport, error) {
	att := new(common.AttestationReport)
	err := rlp.DecodeBytes(encoded, att)

	return att, err
}
