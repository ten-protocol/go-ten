package rpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"

	gethlog "github.com/ethereum/go-ethereum/log"
)

const (
	ws   = "ws://"
	http = "http://"
)

// NetworkClient is a Client implementation that wraps Geth's rpc.Client to make calls to the obscuro node
type NetworkClient struct {
	RpcClient *rpc.Client
}

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
	encClient, err := NewEncRPCClient(&NetworkClient{RpcClient: connection}, viewingKey, logger)
	if err != nil {
		return nil, err
	}
	return encClient, nil
}

// NewNetworkClient returns a client that can make RPC calls to an Obscuro node
func NewNetworkClient(address string) (Client, error) {
	if !strings.HasPrefix(address, http) && !strings.HasPrefix(address, ws) {
		return nil, fmt.Errorf("clients for Obscuro only support the %s and %s protocols", http, ws)
	}

	rpcClient, err := rpc.Dial(address)
	if err != nil {
		return nil, fmt.Errorf("could not create RPC client on %s. Cause: %w", address, err)
	}

	return &NetworkClient{
		RpcClient: rpcClient,
	}, nil
}

// Call handles JSON rpc requests, delegating to the geth RPC client
// The result must be a pointer so that package json can unmarshal into it. You can also pass nil, in which case the result is ignored.
func (c *NetworkClient) Call(result interface{}, method string, args ...interface{}) error {
	return c.RpcClient.Call(result, method, args...)
}

func (c *NetworkClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	return c.RpcClient.CallContext(ctx, result, method, args...)
}

func (c *NetworkClient) Subscribe(ctx context.Context, _ interface{}, namespace string, channel interface{}, args ...interface{}) (*gethrpc.ClientSubscription, error) {
	return c.RpcClient.Subscribe(ctx, namespace, channel, args...)
}

func (c *NetworkClient) Stop() {
	c.RpcClient.Close()
}
