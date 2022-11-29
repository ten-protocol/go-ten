package core

import (
	"github.com/ethereum/go-ethereum/common"
)

// BlockState pairs the hash of an L1 block with the hash of the head rollup in the L2 chain after processing that block.
type BlockState struct {
	Block             common.Hash // The hash of an L1 block.
	HeadRollup        common.Hash // The head rollup after processing the L1 block.
	UpdatedHeadRollup bool        // Indicates whether ingesting this block updated the head rollup.
}
