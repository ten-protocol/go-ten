package ethereummock

import (
	"sync"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/core"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/db"

	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

type blockAndHeight struct {
	b      *types.Block
	height uint64
}

// Received blocks ar stored here
type blockResolverInMem struct {
	blockCache map[obscurocommon.L1RootHash]blockAndHeight
	m          sync.RWMutex
}

func (n *blockResolverInMem) ProofHeight(_ *core.Rollup) int64 {
	panic("implement me")
}

func (n *blockResolverInMem) Proof(_ *core.Rollup) *types.Block {
	panic("implement me")
}

func NewResolver() db.BlockResolver {
	return &blockResolverInMem{
		blockCache: map[obscurocommon.L1RootHash]blockAndHeight{},
		m:          sync.RWMutex{},
	}
}

func (n *blockResolverInMem) StoreBlock(block *types.Block) bool {
	n.m.Lock()
	defer n.m.Unlock()
	if block.ParentHash() == obscurocommon.GenesisHash {
		n.blockCache[block.Hash()] = blockAndHeight{block, obscurocommon.L1GenesisHeight}
		return true
	}

	p, f := n.blockCache[block.ParentHash()]
	if !f {
		log.Log("Trying to store block but haven't yet stored its parent. Trying increasing the simulation's block " +
			"time or reducing the number of nodes")
		return false
	}
	n.blockCache[block.Hash()] = blockAndHeight{block, p.height + 1}
	return true
}

func (n *blockResolverInMem) FetchBlock(hash obscurocommon.L1RootHash) (*types.Block, bool) {
	n.m.RLock()
	defer n.m.RUnlock()
	block, f := n.blockCache[hash]

	return block.b, f
}

func (n *blockResolverInMem) FetchHeadBlock() *types.Block {
	n.m.RLock()
	defer n.m.RUnlock()
	var max blockAndHeight
	for _, bh := range n.blockCache {
		if max.height < bh.height {
			max = bh
		}
	}
	return max.b
}

func (n *blockResolverInMem) HeightBlock(block *types.Block) uint64 {
	n.m.RLock()
	defer n.m.RUnlock()
	b, f := n.blockCache[block.Hash()]
	if f {
		return b.height
	}
	panic("block not stored")
}

func (n *blockResolverInMem) ParentBlock(b *types.Block) (*types.Block, bool) {
	return n.FetchBlock(b.Header().ParentHash)
}

func (n *blockResolverInMem) IsAncestor(block *types.Block, maybeAncestor *types.Block) bool {
	if maybeAncestor.Hash() == block.Hash() {
		return true
	}

	if n.HeightBlock(maybeAncestor) >= n.HeightBlock(block) {
		return false
	}

	p, f := n.ParentBlock(block)
	if !f {
		return false
	}

	return n.IsAncestor(p, maybeAncestor)
}

func (n *blockResolverInMem) IsBlockAncestor(block *types.Block, maybeAncestor obscurocommon.L1RootHash) bool {
	if maybeAncestor == block.Hash() {
		return true
	}

	if maybeAncestor == obscurocommon.GenesisBlock.Hash() {
		return true
	}

	if n.HeightBlock(block) == obscurocommon.L1GenesisHeight {
		return false
	}

	resolvedBlock, found := n.FetchBlock(maybeAncestor)
	if found {
		if n.HeightBlock(resolvedBlock) >= n.HeightBlock(block) {
			return false
		}
	}

	p, f := n.ParentBlock(block)
	if !f {
		return false
	}

	return n.IsBlockAncestor(p, maybeAncestor)
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
	resolver db.BlockResolver,
	db TxDB,
) []*obscurocommon.L1Tx {
	if cb.NumberU64() <= obscurocommon.HeightCommittedBlocks {
		return mempool
	}

	b := cb
	i := 0

	for {
		if i == obscurocommon.HeightCommittedBlocks {
			break
		}

		p, f := resolver.ParentBlock(b)
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
