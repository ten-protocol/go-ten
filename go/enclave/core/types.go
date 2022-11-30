package core

import (
	"github.com/ethereum/go-ethereum/common"
)

// HeadsAfterL1Block is the heads of the L1 and L2 chains, after processing a given L1 block.
type HeadsAfterL1Block struct {
	HeadBlock  common.Hash // The hash of the processed L1 block.
	HeadRollup common.Hash // The corresponding head rollup.
}
