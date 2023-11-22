package common

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"github.com/ten-protocol/go-ten/go/rpc"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
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

// CalculateUserID calculates userID from a public key
func CalculateUserID(publicKeyBytes []byte) []byte {
	return crypto.Keccak256Hash(publicKeyBytes).Bytes()
}

// GetUserIDAndAddressFromMessage checks if message is in correct format and extracts userID and address from it
func GetUserIDAndAddressFromMessage(message string) (string, string, error) {
	if authenticateMessageRegex.MatchString(message) {
		params := authenticateMessageRegex.FindStringSubmatch(message)
		if len(params) >= 3 {
			return params[1], params[2], nil
		}
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

func CreateEncClient(
	hostRPCBindAddr string,
	addressBytes []byte,
	privateKeyBytes []byte,
	signature []byte,
	logger gethlog.Logger,
) (*rpc.EncRPCClient, error) {
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
	encClient, err := rpc.NewEncNetworkClient(hostRPCBindAddr, vk, logger)
	if err != nil {
		return nil, fmt.Errorf("unable to create EncRPCClient: %w", err)
	}
	return encClient, nil
}

type RPCRequest struct {
	ID     json.RawMessage
	Method string
	Params []interface{}
}

// Clone returns a new instance of the *RPCRequest
func (r *RPCRequest) Clone() *RPCRequest {
	return &RPCRequest{
		ID:     r.ID,
		Method: r.Method,
		Params: r.Params,
	}
}
