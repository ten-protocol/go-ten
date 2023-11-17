package erc20contractlib

import (
	"math/big"
	"strings"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/log"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func init() { //nolint:gochecknoinits
	var err error
	obscuroERC20ContractABIJSON, err = abi.JSON(strings.NewReader(ERC20ContractABI))
	if err != nil {
		panic(err)
	}
}

var obscuroERC20ContractABIJSON = abi.ABI{}

const (
	TransferFunction  = "transfer"
	BalanceOfFunction = "balanceOf"
	AmountField       = "amount"
	ToField           = "to"
)

func DecodeTransferTx(t *types.Transaction, logger gethlog.Logger) (bool, *gethcommon.Address, *big.Int) {
	if len(t.Data()) < methodBytesLen {
		return false, nil, nil
	}
	method, err := obscuroERC20ContractABIJSON.MethodById(t.Data()[:methodBytesLen])
	if err != nil {
		logger.Trace("Could not decode tx", log.TxKey, t.Hash(), log.ErrKey, err)
		return false, nil, nil
	}

	if method.Name != TransferFunction {
		return false, nil, nil
	}

	args := map[string]interface{}{}
	if err := method.Inputs.UnpackIntoMap(args, t.Data()[4:]); err != nil {
		panic(err)
	}

	address := args[ToField].(gethcommon.Address)
	amount := args[AmountField].(*big.Int)
	return true, &address, amount
}

func CreateTransferTxData(address gethcommon.Address, amount *big.Int) []byte {
	transferERC20data, err := obscuroERC20ContractABIJSON.Pack(TransferFunction, address, amount)
	if err != nil {
		panic(err)
	}
	return transferERC20data
}

func CreateBalanceOfData(address gethcommon.Address) []byte {
	balanceData, err := obscuroERC20ContractABIJSON.Pack(BalanceOfFunction, address)
	if err != nil {
		panic(err)
	}
	return balanceData
}
