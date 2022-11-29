package core

import (
	"github.com/ethereum/go-ethereum/common"
)

// ChainHeads pairs the hash of the head L1 block with the hash of the head rollup after processing that block.
type ChainHeads struct {
	HeadBlock         common.Hash // The hash of an L1 block.
	HeadRollup        common.Hash // The head rollup after processing the L1 block.
	UpdatedHeadRollup bool        // Indicates whether ingesting this block updated the head rollup.
}
