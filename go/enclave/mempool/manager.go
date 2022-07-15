package mempool

import (
	"sort"
	"sync"

	gethcommon "github.com/ethereum/go-ethereum/common"

	obscurocore "github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/db"

	"github.com/obscuronet/go-obscuro/go/common"
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
	mempool        map[gethcommon.Hash]*common.L2Tx
}

func New(chainID int64) Manager {
	return &mempoolManager{
		mempool:        make(map[gethcommon.Hash]*common.L2Tx),
		obscuroChainID: chainID,
		mpMutex:        sync.RWMutex{},
	}
}

func (db *mempoolManager) AddMempoolTx(tx *common.L2Tx) error {
	db.mpMutex.Lock()
	defer db.mpMutex.Unlock()
	err := obscurocore.VerifySignature(db.obscuroChainID, tx)
	if err != nil {
		return err
	}
	db.mempool[tx.Hash()] = tx
	return nil
}

func (db *mempoolManager) FetchMempoolTxs() []*common.L2Tx {
	db.mpMutex.RLock()
	defer db.mpMutex.RUnlock()

	mpCopy := make([]*common.L2Tx, 0)
	for _, tx := range db.mempool {
		mpCopy = append(mpCopy, tx)
	}
	return mpCopy
}

func (db *mempoolManager) RemoveMempoolTxs(rollup *obscurocore.Rollup, resolver db.RollupResolver) {
	db.mpMutex.Lock()
	defer db.mpMutex.Unlock()

	toRemove := historicTxs(rollup, resolver)
	r := make(map[gethcommon.Hash]*common.L2Tx)
	for id, t := range db.mempool {
		_, f := toRemove[id]
		if !f {
			r[id] = t
		}
	}
	db.mempool = r
}

// Returns all transactions found 20 levels below
func historicTxs(r *obscurocore.Rollup, resolver db.RollupResolver) map[gethcommon.Hash]gethcommon.Hash {
	i := common.HeightCommittedBlocks
	c := r
	// todo - create method to return the canonical rollup from height N
	for {
		if i == 0 || c.Header.Number.Uint64() == common.L2GenesisHeight {
			return obscurocore.ToMap(c.Transactions)
		}
		i--
		c = resolver.ParentRollup(c)
	}
}

// CurrentTxs - Calculate transactions to be included in the current rollup
func (db *mempoolManager) CurrentTxs(head *obscurocore.Rollup, resolver db.RollupResolver) obscurocore.L2Txs {
	txs := findTxsNotIncluded(head, db.FetchMempoolTxs(), resolver)
	sort.Sort(sortByNonce(txs))
	return txs
}

// findTxsNotIncluded - given a list of transactions, it keeps only the ones that were not included in the block
// todo - inefficient
func findTxsNotIncluded(head *obscurocore.Rollup, txs []*common.L2Tx, s db.RollupResolver) []*common.L2Tx {
	// go back only HeightCommittedBlocks blocks to accumulate transactions to be diffed against the mempool
	startAt := uint64(0)
	if head.NumberU64() > common.HeightCommittedBlocks {
		startAt = head.NumberU64() - common.HeightCommittedBlocks
	}
	included := allIncludedTransactions(head, s, startAt)
	return removeExisting(txs, included)
}

func allIncludedTransactions(r *obscurocore.Rollup, s db.RollupResolver, stopAtHeight uint64) map[gethcommon.Hash]*common.L2Tx {
	if r.Header.Number.Uint64() == stopAtHeight {
		return obscurocore.MakeMap(r.Transactions)
	}
	newMap := make(map[gethcommon.Hash]*common.L2Tx)
	for k, v := range allIncludedTransactions(s.ParentRollup(r), s, stopAtHeight) {
		newMap[k] = v
	}
	for _, tx := range r.Transactions {
		newMap[tx.Hash()] = tx
	}
	return newMap
}

func removeExisting(base []*common.L2Tx, toRemove map[gethcommon.Hash]*common.L2Tx) (r []*common.L2Tx) {
	for _, t := range base {
		_, f := toRemove[t.Hash()]
		if !f {
			r = append(r, t)
		}
	}
	return
}
