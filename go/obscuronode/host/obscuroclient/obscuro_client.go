package obscuroclient

import (
	"github.com/ethereum/go-ethereum/rpc"
)

// todo - joel - describe
type Client interface {
	// todo - joel - describe
	Call(result interface{}, method string, args ...interface{}) error
}

// todo - joel - describe
type ClientImpl struct {
	rpcClient *rpc.Client
}

func NewClient() Client {
	// todo - joel - parameterise this
	rpcClient, err := rpc.Dial("http://127.0.0.1:12000")
	if err != nil {
		panic(err)
	}

	return ClientImpl{rpcClient: rpcClient}
}

// todo - joel - describe
func (c ClientImpl) Call(result interface{}, method string, args ...interface{}) error {
	return c.rpcClient.Call(result, method, args)
}

// todo - joel - describe
type ClientDummy struct{}

// todo - joel - describe
func (c ClientDummy) Call(interface{}, string, ...interface{}) error {
	println("did nothing")
	return nil
}
