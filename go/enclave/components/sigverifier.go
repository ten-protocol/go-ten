package components

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	"github.com/ten-protocol/go-ten/go/common"

	"github.com/ten-protocol/go-ten/go/common/log"

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

	// loop through sequencer keys and exit early if one of them matches
	for _, seqID := range sequencerIDs {
		err := sigChecker.verifyForSequencer(seqID, hash, sig)
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf("could not verify the signature of batch %s against any of the stored sequencer enclave keys", hash.Hex())
}

func (sigChecker *SignatureValidator) verifyForSequencer(seqID common.EnclaveID, hash gethcommon.Hash, sig []byte) error {
	attestedEnclave, err := sigChecker.storage.GetEnclavePubKey(context.Background(), seqID)
	if err != nil {
		sigChecker.logger.Error("Could not get public key for sequencer. Should not happen", "sequencerID", seqID, log.ErrKey, err)
		return fmt.Errorf("could not load public key for sequencer %s", seqID.Hex())
	}

	err = signature.VerifySignature(attestedEnclave.PubKey, hash.Bytes(), sig)
	if err != nil {
		sigChecker.logger.Debug("Could not verify signature", "sequencerID", seqID, log.ErrKey, err)
		return fmt.Errorf("could not verify the signature of batch %s against the stored sequencer enclave key: %s", hash.Hex(), seqID.Hex())
	}

	return nil
}
