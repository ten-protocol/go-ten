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
	Amount        uint64          // Amount to be deposited
	To            *common.Address // Address the ERC20 Transfer was made to (always be the Management Contract Addr)
	Sender        *common.Address // Address that issued the ERC20, the token holder or tx.origin
	TokenContract *common.Address // Address of the ERC20 Contract address that was executed
}

type L1RespondSecretTx struct {
	Secret          EncryptedSharedEnclaveSecret
	Attestation     EncodedAttestationReport
	RequesterPubKey []byte
	RequesterID     common.Address
}

type L1RequestSecretTx struct {
	Attestation EncodedAttestationReport
}
