package core

import "github.com/ethereum/go-ethereum/common"

// BlockState - Represents the state after an L1 Block was processed.
type BlockState struct {
	Block          common.Hash
	HeadRollup     common.Hash
	FoundNewRollup bool
}
