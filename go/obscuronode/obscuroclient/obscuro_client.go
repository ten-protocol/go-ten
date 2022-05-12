package obscuroclient

import (
	"github.com/ethereum/go-ethereum/rpc"
)

const RPCSendTransactionEncrypted = "obscuro_sendTransactionEncrypted"

// todo - joel - describe
type Client interface {
	// todo - joel - describe
	Call(result interface{}, method string, args ...interface{}) error
}

// todo - joel - describe
type clientImpl struct {
	rpcClient *rpc.Client
}

func NewClient(address string) Client {
	// todo - joel - use constant
	rpcClient, err := rpc.Dial("http://" + address)
	if err != nil {
		panic(err)
	}

	return clientImpl{rpcClient: rpcClient}
}

// todo - joel - describe
func (c clientImpl) Call(result interface{}, method string, args ...interface{}) error {
	return c.rpcClient.Call(result, method, args...)
}

// todo - joel - pull dummy into host package
// todo - joel - use dummy in in-mem tests

// todo - joel - describe
type ClientDummy struct{}

// todo - joel - describe
func (c ClientDummy) Call(interface{}, string, ...interface{}) error {
	println("did nothing")
	return nil
}
