package erc20contractlib

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/log"
	"github.com/obscuronet/go-obscuro/go/obscurocommon"

	"github.com/ethereum/go-ethereum/accounts/abi"
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

func DecodeTransferTx(t types.Transaction) (bool, *common.Address, *big.Int) {
	method, err := obscuroERC20ContractABIJSON.MethodById(t.Data()[:methodBytesLen])
	if err != nil {
		log.Info("Could not decode tx %d, Err: %s", obscurocommon.ShortHash(t.Hash()), err)
		return false, nil, nil
	}

	if method.Name != TransferFunction {
		return false, nil, nil
	}

	args := map[string]interface{}{}
	if err := method.Inputs.UnpackIntoMap(args, t.Data()[4:]); err != nil {
		panic(err)
	}

	address := args[ToField].(common.Address)
	amount := args[AmountField].(*big.Int)
	return true, &address, amount
}

func CreateTransferTxData(address common.Address, amount uint64) []byte {
	transferERC20data, err := obscuroERC20ContractABIJSON.Pack(TransferFunction, address, big.NewInt(int64(amount)))
	if err != nil {
		panic(err)
	}
	return transferERC20data
}

func CreateBalanceOfData(address common.Address) []byte {
	balanceData, err := obscuroERC20ContractABIJSON.Pack(BalanceOfFunction, address)
	if err != nil {
		panic(err)
	}
	return balanceData
}
