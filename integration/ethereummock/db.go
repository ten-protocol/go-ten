package ethereummock

import (
	"sync"

	common2 "github.com/obscuronet/obscuro-playground/go/common"
)

type blockAndHeight struct {
	b      *common2.Block
	height int
}

// Received blocks ar stored here
type blockResolverInMem struct {
	blockCache map[common2.L1RootHash]blockAndHeight
	m          sync.RWMutex
}

func NewResolver() common2.BlockResolver {
	return &blockResolverInMem{
		blockCache: map[common2.L1RootHash]blockAndHeight{},
		m:          sync.RWMutex{},
	}
}

func (n *blockResolverInMem) Store(block *common2.Block) {
	n.m.Lock()
	p, f := n.blockCache[block.ParentHash()]
	if !f {
		panic("Parent not found. Should not happen")
	}
	n.blockCache[block.Hash()] = blockAndHeight{block, p.height + 1}
	n.m.Unlock()
}

func (n *blockResolverInMem) Resolve(hash common2.L1RootHash) (*common2.Block, bool) {
	n.m.RLock()
	defer n.m.RUnlock()
	block, f := n.blockCache[hash]

	return block.b, f
}

func (n *blockResolverInMem) Height(block *common2.Block) int {
	n.m.RLock()
	defer n.m.RUnlock()
	return n.blockCache[block.Hash()].height
}

func (n *blockResolverInMem) Parent(block *common2.Block) (*common2.Block, bool) {
	return common2.Parent(n, block)
}

// The cache of included transactions
type txDBInMem struct {
	transactionsPerBlockCache map[common2.L1RootHash]map[common2.TxHash]*common2.L1Tx
	rpbcM                     *sync.RWMutex
}

func NewTxDB() TxDB {
	return &txDBInMem{
		transactionsPerBlockCache: make(map[common2.L1RootHash]map[common2.TxHash]*common2.L1Tx),
		rpbcM:                     &sync.RWMutex{},
	}
}

func (n *txDBInMem) Txs(b *common2.Block) (map[common2.TxHash]*common2.L1Tx, bool) {
	n.rpbcM.RLock()
	val, found := n.transactionsPerBlockCache[b.Hash()]
	n.rpbcM.RUnlock()

	return val, found
}

func (n *txDBInMem) AddTxs(b *common2.Block, newMap map[common2.TxHash]*common2.L1Tx) {
	n.rpbcM.Lock()
	n.transactionsPerBlockCache[b.Hash()] = newMap
	n.rpbcM.Unlock()
}

// removeCommittedTransactions returns a copy of `mempool` where all transactions that are exactly `committedBlocks`
// deep have been removed.
func removeCommittedTransactions(
	cb *common2.Block,
	mempool []*common2.L1Tx,
	resolver common2.BlockResolver,
	db TxDB,
) []*common2.L1Tx {
	if resolver.Height(cb) <= common2.HeightCommittedBlocks {
		return mempool
	}

	b := cb
	i := 0

	for {
		if i == common2.HeightCommittedBlocks {
			break
		}

		p, f := resolver.Parent(b)
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
