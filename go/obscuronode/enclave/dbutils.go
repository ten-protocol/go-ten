package enclave

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

func Parent(r BlockResolver, b *types.Block) (*types.Block, bool) {
	return r.FetchBlock(b.Header().ParentHash)
}

// IsAncestor return true if a is the ancestor of b
func IsAncestor(blockA *types.Block, blockB *types.Block, r BlockResolver) bool {
	if blockA.Hash() == blockB.Hash() {
		return true
	}

	if r.HeightBlock(blockA) >= r.HeightBlock(blockB) {
		return false
	}

	p, f := Parent(r, blockB)
	if !f {
		return false
	}

	return IsAncestor(blockA, p, r)
}

// IsBlockAncestor - takes into consideration that the block to verify might be on a branch we haven't received yet
func IsBlockAncestor(l1BlockHash obscurocommon.L1RootHash, block *types.Block, resolver BlockResolver) bool {
	if l1BlockHash == block.Hash() {
		return true
	}

	if l1BlockHash == obscurocommon.GenesisBlock.Hash() {
		return true
	}

	if resolver.HeightBlock(block) == 0 {
		return false
	}

	resolvedBlock, found := resolver.FetchBlock(l1BlockHash)
	if found {
		if resolver.HeightBlock(resolvedBlock) >= resolver.HeightBlock(block) {
			return false
		}
	}

	p, f := Parent(resolver, block)
	if !f {
		return false
	}

	return IsBlockAncestor(l1BlockHash, p, resolver)
}

func parentRollup(db DB, r *Rollup) *Rollup {
	parent, found := db.FetchRollup(r.Header.ParentHash)
	if !found {
		panic(fmt.Sprintf("Could not find rollup: r_%s", r.Hash()))
	}
	return parent
}

func heightRollup(db DB, r *Rollup) int {
	if height := r.Height.Load(); height != nil {
		return height.(int)
	}
	if r.Hash() == GenesisRollup.Hash() {
		r.Height.Store(obscurocommon.L2GenesisHeight)
		return obscurocommon.L2GenesisHeight
	}
	v := heightRollup(db, parentRollup(db, r)) + 1
	r.Height.Store(v)
	return v
}
