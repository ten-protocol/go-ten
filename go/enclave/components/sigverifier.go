package components

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/obscuronet/go-obscuro/go/enclave/storage"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

type SignatureValidator struct {
	SequencerID gethcommon.Address
	attestedKey *ecdsa.PublicKey
	storage     storage.Storage
}

func NewSignatureValidator(seqID gethcommon.Address, storage storage.Storage) (*SignatureValidator, error) {
	// todo (#718) - sequencer identities should be retrieved from the L1 management contract
	return &SignatureValidator{
		SequencerID: seqID,
		storage:     storage,
		attestedKey: nil,
	}, nil
}

// CheckSequencerSignature - verifies the signature against the registered sequencer
func (sigChecker *SignatureValidator) CheckSequencerSignature(headerHash gethcommon.Hash, sigR *big.Int, sigS *big.Int) error {
	if sigR == nil || sigS == nil {
		return fmt.Errorf("missing signature on batch")
	}

	if sigChecker.attestedKey == nil {
		attestedKey, err := sigChecker.storage.FetchAttestedKey(sigChecker.SequencerID)
		if err != nil {
			return fmt.Errorf("could not retrieve attested key for aggregator %s. Cause: %w", sigChecker.SequencerID, err)
		}
		sigChecker.attestedKey = attestedKey
	}

	if !ecdsa.Verify(sigChecker.attestedKey, headerHash.Bytes(), sigR, sigS) {
		return fmt.Errorf("could not verify ECDSA signature")
	}
	return nil
}
