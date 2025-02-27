package contractlib

import (
	"github.com/ten-protocol/go-ten/go/ethadapter/contractlib"
)

type mockContractRegistry struct {
	rollupLib         mockRollupContractLib
	networkEnclaveLib MockNetworkEnclaveRegistryLib
	networkConfigLib  contractlib.NetworkConfigLib
}

func (m *mockContractRegistry) GetContractAddresses() *contractlib.NetworkAddresses {
	return &contractlib.NetworkAddresses{
		CrossChain:             RollupTxAddr,
		MessageBus:             MessageBusAddr,
		NetworkEnclaveRegistry: StoreSecretTxAddr,
		RollupContract:         RollupTxAddr,
	}
}

func NewContractRegistryMock() contractlib.ContractRegistry {
	return &mockContractRegistry{}
}

func (m *mockContractRegistry) RollupLib() contractlib.RollupContractLib {
	return NewRollupContractLibMock()
}

func (m *mockContractRegistry) NetworkEnclaveLib() contractlib.NetworkEnclaveRegistryLib {
	return NewNetworkEnclaveRegistryLibMock()
}

func (m *mockContractRegistry) NetworkConfigLib() contractlib.NetworkConfigLib {
	return NewNetworkConfigLibMock()
}
