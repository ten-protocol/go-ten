package viewingkey

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/ten-protocol/go-ten/go/wallet"
)

const (
	EIP712Signature SignatureType = 0
	PersonalSign    SignatureType = 1
	Legacy          SignatureType = 2
)

// SignatureType is used to differentiate between different signature types (string is used, because int is not RLP-serializable)
type SignatureType uint8

const (
	EIP712Domain              = "EIP712Domain"
	EIP712Type                = "Authentication"
	EIP712DomainName          = "name"
	EIP712DomainVersion       = "version"
	EIP712DomainChainID       = "chainId"
	EIP712EncryptionToken     = "Encryption Token"
	EIP712DomainNameValue     = "Ten"
	EIP712DomainVersionValue  = "1.0"
	UserIDHexLength           = 40
	PersonalSignMessageFormat = "Token: %s on chain: %d version: %d"
	SignedMsgPrefix           = "vk" // prefix for legacy signed messages (remove when legacy signature type is removed)
	PersonalSignVersion       = 1
)

// EIP712EncryptionTokens is a list of all possible options for Encryption token name
var EIP712EncryptionTokens = [...]string{
	EIP712EncryptionToken,
}

type MessageGenerator interface {
	generateMessage(encryptionToken string, chainID int64, version int, hash bool) ([]byte, error)
}

type (
	PersonalMessageGenerator struct{}
	EIP712MessageGenerator   struct{}
)

var messageGenerators = map[SignatureType]MessageGenerator{
	PersonalSign:    PersonalMessageGenerator{},
	EIP712Signature: EIP712MessageGenerator{},
}

// GenerateMessage generates a message for the given encryptionToken, chainID, version and signatureType
func (p PersonalMessageGenerator) generateMessage(encryptionToken string, chainID int64, version int, hash bool) ([]byte, error) {
	textMessage := fmt.Sprintf(PersonalSignMessageFormat, encryptionToken, chainID, version)
	if hash {
		return accounts.TextHash([]byte(textMessage)), nil
	}
	return []byte(textMessage), nil
}

func (e EIP712MessageGenerator) generateMessage(encryptionToken string, chainID int64, _ int, hash bool) ([]byte, error) {
	if len(encryptionToken) != UserIDHexLength {
		return nil, fmt.Errorf("userID hex length must be %d, received %d", UserIDHexLength, len(encryptionToken))
	}
	encryptionToken = "0x" + encryptionToken

	domain := apitypes.TypedDataDomain{
		Name:    EIP712DomainNameValue,
		Version: EIP712DomainVersionValue,
		ChainId: (*math.HexOrDecimal256)(big.NewInt(chainID)),
	}

	message := map[string]interface{}{
		EIP712EncryptionToken: encryptionToken,
	}

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

	newTypeElement := apitypes.TypedData{
		Types:       types,
		PrimaryType: EIP712Type,
		Domain:      domain,
		Message:     message,
	}

	rawData, err := getBytesFromTypedData(newTypeElement)
	if err != nil {
		return nil, err
	}

	// add the JSON message to the list of messages
	jsonData, err := json.Marshal(newTypeElement)
	if err != nil {
		return nil, err
	}

	if hash {
		return crypto.Keccak256(rawData), nil
	}
	return jsonData, nil
}

// GenerateMessage generates a message for the given encryptionToken, chainID, version and signatureType
// hashed parameter is used to determine if the returned message should be hashed or returned in plain text
func GenerateMessage(encryptionToken string, chainID int64, version int, signatureType SignatureType, hashed bool) ([]byte, error) {
	generator, exists := messageGenerators[signatureType]
	if !exists {
		return nil, fmt.Errorf("unsupported signature type")
	}
	return generator.generateMessage(encryptionToken, chainID, version, hashed)
}

// GenerateSignMessage creates the message to be signed
// vkPubKey is expected to be a []byte("0x....") to create the signing message
// todo (@ziga) Remove this method once old WE endpoints are removed
func GenerateSignMessage(vkPubKey []byte) string {
	return SignedMsgPrefix + hex.EncodeToString(vkPubKey)
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

// CalculateUserIDHex CalculateUserID calculates userID from a public key
// (we truncate it, because we want it to have length 20) and encode to hex strings
func CalculateUserIDHex(publicKeyBytes []byte) string {
	return hex.EncodeToString(CalculateUserID(publicKeyBytes))
}

// CalculateUserID calculates userID from a public key (we truncate it, because we want it to have length 20)
func CalculateUserID(publicKeyBytes []byte) []byte {
	return crypto.Keccak256Hash(publicKeyBytes).Bytes()[:20]
}

// GetBestFormat returns the best format for a message based on available formats that are supported by the user
func GetBestFormat(formatsSlice []string) SignatureType {
	// If "Personal" is the only format available, choose it
	if len(formatsSlice) == 1 && formatsSlice[0] == "Personal" {
		return PersonalSign
	}

	// otherwise, choose EIP712
	return EIP712Signature
}

func GetSignatureTypeString(expectedSignatureType SignatureType) string {
	for key, value := range SignatureTypeMap {
		if value == expectedSignatureType {
			return key
		}
	}
	return ""
}

var SignatureTypeMap = map[string]SignatureType{
	"EIP712":   EIP712Signature,
	"Personal": PersonalSign,
}
