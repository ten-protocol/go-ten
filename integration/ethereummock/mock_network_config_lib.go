package ethereummock

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/ethadapter/contractlib"
)

type mockNetworkConfigLib struct{}

func NewNetworkConfigLibMock() contractlib.NetworkConfigLib {
	return &mockNetworkConfigLib{}
}

func (m *mockNetworkConfigLib) GetContractAddr() *gethcommon.Address {
	return &NetworkConfigAddr
}

func (m *mockNetworkConfigLib) GetContractAddresses() (*common.NetworkConfigAddresses, error) {
	return &common.NetworkConfigAddresses{
		CrossChain:             CrossChainAddr,
		MessageBus:             MessageBusAddr,
		NetworkEnclaveRegistry: StoreSecretTxAddr,
		RollupContract:         RollupTxAddr,
	}, nil
}

func (m *mockNetworkConfigLib) IsMock() bool {
	return true
}
