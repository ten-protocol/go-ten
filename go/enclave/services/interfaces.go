package services

import (
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/enclave/components"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
)

// ObscuroService - the interface for any service type running in Obscuro nodes.
// Should only contain the shared functionality that every service type needs to have.
type ObscuroService interface {
	//	ReceiveBlock - function that accepts L1 blocks and their receipts along with a flag
	// that signals if the block is the latest one or not. Processing of those blocks and
	// resulting actions differ between the service types.
	ReceiveBlock(*common.BlockAndReceipts, bool) (*components.BlockIngestionType, error)

	// SubmitTransaction - L2 obscuro transactions need to be passed here. Sequencers
	// will put them in the mempool while validators might put them in a queue and monitor
	// for censorship.
	SubmitTransaction(*common.L2Tx) error
}

type Sequencer interface {
	// CreateBatch - creates a new head batch for the parameter block or if there is no block
	// provided, for the latest known L1 head block.
	CreateBatch(*common.L1Block) (*core.Batch, error)

	// CreateRollup - creates a new rollup from the latest recorded rollup in the head l1 chain
	// and adds as many batches to it as possible.
	CreateRollup() (*common.ExtRollup, error)

	ObscuroService
}

type ObsValidator interface {
	// ValidateAndStoreBatch - if all the prerequisites are available (parent batch, l1 block) then
	// this function recomputes the batch using the exact same context and compares the results.
	// If the batch is valid it will be stored.
	ValidateAndStoreBatch(*core.Batch) error

	ObscuroService
}
