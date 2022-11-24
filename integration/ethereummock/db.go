package ethereummock

import (
	"bytes"
	"sync"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/obscuronet/go-obscuro/go/common/errutil"

	"github.com/obscuronet/go-obscuro/go/common"

	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/db"

	"github.com/ethereum/go-ethereum/core/types"
)

// Received blocks ar stored here
type blockResolverInMem struct {
	blockCache map[common.L1RootHash]*types.Block
	m          sync.RWMutex
}

func (n *blockResolverInMem) ProofHeight(_ *core.Rollup) int64 {
	panic("implement me")
}

func (n *blockResolverInMem) Proof(_ *core.Rollup) (*types.Block, error) {
	panic("implement me")
}

func NewResolver() db.BlockResolver {
	return &blockResolverInMem{
		blockCache: map[common.L1RootHash]*types.Block{},
		m:          sync.RWMutex{},
	}
}

func (n *blockResolverInMem) StoreBlock(block *types.Block) {
	n.m.Lock()
	defer n.m.Unlock()
	n.blockCache[block.Hash()] = block
}

func (n *blockResolverInMem) FetchBlock(hash common.L1RootHash) (*types.Block, error) {
	n.m.RLock()
	defer n.m.RUnlock()
	block, f := n.blockCache[hash]

	if !f {
		return nil, errutil.ErrNotFound
	}
	return block, nil
}

func (n *blockResolverInMem) FetchHeadBlock() (*types.Block, error) {
	n.m.RLock()
	defer n.m.RUnlock()
	var max *types.Block
	for k := range n.blockCache {
		bh := n.blockCache[k]
		if max == nil || max.NumberU64() < bh.NumberU64() {
			max = bh
		}
	}
	if max == nil {
		return nil, errutil.ErrNotFound
	}
	return max, nil
}

func (n *blockResolverInMem) ParentBlock(b *types.Block) (*types.Block, error) {
	return n.FetchBlock(b.Header().ParentHash)
}

func (n *blockResolverInMem) IsAncestor(block *types.Block, maybeAncestor *types.Block) bool {
	if bytes.Equal(maybeAncestor.Hash().Bytes(), block.Hash().Bytes()) {
		return true
	}

	if maybeAncestor.NumberU64() >= block.NumberU64() {
		return false
	}

	p, err := n.ParentBlock(block)
	if err != nil {
		return false
	}

	return n.IsAncestor(p, maybeAncestor)
}

func (n *blockResolverInMem) IsBlockAncestor(block *types.Block, maybeAncestor common.L1RootHash) bool {
	if bytes.Equal(maybeAncestor.Bytes(), block.Hash().Bytes()) {
		return true
	}

	if bytes.Equal(maybeAncestor.Bytes(), common.GenesisBlock.Hash().Bytes()) {
		return true
	}

	if block.NumberU64() == common.L1GenesisHeight {
		return false
	}

	resolvedBlock, err := n.FetchBlock(maybeAncestor)
	if err == nil {
		if resolvedBlock.NumberU64() >= block.NumberU64() {
			return false
		}
	}

	p, err := n.ParentBlock(block)
	if err != nil {
		// TODO - If error is not `errutil.ErrNotFound`, throw.
		return false
	}

	return n.IsBlockAncestor(p, maybeAncestor)
}

// The cache of included transactions
type txDBInMem struct {
	transactionsPerBlockCache map[common.L1RootHash]map[common.TxHash]*types.Transaction
	rpbcM                     *sync.RWMutex
}

func NewTxDB() TxDB {
	return &txDBInMem{
		transactionsPerBlockCache: make(map[common.L1RootHash]map[common.TxHash]*types.Transaction),
		rpbcM:                     &sync.RWMutex{},
	}
}

func (n *txDBInMem) Txs(b *types.Block) (map[common.TxHash]*types.Transaction, bool) {
	n.rpbcM.RLock()
	val, found := n.transactionsPerBlockCache[b.Hash()]
	n.rpbcM.RUnlock()

	return val, found
}

func (n *txDBInMem) AddTxs(b *types.Block, newMap map[common.TxHash]*types.Transaction) {
	n.rpbcM.Lock()
	n.transactionsPerBlockCache[b.Hash()] = newMap
	n.rpbcM.Unlock()
}

// removeCommittedTransactions returns a copy of `mempool` where all transactions that are exactly `committedBlocks`
// deep have been removed.
func (m *Node) removeCommittedTransactions(
	cb *types.Block,
	mempool []*types.Transaction,
	resolver db.BlockResolver,
	db TxDB,
) []*types.Transaction {
	if cb.NumberU64() <= common.HeightCommittedBlocks {
		return mempool
	}

	b := cb
	i := 0

	for {
		if i == common.HeightCommittedBlocks {
			break
		}

		p, err := resolver.ParentBlock(b)
		if err != nil {
			m.logger.Crit("Could not retrieve parent block.", log.ErrKey, err)
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
