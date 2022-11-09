package gethutil

import (
	"bytes"
	"errors"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
)

// Utilities for working with geth structures

var errNoCommonAncestor = errors.New("no common ancestor found")

// LCA - returns the least common ancestor of the 2 blocks or an error if no common ancestor is found
func LCA(blockA *types.Block, blockB *types.Block, resolver db.BlockResolver) (*types.Block, error) {
	if blockA.NumberU64() == common.L1GenesisHeight || blockB.NumberU64() == common.L1GenesisHeight {
		return blockA, nil
	}
	if bytes.Equal(blockA.Hash().Bytes(), blockB.Hash().Bytes()) {
		return blockA, nil
	}
	if blockA.NumberU64() > blockB.NumberU64() {
		p, f := resolver.ParentBlock(blockA)
		if !f {
			return nil, errNoCommonAncestor
		}
		return LCA(p, blockB, resolver)
	}
	if blockB.NumberU64() > blockA.NumberU64() {
		p, f := resolver.ParentBlock(blockB)
		if !f {
			return nil, errNoCommonAncestor
		}

		return LCA(blockA, p, resolver)
	}
	parentBlockA, f := resolver.ParentBlock(blockA)
	if !f {
		return nil, errNoCommonAncestor
	}
	parentBlockB, f := resolver.ParentBlock(blockB)
	if !f {
		return nil, errNoCommonAncestor
	}

	return LCA(parentBlockA, parentBlockB, resolver)
}
