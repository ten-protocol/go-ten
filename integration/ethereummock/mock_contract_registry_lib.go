package ethereummock

import (
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/ethadapter/contractlib"
)

type mockContractRegistryLib struct{}

func (m *mockContractRegistryLib) GetContractAddresses() *common.NetworkConfigAddresses {
	addresses := &common.NetworkConfigAddresses{
		CrossChain:             RollupTxAddr,
		MessageBus:             MessageBusAddr,
		NetworkEnclaveRegistry: RespondSecretTxAddr,
		RollupContract:         RollupTxAddr,
	}
	return addresses
}

func NewContractRegistryLibMock() contractlib.ContractRegistryLib {
	return &mockContractRegistryLib{}
}

func (m *mockContractRegistryLib) RollupLib() contractlib.RollupContractLib {
	return NewRollupContractLibMock()
}

func (m *mockContractRegistryLib) NetworkEnclaveLib() contractlib.EnclaveRegistryLib {
	return NewNetworkEnclaveRegistryLibMock()
}

func (m *mockContractRegistryLib) NetworkConfigLib() contractlib.NetworkConfigLib {
	return NewNetworkConfigLibMock()
}

func (m *mockContractRegistryLib) IsMock() bool { return true }
