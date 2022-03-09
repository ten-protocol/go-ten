package common

// BlockResolver -database of blocks indexed by the root hash
type BlockResolver interface {
	Resolve(hash L1RootHash) (*Block, bool)
	Store(block *Block)
	Height(block *Block) int
	Parent(block *Block) (*Block, bool)
}

func Parent(r BlockResolver, b *Block) (*Block, bool) {
	return r.Resolve(b.Header().ParentHash)
}

// IsAncestor return true if a is the ancestor of b
func IsAncestor(blockA *Block, blockB *Block, r BlockResolver) bool {
	if blockA.Hash() == blockB.Hash() {
		return true
	}

	if r.Height(blockA) >= r.Height(blockB) {
		return false
	}

	p, f := r.Parent(blockB)
	if !f {
		return false
	}

	return IsAncestor(blockA, p, r)
}

// IsBlockAncestor - takes into consideration that the block to verify might be on a branch we haven't received yet
func IsBlockAncestor(l1BlockHash L1RootHash, block *Block, resolver BlockResolver) bool {
	if l1BlockHash == block.Hash() {
		return true
	}

	if l1BlockHash == GenesisBlock.Hash() {
		return true
	}

	if resolver.Height(block) == 0 {
		return false
	}

	resolvedBlock, found := resolver.Resolve(l1BlockHash)
	if found {
		if resolver.Height(resolvedBlock) >= resolver.Height(block) {
			return false
		}
	}

	p, f := resolver.Parent(block)
	if !f {
		return false
	}

	return IsBlockAncestor(l1BlockHash, p, resolver)
}
