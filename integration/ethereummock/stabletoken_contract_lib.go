package ethereummock

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/ethclient/erc20contractlib"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

type contractLib struct{}

func (c contractLib) DecodeTx(tx *types.Transaction) obscurocommon.L1Transaction {
	// TODO implement me
	panic("implement me")
}

// NewStableTokenContractLibMock Mock package, should not be used
func NewStableTokenContractLibMock() erc20contractlib.ERC20ContractLib {
	return &contractLib{}
}
