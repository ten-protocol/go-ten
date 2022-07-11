package common

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// EncodedBlock the encoded version of an ExtBlock
type EncodedBlock []byte

func EncodeBlock(b *types.Block) (EncodedBlock, error) {
	encoded, err := rlp.EncodeToBytes(b)
	if err != nil {
		return nil, fmt.Errorf("could not encode block to bytes. Cause: %w", err)
	}
	return encoded, nil
}

func (eb EncodedBlock) DecodeBlock() (*types.Block, error) {
	b := types.Block{}
	if err := rlp.DecodeBytes(eb, &b); err != nil {
		return nil, fmt.Errorf("could not decode block from bytes. Cause: %w", err)
	}
	return &b, nil
}

func EncodeRollup(r *EncryptedRollup) (EncodedRollup, error) {
	return rlp.EncodeToBytes(r)
}

func DecodeRollup(encoded EncodedRollup) (*EncryptedRollup, error) {
	r := new(EncryptedRollup)
	err := rlp.DecodeBytes(encoded, r)
	return r, err
}

func EncodeAttestation(att *AttestationReport) (EncodedAttestationReport, error) {
	return rlp.EncodeToBytes(att)
}

func DecodeAttestation(encoded EncodedAttestationReport) (*AttestationReport, error) {
	att := new(AttestationReport)
	err := rlp.DecodeBytes(encoded, att)
	return att, err
}
