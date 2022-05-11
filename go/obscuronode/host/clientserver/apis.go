package clientserver

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

const chainID = 1337 // TODO - Retrieve this value from the config service.

var bigChainID = (*hexutil.Big)(big.NewInt(chainID))

// EthAPI implements specific Ethereum JSON RPC operations in the "eth" namespace.
type EthAPI struct{}

func NewEthAPI() *EthAPI {
	return &EthAPI{}
}

func (api *EthAPI) ChainId() (*hexutil.Big, error) { //nolint
	return bigChainID, nil
}

// ObscuroAPI implements Obscuro-specific JSON RPC operations.
type ObscuroAPI struct{}

func NewObscuroAPI() *ObscuroAPI {
	return &ObscuroAPI{}
}

// todo - joel - want to receive a string here instead - how?
func (api *ObscuroAPI) SendTransactionEncrypted([]interface{}) string { //nolint
	println("received encrypted tx")
	return "hello joel"
}
