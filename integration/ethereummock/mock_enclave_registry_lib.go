package ethereummock

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/ethadapter/contractlib"
)

type MockEnclaveRegistryLib struct{}

func NewEnclaveRegistryLibMock() contractlib.EnclaveRegistryLib {
	return &MockEnclaveRegistryLib{}
}

func (m *MockEnclaveRegistryLib) IsMock() bool {
	return true
}

func (m *MockEnclaveRegistryLib) GetContractAddr() *gethcommon.Address {
	return &InitializeSecretTxAddr
}

func (m *MockEnclaveRegistryLib) DecodeTx(tx *types.Transaction) (common.L1TenTransaction, error) {
	if tx.To() == nil || len(tx.Data()) == 0 {
		return nil, nil
	}
	switch tx.To().Hex() {
	case InitializeSecretTxAddr.Hex():
		return DecodeTx(tx), nil
	case RequestSecretTxAddr.Hex():
		return DecodeTx(tx), nil
	case RespondSecretTxAddr.Hex():
		return DecodeTx(tx), nil
	case GrantSeqTxAddr.Hex():
		return DecodeTx(tx), nil
	default:
		return nil, nil
	}
}

func (m *MockEnclaveRegistryLib) CreateInitializeSecret(tx *common.L1InitializeSecretTx) (types.TxData, error) {
	return EncodeTx(tx, InitializeSecretTxAddr), nil
}

func (m *MockEnclaveRegistryLib) CreateRequestSecret(tx *common.L1RequestSecretTx) (types.TxData, error) {
	return EncodeTx(tx, RequestSecretTxAddr), nil
}

func (m *MockEnclaveRegistryLib) CreateRespondSecret(tx *common.L1RespondSecretTx) (types.TxData, error) {
	return EncodeTx(tx, RespondSecretTxAddr), nil
}
