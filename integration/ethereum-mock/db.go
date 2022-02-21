package ethereum_mock

import (
	common2 "github.com/otherview/obscuro-playground/go/common"
	"sync"
)

// Received blocks ar stored here
type blockResolverInMem struct {
	blockCache map[common2.L1RootHash]*common2.Block
	m          sync.RWMutex
}

func NewResolver() common2.BlockResolver {
	return &blockResolverInMem{
		blockCache: map[common2.L1RootHash]*common2.Block{},
		m:          sync.RWMutex{},
	}
}

func (n *blockResolverInMem) Store(node *common2.Block) {
	n.m.Lock()
	n.blockCache[node.Hash()] = node
	n.m.Unlock()
}

func (n *blockResolverInMem) Resolve(hash common2.L1RootHash) (*common2.Block, bool) {
	n.m.RLock()
	defer n.m.RUnlock()
	block, f := n.blockCache[hash]
	return block, f
}

// The cache of included transactions
type txDbInMem struct {
	transactionsPerBlockCache map[common2.L1RootHash]map[common2.TxHash]*common2.L1Tx
	rpbcM                     *sync.RWMutex
}

func NewTxDb() TxDb {
	return &txDbInMem{
		transactionsPerBlockCache: make(map[common2.L1RootHash]map[common2.TxHash]*common2.L1Tx),
		rpbcM:                     &sync.RWMutex{},
	}
}

func (n *txDbInMem) Txs(b *common2.Block) (map[common2.TxHash]*common2.L1Tx, bool) {
	n.rpbcM.RLock()
	val, found := n.transactionsPerBlockCache[b.Hash()]
	n.rpbcM.RUnlock()
	return val, found
}

func (n *txDbInMem) AddTxs(b *common2.Block, newMap map[common2.TxHash]*common2.L1Tx) {
	n.rpbcM.Lock()
	n.transactionsPerBlockCache[b.Hash()] = newMap
	n.rpbcM.Unlock()
}

// removeCommittedTransactions returns a copy of `mempool` where all transactions that are exactly `committedBlocks`
// deep have been removed.
func removeCommittedTransactions(cb *common2.Block, mempool []*common2.L1Tx, r common2.BlockResolver, db TxDb) []*common2.L1Tx {
	if cb.Height(r) <= common2.HeightCommittedBlocks {
		return mempool
	}

	b := cb
	i := 0
	for {
		if i == common2.HeightCommittedBlocks {
			break
		}
		p, f := b.Parent(r)
		if !f {
			panic("wtf")
		}
		b = p
		i++
	}
	val, _ := db.Txs(b)
	//if !found {
	//	panic("should not fail here")
	//}
	return removeExisting(mempool, val)
}
