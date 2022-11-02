package datagenerator

import (
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/go-obscuro/go/common"
)

func RandomRollup() common.ExtRollupWithHash {
	extRollup := common.ExtRollup{
		Header: &common.Header{
			ParentHash:  randomHash(),
			Agg:         RandomAddress(),
			RollupNonce: randomUInt64(),
			L1Proof:     randomHash(),
			Root:        randomHash(),
			Number:      big.NewInt(int64(randomUInt64())),
			Withdrawals: randomWithdrawals(10),
		},
		TxHashes:        []gethcommon.Hash{randomHash()},
		EncryptedTxBlob: RandomBytes(10),
	}
	return common.ExtRollupWithHash{
		ExtRollup: extRollup,
	}
}
