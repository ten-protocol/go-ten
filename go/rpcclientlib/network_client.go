package rpcclientlib

import (
	"fmt"

	"github.com/ethereum/go-ethereum/rpc"
)

// NewNetworkClient returns a client that can make RPC calls to an Obscuro node
func NewNetworkClient(address string) (Client, error) {
	rpcClient, err := rpc.Dial(http + address)
	if err != nil {
		return nil, fmt.Errorf("could not create RPC client on %s. Cause: %w", http+address, err)
	}

	return &networkClient{
		rpcClient: rpcClient,
	}, nil
}

// networkClient is a Client implementation that wraps Geth's rpc.Client to make calls to the obscuro node
type networkClient struct {
	rpcClient *rpc.Client
}

// Call handles JSON rpc requests - if the method is sensitive it will encrypt the args before sending the request and
//	then decrypts the response before returning.
// The result must be a pointer so that package json can unmarshal into it. You can also pass nil, in which case the result is ignored.
func (c *networkClient) Call(result interface{}, method string, args ...interface{}) error {
	return c.rpcClient.Call(&result, method, args...)
}

func (c *networkClient) Stop() {
	c.rpcClient.Close()
}
