package contractlib

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/ethadapter/contractlib"
)

type MockNetworkEnclaveRegistryLib struct{}

func NewNetworkEnclaveRegistryLibMock() contractlib.NetworkEnclaveRegistryLib {
	return &MockNetworkEnclaveRegistryLib{}
}

func (m *MockNetworkEnclaveRegistryLib) IsMock() bool {
	return true
}

func (m *MockNetworkEnclaveRegistryLib) GetContractAddr() *gethcommon.Address {
	return &StoreSecretTxAddr
}

func (m *MockNetworkEnclaveRegistryLib) DecodeTx(tx *types.Transaction) (common.L1TenTransaction, error) {
	if tx.To() == nil || len(tx.Data()) == 0 {
		return nil, nil
	}

	switch tx.To().Hex() {
	case StoreSecretTxAddr.Hex():
		return DecodeTx(tx), nil
	case RequestSecretTxAddr.Hex():
		return DecodeTx(tx), nil
	case InitializeSecretTxAddr.Hex():
		return DecodeTx(tx), nil
	default:
		return nil, nil
	}
}

func (m *MockNetworkEnclaveRegistryLib) CreateInitializeSecret(tx *common.L1InitializeSecretTx) (types.TxData, error) {
	return EncodeTx(tx, InitializeSecretTxAddr), nil
}

func (m *MockNetworkEnclaveRegistryLib) CreateRequestSecret(tx *common.L1RequestSecretTx) (types.TxData, error) {
	return EncodeTx(tx, RequestSecretTxAddr), nil
}

func (m *MockNetworkEnclaveRegistryLib) CreateRespondSecret(tx *common.L1RespondSecretTx, _ bool) (types.TxData, error) {
	return EncodeTx(tx, StoreSecretTxAddr), nil
}
