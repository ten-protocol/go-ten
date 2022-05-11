package clientserver

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
)

const chainID = 1337 // TODO - Retrieve this value from the config service.

var bigChainID = (*hexutil.Big)(big.NewInt(chainID))

type EthAPI struct{}

func NewEthAPI() *EthAPI {
	return &EthAPI{}
}

func (api *EthAPI) ChainId() (*hexutil.Big, error) {
	return bigChainID, nil
}
