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
func IsAncestor(a *Block, b *Block, r BlockResolver) bool {
	if a.Hash() == b.Hash() {
		return true
	}
	if a.Height(r) >= b.Height(r) {
		return false
	}
	p, f := b.Parent(r)
	if !f {
		return false
	}
	return IsAncestor(a, p, r)
}

// IsBlockAncestor - takes into consideration that the block to verify might be on a branch we haven't received yet
func IsBlockAncestor(a L1RootHash, b *Block, r BlockResolver) bool {
	if a == b.Hash() {
		return true
	}
	if a == GenesisBlock.Hash() {
		return true
	}
	if b.Height(r) == 0 {
		return false
	}
	block, found := r.Resolve(a)
	if found {
		if block.Height(r) >= b.Height(r) {
			return false
		}
	}
	p, f := b.Parent(r)
	if !f {
		return false
	}
	return IsBlockAncestor(a, p, r)
}
