package nodetype

import (
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
)

// NodeType - the interface for any service type running in Obscuro nodes.
// Should only contain the shared functionality that every service type needs to have.
type NodeType interface {
	// SubmitTransaction - L2 obscuro transactions need to be passed here. Sequencers
	// will put them in the mempool while validators might put them in a queue and monitor
	// for censorship.
	SubmitTransaction(*common.L2Tx) error

	// OnL1Fork - logic to be performed when there is an L1 Fork
	OnL1Fork(fork *common.ChainFork) error
}

type Sequencer interface {
	// CreateBatch - creates a new head batch for the latest known L1 head block.
	CreateBatch() error

	// CreateRollup - creates a new rollup from the latest recorded rollup in the head l1 chain
	// and adds as many batches to it as possible.
	CreateRollup(lastBatchNo uint64) (*common.ExtRollup, error)

	NodeType
}

type ObsValidator interface {
	// ValidateAndStoreBatch - if all the prerequisites are available (parent batch, l1 block) then
	// this function recomputes the batch using the exact same context and compares the results.
	// If the batch is valid it will be stored.
	ValidateAndStoreBatch(*core.Batch) error

	NodeType
}
