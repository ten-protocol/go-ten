package host

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
)

const (
	// service names - these are the keys used to register known services with the host
	L1BlockRepositoryName = "l1-block-repo"
)

// The host has a number of services that encapsulate the various responsibilities of the host.
// This file contains service-level abstractions and utilities, as well as all the interfaces for these services, code
// should depend on these interfaces rather than the concrete implementations.

// Service interface allows the host to manage all services in a generic way
type Service interface {
	Start() error
	Stop() error
	HealthStatus() HealthStatus
}

type L1BlockRepository interface {
	// Subscribe will register a block handler to receive new blocks as they arrive
	Subscribe(handler L1BlockHandler)

	FetchBlockByHeight(height int) (*types.Block, error)
	// FetchNextBlock returns the next canonical block after a given block hash
	// It returns the new block and a bool which is true if the block is the current L1 head
	FetchNextBlock(prevBlock gethcommon.Hash) (*types.Block, bool, error)
	// FetchReceipts returns the receipts for a given L1 block
	FetchReceipts(block *common.L1Block) types.Receipts
}

// L1BlockHandler is an interface for receiving new blocks from the repository as they arrive
type L1BlockHandler interface {
	// HandleBlock will be called in a new goroutine for each new block as it arrives
	HandleBlock(block *types.Block)
}
