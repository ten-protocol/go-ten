package core

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
)

const SharedSecretLen = 32

// SharedEnclaveSecret - the entropy
type SharedEnclaveSecret [SharedSecretLen]byte

// BlockState - Represents the state after an L1 Block was processed.
type BlockState struct {
	Block          gethcommon.Hash
	HeadRollup     gethcommon.Hash
	FoundNewRollup bool
}
