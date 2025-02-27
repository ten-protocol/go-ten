package contractlib

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/ethadapter/contractlib"
)

type mockContractRegistry struct {
	rollupLib         mockRollupContractLib
	networkEnclaveLib MockNetworkEnclaveRegistryLib
	networkConfigLib  *contractlib.NetworkConfigLib
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

func (m *mockContractRegistry) NetworkConfigLib() *contractlib.NetworkConfigLib {
	networkConfigMock := NewNetworkConfigLibMock()
	return &networkConfigMock
}

func (m *mockContractRegistry) GetContractAddresses() *contractlib.NetworkAddresses {
	addresses, _ := m.networkConfigLib.GetContractAddresses()
	return &addresses
}

func (m *mockContractRegistry) GetContractByAddress(addr gethcommon.Address) contractlib.ContractLib {
	switch addr {
	case RollupTxAddr:
		return m.rollupLib
	case StoreSecretTxAddr, RequestSecretTxAddr, InitializeSecretTxAddr:
		return m.networkEnclaveLib
	default:
		return nil
	}
}
