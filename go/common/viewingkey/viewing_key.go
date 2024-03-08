package viewingkey

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

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
	EIP712Domain          = "EIP712Domain"
	EIP712Type            = "Authentication"
	EIP712DomainName      = "name"
	EIP712DomainVersion   = "version"
	EIP712DomainChainID   = "chainId"
	EIP712EncryptionToken = "Encryption Token"
	// EIP712EncryptionTokenV2 is used to support older versions of third party libraries
	// that don't have the support for spaces in type names
	EIP712EncryptionTokenV2      = "EncryptionToken"
	EIP712DomainNameValue        = "Ten"
	EIP712DomainVersionValue     = "1.0"
	UserIDHexLength              = 40
	PersonalSignMessagePrefix    = "\x19Ethereum Signed Message:\n%d%s"
	PersonalSignMessageFormat    = "Token: %s on chain: %d version:%d"
	EIP712SignatureTypeInt       = 0
	PersonalSignSignatureTypeInt = 1
	LegacySignatureTypeInt       = 2
)

const (
	EIP712Signature SignatureType = "EIP712"
	PersonalSign    SignatureType = "PersonalSign"
	Legacy          SignatureType = "Legacy"
)

// EIP712EncryptionTokens is a list of all possible options for Encryption token name
var EIP712EncryptionTokens = [...]string{
	EIP712EncryptionToken,
	EIP712EncryptionTokenV2,
}

// PersonalSignMessageSupportedVersions is a list of supported versions for the personal sign message
var PersonalSignMessageSupportedVersions = []int{1}

// SignatureType is used to differentiate between different signature types (string is used, because int is not RLP-serializable)
type SignatureType string

// IntToSignatureType converts an int to a SignatureType
func IntToSignatureType(signatureType int) SignatureType {
	switch signatureType {
	case EIP712SignatureTypeInt:
		return EIP712Signature
	case PersonalSignSignatureTypeInt:
		return PersonalSign
	case LegacySignatureTypeInt:
		return Legacy
	default:
		return ""
	}
}

// ViewingKey encapsulates the signed viewing key for an account for use in encrypted communication with an enclave.
// It is th client-side perspective of the viewing key used for decrypting incoming traffic.
type ViewingKey struct {
	Account                 *gethcommon.Address // Account address that this Viewing Key is bound to - Users Pubkey address
	PrivateKey              *ecies.PrivateKey   // ViewingKey private key to encrypt data to the enclave
	PublicKey               []byte              // ViewingKey public key in decrypt data from the enclave
	SignatureWithAccountKey []byte              // ViewingKey public key signed by the Accounts Private key - Allows to retrieve the Account address
	SignatureType           SignatureType       // Type of signature used to sign the public key
}

// RPCSignedViewingKey - used for transporting a minimalist viewing key via
// every RPC request to a sensitive method, including Log subscriptions.
// only the public key and the signature are required
// the account address is sent as well to aid validation
type RPCSignedViewingKey struct {
	PublicKey               []byte
	SignatureWithAccountKey []byte
	SignatureType           SignatureType
}

// GenerateViewingKeyForWallet takes an account wallet, generates a viewing key and signs the key with the acc's private key
// uses the same method of signature handling as Metamask/geth
// TODO @Ziga - update this method to use the new EIP-712 signature format / personal sign after the removal of the legacy format
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
		Account:                 &accAddress,
		PrivateKey:              viewingPrivateKeyECIES,
		PublicKey:               viewingPubKeyBytes,
		SignatureWithAccountKey: signature,
		SignatureType:           Legacy,
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

func GeneratePersonalSignMessage(encryptionToken string, chainID int64, version int) string {
	return fmt.Sprintf(PersonalSignMessageFormat, encryptionToken, chainID, version)
}

// getBytesFromTypedData creates EIP-712 compliant hash from typedData.
// It involves hashing the message with its structure, hashing domain separator,
// and then encoding both hashes with specific EIP-712 bytes to construct the final message format.
func getBytesFromTypedData(typedData apitypes.TypedData) ([]byte, error) {
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

// GenerateAuthenticationEIP712RawDataOptions generates all the options or raw data messages (bytes)
// for an EIP-712 message used to authenticate an address with user
func GenerateAuthenticationEIP712RawDataOptions(userID string, chainID int64) ([][]byte, error) {
	if len(userID) != UserIDHexLength {
		return nil, fmt.Errorf("userID hex length must be %d, received %d", UserIDHexLength, len(userID))
	}
	encryptionToken := "0x" + userID

	domain := apitypes.TypedDataDomain{
		Name:    EIP712DomainNameValue,
		Version: EIP712DomainVersionValue,
		ChainId: (*math.HexOrDecimal256)(big.NewInt(chainID)),
	}

	typedDataList := make([]apitypes.TypedData, 0, len(EIP712EncryptionTokens))
	for _, encTokenName := range EIP712EncryptionTokens {
		message := map[string]interface{}{
			encTokenName: encryptionToken,
		}

		types := apitypes.Types{
			EIP712Domain: {
				{Name: EIP712DomainName, Type: "string"},
				{Name: EIP712DomainVersion, Type: "string"},
				{Name: EIP712DomainChainID, Type: "uint256"},
			},
			EIP712Type: {
				{Name: encTokenName, Type: "address"},
			},
		}

		newTypeElement := apitypes.TypedData{
			Types:       types,
			PrimaryType: EIP712Type,
			Domain:      domain,
			Message:     message,
		}
		typedDataList = append(typedDataList, newTypeElement)
	}

	rawDataOptions := make([][]byte, 0, len(typedDataList))
	for _, typedDataItem := range typedDataList {
		rawData, err := getBytesFromTypedData(typedDataItem)
		if err != nil {
			return nil, err
		}
		rawDataOptions = append(rawDataOptions, rawData)
	}
	return rawDataOptions, nil
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

// checkEIP712Signature checks if signature is valid for provided userID and chainID and return address or nil if not valid
func checkEIP712Signature(userID string, signature []byte, chainID int64) (*gethcommon.Address, error) {
	if len(signature) != 65 {
		return nil, fmt.Errorf("invalid signaure length: %d", len(signature))
	}

	rawDataOptions, err := GenerateAuthenticationEIP712RawDataOptions(userID, chainID)
	if err != nil {
		return nil, fmt.Errorf("cannot generate eip712 message. Cause %w", err)
	}

	// We transform the V from 27/28 to 0/1. This same change is made in Geth internals, for legacy reasons to be able
	// to recover the address: https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L452-L459
	if signature[64] == 27 || signature[64] == 28 {
		signature[64] -= 27
	}

	for _, rawData := range rawDataOptions {
		// create a hash of structured message (needed for signature verification)
		hashBytes := crypto.Keccak256(rawData)

		// current signature is valid - return account address
		address, err := CheckSignatureAndReturnAccountAddress(hashBytes, signature)
		if err == nil {
			return address, nil
		}
	}
	return nil, errors.New("EIP 712 signature verification failed")
}

// checkPersonalSignSignature checks if signature is valid for provided encryptionToken and chainID and return address or nil if not valid
func checkPersonalSignSignature(encryptionToken string, signature []byte, chainID int64) (*gethcommon.Address, error) {
	if len(signature) != 65 {
		return nil, fmt.Errorf("invalid signaure length: %d", len(signature))
	}
	// We transform the V from 27/28 to 0/1. This same change is made in Geth internals, for legacy reasons to be able
	// to recover the address: https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L452-L459
	if signature[64] == 27 || signature[64] == 28 {
		signature[64] -= 27
	}

	// create all possible hashes (for all the supported versions) of the message (needed for signature verification)
	for _, version := range PersonalSignMessageSupportedVersions {
		message := GeneratePersonalSignMessage(encryptionToken, chainID, version)
		prefixedMessage := fmt.Sprintf(PersonalSignMessagePrefix, len(message), message)
		messageHash := crypto.Keccak256([]byte(prefixedMessage))

		// current signature is valid - return account address
		address, err := CheckSignatureAndReturnAccountAddress(messageHash, signature)
		if err == nil {
			return address, nil
		}
	}

	return nil, fmt.Errorf("signature verification failed")
}

// CheckSignatureWithType TODO @Ziga - Refactor and simplify this function
func CheckSignatureWithType(encryptionToken string, signature []byte, chainID int64, signatureType SignatureType) (*gethcommon.Address, error) {
	if signatureType == PersonalSign {
		addr, err := checkPersonalSignSignature(encryptionToken, signature, chainID)
		if err == nil {
			return addr, nil
		}
	} else if signatureType == EIP712Signature {
		addr, err := checkEIP712Signature(encryptionToken, signature, chainID)
		if err == nil {
			return addr, nil
		}
	} else if signatureType == Legacy {
	}
	return nil, fmt.Errorf("signature verification failed")
}
