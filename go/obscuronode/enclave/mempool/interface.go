package mempool

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/obscuronode/nodecommon"
)

type Manager interface {
	// FetchMempoolTxs returns all L2 transactions in the mempool
	FetchMempoolTxs() []nodecommon.L2Tx
	// AddMempoolTx adds an L2 transaction to the mempool
	AddMempoolTx(tx nodecommon.L2Tx)
	// RemoveMempoolTxs removes any L2 transactions whose hash is keyed in the map from the mempool
	RemoveMempoolTxs(toRemove map[common.Hash]common.Hash)
}
