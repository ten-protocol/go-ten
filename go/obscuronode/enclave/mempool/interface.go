package mempool

import (
	obscurocore "github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/core"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/db"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

type Manager interface {
	// FetchMempoolTxs returns all transactions in the mempool
	FetchMempoolTxs() []*nodecommon.L2Tx
	// AddMempoolTx adds an transaction to the mempool
	AddMempoolTx(tx *nodecommon.L2Tx) error
	// RemoveMempoolTxs removes transactions that are considered immune to re-orgs
	RemoveMempoolTxs(r *obscurocore.Rollup, resolver db.RollupResolver)
	// CurrentTxs Returns the transactions that should be included in the current rollup
	CurrentTxs(head *obscurocore.Rollup, resolver db.RollupResolver) obscurocore.L2Txs
}
