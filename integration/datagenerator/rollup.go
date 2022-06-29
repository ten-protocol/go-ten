package datagenerator

import (
	"math/big"

	"github.com/obscuronet/obscuro-playground/go/common"
)

func RandomRollup() common.EncryptedRollup {
	return common.EncryptedRollup{
		Header: &common.Header{
			ParentHash:  randomHash(),
			Agg:         RandomAddress(),
			Nonce:       randomUInt64(),
			L1Proof:     randomHash(),
			Root:        randomHash(),
			Number:      big.NewInt(int64(randomUInt64())),
			Withdrawals: randomWithdrawals(10),
		},
		Transactions: RandomBytes(10),
	}
}
