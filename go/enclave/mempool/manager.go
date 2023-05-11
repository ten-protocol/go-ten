package mempool

import (
	"errors"
	"sort"
	"sync"

	"github.com/obscuronet/go-obscuro/go/common/errutil"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/go-obscuro/go/enclave/core"
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
	err := core.VerifySignature(db.obscuroChainID, tx)
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

func (db *mempoolManager) RemoveTxs(transactions types.Transactions) error {
	db.mpMutex.Lock()
	defer db.mpMutex.Unlock()

	for _, tx := range transactions {
		delete(db.mempool, tx.Hash())
	}

	return nil
}

// CurrentTxs - Calculate transactions to be included in the current batch
func (db *mempoolManager) CurrentTxs(head *core.Batch, resolver db.BatchResolver) ([]*common.L2Tx, error) {
	txs, err := findTxsNotIncluded(head, db.FetchMempoolTxs(), resolver)
	if err != nil {
		return nil, err
	}
	sort.Sort(sortByNonce(txs))
	return txs, nil
}

// findTxsNotIncluded - given a list of transactions, it keeps only the ones that were not included in the block
// todo (#1491) - inefficient
func findTxsNotIncluded(head *core.Batch, txs []*common.L2Tx, s db.BatchResolver) ([]*common.L2Tx, error) {
	// go back only HeightCommittedBlocks blocks to accumulate transactions to be diffed against the mempool
	startAt := uint64(0)
	if head.NumberU64() > common.HeightCommittedBlocks {
		startAt = head.NumberU64() - common.HeightCommittedBlocks
	}
	included, err := allIncludedTransactions(head, s, startAt)
	if err != nil {
		return nil, err
	}
	return filterOutTransactions(txs, included), nil
}

// Recursively finds all transactions included in the past stopAtHeight batches.
func allIncludedTransactions(batch *core.Batch, s db.BatchResolver, stopAtHeight uint64) (map[gethcommon.Hash]*common.L2Tx, error) {
	if batch.Header.Number.Uint64() == stopAtHeight {
		return core.MakeMap(batch.Transactions), nil
	}

	// We add this batch's transactions to the included transactions.
	newMap := make(map[gethcommon.Hash]*common.L2Tx)
	for _, tx := range batch.Transactions {
		newMap[tx.Hash()] = tx
	}

	// If the batch has a parent (i.e. it is not the genesis block), we recurse.
	parentBatch, err := s.FetchBatch(batch.Header.ParentHash)
	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return nil, err
	}
	if err == nil {
		txsMap, err := allIncludedTransactions(parentBatch, s, stopAtHeight)
		if err != nil {
			return nil, err
		}
		for hash, tx := range txsMap {
			newMap[hash] = tx
		}
	}

	return newMap, nil
}

func filterOutTransactions(txs []*common.L2Tx, txsToRemove map[gethcommon.Hash]*common.L2Tx) (r []*common.L2Tx) {
	for _, tx := range txs {
		_, f := txsToRemove[tx.Hash()]
		if !f {
			r = append(r, tx)
		}
	}
	return
}
