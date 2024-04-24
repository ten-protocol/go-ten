package nodetype

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/enclave/components"
	"github.com/ten-protocol/go-ten/go/enclave/core"
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

	// OnL1Block - performed after the block was processed
	OnL1Block(block types.Block, result *components.BlockIngestionType) error

	Close() error
}

type Sequencer interface {
	// CreateBatch - creates a new head batch for the latest known L1 head block.
	CreateBatch(skipBatchIfEmpty bool) error

	// CreateRollup - creates a new rollup from the latest recorded rollup in the head l1 chain
	// and adds as many batches to it as possible.
	CreateRollup(lastBatchNo uint64) (*common.ExtRollup, error)

	ExportCrossChainData(uint64, uint64) (*common.ExtCrossChainBundle, error)

	NodeType
}

type ObsValidator interface {
	// ExecuteStoredBatches - try to execute all stored by unexecuted batches
	ExecuteStoredBatches() error

	VerifySequencerSignature(*core.Batch) error

	NodeType
}
