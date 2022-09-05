package core

import (
	"github.com/ethereum/go-ethereum/common"
)

// BlockState pairs a block with its corresponding rollup.
type BlockState struct {
	Block          common.Hash
	HeadRollup     common.Hash
	FoundNewRollup bool
}
