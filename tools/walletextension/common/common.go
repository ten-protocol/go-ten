package common

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"github.com/ten-protocol/go-ten/go/rpc"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
)

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

// GetUserIDbyte converts userID from string to correct byte format
func GetUserIDbyte(userID string) ([]byte, error) {
	return hex.DecodeString(userID)
}

func CreateEncClient(
	hostRPCBindAddr string,
	addressBytes []byte,
	privateKeyBytes []byte,
	signature []byte,
	signatureType int,
	logger gethlog.Logger,
) (*rpc.EncRPCClient, error) {
	privateKey, err := BytesToPrivateKey(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to convert bytes to ecies private key: %w", err)
	}

	address := gethcommon.BytesToAddress(addressBytes)

	vk := &viewingkey.ViewingKey{
		Account:                 &address,
		PrivateKey:              privateKey,
		PublicKey:               PrivateKeyToCompressedPubKey(privateKey),
		SignatureWithAccountKey: signature,
		SignatureType:           viewingkey.IntToSignatureType(signatureType),
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

// NewFileLogger is a logger factory function
func NewFileLogger() gethlog.Logger {
	// Open or create your log file
	file, err := os.OpenFile("gateway_logs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		panic(err)
	}

	// Create a new logger instance
	logger := gethlog.New()

	// Set the handler to the file
	logger.SetHandler(gethlog.StreamHandler(file, log.TenLogFormat()))

	return logger
}
