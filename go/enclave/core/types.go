package core

import (
	"github.com/ethereum/go-ethereum/common"
)

const SharedSecretLen = 32

// SharedEnclaveSecret - the entropy
type SharedEnclaveSecret [SharedSecretLen]byte

// BlockState - Represents the state after an L1 Block was processed.
type BlockState struct {
	Block          common.Hash
	HeadRollup     common.Hash
	FoundNewRollup bool
}
