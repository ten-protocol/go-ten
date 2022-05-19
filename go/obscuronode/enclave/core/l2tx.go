package core

import (
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

type (
	L2Txs []nodecommon.L2Tx
	// L2TxType indicates the type of L2 transaction - either a transfer or a withdrawal for now
	L2TxType uint8
)
