package viewingkey

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/ten-protocol/go-ten/go/wallet"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// SignedMsgPrefix is the prefix added when signing the viewing key in MetaMask using the personal_sign
// API. Why is this needed? MetaMask has a security feature whereby if you ask it to sign something that looks like
// a transaction using the personal_sign API, it modifies the data being signed. The goal is to prevent hackers
// from asking a visitor to their website to personal_sign something that is actually a malicious transaction (e.g.
// theft of funds). By adding a prefix, the viewing key bytes no longer looks like a transaction hash, and thus get
// signed as-is.
const SignedMsgPrefix = "vk"

const (
	EIP712Domain             = "EIP712Domain"
	EIP712Type               = "Authentication"
	EIP712DomainName         = "name"
	EIP712DomainVersion      = "version"
	EIP712DomainChainID      = "chainId"
	EIP712EncryptionToken    = "Encryption Token"
	EIP712DomainNameValue    = "Ten"
	EIP712DomainVersionValue = "1.0"
	UserIDHexLength          = 40
)

// ViewingKey encapsulates the signed viewing key for an account for use in encrypted communication with an enclave
type ViewingKey struct {
	Account    *gethcommon.Address // Account address that this Viewing Key is bound to - Users Pubkey address
	PrivateKey *ecies.PrivateKey   // ViewingKey private key to encrypt data to the enclave
	PublicKey  []byte              // ViewingKey public key in decrypt data from the enclave
	Signature  []byte              // ViewingKey public key signed by the Accounts Private key - Allows to retrieve the Account address
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
		Signature:  signature,
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
// todo (@ziga) Remove this method once old WE endpoints are removed
func GenerateSignMessage(vkPubKey []byte) string {
	return SignedMsgPrefix + hex.EncodeToString(vkPubKey)
}

// GenerateSignMessageOG creates the message to be signed by Obscuro Gateway (new format)
// format is expected to be "Register <userID> for <Account>" (with the account in lowercase)
func GenerateSignMessageOG(vkPubKey []byte, addr *gethcommon.Address) string {
	userID := crypto.Keccak256Hash(vkPubKey).Bytes()
	return fmt.Sprintf("Register %s for %s", hex.EncodeToString(userID), strings.ToLower(addr.Hex()))
}

// GenerateAuthenticationEIP712RawData generates raw data (bytes)
// for an EIP-712 message used to authenticate an address with user
func GenerateAuthenticationEIP712RawData(userID string, chainID int64) ([]byte, error) {
	if len(userID) != UserIDHexLength {
		return nil, fmt.Errorf("userID hex length must be %d, received %d", UserIDHexLength, len(userID))
	}
	encryptionToken := "0x" + userID

	types := apitypes.Types{
		EIP712Domain: {
			{Name: EIP712DomainName, Type: "string"},
			{Name: EIP712DomainVersion, Type: "string"},
			{Name: EIP712DomainChainID, Type: "uint256"},
		},
		EIP712Type: {
			{Name: EIP712EncryptionToken, Type: "address"},
		},
	}

	domain := apitypes.TypedDataDomain{
		Name:    EIP712DomainNameValue,
		Version: EIP712DomainVersionValue,
		ChainId: (*math.HexOrDecimal256)(big.NewInt(chainID)),
	}

	message := map[string]interface{}{
		EIP712EncryptionToken: encryptionToken,
	}

	typedData := apitypes.TypedData{
		Types:       types,
		PrimaryType: EIP712Type,
		Domain:      domain,
		Message:     message,
	}

	// Now we need to create EIP-712 compliant hash.
	// It involves hashing the message with its structure, hashing domain separator,
	// and then encoding both hashes with specific EIP-712 bytes to construct the final message format.

	// Hash the EIP-712 message using its type and content
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, err
	}
	// Create the domain separator hash for EIP-712 message context
	domainSeparator, err := typedData.HashStruct(EIP712Domain, typedData.Domain.Map())
	if err != nil {
		return nil, err
	}
	// Prefix domain and message hashes with EIP-712 version and encoding bytes
	rawData := append([]byte("\x19\x01"), append(domainSeparator, typedDataHash...)...)
	return rawData, nil
}

// CalculateUserIDHex CalculateUserID calculates userID from a public key
// (we truncate it, because we want it to have length 20) and encode to hex strings
func CalculateUserIDHex(publicKeyBytes []byte) string {
	return hex.EncodeToString(CalculateUserID(publicKeyBytes))
}

// CalculateUserID calculates userID from a public key (we truncate it, because we want it to have length 20)
func CalculateUserID(publicKeyBytes []byte) []byte {
	return crypto.Keccak256Hash(publicKeyBytes).Bytes()[:20]
}

func VerifySignatureEIP712(userID string, address *gethcommon.Address, signature []byte, chainID int64) (bool, error) {
	// get raw data for structured message
	rawData, err := GenerateAuthenticationEIP712RawData(userID, chainID)
	if err != nil {
		return false, err
	}

	// create a hash of structured message (needed for signature verification)
	hashBytes := crypto.Keccak256(rawData)
	hash := gethcommon.BytesToHash(hashBytes)

	if len(signature) != 65 {
		return false, fmt.Errorf("invalid signature length: %d", len(signature))
	}

	// We transform the V from 27/28 to 0/1. This same change is made in Geth internals, for legacy reasons to be able
	// to recover the address: https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L452-L459
	if signature[64] == 27 || signature[64] == 28 {
		signature[64] -= 27
	}

	pubKeyBytes, err := crypto.Ecrecover(hash[:], signature)
	if err != nil {
		return false, fmt.Errorf("invalid signature: %w", err)
	}

	pubKey, err := crypto.UnmarshalPubkey(pubKeyBytes)
	if err != nil {
		return false, fmt.Errorf("cannot unmarshal public key: %w", err)
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	if !bytes.Equal(recoveredAddr.Bytes(), address.Bytes()) {
		return false, errors.New("address from signature not the same as expected")
	}

	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:64])

	// Verify the signature
	isValid := ecdsa.Verify(pubKey, hashBytes, r, s)

	if !isValid {
		return false, errors.New("signature is not valid")
	}

	return true, nil
}
