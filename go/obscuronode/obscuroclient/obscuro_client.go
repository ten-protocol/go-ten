package obscuroclient

import (
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	http = "http://"

	RPCSendTransactionEncrypted = "obscuro_sendTransactionEncrypted"
)

// Client is used by client applications to interact with the Obscuro node.
type Client interface {
	// Call executes the named method via RPC.
	Call(result interface{}, method string, args ...interface{}) error
	// Stop closes the client.
	Stop()
}

// A Client implementation that wraps rpc.Client to make calls.
type clientImpl struct {
	rpcClient *rpc.Client
}

func NewClient(address string) Client {
	rpcClient, err := rpc.Dial(http + address)
	if err != nil {
		panic(err)
	}

	return clientImpl{rpcClient: rpcClient}
}

func (c clientImpl) Call(result interface{}, method string, args ...interface{}) error {
	return c.rpcClient.Call(result, method, args...)
}

func (c clientImpl) Stop() {
	c.rpcClient.Close()
}
