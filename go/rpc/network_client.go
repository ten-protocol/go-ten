package rpc

import (
	"context"
	"fmt"

	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"

	gethlog "github.com/ethereum/go-ethereum/log"
)

// NewEncNetworkClient returns a network RPC client with Viewing Key encryption/decryption
func NewEncNetworkClient(rpcAddress string, viewingKey *viewingkey.ViewingKey, logger gethlog.Logger) (*EncRPCClient, error) {
	rpcClient, err := NewNetworkClient(rpcAddress)
	if err != nil {
		return nil, err
	}
	enclavePublicKeyBytes, err := ReadEnclaveKey(rpcClient)
	if err != nil {
		return nil, fmt.Errorf("error reading enclave public key: %v", err)
	}
	encClient, err := NewEncRPCClient(rpcClient, viewingKey, enclavePublicKeyBytes, logger)
	if err != nil {
		return nil, err
	}
	return encClient, nil
}

func NewEncNetworkClientFromConn(connection *gethrpc.Client, encKey []byte, viewingKey *viewingkey.ViewingKey, logger gethlog.Logger) (*EncRPCClient, error) {
	encClient, err := NewEncRPCClient(connection, viewingKey, encKey, logger)
	if err != nil {
		return nil, err
	}
	return encClient, nil
}

// NewNetworkClient returns a client that can make RPC calls to an Obscuro node
func NewNetworkClient(address string) (Client, error) {
	return rpc.Dial(address)
}

func ReadEnclaveKey(connection Client) ([]byte, error) {
	var enclavePublicKeyBytes []byte
	err := connection.CallContext(context.Background(), &enclavePublicKeyBytes, RPCKey)
	if err != nil {
		return nil, err
	}
	return enclavePublicKeyBytes, nil
}
