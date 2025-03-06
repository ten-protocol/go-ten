package erc20contract

import (
	"math/big"

	"github.com/ten-protocol/go-ten/integration/erc20contract/generated/EthERC20"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ten-protocol/go-ten/integration/erc20contract/generated/ObsERC20"
)

func L2BytecodeWithDefaultSupply(tokenName string, busAddress common.Address) []byte {
	return L2Bytecode(tokenName, tokenName, "1000000000000000000000000000000000000000", busAddress)
}

func L2Bytecode(tokenName string, tokenSymbol string, initialSupply string, busAddress common.Address) []byte {
	parsed, err := ObsERC20.ObsERC20MetaData.GetAbi()
	if err != nil {
		panic(err)
	}
	supply, _ := big.NewInt(0).SetString(initialSupply, 10)
	input, err := parsed.Pack("", tokenName, tokenSymbol, supply, busAddress)
	if err != nil {
		panic(err)
	}
	bytecode := common.FromHex(ObsERC20.ObsERC20MetaData.Bin)
	return append(bytecode, input...)
}

func L1BytecodeWithDefaultSupply(tokenName string, crossChainContractAddress common.Address) []byte {
	return L1Bytecode(tokenName, tokenName, "1000000000000000000000000000000000000000", crossChainContractAddress)
}

// FIXME what does managment addr do here?
func L1Bytecode(tokenName string, tokenSymbol string, initialSupply string, crossChainContractAddress common.Address) []byte {
	parsed, err := EthERC20.EthERC20MetaData.GetAbi()
	if err != nil {
		panic(err)
	}
	supply, _ := big.NewInt(0).SetString(initialSupply, 10)
	input, err := parsed.Pack("", tokenName, tokenSymbol, supply, crossChainContractAddress)
	if err != nil {
		panic(err)
	}
	bytecode := common.FromHex(EthERC20.EthERC20MetaData.Bin)
	return append(bytecode, input...)
}
