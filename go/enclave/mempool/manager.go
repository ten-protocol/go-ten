package mempool

import (
	"fmt"
	"sort"
	"sync"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/go-obscuro/go/enclave/core"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common"
)

// sortByNonce a very primitive way to implement mempool logic that
// adds transactions sorted by the nonce in the rollup
// which is what the EVM expects
type sortByNonce []*common.L2Tx

func (c sortByNonce) Len() int           { return len(c) }
func (c sortByNonce) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c sortByNonce) Less(i, j int) bool { return c[i].Nonce() < c[j].Nonce() }

// todo - optimize this to use a different data structure that does not require a global lock.
type mempoolManager struct {
	mpMutex        sync.RWMutex // Controls access to `mempool`
	obscuroChainID int64
	logger         gethlog.Logger
	mempool        map[gethcommon.Hash]*common.L2Tx
}

func New(chainID int64, logger gethlog.Logger) Manager {
	return &mempoolManager{
		mempool:        make(map[gethcommon.Hash]*common.L2Tx),
		obscuroChainID: chainID,
		mpMutex:        sync.RWMutex{},
		logger:         logger,
	}
}

func (db *mempoolManager) AddMempoolTx(tx *common.L2Tx) error {
	db.mpMutex.Lock()
	defer db.mpMutex.Unlock()

	// We do not care about the sender return value at this point, only that
	// there is no error coming from validating the signature of said sender.
	_, err := core.GetAuthenticatedSender(db.obscuroChainID, tx)
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
func (db *mempoolManager) CurrentTxs(stateDB *state.StateDB) ([]*common.L2Tx, error) {
	txes := db.FetchMempoolTxs()
	sort.Sort(sortByNonce(txes))

	applicableTransactions := make(common.L2Transactions, 0)
	nonceTracker := NewNonceTracker(stateDB)

	for _, tx := range txes {
		sender, _ := core.GetAuthenticatedSender(db.obscuroChainID, tx)
		if tx.Nonce() == nonceTracker.GetNonce(*sender) {
			applicableTransactions = append(applicableTransactions, tx)
			nonceTracker.IncrementNonce(*sender)
			db.logger.Info(fmt.Sprintf("Including transaction %s with nonce: %d", tx.Hash().Hex(), tx.Nonce()))
		}
	}

	return applicableTransactions, nil
}
