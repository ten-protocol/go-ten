package obscuroclient

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	http = "http://"

	RPCSendTransactionEncrypted = "obscuro_sendTransactionEncrypted"
)

// Client is used by client applications to interact with the Obscuro node.
type Client interface {
	// ID returns the ID of the Obscuro node the client is for.
	ID() common.Address
	// Call executes the named method via RPC.
	Call(result interface{}, method string, args ...interface{}) error
	// Stop closes the client.
	Stop()
}

// A Client implementation that wraps rpc.Client to make calls.
type clientImpl struct {
	nodeID    common.Address
	rpcClient *rpc.Client
}

func NewClient(nodeID int64, address string) Client {
	rpcClient, err := rpc.Dial(http + address)
	if err != nil {
		panic(err)
	}

	return &clientImpl{
		nodeID:    common.BigToAddress(big.NewInt(nodeID)),
		rpcClient: rpcClient,
	}
}

func (c *clientImpl) ID() common.Address {
	return c.nodeID
}

func (c *clientImpl) Call(result interface{}, method string, args ...interface{}) error {
	return c.rpcClient.Call(result, method, args...)
}

func (c *clientImpl) Stop() {
	c.rpcClient.Close()
}
