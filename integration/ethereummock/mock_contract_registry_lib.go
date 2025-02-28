package ethereummock

import (
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/ethadapter/contractlib"
)

type mockContractRegistryLib struct {
	rollupLib         mockRollupContractLib
	networkEnclaveLib MockNetworkEnclaveRegistryLib
	networkConfigLib  contractlib.NetworkConfigLib
}

func (m *mockContractRegistryLib) GetContractAddresses() *common.NetworkAddresses {
	return &common.NetworkAddresses{
		CrossChain:             RollupTxAddr,
		MessageBus:             MessageBusAddr,
		NetworkEnclaveRegistry: StoreSecretTxAddr,
		RollupContract:         RollupTxAddr,
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
