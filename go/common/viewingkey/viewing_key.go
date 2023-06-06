package viewingkey

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/obscuronet/go-obscuro/go/wallet"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// SignedMsgPrefix is the prefix added when signing the viewing key in MetaMask using the personal_sign
// API. Why is this needed? MetaMask has a security feature whereby if you ask it to sign something that looks like
// a transaction using the personal_sign API, it modifies the data being signed. The goal is to prevent hackers
// from asking a visitor to their website to personal_sign something that is actually a malicious transaction (e.g.
// theft of funds). By adding a prefix, the viewing key bytes no longer looks like a transaction hash, and thus get
// signed as-is.
const SignedMsgPrefix = "vk"

// ViewingKey encapsulates the signed viewing key for an account for use in encrypted communication with an enclave
type ViewingKey struct {
	Account    *gethcommon.Address // Account address that this Viewing Key is bound to - Users Pubkey address
	PrivateKey *ecies.PrivateKey   // ViewingKey private key to encrypt data to the enclave
	PublicKey  []byte              // ViewingKey public key in decrypt data from the enclave
	SignedKey  []byte              // User Public Key signed by the Users private key - Matches the Account address
}

// GenerateViewingKeyForWallet takes an account wallet, generates a viewing key and signs the key with the acc's private key
// uses the same method of signature handling as Metamask/geth
func GenerateViewingKeyForWallet(wal wallet.Wallet) (*ViewingKey, error) {
	// generate an ECDSA key pair to encrypt sensitive communications with the obscuro enclave
	vk, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate viewing key for RPC client: %w", err)
	}

	// get key in ECIES format
	viewingPrivateKeyECIES := ecies.ImportECDSA(vk)

	// encode public key as bytes
	viewingPubKeyBytes := crypto.CompressPubkey(&vk.PublicKey)

	// sign public key bytes with the wallet's private key
	signature, err := mmSignViewingKey(viewingPubKeyBytes, wal.PrivateKey())
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

// mmSignViewingKey takes a public key bytes as hex and the private key for a wallet, it simulates the back-and-forth to
// MetaMask and returns the signature bytes to register with the enclave
func mmSignViewingKey(viewingPubKeyBytes []byte, signerKey *ecdsa.PrivateKey) ([]byte, error) {
	signature, err := Sign(signerKey, viewingPubKeyBytes)
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

// Sign takes a users Private key and signs the public viewingKey hex
func Sign(userPrivKey *ecdsa.PrivateKey, vkPubKey []byte) ([]byte, error) {
	msgToSign := GenerateSignMessage(vkPubKey)
	signature, err := crypto.Sign(accounts.TextHash([]byte(msgToSign)), userPrivKey)
	if err != nil {
		return nil, fmt.Errorf("unable to sign messages - %w", err)
	}
	return signature, nil
}

// GenerateSignMessage creates the message to be signed
// vkPubKey is expected to be a []byte("0x....") to create the signing message
func GenerateSignMessage(vkPubKey []byte) string {
	return SignedMsgPrefix + hex.EncodeToString(vkPubKey)
}
