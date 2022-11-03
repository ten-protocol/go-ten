package core

import (
	"github.com/ethereum/go-ethereum/common"
)

// BlockState pairs a block with the rollup it contains.
type BlockState struct {
	Block          common.Hash
	FoundNewRollup bool        // Whether the ingested block contains a new rollup.
	NewRollup      common.Hash // If the block contains a new rollup, the new rollup.
}
