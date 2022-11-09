package obsclient

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/obscuronet/go-obscuro/go/rpc"
)

// ObsClient provides access to general Obscuro functionality that doesn't require viewing keys.
//
// The methods in this client are analogous to the methods in geth's EthClient and should behave the same unless noted otherwise.
type ObsClient struct {
	rpcClient rpc.Client
}

func Dial(rawurl string) (*ObsClient, error) {
	rc, err := rpc.NewNetworkClient(rawurl)
	if err != nil {
		return nil, err
	}
	return NewObsClient(rc), nil
}

func NewObsClient(c rpc.Client) *ObsClient {
	return &ObsClient{c}
}

func (oc *ObsClient) Close() {
	oc.rpcClient.Stop()
}

// Blockchain Access

// ChainID retrieves the current chain ID for transaction replay protection.
func (oc *ObsClient) ChainID() (*big.Int, error) {
	var result hexutil.Big
	err := oc.rpcClient.Call(&result, rpc.ChainID)
	if err != nil {
		return nil, err
	}
	return (*big.Int)(&result), err
}

// RollupNumber returns the height of the head rollup
func (oc *ObsClient) RollupNumber() (uint64, error) {
	var result hexutil.Uint64
	err := oc.rpcClient.Call(&result, rpc.BlockNumber)
	return uint64(result), err
}
