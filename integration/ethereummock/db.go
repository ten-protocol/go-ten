package ethereummock

import (
	"sync"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

type blockAndHeight struct {
	b      *types.Block
	height int
}

// Received blocks ar stored here
type blockResolverInMem struct {
	blockCache map[obscurocommon.L1RootHash]blockAndHeight
	m          sync.RWMutex
}

func NewResolver() enclave.BlockResolver {
	return &blockResolverInMem{
		blockCache: map[obscurocommon.L1RootHash]blockAndHeight{},
		m:          sync.RWMutex{},
	}
}

func (n *blockResolverInMem) StoreBlock(block *types.Block) {
	n.m.Lock()
	defer n.m.Unlock()
	if block.ParentHash() == obscurocommon.GenesisHash {
		n.blockCache[block.Hash()] = blockAndHeight{block, 0}
		return
	}

	p, f := n.blockCache[block.ParentHash()]
	if !f {
		panic("Parent not found. Should not happen")
	}
	n.blockCache[block.Hash()] = blockAndHeight{block, p.height + 1}
}

func (n *blockResolverInMem) ResolveBlock(hash obscurocommon.L1RootHash) (*types.Block, bool) {
	n.m.RLock()
	defer n.m.RUnlock()
	block, f := n.blockCache[hash]

	return block.b, f
}

func (n *blockResolverInMem) HeightBlock(block *types.Block) int {
	n.m.RLock()
	defer n.m.RUnlock()
	return n.blockCache[block.Hash()].height
}

// The cache of included transactions
type txDBInMem struct {
	transactionsPerBlockCache map[obscurocommon.L1RootHash]map[obscurocommon.TxHash]*obscurocommon.L1Tx
	rpbcM                     *sync.RWMutex
}

func NewTxDB() TxDB {
	return &txDBInMem{
		transactionsPerBlockCache: make(map[obscurocommon.L1RootHash]map[obscurocommon.TxHash]*obscurocommon.L1Tx),
		rpbcM:                     &sync.RWMutex{},
	}
}

func (n *txDBInMem) Txs(b *types.Block) (map[obscurocommon.TxHash]*obscurocommon.L1Tx, bool) {
	n.rpbcM.RLock()
	val, found := n.transactionsPerBlockCache[b.Hash()]
	n.rpbcM.RUnlock()

	return val, found
}

func (n *txDBInMem) AddTxs(b *types.Block, newMap map[obscurocommon.TxHash]*obscurocommon.L1Tx) {
	n.rpbcM.Lock()
	n.transactionsPerBlockCache[b.Hash()] = newMap
	n.rpbcM.Unlock()
}

// removeCommittedTransactions returns a copy of `mempool` where all transactions that are exactly `committedBlocks`
// deep have been removed.
func removeCommittedTransactions(
	cb *types.Block,
	mempool []*obscurocommon.L1Tx,
	resolver enclave.BlockResolver,
	db TxDB,
) []*obscurocommon.L1Tx {
	if resolver.HeightBlock(cb) <= obscurocommon.HeightCommittedBlocks {
		return mempool
	}

	b := cb
	i := 0

	for {
		if i == obscurocommon.HeightCommittedBlocks {
			break
		}

		p, f := enclave.Parent(resolver, b)
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
