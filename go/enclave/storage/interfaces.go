package storage

import (
	"crypto/ecdsa"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/core/state"

	"github.com/ethereum/go-ethereum/trie"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/tracers"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/crypto"
)

// BlockResolver stores new blocks and returns information on existing blocks
type BlockResolver interface {
	// FetchBlock returns the L1 Block with the given hash.
	FetchBlock(blockHash common.L1BlockHash) (*types.Block, error)
	// FetchCanonicaBlockByHeight - self explanatory
	FetchCanonicaBlockByHeight(height *big.Int) (*types.Block, error)
	// FetchHeadBlock - returns the head of the current chain.
	FetchHeadBlock() (*types.Block, error)
	// StoreBlock persists the L1 Block and updates the canonical ancestors if there was a fork
	StoreBlock(block *types.Block, fork *common.ChainFork) error
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
	// FetchBatchesByBlock returns all batches with the block hash as the L1 proof
	FetchBatchesByBlock(common.L1BlockHash) ([]*core.Batch, error)

	// FetchCanonicalUnexecutedBatches - return the list of the unexecuted batches that are canonical
	FetchCanonicalUnexecutedBatches(*big.Int) ([]*core.Batch, error)

	// BatchWasExecuted - return true if the batch was executed
	BatchWasExecuted(hash common.L2BatchHash) (bool, error)

	// FetchHeadBatchForBlock returns the hash of the head batch at a given L1 block.
	FetchHeadBatchForBlock(blockHash common.L1BlockHash) (*core.Batch, error)

	// StoreBatch stores an un-executed batch.
	StoreBatch(batch *core.Batch) error
	// StoreExecutedBatch - store the batch after it was executed
	StoreExecutedBatch(batch *core.Batch, receipts []*types.Receipt) error

	// StoreRollup
	StoreRollup(rollup *common.ExtRollup, header *common.CalldataRollupHeader) error
	FetchReorgedRollup(reorgedBlocks []common.L1BlockHash) (*common.L2BatchHash, error)
}

type GethStateDB interface {
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
	// GetReceiptsByBatchHash retrieves the receipts for all transactions in a given rollup.
	GetReceiptsByBatchHash(hash common.L2BatchHash) (types.Receipts, error)
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

	StoreValueTransfers(blockHash common.L1BlockHash, transfers common.ValueTransferEvents) error
	GetL1Transfers(blockHash common.L1BlockHash) (common.ValueTransferEvents, error)
}

type EnclaveKeyStorage interface {
	StoreEnclaveKey(enclaveKey *ecdsa.PrivateKey) error
	GetEnclaveKey() (*ecdsa.PrivateKey, error)
}

// Storage is the enclave's interface for interacting with the enclave's datastore
type Storage interface {
	BlockResolver
	BatchResolver
	GethStateDB
	SharedSecretStorage
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

	// TrieDB - return the underlying trie database
	TrieDB() *trie.Database
}

type ScanStorage interface {
	GetContractCount() (*big.Int, error)
	GetReceiptsPerAddress(address *gethcommon.Address, pagination *common.QueryPagination) (types.Receipts, error)
	GetPublicTransactionData(pagination *common.QueryPagination) ([]common.PublicTransaction, error)
	GetPublicTransactionCount() (uint64, error)

	GetReceiptsPerAddressCount(addr *gethcommon.Address) (uint64, error)
}
