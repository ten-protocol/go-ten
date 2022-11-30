package common

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// EncodedL1Block the encoded version of an L1 block.
type EncodedL1Block []byte

func EncodeBlock(b *types.Block) (EncodedL1Block, error) {
	encoded, err := rlp.EncodeToBytes(b)
	if err != nil {
		return nil, fmt.Errorf("could not encode block to bytes. Cause: %w", err)
	}
	return encoded, nil
}

func (eb EncodedL1Block) DecodeBlock() (*types.Block, error) {
	b := types.Block{}
	if err := rlp.DecodeBytes(eb, &b); err != nil {
		return nil, fmt.Errorf("could not decode block from bytes. Cause: %w", err)
	}
	return &b, nil
}

func EncodeRollup(r *ExtRollup) (EncodedRollup, error) {
	return rlp.EncodeToBytes(r)
}

func DecodeRollup(encoded EncodedRollup) (*ExtRollup, error) {
	r := new(ExtRollup)
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
