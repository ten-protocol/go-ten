package mempool

import (
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/limiters"
)

type Manager interface {
	// FetchMempoolTxs returns all transactions in the mempool
	FetchMempoolTxs() []*common.L2Tx
	// AddMempoolTx adds a transaction to the mempool
	AddMempoolTx(tx *common.L2Tx) error
	// RemoveTxs removes transactions that are considered immune to re-orgs (i.e. over X batches deep).
	RemoveTxs(transactions types.Transactions) error

	// CurrentTxs Returns the transactions that should be included in the current batch
	CurrentTxs(stateDB *state.StateDB, limiter limiters.BatchSizeLimiter) ([]*common.L2Tx, error)
}
