package core

import (
	"github.com/ethereum/go-ethereum/common"
)

// HeadsAfterL1Block pairs the heads of the L1 and L2 chains, at a point in time.
type HeadsAfterL1Block struct {
	HeadBlock         common.Hash // The hash of an L1 block.
	HeadRollup        common.Hash // The head rollup after processing the L1 block.
	UpdatedHeadRollup bool        // Indicates whether ingesting this block updated the head rollup.
}
