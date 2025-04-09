package viewingkey

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/status-im/keycard-go/hexutils"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

const (
	EIP712Signature SignatureType = 0
	PersonalSign    SignatureType = 1
)

// SignatureType is used to differentiate between different signature types (string is used, because int is not RLP-serializable)
type SignatureType uint8

const (
	EIP712Domain                 = "EIP712Domain"
	EIP712Type                   = "Authentication"
	EIP712DomainName             = "name"
	EIP712DomainVersion          = "version"
	EIP712DomainChainID          = "chainId"
	EIP712VerifyingContract      = "verifyingContract"
	EIP712EncryptionToken        = "Encryption Token"
	EIP712DomainNameValue        = "Ten"
	EIP712DomainVersionValue     = "1.0"
	EIP712VerifyingContractValue = "0x0000000000000000000000000000000000000000"
	UserIDLength                 = 20
	PersonalSignMessageFormat    = "Token: %s on chain: %d version: %d"
	PersonalSignVersion          = 1
)

type MessageGenerator interface {
	generateMessage(encryptionToken []byte, chainID int64, version int) ([]byte, error)
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
func (p PersonalMessageGenerator) generateMessage(encryptionToken []byte, chainID int64, version int) ([]byte, error) {
	return []byte(fmt.Sprintf(PersonalSignMessageFormat, hexutils.BytesToHex(encryptionToken), chainID, version)), nil
}

func (e EIP712MessageGenerator) generateMessage(encryptionToken []byte, chainID int64, _ int) ([]byte, error) {
	if len(encryptionToken) != UserIDLength {
		return nil, fmt.Errorf("userID must be %d bytes, received %d", UserIDLength, len(encryptionToken))
	}
	EIP712TypedData := createTypedDataForEIP712Message(encryptionToken, chainID)

	// add the JSON message to the list of messages
	jsonData, err := json.Marshal(EIP712TypedData)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

// GenerateMessage generates a message for the given encryptionToken, chainID, version and signatureType
func GenerateMessage(encryptionToken []byte, chainID int64, version int, signatureType SignatureType) ([]byte, error) {
	generator, exists := messageGenerators[signatureType]
	if !exists {
		return nil, fmt.Errorf("unsupported signature type")
	}
	return generator.generateMessage(encryptionToken, chainID, version)
}

// MessageHash is an interface for getting the hash of the message
type MessageHash interface {
	getMessageHash(message []byte) []byte
}

type (
	PersonalMessageHash struct{}
	EIP712MessageHash   struct{}
)

var messageHash = map[SignatureType]MessageHash{
	PersonalSign:    PersonalMessageHash{},
	EIP712Signature: EIP712MessageHash{},
}

// getMessageHash returns the hash for the personal message
func (p PersonalMessageHash) getMessageHash(message []byte) []byte {
	return accounts.TextHash(message)
}

// getMessageHash returns the hash for the EIP712 message
func (E EIP712MessageHash) getMessageHash(message []byte) []byte {
	var EIP712TypedData apitypes.TypedData
	err := json.Unmarshal(message, &EIP712TypedData)
	if err != nil {
		return nil
	}

	rawData, err := getBytesFromTypedData(EIP712TypedData)
	if err != nil {
		return nil
	}
	return crypto.Keccak256(rawData)
}

// GetMessageHash returns the hash of the message based on the signature type
func GetMessageHash(message []byte, signatureType SignatureType) ([]byte, error) {
	hashFunction, exists := messageHash[signatureType]
	if !exists {
		return nil, fmt.Errorf("unsupported signature type")
	}
	return hashFunction.getMessageHash(message), nil
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

// createTypedDataForEIP712Message creates typed data for EIP712 message
func createTypedDataForEIP712Message(encryptionToken []byte, chainID int64) apitypes.TypedData {
	hexToken := hexutils.BytesToHex(encryptionToken)

	domain := apitypes.TypedDataDomain{
		Name:              EIP712DomainNameValue,
		Version:           EIP712DomainVersionValue,
		ChainId:           (*math.HexOrDecimal256)(big.NewInt(chainID)),
		VerifyingContract: EIP712VerifyingContractValue,
	}

	message := map[string]interface{}{
		EIP712EncryptionToken: hexToken,
	}

	types := apitypes.Types{
		EIP712Domain: {
			{Name: EIP712DomainName, Type: "string"},
			{Name: EIP712DomainVersion, Type: "string"},
			{Name: EIP712DomainChainID, Type: "uint256"},
			{Name: EIP712VerifyingContract, Type: "address"},
		},
		EIP712Type: {
			{Name: EIP712EncryptionToken, Type: "address"},
		},
	}

	typedData := apitypes.TypedData{
		Types:       types,
		PrimaryType: EIP712Type,
		Domain:      domain,
		Message:     message,
	}
	return typedData
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
