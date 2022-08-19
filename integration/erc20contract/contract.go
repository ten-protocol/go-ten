package erc20contract

import (
	"math/big"

	"github.com/obscuronet/go-obscuro/integration/erc20contract/generated/EthERC20"

	"github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/go-obscuro/integration/erc20contract/generated/ObsERC20"
)

func L2BytecodeWithDefaultSupply(tokenName string) []byte {
	return L2Bytecode(tokenName, "1000000000000000000000000000000000000000")
}

func L2Bytecode(tokenName string, initialSupply string) []byte {
	parsed, err := ObsERC20.ObsERC20MetaData.GetAbi()
	if err != nil {
		panic(err)
	}
	supply, _ := big.NewInt(0).SetString(initialSupply, 10)
	input, err := parsed.Pack("", tokenName, tokenName, supply)
	if err != nil {
		panic(err)
	}
	bytecode := common.FromHex(ObsERC20.ObsERC20MetaData.Bin)
	return append(bytecode, input...)
}

func L1BytecodeWithDefaultSupply(tokenName string) []byte {
	return L1Bytecode(tokenName, "1000000000000000000000000000000000000000")
}

func L1Bytecode(tokenName string, initialSupply string) []byte {
	parsed, err := EthERC20.EthERC20MetaData.GetAbi()
	if err != nil {
		panic(err)
	}
	supply, _ := big.NewInt(0).SetString(initialSupply, 10)
	input, err := parsed.Pack("", tokenName, tokenName, supply)
	if err != nil {
		panic(err)
	}
	bytecode := common.FromHex(EthERC20.EthERC20MetaData.Bin)
	return append(bytecode, input...)
}
