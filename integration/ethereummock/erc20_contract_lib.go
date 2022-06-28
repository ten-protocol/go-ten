package ethereummock

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/ethclient"
	"github.com/obscuronet/obscuro-playground/go/ethclient/erc20contractlib"
)

type contractLib struct{}

func (c *contractLib) CreateDepositTx(tx *ethclient.L1DepositTx, nonce uint64) types.TxData {
	return encodeTx(tx, nonce, depositTxAddr)
}

func (c *contractLib) DecodeTx(tx *types.Transaction) ethclient.L1Transaction {
	return decodeTx(tx)
}

// NewERC20ContractLibMock is an implementation of the erc20contractlib.ERC20ContractLib
func NewERC20ContractLibMock() erc20contractlib.ERC20ContractLib {
	return &contractLib{}
}
