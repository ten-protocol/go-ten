package common

import (
	"context"
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto/kzg4844"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ten-protocol/go-ten/go/common/errutil"

	"github.com/ten-protocol/go-ten/go/common/tracers"
	"github.com/ten-protocol/go-ten/go/responses"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

type StatusCode int

// Status represents the enclave's current state - whether the enclave is healthy and ready to process requests, as well
// as its latest known heads for the L1 and L2 chains and enclave ID derived from the public key
type Status struct {
	StatusCode StatusCode
	L1Head     gethcommon.Hash
	L2Head     *big.Int
	EnclaveID  EnclaveID
}

type TxAndReceiptAndBlobs struct {
	Tx      *types.Transaction
	Receipt *types.Receipt
	Blobs   []*kzg4844.Blob
}

const (
	Running        StatusCode = iota // the enclave is running, accepting L1 blocks
	AwaitingSecret                   // the enclave has not received the network secret and cannot process L1 blocks
	Unavailable                      // the enclave is unavailable (no guarantee it will self-recover)
)

// EnclaveInit defines methods for initializing and managing the state of an enclave.
type EnclaveInit interface {
	// Status checks whether the enclave is ready to process requests - only implemented by the RPC layer
	Status(context.Context) (Status, SystemError)

	// Attestation - Produces an attestation report which will be used to request the shared secret from another enclave.
	Attestation(context.Context) (*AttestationReport, SystemError)

	// GenerateSecret - the genesis enclave is responsible with generating the secret entropy
	GenerateSecret(context.Context) (EncryptedSharedEnclaveSecret, SystemError)

	// InitEnclave - initialise an enclave with a seed received by another enclave
	InitEnclave(ctx context.Context, secret EncryptedSharedEnclaveSecret) SystemError

	// EnclaveID - returns the enclave's ID
	EnclaveID(context.Context) (EnclaveID, SystemError)

	// RPCEncryptionKey - returns the key used
	RPCEncryptionKey(context.Context) ([]byte, SystemError)
}

// EnclaveAdmin provides administrative functions for managing an enclave.
type EnclaveAdmin interface {
	// MakeActive - backup sequencer enclave can become active at the command of the host
	MakeActive() SystemError

	// SubmitL1Block - Used for the host to submit L1 pre-processed blocks to the enclave, these may be:
	//  a. historic block - if the enclave is behind and in the process of catching up with the L1 state
	//  b. the latest block published by the L1, to which the enclave should respond with a rollup
	// It is the responsibility of the host to gossip the returned rollup
	// For good functioning the caller should always submit blocks ordered by height
	// submitting a block before receiving ancestors of it, will result in it being ignored
	SubmitL1Block(ctx context.Context, processed *ProcessedL1Data) (*BlockSubmissionResponse, SystemError)

	// SubmitBatch submits a batch received from the sequencer for processing.
	SubmitBatch(ctx context.Context, batch *ExtBatch) SystemError

	// HealthCheck returns whether the enclave is in a healthy state
	HealthCheck(context.Context) (bool, SystemError)

	// GetBatch - retrieve a batch if existing within the enclave db.
	GetBatch(ctx context.Context, hash L2BatchHash) (*ExtBatch, SystemError)

	// GetBatchBySeqNo - retrieve batch by sequencer number if it's in the db.
	GetBatchBySeqNo(ctx context.Context, seqNo uint64) (*ExtBatch, SystemError)

	// GetRollupData - retrieve the first batch sequence and start time for a given rollup.
	GetRollupData(ctx context.Context, hash L2RollupHash) (*PublicRollupMetadata, SystemError)

	// CreateBatch - creates a new head batch extending the previous one for the latest known L1 head if the node is
	// a sequencer. Will panic otherwise.
	CreateBatch(ctx context.Context, skipIfEmpty bool) SystemError

	// CreateRollup - will create a new rollup by going through the sequencer if the node is a sequencer
	// or panic otherwise.
	CreateRollup(ctx context.Context, fromSeqNo uint64) (*ExtRollup, []*kzg4844.Blob, SystemError)

	// StreamL2Updates - will stream any new batches as they are created/detected
	// All will be queued in the channel that has been returned.
	StreamL2Updates() (chan StreamL2UpdatesResponse, func())

	ExportCrossChainData(context.Context, uint64, uint64) (*ExtCrossChainBundle, SystemError)

	// Stop gracefully stops the enclave
	Stop() SystemError
	// StopClient stops the enclave client if one exists - only implemented by the RPC layer
	StopClient() SystemError
}

// EnclaveClientRPC - exposes RPC methods called by users
type EnclaveClientRPC interface {
	EncryptedRPC(ctx context.Context, encryptedReq EncryptedRequest) (*responses.EnclaveResponse, SystemError)

	// Subscribe adds a log subscription to the enclave under the given ID, provided the request is authenticated
	// correctly. The events will be populated in the BlockSubmissionResponse. If there is an existing subscription
	// with the given ID, it is overwritten.
	Subscribe(ctx context.Context, id rpc.ID, encryptedParams EncryptedParamsLogSubscription) SystemError

	// Unsubscribe removes the log subscription with the given ID from the enclave. If there is no subscription with
	// the given ID, nothing is deleted.
	Unsubscribe(id rpc.ID) SystemError

	// DebugTraceTransaction returns the trace of a transaction
	DebugTraceTransaction(ctx context.Context, hash gethcommon.Hash, config *tracers.TraceConfig) (json.RawMessage, SystemError)

	GetCode(ctx context.Context, address gethcommon.Address, blockNrOrHash rpc.BlockNumberOrHash) ([]byte, SystemError)

	// GetTotalContractCount returns the total number of contracts that have been deployed
	GetTotalContractCount(context.Context) (*big.Int, SystemError)

	// EnclavePublicConfig returns network data that is known to the enclave but can be shared publicly
	EnclavePublicConfig(context.Context) (*EnclavePublicConfig, SystemError)
}

// Enclave represents the API of the service that runs inside the TEE.
type Enclave interface {
	EnclaveInit
	EnclaveAdmin
	EnclaveClientRPC
}

// BlockSubmissionResponse is the response sent from the enclave back to the node after ingesting a block
type BlockSubmissionResponse struct {
	ProducedSecretResponses []*ProducedSecretResponse // The responses to any secret requests in the ingested L1 block.
	RejectError             *errutil.BlockRejectError // If block was rejected, contains information about what block to submit next.
	RollupMetadata          []ExtRollupMetadata       // Metadata for each rollup that the host needs to prepare APIs
}

// ProducedSecretResponse contains the data to publish to L1 in response to a secret request discovered while processing an L1 block
type ProducedSecretResponse struct {
	Secret      []byte
	RequesterID gethcommon.Address // enclaveID of the enclave that requested the secret
	AttesterID  gethcommon.Address // enclaveID of the enclave that produced the secret
	HostAddress string
}

type EnclavePublicConfig struct {
	L2MessageBusAddress             gethcommon.Address
	TransactionPostProcessorAddress gethcommon.Address
	PublicSystemContracts           map[string]gethcommon.Address
}
