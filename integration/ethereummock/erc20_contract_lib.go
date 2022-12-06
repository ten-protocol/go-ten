package ethereummock

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/ethadapter/erc20contractlib"
)

type contractLib struct{}

func (c *contractLib) CreateDepositTx(tx *ethadapter.L1DepositTx, nonce uint64) types.TxData {
	return encodeTx(tx, nonce, depositTxAddr)
}

func (c *contractLib) DecodeTx(tx *types.Transaction) ethadapter.L1Transaction {
	if bytes.Equal(tx.To().Bytes(), depositTxAddr.Bytes()) {
		depositTx, ok := decodeTx(tx).(*ethadapter.L1DepositTx)
		if !ok {
			return nil
		}

		// Mock deposits towards the L1 bridge target nil as the management contract address
		// is not set.
		if bytes.Equal(depositTx.To.Bytes(), common.BigToAddress(common.Big0).Bytes()) {
			return depositTx
		}
	}

	return nil
}

// NewERC20ContractLibMock is an implementation of the erc20contractlib.ERC20ContractLib
func NewERC20ContractLibMock() erc20contractlib.ERC20ContractLib {
	return &contractLib{}
}
