package datagenerator

import (
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

func RandomRollup() nodecommon.Rollup {
	return nodecommon.Rollup{
		Header: &nodecommon.Header{
			ParentHash:  randomHash(),
			Agg:         randomAddress(),
			Nonce:       randomUInt64(),
			L1Proof:     randomHash(),
			State:       randomHash(),
			Number:      randomUInt64(),
			Withdrawals: randomWithdrawals(10),
		},
		Transactions: randomBytes(10),
	}
}
