package common

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/obscuronet/go-obscuro/go/common/viewingkey"
	"github.com/obscuronet/go-obscuro/go/rpc"
)

var authenticateMessageRegex = regexp.MustCompile(MessageFormatRegex)

// PrivateKeyToCompressedPubKey converts *ecies.PrivateKey to compressed PubKey ([]byte with length 33)
func PrivateKeyToCompressedPubKey(prvKey *ecies.PrivateKey) []byte {
	ecdsaPublicKey := prvKey.PublicKey.ExportECDSA()
	compressedPubKey := crypto.CompressPubkey(ecdsaPublicKey)
	return compressedPubKey
}

// BytesToPrivateKey converts []bytes to *ecies.PrivateKey
func BytesToPrivateKey(keyBytes []byte) (*ecies.PrivateKey, error) {
	ecdsaPrivateKey, err := crypto.ToECDSA(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to convert bytes to ECDSA private key: %w", err)
	}

	eciesPrivateKey := ecies.ImportECDSA(ecdsaPrivateKey)
	return eciesPrivateKey, nil
}

// CalculateUserID calculates userID from public key
func CalculateUserID(publicKeyBytes []byte) []byte {
	return crypto.Keccak256Hash(publicKeyBytes).Bytes()
}

// GetUserIDAndAddressFromMessage checks if message is in correct format and extracts userID and address from it
func GetUserIDAndAddressFromMessage(message string) (string, string, error) {
	if authenticateMessageRegex.MatchString(message) {
		params := authenticateMessageRegex.FindStringSubmatch(message)
		return params[1], params[2], nil
	}
	return "", "", errors.New("invalid message format")
}

// GetAddressAndPubKeyFromSignature gets an address that was used to sign given signature
func GetAddressAndPubKeyFromSignature(messageHash []byte, signature []byte) (gethcommon.Address, *ecdsa.PublicKey, error) {
	pubKey, err := crypto.SigToPub(messageHash, signature)
	if err != nil {
		return gethcommon.Address{}, nil, err
	}

	return crypto.PubkeyToAddress(*pubKey), pubKey, nil
}

// GetUserIDbyte converts userID from string to correct byte format
func GetUserIDbyte(userID string) ([]byte, error) {
	return hex.DecodeString(userID)
}

func CreateEncClient(hostRPCBindAddr string, addressBytes []byte, privateKeyBytes []byte, signature []byte) (*rpc.EncRPCClient, error) {
	privateKey, err := BytesToPrivateKey(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to convert bytes to ecies private key: %w", err)
	}

	address := gethcommon.BytesToAddress(addressBytes)

	vk := &viewingkey.ViewingKey{
		Account:    &address,
		PrivateKey: privateKey,
		PublicKey:  PrivateKeyToCompressedPubKey(privateKey),
		Signature:  signature,
	}
	encClient, err := rpc.NewEncNetworkClient(hostRPCBindAddr, vk, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create EncRPCClient: %w", err)
	}
	return encClient, nil
}
