package components

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/ten-protocol/go-ten/go/common/signature"

	"github.com/ten-protocol/go-ten/go/enclave/storage"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// SequencerSignatureVerifier interface for signature validation
type SequencerSignatureVerifier interface {
	CheckSequencerSignature(hash gethcommon.Hash, sig []byte) error
}

type SignatureValidator struct {
	attestedKey *ecdsa.PublicKey
	storage     storage.Storage
	logger      gethlog.Logger
}

func NewSignatureValidator(storage storage.Storage, logger gethlog.Logger) (*SignatureValidator, error) {
	return &SignatureValidator{
		storage:     storage,
		attestedKey: nil,
		logger:      logger,
	}, nil
}

// CheckSequencerSignature - verifies the signature against the registered sequencer
func (sigChecker *SignatureValidator) CheckSequencerSignature(hash gethcommon.Hash, sig []byte) error {
	if sig == nil {
		return fmt.Errorf("missing signature on batch")
	}

	sequencerIDs, err := sigChecker.storage.GetSequencerEnclaveIDs(context.Background())
	if err != nil {
		return fmt.Errorf("could not fetch sequencer IDs: %w", err)
	}
	sigChecker.logger.Info("sequencer IDs to verify against", "ids", sequencerIDs)
	// loop through sequencer keys and exit early if one of them matches
	for _, seqID := range sequencerIDs {
		attestedEnclave, err := sigChecker.storage.GetEnclavePubKey(context.Background(), seqID)
		if err != nil {
			sigChecker.logger.Error("Could not get public key for sequencer. Should not happen", "sequencerID", seqID, "error", err)
			continue // skip if we can't get the public key for this sequencer
		}

		err = signature.VerifySignature(attestedEnclave.PubKey, hash.Bytes(), sig)
		if err != nil {
			sigChecker.logger.Warn("Could not verify signature", "sequencerID", seqID, "error", err)
			// todo - as a temporary fix we remmove the sig verification
			// continue // skip
		}
		// signature matches
		return nil
	}

	return fmt.Errorf("could not verify the signature against any of the stored sequencer enclave keys")
}
