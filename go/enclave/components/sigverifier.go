package components

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	"github.com/ten-protocol/go-ten/go/common/signature"

	"github.com/ten-protocol/go-ten/go/enclave/storage"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

type SignatureValidator struct {
	attestedKey *ecdsa.PublicKey
	storage     storage.Storage
}

func NewSignatureValidator(storage storage.Storage) (*SignatureValidator, error) {
	// todo (#718) - sequencer identities should be retrieved from the L1 management contract
	return &SignatureValidator{
		storage:     storage,
		attestedKey: nil,
	}, nil
}

// CheckSequencerSignature - verifies the signature against the registered sequencer
func (sigChecker *SignatureValidator) CheckSequencerSignature(headerHash gethcommon.Hash, sig []byte) error {
	if sig == nil {
		return fmt.Errorf("missing signature on batch")
	}

	// Get all sequencer enclave IDs
	sequencerIDs, err := sigChecker.storage.GetSequencerEnclaveIDs(context.Background())
	if err != nil {
		return fmt.Errorf("could not fetch sequencer IDs: %w", err)
	}

	if len(sequencerIDs) == 0 {
		println("NO SEQ IDS")
		return nil
	}

	// Try to verify the signature against each sequencer's public key
	for _, seqID := range sequencerIDs {
		attestedEnclave, err := sigChecker.storage.GetEnclavePubKey(context.Background(), seqID)
		if err != nil {
			continue // Skip if we can't get the public key for this sequencer
		}

		// Verify signature using this sequencer's public key
		err = signature.VerifySignature(attestedEnclave.PubKey, headerHash.Bytes(), sig)
		if err == nil {
			// Signature verified successfully
			return nil
		}
	}

	return fmt.Errorf("invalid sequencer signature")
}
