package rollupchain

import (
	"bytes"

	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/db"
)

// Returns the new head rollup, and a boolean indicating whether this is a new rollup or the existing head rollup.
// todo - add statistics to determine why there are conflicts.
func selectNextRollup(currentHeadRollup *core.Rollup, rollups []*core.Rollup, blockResolver db.BlockResolver) (*core.Rollup, bool) { //nolint:unused
	var newHeadRollup *core.Rollup

	// We iterate over the proposed rollups to select the best new head rollup.
	for _, rollup := range rollups {
		// We ignore rollups from L2 forks, or that are not newer than the parent rollup.
		isFromFork := !bytes.Equal(rollup.Header.ParentHash.Bytes(), currentHeadRollup.Hash().Bytes())
		isNotNewerThanParent := rollup.Header.Number.Int64() <= currentHeadRollup.Header.Number.Int64()
		if isFromFork || isNotNewerThanParent {
			continue
		}

		// If this is the first rollup to pass the checks above, or it is produced from a newer L1 block, we make it
		// the candidate new head rollup.
		if newHeadRollup == nil || blockResolver.ProofHeight(rollup) > blockResolver.ProofHeight(newHeadRollup) {
			newHeadRollup = rollup
		}
	}

	if newHeadRollup != nil {
		return newHeadRollup, true
	}
	// We remain with the existing head rollup.
	return currentHeadRollup, false
}
