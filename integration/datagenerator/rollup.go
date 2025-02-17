package datagenerator

import (
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ten-protocol/go-ten/go/common"
)

// RandomRollup - block is needed in order to pass the smart contract check
// when submitting cross chain messages.
func RandomRollup(block *types.Header) common.ExtRollup {
	var l1Head gethcommon.Hash
	var l1Number *big.Int

	if block != nil {
		l1Head = block.Hash()
		l1Number = block.Number
	}

	extRollup := common.ExtRollup{
		Header: &common.RollupHeader{
			CompressionL1Head:   l1Head,
			CompressionL1Number: l1Number,
			LastBatchSeqNo:      1,
			CrossChainRoot:      gethcommon.Hash{}, // Empty root for test
		},
	}
	return extRollup
}
