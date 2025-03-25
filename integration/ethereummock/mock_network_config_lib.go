package ethereummock

import (
	"github.com/ethereum/go-ethereum"
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
	addresses := &common.NetworkConfigAddresses{
		CrossChain:               CrossChainAddr,
		L1MessageBus:             MessageBusAddr,
		EnclaveRegistry:          RespondSecretTxAddr,
		DataAvailabilityRegistry: RollupTxAddr,
	}
	return addresses, nil
}

func (m *mockNetworkConfigLib) AddAdditionalAddress(name string, address gethcommon.Address) (ethereum.CallMsg, error) {
	panic("no-op")
}

func (m *mockNetworkConfigLib) IsMock() bool {
	return true
}
