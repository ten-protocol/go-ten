package contractlib

import (
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
)

type ContractRegistryLib interface {
	DARegistryLib() DataAvailabilityRegistryLib
	EnclaveRegistryLib() EnclaveRegistryLib
	NetworkConfigLib() NetworkConfigLib
	GetContractAddresses() *common.NetworkConfigAddresses
	IsMock() bool
}

type ContractRegistryImpl struct {
	daRegistryLib     DataAvailabilityRegistryLib
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

	daRegistryLib := NewDataAvailabilityRegistryLib(&addresses.DataAvailabilityRegistry, logger)
	networkEnclaveLib := NewEnclaveRegistryLib(&addresses.EnclaveRegistry, logger)

	registry := &ContractRegistryImpl{
		daRegistryLib:     daRegistryLib,
		networkEnclaveLib: networkEnclaveLib,
		networkConfig:     networkConfig,
		addresses:         addresses,
		logger:            logger,
	}

	return registry, nil
}

// NewContractRegistryFromLibs - helper function when creating the contract registry on the enclave
func NewContractRegistryFromLibs(daRegistryLib DataAvailabilityRegistryLib, enclaveRegistryLib EnclaveRegistryLib, logger gethlog.Logger) *ContractRegistryImpl {
	registry := &ContractRegistryImpl{
		daRegistryLib:     daRegistryLib,
		networkEnclaveLib: enclaveRegistryLib,
		logger:            logger,
	}

	return registry
}

func (r *ContractRegistryImpl) GetContractAddresses() *common.NetworkConfigAddresses {
	return r.addresses
}

func (r *ContractRegistryImpl) DARegistryLib() DataAvailabilityRegistryLib {
	return r.daRegistryLib
}

func (r *ContractRegistryImpl) EnclaveRegistryLib() EnclaveRegistryLib {
	return r.networkEnclaveLib
}

func (r *ContractRegistryImpl) NetworkConfigLib() NetworkConfigLib {
	return r.networkConfig
}

func (r *ContractRegistryImpl) IsMock() bool { return false }
