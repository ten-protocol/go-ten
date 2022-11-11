package contracts

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func Bytecode(contractMetaData *bind.MetaData) ([]byte, error) {
	parsed, err := contractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	input, err := parsed.Pack("")
	if err != nil {
		return nil, err
	}
	bytecode := common.FromHex(contractMetaData.Bin)
	return append(bytecode, input...), nil
}
