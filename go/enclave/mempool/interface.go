package mempool

import (
	"github.com/obscuronet/obscuro-playground/go/common"
	obscurocore "github.com/obscuronet/obscuro-playground/go/enclave/core"
	"github.com/obscuronet/obscuro-playground/go/enclave/db"
)

type Manager interface {
	// FetchMempoolTxs returns all transactions in the mempool
	FetchMempoolTxs() []*common.L2Tx
	// AddMempoolTx adds an transaction to the mempool
	AddMempoolTx(tx *common.L2Tx) error
	// RemoveMempoolTxs removes transactions that are considered immune to re-orgs
	RemoveMempoolTxs(r *obscurocore.Rollup, resolver db.RollupResolver)
	// CurrentTxs Returns the transactions that should be included in the current rollup
	CurrentTxs(head *obscurocore.Rollup, resolver db.RollupResolver) obscurocore.L2Txs
}
