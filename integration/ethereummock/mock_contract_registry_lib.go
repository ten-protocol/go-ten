package ethereummock

import (
	"github.com/ten-protocol/go-ten/go/ethadapter/contractlib"
	"github.com/ten-protocol/go-ten/integration"
)

type mockContractRegistryLib struct {
	rollupLib         mockRollupContractLib
	networkEnclaveLib MockNetworkEnclaveRegistryLib
	networkConfigLib  contractlib.NetworkConfigLib
}

func (m *mockContractRegistryLib) GetContractAddresses() *contractlib.NetworkAddresses {
	return &contractlib.NetworkAddresses{
		CrossChain:             integration.RollupTxAddr,
		MessageBus:             integration.MessageBusAddr,
		NetworkEnclaveRegistry: integration.StoreSecretTxAddr,
		RollupContract:         integration.RollupTxAddr,
	}
}

func NewContractRegistryLibMock() contractlib.ContractRegistryLib {
	return &mockContractRegistryLib{}
}

func (m *mockContractRegistryLib) RollupLib() contractlib.RollupContractLib {
	return NewRollupContractLibMock()
}

func (m *mockContractRegistryLib) NetworkEnclaveLib() contractlib.NetworkEnclaveRegistryLib {
	return NewNetworkEnclaveRegistryLibMock()
}

func (m *mockContractRegistryLib) NetworkConfigLib() contractlib.NetworkConfigLib {
	return NewNetworkConfigLibMock()
}

func (m *mockContractRegistryLib) IsMock() bool { return true }
