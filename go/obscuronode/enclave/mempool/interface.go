package mempool

import (
	"github.com/ethereum/go-ethereum/common"
	obscurocore "github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/core"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/db"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

type Manager interface {
	// FetchMempoolTxs returns all L2 transactions in the mempool
	FetchMempoolTxs() []*nodecommon.L2Tx
	// AddMempoolTx adds an L2 transaction to the mempool
	AddMempoolTx(tx *nodecommon.L2Tx) error
	// RemoveMempoolTxs removes any L2 transactions whose hash is keyed in the map from the mempool
	RemoveMempoolTxs(toRemove map[common.Hash]common.Hash)
	// CurrentTxs Returns the transactins that should be included in the curent rollup
	CurrentTxs(head *obscurocore.Rollup, resolver db.RollupResolver) obscurocore.L2Txs
}
