package gethutil

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
)

// Utilities for working with geth structures

// LCA - returns the least common ancestor of the 2 blocks or an error if no common ancestor is found
func LCA(blockA *types.Block, blockB *types.Block, resolver db.BlockResolver) (*types.Block, error) {
	if blockA.NumberU64() == common.L1GenesisHeight || blockB.NumberU64() == common.L1GenesisHeight {
		return blockA, nil
	}
	if bytes.Equal(blockA.Hash().Bytes(), blockB.Hash().Bytes()) {
		return blockA, nil
	}
	if blockA.NumberU64() > blockB.NumberU64() {
		p, err := resolver.ParentBlock(blockA)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve parent block. Cause: %w", err)
		}
		return LCA(p, blockB, resolver)
	}
	if blockB.NumberU64() > blockA.NumberU64() {
		p, err := resolver.ParentBlock(blockB)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve parent block. Cause: %w", err)
		}

		return LCA(blockA, p, resolver)
	}
	parentBlockA, err := resolver.ParentBlock(blockA)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve parent block. Cause: %w", err)
	}
	parentBlockB, err := resolver.ParentBlock(blockB)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve parent block. Cause: %w", err)
	}

	return LCA(parentBlockA, parentBlockB, resolver)
}
