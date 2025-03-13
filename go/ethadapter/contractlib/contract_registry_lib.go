package contractlib

import (
	"fmt"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
)

type ContractRegistryLib interface {
	RollupLib() RollupContractLib
	NetworkEnclaveLib() EnclaveRegistryLib
	NetworkConfigLib() NetworkConfigLib
	GetContractAddresses() *common.NetworkConfigAddresses
	IsMock() bool
}

type ContractRegistryImpl struct {
	rollupLib         RollupContractLib
	networkEnclaveLib EnclaveRegistryLib
	networkConfig     NetworkConfigLib
	addresses         *common.NetworkConfigAddresses
	logger            gethlog.Logger
}

func NewContractRegistryLib(networkConfigAddr gethcommon.Address, ethClient ethclient.Client, logger gethlog.Logger) (ContractRegistryLib, error) {
	networkConfig, err := NewNetworkConfigLib(networkConfigAddr, ethClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create NetworkConfig: %w", err)
	}

	addresses, err := networkConfig.GetContractAddresses()
	if err != nil {
		return nil, fmt.Errorf("failed to get contract addresses: %w", err)
	}

	rollupLib := NewRollupContractLib(&addresses.RollupContract, logger)
	networkEnclaveLib := NewEnclaveRegistryLib(&addresses.NetworkEnclaveRegistry, logger)

	registry := &ContractRegistryImpl{
		rollupLib:         rollupLib,
		networkEnclaveLib: networkEnclaveLib,
		networkConfig:     networkConfig,
		addresses:         addresses,
		logger:            logger,
	}

	return registry, nil
}

// NewContractRegistryFromLibs - helper function when creating the contract registry on the enclave
func NewContractRegistryFromLibs(rolluplib RollupContractLib, enclaveRegistryLib EnclaveRegistryLib, logger gethlog.Logger) *ContractRegistryImpl {
	registry := &ContractRegistryImpl{
		rollupLib:         rolluplib,
		networkEnclaveLib: enclaveRegistryLib,
		logger:            logger,
	}

	return registry
}

func (r *ContractRegistryImpl) GetContractAddresses() *common.NetworkConfigAddresses {
	return r.addresses
}

func (r *ContractRegistryImpl) RollupLib() RollupContractLib {
	return r.rollupLib
}

func (r *ContractRegistryImpl) NetworkEnclaveLib() EnclaveRegistryLib {
	return r.networkEnclaveLib
}

func (r *ContractRegistryImpl) NetworkConfigLib() NetworkConfigLib {
	return r.networkConfig
}

func (r *ContractRegistryImpl) IsMock() bool { return false }
