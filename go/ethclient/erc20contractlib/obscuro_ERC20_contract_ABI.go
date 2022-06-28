package erc20contractlib

import (
	"math/big"
	"strings"

	"github.com/obscuronet/obscuro-playground/go/common/log"

	"github.com/obscuronet/obscuro-playground/go/common"

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

func DecodeTransferTx(t *types.Transaction) (bool, *gethcommon.Address, *big.Int) {
	method, err := obscuroERC20ContractABIJSON.MethodById(t.Data()[:methodBytesLen])
	if err != nil {
		log.Info("Could not decode tx %d, Err: %s", common.ShortHash(t.Hash()), err)
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

func CreateTransferTxData(address gethcommon.Address, amount uint64) []byte {
	transferERC20data, err := obscuroERC20ContractABIJSON.Pack(TransferFunction, address, big.NewInt(int64(amount)))
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
