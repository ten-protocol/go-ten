package enclave

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

func AssertSecretAvailable(db DB) {
	if db.FetchSecret() == nil {
		panic("Enclave not initialized")
	}
}

func ParentRollup(db DB, r *Rollup) *Rollup {
	AssertSecretAvailable(db)
	parent, found := db.FetchRollup(r.Header.ParentHash)
	if !found {
		panic(fmt.Sprintf("Could not find rollup: r_%s", r.Hash()))
	}
	return parent
}

func HeightRollup(db DB, r *Rollup) int {
	AssertSecretAvailable(db)
	if height := r.Height.Load(); height != nil {
		return height.(int)
	}
	if r.Hash() == GenesisRollup.Hash() {
		r.Height.Store(obscurocommon.L2GenesisHeight)
		return obscurocommon.L2GenesisHeight
	}
	v := HeightRollup(db, ParentRollup(db, r)) + 1
	r.Height.Store(v)
	return v
}

func Balance(db DB, address common.Address) uint64 {
	AssertSecretAvailable(db)
	return db.Head().State[address]
}

func ExistRollup(db DB, hash obscurocommon.L2RootHash) bool {
	AssertSecretAvailable(db)
	_, f := db.FetchRollup(hash)
	return f
}

func ParentBlock(db DB, block *types.Block) (*types.Block, bool) {
	return db.ResolveBlock(block.ParentHash())
}
