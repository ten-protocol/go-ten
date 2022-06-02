package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/obscuronode/nodecommon"
)

// L2Txs Todo - is this type useful?
type L2Txs []nodecommon.L2Tx

// BlockState - Represents the state after an L1 Block was processed.
type BlockState struct {
	Block          common.Hash
	HeadRollup     common.Hash
	FoundNewRollup bool
}
