package contractlib

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/ethadapter/contractlib"
)

type mockNetworkConfigLib struct{}

func NewNetworkConfigLibMock() contractlib.NetworkConfigLib {
	return &mockNetworkConfigLib{}
}

func (m *mockNetworkConfigLib) GetContractAddr() *gethcommon.Address {
	return &NetworkConfigAddr
}

func (m *mockNetworkConfigLib) GetContractAddresses() (*contractlib.NetworkAddresses, error) {
	return &contractlib.NetworkAddresses{
		CrossChain:             CrossChainAddr,
		MessageBus:             MessageBusAddr,
		NetworkEnclaveRegistry: StoreSecretTxAddr,
		RollupContract:         RollupTxAddr,
	}, nil
}
