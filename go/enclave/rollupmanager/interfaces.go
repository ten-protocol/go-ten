package rollupmanager

import (
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
)

type RollupManager interface {
	// CreateRollup - creates a rollup encapsulating the state from the
	// latest published head batch to the most current headbatch.
	CreateRollup() (*core.Rollup, error)
	// ProcessL1Block - extracts the rollups from the block's transactions
	// and verifies their integrity, saving and processing any batches that have
	// not been seenp previously.
	ProcessL1Block(b *common.BlockAndReceipts) ([]*core.Rollup, error)
}
