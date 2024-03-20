package viewingkey

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// SignatureChecker is an interface for checking
// if signature is valid for provided encryptionToken and chainID and return singing address or nil if not valid
type SignatureChecker interface {
	CheckSignature(encryptionToken string, signature []byte, chainID int64) (*gethcommon.Address, error)
}

type (
	PersonalSignChecker struct{}
	EIP712Checker       struct{}
	LegacyChecker       struct{}
)

// CheckSignature checks if signature is valid for provided encryptionToken and chainID and return address or nil if not valid
func (psc PersonalSignChecker) CheckSignature(encryptionToken string, signature []byte, chainID int64) (*gethcommon.Address, error) {
	if len(signature) != 65 {
		return nil, fmt.Errorf("invalid signaure length: %d", len(signature))
	}
	// We transform the V from 27/28 to 0/1. This same change is made in Geth internals, for legacy reasons to be able
	// to recover the address: https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L452-L459
	if signature[64] == 27 || signature[64] == 28 {
		signature[64] -= 27
	}

	messageHash, err := GenerateMessage(encryptionToken, chainID, PersonalSignVersion, PersonalSign, true)
	if err != nil {
		return nil, fmt.Errorf("cannot generate message. Cause %w", err)
	}

	// signature is valid - return account address
	address, err := CheckSignatureAndReturnAccountAddress(messageHash, signature)
	if err == nil {
		return address, nil
	}

	return nil, fmt.Errorf("signature verification failed")
}

func (e EIP712Checker) CheckSignature(encryptionToken string, signature []byte, chainID int64) (*gethcommon.Address, error) {
	if len(signature) != 65 {
		return nil, fmt.Errorf("invalid signaure length: %d", len(signature))
	}

	messageHash, err := GenerateMessage(encryptionToken, chainID, 1, EIP712Signature, true)
	if err != nil {
		return nil, fmt.Errorf("cannot generate message. Cause %w", err)
	}

	// We transform the V from 27/28 to 0/1. This same change is made in Geth internals, for legacy reasons to be able
	// to recover the address: https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L452-L459
	if signature[64] == 27 || signature[64] == 28 {
		signature[64] -= 27
	}

	// current signature is valid - return account address
	address, err := CheckSignatureAndReturnAccountAddress(messageHash, signature)
	if err == nil {
		return address, nil
	}

	return nil, errors.New("EIP 712 signature verification failed")
}

// CheckSignature checks if signature is valid for provided encryptionToken and chainID and return address or nil if not valid
// todo (@ziga) Remove this method once old WE endpoints are removed
// encryptionToken is expected to be a public key and not encrypted token as with other signature types
// (since this is only temporary fix and legacy format will be removed soon)
func (lsc LegacyChecker) CheckSignature(encryptionToken string, signature []byte, _ int64) (*gethcommon.Address, error) {
	publicKey := []byte(encryptionToken)
	msgToSignLegacy := GenerateSignMessage(publicKey)

	recoveredAccountPublicKeyLegacy, err := crypto.SigToPub(accounts.TextHash([]byte(msgToSignLegacy)), signature)
	if err != nil {
		return nil, fmt.Errorf("failed to recover account public key from legacy signature: %w", err)
	}
	recoveredAccountAddressLegacy := crypto.PubkeyToAddress(*recoveredAccountPublicKeyLegacy)
	return &recoveredAccountAddressLegacy, nil
}

// SignatureChecker is a map of SignatureType to SignatureChecker
var signatureCheckers = map[SignatureType]SignatureChecker{
	PersonalSign:    PersonalSignChecker{},
	EIP712Signature: EIP712Checker{},
	Legacy:          LegacyChecker{},
}

// CheckSignature checks if signature is valid for provided encryptionToken and chainID and return address or nil if not valid
func CheckSignature(encryptionToken string, signature []byte, chainID int64, signatureType SignatureType) (*gethcommon.Address, error) {
	checker, exists := signatureCheckers[signatureType]
	if !exists {
		return nil, fmt.Errorf("unsupported signature type")
	}
	return checker.CheckSignature(encryptionToken, signature, chainID)
}

// CheckSignatureAndReturnAccountAddress checks if the signature is valid for hash of the message and checks if
// signer is an address provided to the function.
// It returns an address if the signature is valid and nil otherwise
func CheckSignatureAndReturnAccountAddress(hashBytes []byte, signature []byte) (*gethcommon.Address, error) {
	pubKeyBytes, err := crypto.Ecrecover(hashBytes, signature)
	if err != nil {
		return nil, err
	}

	pubKey, err := crypto.UnmarshalPubkey(pubKeyBytes)
	if err != nil {
		return nil, err
	}

	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:64])

	// Verify the signature and return the result (all the checks above passed)
	isSigValid := ecdsa.Verify(pubKey, hashBytes, r, s)
	if isSigValid {
		address := crypto.PubkeyToAddress(*pubKey)
		return &address, nil
	}
	return nil, fmt.Errorf("invalid signature")
}
