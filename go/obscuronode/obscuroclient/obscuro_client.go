package obscuroclient

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/obscuro-playground/go/log"
)

type RPCMethod uint8

const (
	http = "http://"

	RPCSendTransactionEncrypted = "obscuro_sendTransactionEncrypted"
	RPCGetCurrentBlockHead      = "obscuro_getCurrentBlockHead"
	RPCGetCurrentRollupHead     = "obscuro_getCurrentRollupHead"
	RPCGetRollupHeader          = "obscuro_getRollupHeader"
	RPCGetTransaction           = "obscuro_getTransaction"
	RPCBalance                  = "obscuro_balance"
	RPCNonce                     = "obscuro_nonce"
	RPCStopHost                 = "obscuro_stopHost"
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

func NewClient(nodeID common.Address, address string) Client {
	rpcClient, err := rpc.Dial(http + address)
	if err != nil {
		log.Panic("could not create RPC client on %s. Cause: %s", http+address, err)
	}

	return &clientImpl{
		nodeID:    nodeID,
		rpcClient: rpcClient,
	}
}

func (c *clientImpl) ID() common.Address {
	return c.nodeID
}

func (c *clientImpl) Call(result interface{}, method string, args ...interface{}) error {
	return c.rpcClient.Call(&result, method, args...)
}

func (c *clientImpl) Stop() {
	c.rpcClient.Close()
}
