package host

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
)

// service names - these are the keys used to register known services with the host
const (
	L1BlockRepositoryName = "l1-block-repo"
	L1PublisherName       = "l1-publisher"
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

// L1BlockRepository provides an interface for the host to request L1 block data (live-streaming and historical)
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

// L1Publisher provides an interface for the host to interact with Obscuro data (management contract etc.) on L1
type L1Publisher interface {
	// InitializeSecret will send a management contract transaction to initialize the network with the generated secret
	InitializeSecret(attestation *common.AttestationReport, encSecret common.EncryptedSharedEnclaveSecret) error
	// RequestSecret will send a management contract transaction to request a secret from the enclave, returning the L1 head at time of sending
	RequestSecret(report *common.AttestationReport) (gethcommon.Hash, error)
	// ExtractSecretResponses will return all secret response tx from an L1 block
	ExtractSecretResponses(block *types.Block) []*ethadapter.L1RespondSecretTx
	// PublishRollup will create and publish a rollup tx to the management contract - fire and forget we don't wait for receipt
	// todo (#1624) - With a single sequencer, it is problematic if rollup publication fails; handle this case better
	PublishRollup(producedRollup *common.ExtRollup)
	// PublishSecretResponse will create and publish a secret response tx to the management contract - fire and forget we don't wait for receipt
	PublishSecretResponse(secretResponse *common.ProducedSecretResponse) error

	FetchLatestPeersList() ([]string, error)
}
