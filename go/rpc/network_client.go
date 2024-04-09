package rpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"

	gethlog "github.com/ethereum/go-ethereum/log"
)

const (
	ws   = "ws://"
	http = "http://"
)

// networkClient is a Client implementation that wraps Geth's rpc.Client to make calls to the TEN node
type networkClient struct {
	rpcClient *rpc.Client
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

// NewNetworkClient returns a client that can make RPC calls to an TEN node
func NewNetworkClient(address string) (Client, error) {
	if !strings.HasPrefix(address, http) && !strings.HasPrefix(address, ws) {
		return nil, fmt.Errorf("clients for TEN only support the %s and %s protocols", http, ws)
	}

	rpcClient, err := rpc.Dial(address)
	if err != nil {
		return nil, fmt.Errorf("could not create RPC client on %s. Cause: %w", address, err)
	}

	return &networkClient{
		rpcClient: rpcClient,
	}, nil
}

// Call handles JSON rpc requests, delegating to the geth RPC client
// The result must be a pointer so that package json can unmarshal into it. You can also pass nil, in which case the result is ignored.
func (c *networkClient) Call(result interface{}, method string, args ...interface{}) error {
	return c.rpcClient.Call(result, method, args...)
}

func (c *networkClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	return c.rpcClient.CallContext(ctx, result, method, args...)
}

func (c *networkClient) Subscribe(ctx context.Context, _ interface{}, namespace string, channel interface{}, args ...interface{}) (*rpc.ClientSubscription, error) {
	return c.rpcClient.Subscribe(ctx, namespace, channel, args...)
}

func (c *networkClient) Stop() {
	c.rpcClient.Close()
}
