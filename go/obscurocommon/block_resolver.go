package obscurocommon

import "github.com/ethereum/go-ethereum/core/types"

// BlockResolver -database of blocks indexed by the root hash
type BlockResolver interface {
	ResolveBlock(hash L1RootHash) (*types.Block, bool)
	StoreBlock(block *types.Block)
	HeightBlock(block *types.Block) int
	ParentBlock(block *types.Block) (*types.Block, bool)
}

func Parent(r BlockResolver, b *types.Block) (*types.Block, bool) {
	return r.ResolveBlock(b.Header().ParentHash)
}

// IsAncestor return true if a is the ancestor of b
func IsAncestor(blockA *types.Block, blockB *types.Block, r BlockResolver) bool {
	if blockA.Hash() == blockB.Hash() {
		return true
	}

	if r.HeightBlock(blockA) >= r.HeightBlock(blockB) {
		return false
	}

	p, f := r.ParentBlock(blockB)
	if !f {
		return false
	}

	return IsAncestor(blockA, p, r)
}

// IsBlockAncestor - takes into consideration that the block to verify might be on a branch we haven't received yet
func IsBlockAncestor(l1BlockHash L1RootHash, block *types.Block, resolver BlockResolver) bool {
	if l1BlockHash == block.Hash() {
		return true
	}

	if l1BlockHash == GenesisBlock.Hash() {
		return true
	}

	if resolver.HeightBlock(block) == 0 {
		return false
	}

	resolvedBlock, found := resolver.ResolveBlock(l1BlockHash)
	if found {
		if resolver.HeightBlock(resolvedBlock) >= resolver.HeightBlock(block) {
			return false
		}
	}

	p, f := resolver.ParentBlock(block)
	if !f {
		return false
	}

	return IsBlockAncestor(l1BlockHash, p, resolver)
}
