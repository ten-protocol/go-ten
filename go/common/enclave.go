package common

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
)

// Enclave represents the API of the service that runs inside the TEE.
type Enclave interface {
	// IsReady checks whether the enclave is ready to process requests - only implemented by the RPC layer
	IsReady() error

	// Attestation - Produces an attestation report which will be used to request the shared secret from another enclave.
	Attestation() *AttestationReport

	// GenerateSecret - the genesis enclave is responsible with generating the secret entropy
	GenerateSecret() EncryptedSharedEnclaveSecret

	// ShareSecret - verify the attestation and return the shared secret (encrypted with the key from the attestation)
	ShareSecret(report *AttestationReport) (EncryptedSharedEnclaveSecret, error)

	// StoreAttestation - store the public key against the counterparties
	StoreAttestation(report *AttestationReport) error

	// InitEnclave - initialise an enclave with a seed received by another enclave
	InitEnclave(secret EncryptedSharedEnclaveSecret) error

	// IsInitialised - true if the shared secret is available
	IsInitialised() bool

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

	// Nonce returns the nonce of the wallet with the given address.
	Nonce(address gethcommon.Address) uint64

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

	// Subscribe registers a new event subscription. The events will be populated in the BlockSubmissionResponse
	Subscribe(id uuid.UUID, subscription EncryptedEventSubscription) error

	// Unsubscribe - removes a subscription
	Unsubscribe(id uuid.UUID) error

	// StopClient stops the enclave client if one exists - only implemented by the RPC layer
	StopClient() error
}

// BlockSubmissionResponse is the response sent from the enclave back to the node after ingesting a block
type BlockSubmissionResponse struct {
	BlockHeader           *types.Header // the header of the consumed block. Todo - only the hash required
	IngestedBlock         bool          // Whether the Block was ingested or discarded
	BlockNotIngestedCause string        // The reason the block was not ingested. This message has to not disclose anything useful from the enclave.

	ProducedRollup ExtRollup // The new Rollup when ingesting the block produces a new Rollup
	FoundNewHead   bool      // Ingested Block contained a new Rollup - Block, and Rollup heads were updated
	RollupHead     *Header   // If a new header was found, this field will be populated with the header of the rollup.

	SubscribedEvents map[uuid.UUID]EncryptedEvents // For each subscription id, there is an encrypted list of events, which has to be sent back to the requester by the host.
}

// EventSubscription
// From the design - call must take a list of signed owning accounts.
// Each account must be signed with the latest viewing key (to prevent someone from asking random events, just to leak info).
// The call will fail if there are no viewing keys for all those accounts.
type EventSubscription struct {
	ID       uuid.UUID
	Accounts []*SubscriptionAccount
	// todo Filters - the geth log filters
}

func (s EventSubscription) Matches(log *types.Log, db *state.StateDB) (bool, *SubscriptionAccount) {
	// todo
	// transform the log into a useful data structure by extracting addresses from the log (according to the design)
	// identify what type of log it is ( account specific or relevant)
	// if account-specific go through each SubscriptionAccount and check whether the log is relevant
	// note - the above logic has to be reused to filter out the logs when someone requests a transaction receipt
	// for logs that pass the above the step apply the filters
	// return the first account for which the log matches, so it can be used for encryption
	return true, nil
}

// SubscriptionAccount - represents an authenticated account used for subscribing to events
// note - the fields below are just indicative, to showcase the required functionality
type SubscriptionAccount struct {
	Account    gethcommon.Address // the account for which the events are relevant
	ViewingKey []byte             // the viewing key to use for the encryption of events relevant to this account
	SignedKey  []byte             // public viewing key signed by the Account's private key. Useful to authenticate the VK.
	Signature  []byte             // a signature over the account using the private viewing key. To prevent anyone but the account owner to request subscriptions to leak data.
}
