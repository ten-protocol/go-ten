package ManagementContract

import (
	"github.com/ethereum/go-ethereum/common"
)

func Bytecode() ([]byte, error) {
	parsed, err := ManagementContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	input, err := parsed.Pack("")
	if err != nil {
		return nil, err
	}
	bytecode := common.FromHex(ManagementContractMetaData.Bin)
	return append(bytecode, input...), nil
}
