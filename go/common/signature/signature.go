package signature

import (
	"crypto/ecdsa"
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

/**
 * The signature package provides utilities for verifying and recovering ECDSA signatures.
 *
 * These signatures are used by the enclave to sign batches, rollups and other messages.
 * They are also used for the viewing keys included with RPC requests.
 *
 * We create the ethereum standard, non-malleable signature format with a 1-byte recoveryID offset by 27.
 */

const (
	SigLength         = 65
	SigRecoveryOffset = 27
)

func VerifySignature(pubKey *ecdsa.PublicKey, hash, signature []byte) error {
	if len(signature) != SigLength {
		return fmt.Errorf("invalid signature length, expected %d, got %d", SigLength, len(signature))
	}
	if !crypto.VerifySignature(crypto.FromECDSAPub(pubKey), hash, signature[:SigLength-1]) {
		return fmt.Errorf("ECDSA verification failed")
	}
	return nil
}

func RecoverPubKeyBytes(hash, signature []byte) ([]byte, error) {
	if len(signature) != SigLength {
		return nil, fmt.Errorf("invalid signature length, expected %d, got %d", SigLength, len(signature))
	}

	// remove offset from recovery ID if necessary
	if signature[SigLength-1] >= SigRecoveryOffset {
		signature[SigLength-1] -= SigRecoveryOffset
	}

	return crypto.Ecrecover(hash, signature)
}

func RecoverPubKey(hash, signature []byte) (*ecdsa.PublicKey, error) {
	pubKeyBytes, err := RecoverPubKeyBytes(hash, signature)
	if err != nil {
		return nil, fmt.Errorf("failed to recover public key: %w", err)
	}
	key, err := crypto.UnmarshalPubkey(pubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal public key: %w", err)
	}
	return key, nil
}

func RecoverAddress(hash, signature []byte) (*gethcommon.Address, error) {
	pubKey, err := RecoverPubKey(hash, signature)
	if err != nil {
		return nil, err
	}
	address := crypto.PubkeyToAddress(*pubKey)
	return &address, nil
}

func Sign(messageHash []byte, key *ecdsa.PrivateKey) ([]byte, error) {
	sig, err := crypto.Sign(messageHash, key)
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %w", err)
	}
	// add recovery ID offset
	sig[64] += SigRecoveryOffset
	return sig, nil
}
