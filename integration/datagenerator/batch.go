package datagenerator

import (
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/go-obscuro/go/common"
)

// RandomBatch - block is needed in order to pass the smart contract check
// when submitting cross chain messages.
func RandomBatch(block *types.Block) common.ExtBatch {
	extBatch := common.ExtBatch{
		Header: &common.Header{
			ParentHash:  randomHash(),
			Agg:         RandomAddress(),
			L1Proof:     randomHash(),
			Root:        randomHash(),
			Number:      big.NewInt(int64(RandomUInt64())),
			Withdrawals: randomWithdrawals(10),
		},
		TxHashes:        []gethcommon.Hash{randomHash()},
		EncryptedTxBlob: RandomBytes(10),
	}

	if block != nil {
		extBatch.Header.LatestInboundCrossChainHeight = block.Number()
		extBatch.Header.LatestInboudCrossChainHash = block.Hash()
	}

	return extBatch
}
