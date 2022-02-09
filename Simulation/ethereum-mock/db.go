package ethereum_mock

import (
	"simulation/common"
	"sync"
)

// Received blocks ar stored here
type blockResolverInMem struct {
	blockCache map[common.RootHash]common.Block
	m          sync.RWMutex
}

func NewResolver() common.BlockResolver {
	return &blockResolverInMem{
		blockCache: map[common.RootHash]common.Block{},
		m:          sync.RWMutex{},
	}
}

func (n *blockResolverInMem) Store(node common.Block) {
	n.m.Lock()
	n.blockCache[node.RootHash] = node
	n.m.Unlock()
}

func (n *blockResolverInMem) Resolve(hash common.RootHash) (common.Block, bool) {
	n.m.RLock()
	defer n.m.RUnlock()
	block, f := n.blockCache[hash]
	return block, f
}

// The cache of included transactions
type txDbInMem struct {
	transactionsPerBlockCache map[common.RootHash]map[common.TxHash]common.L1Tx
	rpbcM                     *sync.RWMutex
}

func NewTxDb() TxDb {
	return &txDbInMem{
		transactionsPerBlockCache: make(map[common.RootHash]map[common.TxHash]common.L1Tx),
		rpbcM:                     &sync.RWMutex{},
	}
}

func (n *txDbInMem) Txs(b common.Block) (map[common.TxHash]common.L1Tx, bool) {
	n.rpbcM.RLock()
	val, found := n.transactionsPerBlockCache[b.RootHash]
	n.rpbcM.RUnlock()
	return val, found
}

func (n *txDbInMem) AddTxs(b common.Block, newMap map[common.TxHash]common.L1Tx) {
	n.rpbcM.Lock()
	n.transactionsPerBlockCache[b.RootHash] = newMap
	n.rpbcM.Unlock()
}

// removeCommittedTransactions returns a copy of `mempool` where all transactions that are exactly `committedBlocks`
// deep have been removed.
func removeCommittedTransactions(cb common.Block, mempool []common.L1Tx, r common.BlockResolver, db TxDb) []common.L1Tx {
	if cb.Height <= common.HeightCommittedBlocks {
		return mempool
	}

	b := cb
	i := 0
	for {
		if i == common.HeightCommittedBlocks {
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
