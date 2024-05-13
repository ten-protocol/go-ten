package rpc

import (
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
	encClient, err := NewEncRPCClient(rpcClient, viewingKey, logger)
	if err != nil {
		return nil, err
	}
	return encClient, nil
}

func NewEncNetworkClientFromConn(connection *gethrpc.Client, viewingKey *viewingkey.ViewingKey, logger gethlog.Logger) (*EncRPCClient, error) {
	encClient, err := NewEncRPCClient(connection, viewingKey, logger)
	if err != nil {
		return nil, err
	}
	return encClient, nil
}

// NewNetworkClient returns a client that can make RPC calls to an Obscuro node
func NewNetworkClient(address string) (Client, error) {
	return rpc.Dial(address)
}
