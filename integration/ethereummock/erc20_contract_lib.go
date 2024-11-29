package ethereummock

import (
	"bytes"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/go/ethadapter/erc20contractlib"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

type contractLib struct{}

func (c *contractLib) CreateDepositTx(tx *ethadapter.L1DepositTx) types.TxData {
	return encodeTx(tx, depositTxAddr)
}

// DecodeTx - return only deposit transactions to the management contract
func (c *contractLib) DecodeTx(tx *types.Transaction) common.TenTransaction {
	if bytes.Equal(tx.To().Bytes(), depositTxAddr.Bytes()) {
		depositTx, ok := decodeTx(tx).(*ethadapter.L1DepositTx)
		if !ok {
			return nil
		}

		// Mock deposits towards the L1 bridge target nil as the management contract address
		// is not set.
		if bytes.Equal(depositTx.To.Bytes(), gethcommon.BigToAddress(gethcommon.Big0).Bytes()) {
			return depositTx
		}
	}

	return nil
}

// NewERC20ContractLibMock is an implementation of the erc20contractlib.ERC20ContractLib
func NewERC20ContractLibMock() erc20contractlib.ERC20ContractLib {
	return &contractLib{}
}
