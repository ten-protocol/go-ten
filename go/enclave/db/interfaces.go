package db

import (
	"crypto/ecdsa"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/trie"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/tracers"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/crypto"
	"github.com/obscuronet/go-obscuro/go/enclave/db/sql"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// BlockResolver stores new blocks and returns information on existing blocks
type BlockResolver interface {
	// FetchBlock returns the L1 Block with the given hash.
	FetchBlock(blockHash common.L1BlockHash) (*types.Block, error)
	// FetchHeadBlock - returns the head of the current chain.
	FetchHeadBlock() (*types.Block, error)
	// StoreBlock persists the L1 Block
	StoreBlock(block *types.Block)
	// IsAncestor returns true if maybeAncestor is an ancestor of the L1 Block, and false otherwise
	IsAncestor(block *types.Block, maybeAncestor *types.Block) bool
	// IsBlockAncestor returns true if maybeAncestor is an ancestor of the L1 Block, and false otherwise
	// Takes into consideration that the Block to verify might be on a branch we haven't received yet
	// todo (low priority) - this is super confusing, analyze the usage
	IsBlockAncestor(block *types.Block, maybeAncestor common.L1BlockHash) bool
}

type BatchResolver interface {
	// FetchBatch returns the batch with the given hash.
	FetchBatch(hash common.L2BatchHash) (*core.Batch, error)
	// FetchBatchHeader returns the batch header with the given hash.
	FetchBatchHeader(hash common.L2BatchHash) (*common.BatchHeader, error)
	// FetchBatchByHeight returns the batch on the canonical chain with the given height.
	FetchBatchByHeight(height uint64) (*core.Batch, error)
	// FetchBatchBySeqNo returns the batch with the given seq number.
	FetchBatchBySeqNo(seqNum uint64) (*core.Batch, error)
	// FetchHeadBatch returns the current head batch of the canonical chain.
	FetchHeadBatch() (*core.Batch, error)
	// FetchCurrentSequencerNo returns the sequencer number
	FetchCurrentSequencerNo() (*big.Int, error)
}

type BatchUpdater interface {
	// StoreBatch stores a batch.
	StoreBatch(batch *core.Batch, receipts []*types.Receipt, dbBatch *sql.Batch) error
	// UpdateHeadBatch updates the canonical L2 head batch for a given L1 block.
	UpdateHeadBatch(l1Head common.L1BlockHash, l2Head *core.Batch, receipts []*types.Receipt, dbBatch *sql.Batch) error
	// SetHeadBatchPointer updates the canonical L2 head batch for a given L1 block.
	SetHeadBatchPointer(l2Head *core.Batch, dbBatch *sql.Batch) error
}

type HeadsAfterL1BlockStorage interface {
	// FetchHeadBatchForBlock returns the hash of the head batch at a given L1 block.
	FetchHeadBatchForBlock(blockHash common.L1BlockHash) (*core.Batch, error)
	// UpdateL1Head updates the L1 head.
	UpdateL1Head(l1Head common.L1BlockHash) error
	// CreateStateDB creates a database that can be used to execute transactions
	CreateStateDB(hash common.L2BatchHash) (*state.StateDB, error)
	// EmptyStateDB creates the original empty StateDB
	EmptyStateDB() (*state.StateDB, error)
}

type SharedSecretStorage interface {
	// FetchSecret returns the enclave's secret.
	FetchSecret() (*crypto.SharedEnclaveSecret, error)
	// StoreSecret stores a secret in the enclave
	StoreSecret(secret crypto.SharedEnclaveSecret) error
}

type TransactionStorage interface {
	// GetTransaction - returns the positional metadata of the tx by hash
	GetTransaction(txHash common.L2TxHash) (*types.Transaction, gethcommon.Hash, uint64, uint64, error)
	// GetTransactionReceipt - returns the receipt of a tx by tx hash
	GetTransactionReceipt(txHash common.L2TxHash) (*types.Receipt, error)
	// GetReceiptsByHash retrieves the receipts for all transactions in a given rollup.
	GetReceiptsByHash(hash common.L2BatchHash) (types.Receipts, error)
	// GetSender returns the sender of the tx by hash
	GetSender(txHash common.L2TxHash) (gethcommon.Address, error)
	// GetContractCreationTx returns the hash of the tx that created a contract
	GetContractCreationTx(address gethcommon.Address) (*gethcommon.Hash, error)
}

type AttestationStorage interface {
	// FetchAttestedKey returns the public key of an attested aggregator
	FetchAttestedKey(aggregator gethcommon.Address) (*ecdsa.PublicKey, error)
	// StoreAttestedKey - store the public key of an attested aggregator
	StoreAttestedKey(aggregator gethcommon.Address, key *ecdsa.PublicKey) error
}

type CrossChainMessagesStorage interface {
	StoreL1Messages(blockHash common.L1BlockHash, messages common.CrossChainMessages) error
	GetL1Messages(blockHash common.L1BlockHash) (common.CrossChainMessages, error)
}

type EnclaveKeyStorage interface {
	StoreEnclaveKey(enclaveKey *ecdsa.PrivateKey) error
	GetEnclaveKey() (*ecdsa.PrivateKey, error)
}

// Storage is the enclave's interface for interacting with the enclave's datastore
type Storage interface {
	BlockResolver
	BatchResolver
	BatchUpdater
	SharedSecretStorage
	HeadsAfterL1BlockStorage
	TransactionStorage
	AttestationStorage
	CrossChainMessagesStorage
	EnclaveKeyStorage
	ScanStorage
	io.Closer

	// HealthCheck returns whether the storage is deemed healthy or not
	HealthCheck() (bool, error)

	// FilterLogs - applies the properties the relevancy checks for the requestingAccount to all the stored log events
	// nil values will be ignored. Make sure to set all fields to the right values before calling this function
	// the blockHash should always be nil.
	FilterLogs(requestingAccount *gethcommon.Address, fromBlock, toBlock *big.Int, blockHash *common.L2BatchHash, addresses []gethcommon.Address, topics [][]gethcommon.Hash) ([]*types.Log, error)

	// DebugGetLogs returns logs for a given tx hash without any constraints - should only be used for debug purposes
	DebugGetLogs(txHash common.TxHash) ([]*tracers.DebugLogs, error)

	// todo (@stefan) - OpenBatch should return a custom type that hides any methods and properties to outside callers
	// in order to prevent accidental messing up the internal state
	// OpenBatch - returns a batch struct that allows for grouping write calls to the database together.
	OpenBatch() *sql.Batch
	// CommitBatch - finalizes a batch and pushes the changes to the database
	CommitBatch(dbBatch *sql.Batch) error

	// TrieDB - return the underlying trie database
	TrieDB() *trie.Database
}

type ScanStorage interface {
	GetContractCount() (*big.Int, error)
}
