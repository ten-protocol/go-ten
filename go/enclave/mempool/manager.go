package mempool

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/obscuronet/go-obscuro/go/common/errutil"

	gethcommon "github.com/ethereum/go-ethereum/common"

	obscurocore "github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/db"

	"github.com/obscuronet/go-obscuro/go/common"
)

// sortByNonce a very primitive way to implement mempool logic that
// adds transactions sorted by the nonce in the rollup
// which is what the EVM expects
type sortByNonce []*common.L2Tx

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

func (db *mempoolManager) RemoveMempoolTxs(rollup *obscurocore.Rollup, resolver db.RollupResolver) error {
	db.mpMutex.Lock()
	defer db.mpMutex.Unlock()

	toRemove, err := historicTxs(rollup, resolver)
	if err != nil {
		return fmt.Errorf("error retrieiving historic transactions. Cause: %w", err)
	}

	r := make(map[gethcommon.Hash]*common.L2Tx)
	for id, t := range db.mempool {
		_, f := toRemove[id]
		if !f {
			r[id] = t
		}
	}
	db.mempool = r

	return nil
}

// Returns all transactions found 20 levels below
func historicTxs(r *obscurocore.Rollup, resolver db.RollupResolver) (map[gethcommon.Hash]gethcommon.Hash, error) {
	i := common.HeightCommittedBlocks
	c := r
	found := true
	var err error
	// todo - create method to return the canonical rollup from height N
	for {
		if !found || i == 0 || c.Header.Number.Uint64() == common.L2GenesisHeight {
			return obscurocore.ToMap(c.Transactions), nil
		}
		i--
		c, err = resolver.ParentRollup(c)
		if err != nil && !errors.Is(err, errutil.ErrNotFound) {
			return nil, fmt.Errorf("could not retrieve parent rollup. Cause: %w", err)
		}
		found = err != nil
	}
}

// CurrentTxs - Calculate transactions to be included in the current rollup
func (db *mempoolManager) CurrentTxs(head *obscurocore.Rollup, resolver db.RollupResolver) ([]*common.L2Tx, error) {
	txs, err := findTxsNotIncluded(head, db.FetchMempoolTxs(), resolver)
	if err != nil {
		return nil, err
	}
	sort.Sort(sortByNonce(txs))
	return txs, nil
}

// findTxsNotIncluded - given a list of transactions, it keeps only the ones that were not included in the block
// todo - inefficient
func findTxsNotIncluded(head *obscurocore.Rollup, txs []*common.L2Tx, s db.RollupResolver) ([]*common.L2Tx, error) {
	// go back only HeightCommittedBlocks blocks to accumulate transactions to be diffed against the mempool
	startAt := uint64(0)
	if head.NumberU64() > common.HeightCommittedBlocks {
		startAt = head.NumberU64() - common.HeightCommittedBlocks
	}
	included, err := allIncludedTransactions(head, s, startAt)
	if err != nil {
		return nil, err
	}
	return removeExisting(txs, included), nil
}

func allIncludedTransactions(r *obscurocore.Rollup, s db.RollupResolver, stopAtHeight uint64) (map[gethcommon.Hash]*common.L2Tx, error) {
	if r.Header.Number.Uint64() == stopAtHeight {
		return obscurocore.MakeMap(r.Transactions), nil
	}
	newMap := make(map[gethcommon.Hash]*common.L2Tx)
	parent, err := s.ParentRollup(r)
	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return nil, err
	}

	if err == nil {
		txsMap, err := allIncludedTransactions(parent, s, stopAtHeight)
		if err != nil {
			return nil, err
		}
		for k, v := range txsMap {
			newMap[k] = v
		}
		for _, tx := range r.Transactions {
			newMap[tx.Hash()] = tx
		}
	}

	return newMap, nil
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
