package contractlib

import (
	"fmt"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	gethlog "github.com/ethereum/go-ethereum/log"
)

type ContractRegistryLib interface {
	RollupLib() RollupContractLib
	NetworkEnclaveLib() NetworkEnclaveRegistryLib
	NetworkConfigLib() NetworkConfigLib
	GetContractAddresses() *NetworkAddresses
	IsMock() bool
}

type contractRegistryImpl struct {
	rollupLib         RollupContractLib
	networkEnclaveLib NetworkEnclaveRegistryLib
	networkConfig     NetworkConfigLib
	addresses         *NetworkAddresses
	logger            gethlog.Logger
}

func NewContractRegistry(networkConfigAddr gethcommon.Address, ethClient ethclient.Client, logger gethlog.Logger) (ContractRegistryLib, error) {
	networkConfig, err := NewNetworkConfigLib(networkConfigAddr, ethClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create NetworkConfig: %w", err)
	}

	addresses, err := networkConfig.GetContractAddresses()
	if err != nil {
		return nil, fmt.Errorf("failed to get contract addresses: %w", err)
	}

	rollupLib := NewRollupContractLib(&addresses.RollupContract, logger)
	networkEnclaveLib := NewNetworkEnclaveRegistryLib(&addresses.NetworkEnclaveRegistry, logger)

	registry := &contractRegistryImpl{
		rollupLib:         rollupLib,
		networkEnclaveLib: networkEnclaveLib,
		networkConfig:     networkConfig,
		addresses:         addresses,
		logger:            logger,
	}

	return registry, nil
}

// NewContractRegistryFromLibs - helper function when creating the contract registry on the enclave
func NewContractRegistryFromLibs(rolluplib RollupContractLib, enclaveRegistryLib NetworkEnclaveRegistryLib, logger gethlog.Logger) *contractRegistryImpl {
	registry := &contractRegistryImpl{
		rollupLib:         rolluplib,
		networkEnclaveLib: enclaveRegistryLib,
		logger:            logger,
	}

	return registry
}

func (r *contractRegistryImpl) GetContractAddresses() *NetworkAddresses {
	return r.addresses
}

func (r *contractRegistryImpl) RollupLib() RollupContractLib {
	return r.rollupLib
}

func (r *contractRegistryImpl) NetworkEnclaveLib() NetworkEnclaveRegistryLib {
	return r.networkEnclaveLib
}

func (r *contractRegistryImpl) NetworkConfigLib() NetworkConfigLib {
	return r.networkConfig
}

func (r *contractRegistryImpl) IsMock() bool { return false }
