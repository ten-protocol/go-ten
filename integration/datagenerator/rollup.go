package datagenerator

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ten-protocol/go-ten/go/common"
)

// RandomRollup - block is needed in order to pass the smart contract check
// when submitting cross chain messages.
// TODO @will this wont work now since we verify the data on the management contract
func RandomRollup(block *types.Header) common.ExtRollup {
	extRollup := common.ExtRollup{
		Header: &common.RollupHeader{
			CompressionL1Head:   block.Hash(),
			CompressionL1Number: block.Number,
			LastBatchSeqNo:      1,
			CrossChainRoot:      gethcommon.Hash{}, // Empty root for test
			BlobHash:            gethcommon.Hash{}, // Will be set by contract from blob data
		},
	}
	return extRollup
}
