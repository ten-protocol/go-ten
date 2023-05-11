package mempool

import (
	"fmt"
	"sort"
	"sync"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/db"

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

type mempoolManager struct {
	mpMutex        sync.RWMutex // Controls access to `mempool`
	obscuroChainID int64
	logger         gethlog.Logger
	mempool        map[gethcommon.Hash]*common.L2Tx
	pendingAccTxs  map[gethcommon.Address]common.L2Transactions
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

	_, _, err := core.ExtractOrderingInfo(db.obscuroChainID, tx)
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
func (db *mempoolManager) CurrentTxs(head *core.Batch, resolver db.Storage) ([]*common.L2Tx, error) {
	stateDB, err := resolver.CreateStateDB(*head.Hash())
	if err != nil {
		return nil, err
	}

	txes := db.FetchMempoolTxs()
	sort.Sort(sortByNonce(txes))

	applicableTransactions := make(common.L2Transactions, 0)

	addressNonces := make(map[gethcommon.Address]uint64)
	NonceFor := func(address gethcommon.Address) uint64 {
		nonce, found := addressNonces[address]
		if !found {
			nonce = stateDB.GetNonce(address)
		}
		addressNonces[address] = nonce
		return nonce
	}

	IncrementFor := func(address gethcommon.Address) {
		nonce := addressNonces[address]
		addressNonces[address] = nonce + 1
	}

	for _, tx := range txes {
		sender, txNonce, _ := core.ExtractOrderingInfo(db.obscuroChainID, tx)
		addressNonce := NonceFor(*sender)
		if txNonce == addressNonce {
			applicableTransactions = append(applicableTransactions, tx)
			IncrementFor(*sender)
			db.logger.Info(fmt.Sprintf("Including transaction %s with nonce: %d", tx.Hash().Hex(), tx.Nonce()))
		}
	}

	return applicableTransactions, nil
}
