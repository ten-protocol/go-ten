package constants

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/contracts/generated/CrossChain"
	"github.com/ten-protocol/go-ten/contracts/generated/NetworkConfig"
	"github.com/ten-protocol/go-ten/contracts/generated/NetworkEnclaveRegistry"
	"github.com/ten-protocol/go-ten/contracts/generated/RollupContract"
)

func NetworkConfigBytecode() ([]byte, error) {
	parsed, err := NetworkConfig.NetworkConfigMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	input, err := parsed.Pack("")
	if err != nil {
		return nil, err
	}
	bytecode := common.FromHex(NetworkConfig.NetworkConfigMetaData.Bin)
	return append(bytecode, input...), nil
}

func EnclaveRegistryBytecode() ([]byte, error) {
	parsed, err := NetworkEnclaveRegistry.NetworkEnclaveRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	input, err := parsed.Pack("")
	if err != nil {
		return nil, err
	}
	bytecode := common.FromHex(NetworkEnclaveRegistry.NetworkEnclaveRegistryMetaData.Bin)
	return append(bytecode, input...), nil
}

func RollupContractBytecode() ([]byte, error) {
	parsed, err := RollupContract.RollupContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	input, err := parsed.Pack("")
	if err != nil {
		return nil, err
	}
	bytecode := common.FromHex(RollupContract.RollupContractMetaData.Bin)
	return append(bytecode, input...), nil
}

func CrossChainBytecode() ([]byte, error) {
	parsed, err := CrossChain.CrossChainMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	input, err := parsed.Pack("")
	if err != nil {
		return nil, err
	}
	bytecode := common.FromHex(CrossChain.CrossChainMetaData.Bin)
	return append(bytecode, input...), nil
}
