package obsclient

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/obscuronet/go-obscuro/go/rpcclientlib"
	"math/big"
)

// ObsClient requires an RPC Client and provides access to general Obscuro functionality that doesn't require viewing keys.
type ObsClient struct {
	c rpcclientlib.Client
}

func Dial(rawurl string) (*ObsClient, error) {
	rc, err := rpcclientlib.NewNetworkClient(rawurl)
	if err != nil {
		return nil, err
	}
	return NewClient(rc), nil
}

func NewClient(c rpcclientlib.Client) *ObsClient {
	return &ObsClient{c}
}

func (oc *ObsClient) Close() {
	oc.c.Stop()
}

// Blockchain Access

// ChainID retrieves the current chain ID for transaction replay protection.
func (oc *ObsClient) ChainID() (*big.Int, error) {
	var result hexutil.Big
	err := oc.c.Call(&result, "eth_chainId")
	if err != nil {
		return nil, err
	}
	return (*big.Int)(&result), err
}
