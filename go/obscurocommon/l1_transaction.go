package obscurocommon

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// L1Transaction represents how obscuro interprets Transactions that happen on the Layer 1
// These transactions can have different implementations
type L1Transaction interface {
	Bytes() []byte
}

type L1RollupTx struct {
	Rollup EncodedRollup
}

func (t *L1RollupTx) Bytes() []byte { return []byte(fmt.Sprintf("%v", *t)) }

type L1DepositTx struct {
	Amount        uint64
	To            common.Address
	TokenContract common.Address
}

func (t *L1DepositTx) Bytes() []byte { return []byte(fmt.Sprintf("%v", *t)) }

type L1StoreSecretTx struct {
	Secret      EncryptedSharedEnclaveSecret
	Attestation AttestationReport
}

func (t *L1StoreSecretTx) Bytes() []byte { return []byte(fmt.Sprintf("%v", *t)) }

type L1RequestSecretTx struct {
	Attestation AttestationReport
}

func (t *L1RequestSecretTx) Bytes() []byte { return []byte(fmt.Sprintf("%v", *t)) }
