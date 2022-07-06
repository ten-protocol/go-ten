package core

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/common"
)

const SharedSecretLen = 32

// SharedEnclaveSecret - the entropy
type SharedEnclaveSecret [SharedSecretLen]byte

// L2Txs Todo - is this type useful?
type L2Txs []*common.L2Tx

// BlockState - Represents the state after an L1 Block was processed.
type BlockState struct {
	Block          gethcommon.Hash
	HeadRollup     gethcommon.Hash
	FoundNewRollup bool
}
