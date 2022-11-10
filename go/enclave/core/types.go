package core

import (
	"github.com/ethereum/go-ethereum/common"
)

// BlockState pairs a block with the canonical head rollup and the rollup the block contains, if any.
type BlockState struct {
	Block          common.Hash
	HeadRollup     common.Hash // The head rollup of the canonical L2 chain.
	FoundNewRollup bool        // Whether the ingested block contains a new rollup.
}
