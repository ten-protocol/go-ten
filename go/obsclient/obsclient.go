package obsclient

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/obscuronet/go-obscuro/go/rpcclientlib"
)

// ObsClient requires an RPC Client and provides access to general Obscuro functionality that doesn't require viewing keys.
//
// The methods in this client are analogous to the methods in geth's EthClient and should behave the same unless noted otherwise.
type ObsClient struct {
	RPCClient rpcclientlib.Client
}

func Dial(rawurl string) (*ObsClient, error) {
	rc, err := rpcclientlib.NewNetworkClient(rawurl)
	if err != nil {
		return nil, err
	}
	return NewObsClient(rc), nil
}

func NewObsClient(c rpcclientlib.Client) *ObsClient {
	return &ObsClient{c}
}

func (oc *ObsClient) Close() {
	oc.RPCClient.Stop()
}

// Blockchain Access

// ChainID retrieves the current chain ID for transaction replay protection.
func (oc *ObsClient) ChainID() (*big.Int, error) {
	var result hexutil.Big
	err := oc.RPCClient.Call(&result, rpcclientlib.RPCChainID)
	if err != nil {
		return nil, err
	}
	return (*big.Int)(&result), err
}
