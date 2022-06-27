package mempool

import (
	"sort"
	"sync"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
	obscurocore "github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/core"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/db"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

// sortByNonce a very primitive way to implement mempool logic that
// adds transactions sorted by the nonce in the rollup
// which is what the EVM expects
type sortByNonce obscurocore.L2Txs

func (c sortByNonce) Len() int           { return len(c) }
func (c sortByNonce) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c sortByNonce) Less(i, j int) bool { return c[i].Nonce() < c[j].Nonce() }

type mempoolManager struct {
	mpMutex        sync.RWMutex // Controls access to `mempool`
	obscuroChainID int64
	mempool        map[common.Hash]*nodecommon.L2Tx
}

func New(chainID int64) Manager {
	return &mempoolManager{
		mempool:        make(map[common.Hash]*nodecommon.L2Tx),
		obscuroChainID: chainID,
		mpMutex:        sync.RWMutex{},
	}
}

func (db *mempoolManager) AddMempoolTx(tx *nodecommon.L2Tx) error {
	db.mpMutex.Lock()
	defer db.mpMutex.Unlock()
	err := obscurocore.VerifySignature(db.obscuroChainID, tx)
	if err != nil {
		return err
	}
	db.mempool[tx.Hash()] = tx
	return nil
}

func (db *mempoolManager) FetchMempoolTxs() []*nodecommon.L2Tx {
	db.mpMutex.RLock()
	defer db.mpMutex.RUnlock()

	mpCopy := make([]*nodecommon.L2Tx, 0)
	for _, tx := range db.mempool {
		mpCopy = append(mpCopy, tx)
	}
	return mpCopy
}

func (db *mempoolManager) RemoveMempoolTxs(toRemove map[common.Hash]common.Hash) {
	db.mpMutex.Lock()
	defer db.mpMutex.Unlock()

	r := make(map[common.Hash]*nodecommon.L2Tx)
	for id, t := range db.mempool {
		_, f := toRemove[id]
		if !f {
			r[id] = t
		}
	}
	db.mempool = r
}

// Calculate transactions to be included in the current rollup
func (db *mempoolManager) CurrentTxs(head *obscurocore.Rollup, resolver db.RollupResolver) obscurocore.L2Txs {
	txs := findTxsNotIncluded(head, db.FetchMempoolTxs(), resolver)
	sort.Sort(sortByNonce(txs))
	return txs
}

// findTxsNotIncluded - given a list of transactions, it keeps only the ones that were not included in the block
// todo - inefficient
func findTxsNotIncluded(head *obscurocore.Rollup, txs []*nodecommon.L2Tx, s db.RollupResolver) []*nodecommon.L2Tx {
	included := allIncludedTransactions(head, s)
	return removeExisting(txs, included)
}

func allIncludedTransactions(r *obscurocore.Rollup, s db.RollupResolver) map[common.Hash]*nodecommon.L2Tx {
	if r.Header.Number.Uint64() == obscurocommon.L2GenesisHeight {
		return obscurocore.MakeMap(r.Transactions)
	}
	newMap := make(map[common.Hash]*nodecommon.L2Tx)
	for k, v := range allIncludedTransactions(s.ParentRollup(r), s) {
		newMap[k] = v
	}
	for _, tx := range r.Transactions {
		newMap[tx.Hash()] = tx
	}
	return newMap
}

func removeExisting(base []*nodecommon.L2Tx, toRemove map[common.Hash]*nodecommon.L2Tx) (r []*nodecommon.L2Tx) {
	for _, t := range base {
		_, f := toRemove[t.Hash()]
		if !f {
			r = append(r, t)
		}
	}
	return
}
