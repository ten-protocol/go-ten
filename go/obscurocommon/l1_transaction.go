package obscurocommon

import (
	"github.com/ethereum/go-ethereum/common"
)

// L1Transaction represents how obscuro interprets Transactions that happen on the Layer 1
// These transactions can have different implementations
type L1Transaction interface{}

type L1RollupTx struct {
	Rollup EncodedRollup
}

type L1DepositTx struct {
	Amount        uint64
	To            common.Address
	TokenContract *common.Address
}

type L1StoreSecretTx struct {
	Secret      EncryptedSharedEnclaveSecret
	Attestation AttestationReport
}

type L1RequestSecretTx struct {
	Attestation AttestationReport
}
