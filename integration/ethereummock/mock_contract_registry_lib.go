package ethereummock

import (
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/ethadapter/contractlib"
)

type mockContractRegistryLib struct{}

func (m *mockContractRegistryLib) GetContractAddresses() *common.NetworkConfigAddresses {
	addresses := &common.NetworkConfigAddresses{
		CrossChain:      RollupTxAddr,
		L1MessageBus:    MessageBusAddr,
		EnclaveRegistry: RespondSecretTxAddr,
		RollupContract:  RollupTxAddr,
	}
	return addresses
}

func NewContractRegistryLibMock() contractlib.ContractRegistryLib {
	return &mockContractRegistryLib{}
}

func (m *mockContractRegistryLib) RollupLib() contractlib.RollupContractLib {
	return NewRollupContractLibMock()
}

func (m *mockContractRegistryLib) EnclaveRegistryLib() contractlib.EnclaveRegistryLib {
	return NewEnclaveRegistryLibMock()
}

func (m *mockContractRegistryLib) NetworkConfigLib() contractlib.NetworkConfigLib {
	return NewNetworkConfigLibMock()
}

func (m *mockContractRegistryLib) IsMock() bool { return true }
