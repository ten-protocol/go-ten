package common

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// Status represents the enclave's current status - the status and behaviour of the host is a function of the status of the enclave
// since the host's responsibility is to manage communication in and out of the enclave.
type Status int

const (
	Running        Status = iota // the enclave is running, accepting L1 blocks
	AwaitingSecret               // the enclave has not received the network secret and cannot process L1 blocks
	Unavailable                  // the enclave is unavailable (no guarantee it will self-recover)
)

// Enclave represents the API of the service that runs inside the TEE.
type Enclave interface {
	// Status checks whether the enclave is ready to process requests - only implemented by the RPC layer
	Status() (Status, error)

	// Attestation - Produces an attestation report which will be used to request the shared secret from another enclave.
	Attestation() (*AttestationReport, error)

	// GenerateSecret - the genesis enclave is responsible with generating the secret entropy
	GenerateSecret() (EncryptedSharedEnclaveSecret, error)

	// ShareSecret - verify the attestation and return the shared secret (encrypted with the key from the attestation)
	ShareSecret(report *AttestationReport) (EncryptedSharedEnclaveSecret, error)

	// StoreAttestation - store the public key against the counterparties
	StoreAttestation(report *AttestationReport) error

	// InitEnclave - initialise an enclave with a seed received by another enclave
	InitEnclave(secret EncryptedSharedEnclaveSecret) error

	// ProduceGenesis - the genesis enclave produces the genesis rollup
	ProduceGenesis(blkHash gethcommon.Hash) BlockSubmissionResponse

	// IngestBlocks - feed L1 blocks into the enclave to catch up
	IngestBlocks(blocks []*types.Block) []BlockSubmissionResponse

	// Start - start speculative execution
	Start(block types.Block)

	// SubmitBlock - When a new POBI round starts, the host submits a block to the enclave, which responds with a rollup
	// it is the responsibility of the host to gossip the returned rollup
	// For good functioning the caller should always submit blocks ordered by height
	// submitting a block before receiving a parent of it, will result in it being ignored
	SubmitBlock(block types.Block) (BlockSubmissionResponse, error)

	// SubmitRollup - receive gossiped rollups
	SubmitRollup(rollup ExtRollup)

	// SubmitTx - user transactions
	SubmitTx(tx EncryptedTx) (EncryptedResponseSendRawTx, error)

	// ExecuteOffChainTransaction - Execute a smart contract to retrieve data
	// Todo - return the result with a block delay. To prevent frontrunning.
	ExecuteOffChainTransaction(encryptedParams EncryptedParamsCall) (EncryptedResponseCall, error)

	// GetTransactionCount returns the nonce of the wallet with the given address (encrypted with the acc viewing key)
	GetTransactionCount(encryptedParams EncryptedParamsGetTxCount) (EncryptedResponseGetTxCount, error)

	// RoundWinner - calculates and returns the winner for a round, and whether this node is the winner
	RoundWinner(parent L2RootHash) (ExtRollup, bool, error)

	// Stop gracefully stops the enclave
	Stop() error

	// GetTransaction returns a transaction in JSON format, encrypted with the viewing key for the transaction's `from` field.
	GetTransaction(encryptedParams EncryptedParamsGetTxByHash) (EncryptedResponseGetTxByHash, error)

	// GetTransactionReceipt returns a transaction receipt given its signed hash, or nil if the transaction is unknown
	GetTransactionReceipt(encryptedParams EncryptedParamsGetTxReceipt) (EncryptedResponseGetTxReceipt, error)

	// GetRollup returns the rollup with the given hash, or nil if no such rollup exists.
	GetRollup(rollupHash L2RootHash) (*ExtRollup, error)

	// AddViewingKey - Decrypts, verifies and saves viewing keys.
	// Viewing keys are asymmetric keys generated inside the wallet extension, and then signed by the wallet (e.g.
	// MetaMask) in which the user holds the signing keys.
	// The keys are then are sent to the enclave via RPC and processed using this method.
	// The first step is to check the validity of the signature over the viewing key.
	// Then, we need to find the account which has empowered this viewing key. We can do that by retrieving the signing
	// public key from the signature. By hashing the public key, we can then determine the address of the account.
	// At the end, we save the viewing key (which is a public key) against the account, and use it to encrypt any
	// "eth_call" and "eth_getBalance" requests that have that address as a "from" field.
	AddViewingKey(encryptedViewingKeyBytes []byte, signature []byte) error

	// GetBalance returns the balance of the address on the Obscuro network, encrypted with the viewing key for the
	// address.
	GetBalance(encryptedParams EncryptedParamsGetBalance) (EncryptedResponseGetBalance, error)

	// GetCode returns the code stored at the given address in the state for the given rollup hash.
	GetCode(address gethcommon.Address, rollupHash *gethcommon.Hash) ([]byte, error)

	// Subscribe adds a log subscription to the enclave under the given ID, provided the request is authenticated
	// correctly. The events will be populated in the BlockSubmissionResponse. If there is an existing subscription
	// with the given ID, it is overwritten.
	Subscribe(id uuid.UUID, encryptedParams EncryptedParamsLogSubscription) error

	// Unsubscribe removes the log subscription with the given ID from the enclave. If there is no subscription with
	// the given ID, nothing is deleted.
	Unsubscribe(id uuid.UUID) error

	// StopClient stops the enclave client if one exists - only implemented by the RPC layer
	StopClient() error

	// EstimateGas tries to estimate the gas needed to execute a specific transaction based on the pending state.
	EstimateGas(encryptedParams EncryptedParamsEstimateGas) (EncryptedResponseEstimateGas, error)
}

// BlockSubmissionResponse is the response sent from the enclave back to the node after ingesting a block
type BlockSubmissionResponse struct {
	BlockHeader           *types.Header // the header of the consumed block. Todo - only the hash required
	IngestedBlock         bool          // Whether the Block was ingested or discarded
	BlockNotIngestedCause string        // The reason the block was not ingested. This message has to not disclose anything useful from the enclave.

	ProducedRollup ExtRollup // The new Rollup when ingesting the block produces a new Rollup
	FoundNewHead   bool      // Ingested Block contained a new Rollup - Block, and Rollup heads were updated
	RollupHead     *Header   // If a new header was found, this field will be populated with the header of the rollup.

	SubscribedLogs map[uuid.UUID]EncryptedLogs // The logs produced by the block and all its ancestors for each subscription ID.
}
