package constants

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/contracts/generated/ManagementContract"
)

func Bytecode() ([]byte, error) {
	parsed, err := ManagementContract.ManagementContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	input, err := parsed.Pack("")
	if err != nil {
		return nil, err
	}
	bytecode := common.FromHex(ManagementContract.ManagementContractMetaData.Bin)
	return append(bytecode, input...), nil
}
