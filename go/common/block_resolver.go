package common

// BlockResolver -database of blocks indexed by the root hash
type BlockResolver interface {
	Resolve(hash L1RootHash) (*Block, bool)
	Store(node *Block)
}

func (b Block) Parent(r BlockResolver) (*Block, bool) {
	return r.Resolve(b.Header.ParentHash)
}

// IsAncestor return true if a is the ancestor of b
func IsAncestor(blockA *Block, blockB *Block, r BlockResolver) bool {
	if blockA.Hash() == blockB.Hash() {
		return true
	}

	if blockA.Height(r) >= blockB.Height(r) {
		return false
	}

	p, f := blockB.Parent(r)
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

	if block.Height(resolver) == 0 {
		return false
	}

	resolvedBlock, found := resolver.Resolve(l1BlockHash)
	if found {
		if resolvedBlock.Height(resolver) >= block.Height(resolver) {
			return false
		}
	}

	p, f := block.Parent(resolver)
	if !f {
		return false
	}

	return IsBlockAncestor(l1BlockHash, p, resolver)
}
