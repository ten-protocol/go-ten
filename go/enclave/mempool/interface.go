package mempool

import (
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
)

type Manager interface {
	// FetchMempoolTxs returns all transactions in the mempool
	FetchMempoolTxs() []*common.L2Tx
	// AddMempoolTx adds a transaction to the mempool
	AddMempoolTx(tx *common.L2Tx) error
	// RemoveMempoolTxs removes transactions that are considered immune to re-orgs (i.e. over X batches deep).
	RemoveMempoolTxs(batch *core.Batch, resolver db.BatchResolver) error
	// CurrentTxs Returns the transactions that should be included in the current batch
	CurrentTxs(head *core.Batch, resolver db.BatchResolver) ([]*common.L2Tx, error)
}
