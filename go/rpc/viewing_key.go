package rpc

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/go-obscuro/go/enclave/rpc"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/obscuronet/go-obscuro/go/wallet"
)

// ViewingKey encapsulates the signed viewing key for an account for use in encrypted communication with an enclave
type ViewingKey struct {
	Account    *common.Address   // Account address that this private key is bound to
	PrivateKey *ecies.PrivateKey // private viewing key
	PublicKey  []byte            // public viewing key in bytes to share with enclave
	SignedKey  []byte            // public viewing key signed by the Account's private key
}

// GenerateAndSignViewingKey takes an account wallet, it generate a viewing key and signs the key with the acc's private key
func GenerateAndSignViewingKey(wal wallet.Wallet) (*ViewingKey, error) {
	// generate an ECDSA key pair to encrypt sensitive communications with the obscuro enclave
	vk, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate viewing key for RPC client: %w", err)
	}

	// get key in ECIES format
	viewingPrivateKeyECIES := ecies.ImportECDSA(vk)

	// encode public key as bytes
	viewingPubKeyBytes := crypto.CompressPubkey(&vk.PublicKey)

	// sign hex-encoded public key string with the wallet's private key
	viewingKeyHex := hex.EncodeToString(viewingPubKeyBytes)
	signature, err := signViewingKey(viewingKeyHex, wal.PrivateKey())
	if err != nil {
		return nil, err
	}

	accAddress := wal.Address()
	return &ViewingKey{
		Account:    &accAddress,
		PrivateKey: viewingPrivateKeyECIES,
		PublicKey:  viewingPubKeyBytes,
		SignedKey:  signature,
	}, nil
}

// signViewingKey takes a public key bytes as hex and the private key for a wallet, it simulates the back-and-forth to
// MetaMask and returns the signature bytes to register with the enclave
func signViewingKey(viewingKeyHex string, signerKey *ecdsa.PrivateKey) ([]byte, error) {
	msgToSign := rpc.ViewingKeySignedMsgPrefix + viewingKeyHex
	signature, err := crypto.Sign(accounts.TextHash([]byte(msgToSign)), signerKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign viewing key: %w", err)
	}

	// We have to transform the V from 0/1 to 27/28, and add the leading "0".
	signature[64] += 27
	signatureWithLeadBytes := append([]byte("0"), signature...)

	// this string encoded signature is what the wallet extension would receive after it is signed by metamask
	sigStr := hex.EncodeToString(signatureWithLeadBytes)
	// and then we extract the signature bytes in the same way as the wallet extension
	outputSig, err := hex.DecodeString(sigStr[2:])
	if err != nil {
		return nil, fmt.Errorf("failed to decode signature string: %w", err)
	}
	// This same change is made in geth internals, for legacy reasons to be able to recover the address:
	//	https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L452-L459
	outputSig[64] -= 27

	return outputSig, nil
}
