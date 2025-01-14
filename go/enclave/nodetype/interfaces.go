package nodetype

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/enclave/components"
	"github.com/ten-protocol/go-ten/go/enclave/core"
)

// NodeType - the interface for any service type running in Obscuro nodes.
// Should only contain the shared functionality that every service type needs to have.
type NodeType interface {
	// OnL1Fork - logic to be performed when there is an L1 Fork
	OnL1Fork(ctx context.Context, fork *common.ChainFork) error

	// OnL1Block - performed after the block was processed
	OnL1Block(ctx context.Context, block *types.Header, result *components.BlockIngestionType) error

	Close() error
}

type ActiveSequencer interface {
	// CreateBatch - creates a new head batch for the latest known L1 head block.
	CreateBatch(ctx context.Context, skipBatchIfEmpty bool) error

	// CreateRollup - creates a new rollup from the latest recorded rollup in the head l1 chain
	// and adds as many batches to it as possible.
	CreateRollup(ctx context.Context, lastBatchNo uint64) (*common.ExtRollup, *common.ExtRollupMetadata, error)

	NodeType
}

type Validator interface {
	// ExecuteStoredBatches - try to execute all stored by unexecuted batches
	ExecuteStoredBatches(context.Context) error

	VerifySequencerSignature(*core.Batch) error

	NodeType
}
