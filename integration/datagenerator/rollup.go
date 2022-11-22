package datagenerator

import (
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/go-obscuro/go/common"
)

func RandomRollup(block *types.Block) common.ExtRollup {
	extRollup := common.ExtRollup{
		Header: &common.Header{
			ParentHash:                    randomHash(),
			Agg:                           RandomAddress(),
			L1Proof:                       randomHash(),
			Root:                          randomHash(),
			Number:                        big.NewInt(int64(randomUInt64())),
			Withdrawals:                   randomWithdrawals(10),
			LatestInboundCrossChainHeight: block.Number(),
			LatestInboudCrossChainHash:    block.Hash(),
		},
		TxHashes:        []gethcommon.Hash{randomHash()},
		EncryptedTxBlob: RandomBytes(10),
	}
	return extRollup
}
