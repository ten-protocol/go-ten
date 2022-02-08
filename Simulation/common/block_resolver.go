package common

// BlockResolver -database ob blocks indexed by the root hash
type BlockResolver interface {
	Resolve(hash RootHash) (Block, bool)
	Store(node Block)
}

func (b Block) Parent(r BlockResolver) (Block, bool) {
	return r.Resolve(b.ParentHash)
}

// IsAncestor return true if a is the ancestor of b
func IsAncestor(a Block, b Block, r BlockResolver) bool {
	if a.RootHash == b.RootHash {
		return true
	}
	if a.Height >= b.Height {
		return false
	}
	p, f := b.Parent(r)
	if !f {
		return false
	}
	return IsAncestor(a, p, r)
}

// IsBlockAncestor - takes into conssideration that the block to verify might be on a branch we haven't recevied yet
func IsBlockAncestor(a RootHash, b Block, r BlockResolver) bool {
	if a == b.RootHash {
		return true
	}
	if a == GenesisBlock.RootHash {
		return true
	}
	if b.Height == 0 {
		return false
	}
	block, found := r.Resolve(a)
	if found {
		if block.Height >= b.Height {
			return false
		}
	}
	p, f := b.Parent(r)
	if !f {
		return false
	}
	return IsBlockAncestor(a, p, r)
}
