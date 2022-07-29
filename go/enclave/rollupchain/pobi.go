package rollupchain

import (
	"bytes"

	obscurocore "github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
)

// FindWinner - implements the logic of determining the canonical chain rollup.
func FindWinner(parent *obscurocore.Rollup, rollups []*obscurocore.Rollup, blockResolver db.BlockResolver) (*obscurocore.Rollup, bool) {
	win := -1
	// todo - add statistics to determine why there are conflicts.
	for i, r := range rollups {
		switch {
		case !bytes.Equal(r.Header.ParentHash.Bytes(), parent.Hash().Bytes()): // ignore rollups from L2 forks
		case r.Header.Number.Int64() <= parent.Header.Number.Int64(): // ignore rollups that are older than the parent
		case win == -1:
			win = i
		case blockResolver.ProofHeight(r) < blockResolver.ProofHeight(rollups[win]): // ignore rollups generated with an older proof
		case blockResolver.ProofHeight(r) > blockResolver.ProofHeight(rollups[win]): // newer rollups win
			win = i
		case r.Header.RollupNonce < rollups[win].Header.RollupNonce: // for rollups with the same proof, the one with the lowest nonce wins
			win = i
		}
	}
	if win == -1 {
		return nil, false
	}
	return rollups[win], true
}
