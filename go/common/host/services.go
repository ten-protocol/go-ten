package host

import (
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
)

// service names - these are the keys used to register known services with the host
const (
	P2PName               = "p2p"
	L1BlockRepositoryName = "l1-block-repo"
	L1PublisherName       = "l1-publisher"
	L2BatchRepositoryName = "l2-batch-repo"
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

// P2P provides an interface for the host to interact with the P2P network
type P2P interface {
	// BroadcastBatches sends live batch(es) to every other node on the network
	BroadcastBatches(batches []*common.ExtBatch) error
	// SendTxToSequencer sends the encrypted transaction to the sequencer.
	SendTxToSequencer(tx common.EncryptedTx) error

	// RequestBatchesFromSequencer asynchronously requests batches from the sequencer, from the given sequence number
	RequestBatchesFromSequencer(fromSeqNo *big.Int) error
	// RespondToBatchRequest sends the requested batches to the requesting peer
	RespondToBatchRequest(requestID string, batches []*common.ExtBatch) error

	// SubscribeForBatches will register a handler to receive new batches from peers, returns unsubscribe func
	SubscribeForBatches(handler P2PBatchHandler) func()
	// SubscribeForTx will register a handler to receive new transactions from peers, returns unsubscribe func
	SubscribeForTx(handler P2PTxHandler) func()
	// SubscribeForBatchRequests will register a handler to receive new batch requests from peers, returns unsubscribe func
	// todo (@matt) feels a bit weird to have this in this interface since it relates to serving data rather than receiving
	SubscribeForBatchRequests(handler P2PBatchRequestHandler) func()

	// UpdatePeerList allows the host to notify the p2p service of a change in the peer list
	UpdatePeerList([]string)
}

// P2PBatchHandler is an interface for receiving new batches from the P2P network as they arrive
type P2PBatchHandler interface {
	// HandleBatches will be called in a new goroutine for batches that arrive
	HandleBatches(batch []*common.ExtBatch, isLive bool)
}

// P2PTxHandler is an interface for receiving new transactions from the P2P network as they arrive
type P2PTxHandler interface {
	// HandleTransaction will be called in a new goroutine for each new tx as it arrives
	HandleTransaction(tx common.EncryptedTx)
}

type P2PBatchRequestHandler interface {
	// HandleBatchRequest will be called in a new goroutine for each new batch request as it arrives
	HandleBatchRequest(requestID string, fromSeqNo *big.Int)
}

// L1BlockRepository provides an interface for the host to request L1 block data (live-streaming and historical)
type L1BlockRepository interface {
	// Subscribe will register a block handler to receive new blocks as they arrive, returns unsubscribe func
	Subscribe(handler L1BlockHandler) func()

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

	FetchLatestSeqNo() (*big.Int, error)
}

// L2BatchRepository provides an interface for the host to request L2 batch data (live-streaming and historical)
type L2BatchRepository interface {
	// Subscribe will register a batch handler to receive new batches as they arrive
	Subscribe(handler L2BatchHandler)

	FetchBatchBySeqNo(seqNo *big.Int) (*common.ExtBatch, error)

	// AddBatch is used to notify the repository of a new batch, e.g. from the enclave when seq produces one or a rollup is consumed
	// Note: it is fine to add batches that the repo already has, it will just ignore them
	AddBatch(batch *common.ExtBatch) error
}

// L2BatchHandler is an interface for receiving new batches from the publisher as they arrive
type L2BatchHandler interface {
	// HandleBatch will be called in a new goroutine for each new batch as it arrives
	HandleBatch(batch *common.ExtBatch)
}
