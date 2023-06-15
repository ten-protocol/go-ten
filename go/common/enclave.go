package common

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/tracers"
	"github.com/obscuronet/go-obscuro/go/responses"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

type StatusCode int

// Status represents the enclave's current state - whether the enclave is healthy and ready to process requests, as well
// as its latest known heads for the L1 and L2 chains
type Status struct {
	StatusCode StatusCode
	L1Head     gethcommon.Hash
	L2Head     gethcommon.Hash
}

const (
	Running        StatusCode = iota // the enclave is running, accepting L1 blocks
	AwaitingSecret                   // the enclave has not received the network secret and cannot process L1 blocks
	Unavailable                      // the enclave is unavailable (no guarantee it will self-recover)
)

// Enclave represents the API of the service that runs inside the TEE.
type Enclave interface {
	// Status checks whether the enclave is ready to process requests - only implemented by the RPC layer
	Status() (Status, SystemError)

	// Attestation - Produces an attestation report which will be used to request the shared secret from another enclave.
	Attestation() (*AttestationReport, SystemError)

	// GenerateSecret - the genesis enclave is responsible with generating the secret entropy
	GenerateSecret() (EncryptedSharedEnclaveSecret, SystemError)

	// InitEnclave - initialise an enclave with a seed received by another enclave
	InitEnclave(secret EncryptedSharedEnclaveSecret) SystemError

	// SubmitL1Block - Used for the host to submit L1 blocks to the enclave, these may be:
	//  a. historic block - if the enclave is behind and in the process of catching up with the L1 state
	//  b. the latest block published by the L1, to which the enclave should respond with a rollup
	// It is the responsibility of the host to gossip the returned rollup
	// For good functioning the caller should always submit blocks ordered by height
	// submitting a block before receiving ancestors of it, will result in it being ignored
	SubmitL1Block(block L1Block, receipts L1Receipts, isLatest bool) (*BlockSubmissionResponse, SystemError)

	// SubmitTx - user transactions
	SubmitTx(tx EncryptedTx) (*responses.RawTx, SystemError)

	// SubmitBatch submits a batch received from the sequencer for processing.
	SubmitBatch(batch *ExtBatch) SystemError

	// ObsCall - Execute a smart contract to retrieve data. The equivalent of "Eth_call"
	// Todo - return the result with a block delay. To prevent frontrunning.
	ObsCall(encryptedParams EncryptedParamsCall) (*responses.Call, SystemError)

	// GetTransactionCount returns the nonce of the wallet with the given address (encrypted with the acc viewing key)
	GetTransactionCount(encryptedParams EncryptedParamsGetTxCount) (*responses.TxCount, SystemError)

	// Stop gracefully stops the enclave
	Stop() SystemError

	// GetTransaction returns a transaction in JSON format, encrypted with the viewing key for the transaction's `from` field.
	GetTransaction(encryptedParams EncryptedParamsGetTxByHash) (*responses.TxByHash, SystemError)

	// GetTransactionReceipt returns a transaction receipt given its signed hash, or nil if the transaction is unknown
	GetTransactionReceipt(encryptedParams EncryptedParamsGetTxReceipt) (*responses.TxReceipt, SystemError)

	// AddViewingKey - Decrypts, verifies and saves viewing keys.
	// Viewing keys are asymmetric keys generated inside the wallet extension, and then signed by the wallet (e.g.
	// MetaMask) in which the user holds the signing keys.
	// The keys are then are sent to the enclave via RPC and processed using this method.
	// The first step is to check the validity of the signature over the viewing key.
	// Then, we need to find the account which has empowered this viewing key. We can do that by retrieving the signing
	// public key from the signature. By hashing the public key, we can then determine the address of the account.
	// At the end, we save the viewing key (which is a public key) against the account, and use it to encrypt any
	// "eth_call" and "eth_getBalance" requests that have that address as a "from" field.
	AddViewingKey(encryptedViewingKeyBytes []byte, signature []byte) SystemError

	// GetBalance returns the balance of the address on the Obscuro network, encrypted with the viewing key for the
	// address.
	GetBalance(encryptedParams EncryptedParamsGetBalance) (*responses.Balance, SystemError)

	// GetCode returns the code stored at the given address in the state for the given rollup hash.
	GetCode(address gethcommon.Address, rollupHash *gethcommon.Hash) ([]byte, SystemError)

	// Subscribe adds a log subscription to the enclave under the given ID, provided the request is authenticated
	// correctly. The events will be populated in the BlockSubmissionResponse. If there is an existing subscription
	// with the given ID, it is overwritten.
	Subscribe(id rpc.ID, encryptedParams EncryptedParamsLogSubscription) SystemError

	// Unsubscribe removes the log subscription with the given ID from the enclave. If there is no subscription with
	// the given ID, nothing is deleted.
	Unsubscribe(id rpc.ID) SystemError

	// StopClient stops the enclave client if one exists - only implemented by the RPC layer
	StopClient() SystemError

	// EstimateGas tries to estimate the gas needed to execute a specific transaction based on the pending state.
	EstimateGas(encryptedParams EncryptedParamsEstimateGas) (*responses.Gas, SystemError)

	// GetLogs returns all the logs matching the filter.
	GetLogs(encryptedParams EncryptedParamsGetLogs) (*responses.Logs, SystemError)

	// HealthCheck returns whether the enclave is in a healthy state
	HealthCheck() (bool, SystemError)

	// GetBatch - retrieve a batch if existing within the enclave db.
	GetBatch(hash L2BatchHash) (*ExtBatch, error)

	// CreateBatch - creates a new head batch extending the previous one for the latest known L1 head if the node is
	// a sequencer. Will panic otherwise.
	CreateBatch() SystemError

	// CreateRollup - will create a new rollup by going through the sequencer if the node is a sequencer
	// or panic otherwise.
	CreateRollup() (*ExtRollup, SystemError)

	// DebugTraceTransaction returns the trace of a transaction
	DebugTraceTransaction(hash gethcommon.Hash, config *tracers.TraceConfig) (json.RawMessage, SystemError)

	// StreamL2Updates - will stream the batches following the L2 batch hash given along with any newly created batches
	// in the right order. All will be queued in the channel that has been returned
	StreamL2Updates(*L2BatchHash) (chan StreamL2UpdatesResponse, func())
	// DebugEventLogRelevancy returns the logs of a transaction
	DebugEventLogRelevancy(hash gethcommon.Hash) (json.RawMessage, SystemError)
}

// BlockSubmissionResponse is the response sent from the enclave back to the node after ingesting a block
type BlockSubmissionResponse struct {
	ProducedBatch           *ExtBatch                 // The batch produced iff the node is a sequencer and is on the latest block.
	ProducedRollup          *ExtRollup                // The rollup produced iff the node is a sequencer and it is time to produce a new rollup.
	ProducedSecretResponses []*ProducedSecretResponse // The responses to any secret requests in the ingested L1 block.
	SubscribedLogs          map[rpc.ID][]byte         // The logs produced by the L1 block and all its ancestors for each subscription ID.
	RejectError             *errutil.BlockRejectError // If block was rejected, contains information about what block to submit next.
}

// ProducedSecretResponse contains the data to publish to L1 in response to a secret request discovered while processing an L1 block
type ProducedSecretResponse struct {
	Secret      []byte
	RequesterID gethcommon.Address
	HostAddress string
}
