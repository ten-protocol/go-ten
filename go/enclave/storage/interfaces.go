package storage

import (
	"context"
	"crypto/ecdsa"
	"io"
	"math/big"

	"github.com/ten-protocol/go-ten/go/enclave/storage/enclavedb"

	"github.com/ethereum/go-ethereum/triedb"

	"github.com/ethereum/go-ethereum/core/state"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"github.com/ten-protocol/go-ten/go/enclave/crypto"
)

// BlockResolver stores new blocks and returns information on existing blocks
type BlockResolver interface {
	// FetchBlock returns the L1 BlockHeader with the given hash.
	FetchBlock(ctx context.Context, blockHash common.L1BlockHash) (*types.Header, error)
	IsBlockCanonical(ctx context.Context, blockHash common.L1BlockHash) (bool, error)
	// FetchCanonicaBlockByHeight - self explanatory
	FetchCanonicaBlockByHeight(ctx context.Context, height *big.Int) (*types.Header, error)
	// FetchHeadBlock - returns the head of the current chain.
	FetchHeadBlock(ctx context.Context) (*types.Header, error)
	// StoreBlock persists the L1 BlockHeader and updates the canonical ancestors if there was a fork
	StoreBlock(ctx context.Context, block *types.Header, fork *common.ChainFork) error
	// IsAncestor returns true if maybeAncestor is an ancestor of the L1 BlockHeader, and false otherwise
	IsAncestor(ctx context.Context, block *types.Header, maybeAncestor *types.Header) bool
}

type BatchResolver interface {
	// FetchBatch returns the batch with the given hash.
	FetchBatch(ctx context.Context, hash common.L2BatchHash) (*core.Batch, error)
	// FetchBatchHeader returns the batch header with the given hash.
	FetchBatchHeader(ctx context.Context, hash common.L2BatchHash) (*common.BatchHeader, error)
	FetchBatchTransactionsBySeq(ctx context.Context, seqNo uint64) ([]*common.L2Tx, error)
	// FetchBatchByHeight returns the batch on the canonical chain with the given height.
	FetchBatchByHeight(ctx context.Context, height uint64) (*core.Batch, error)
	// FetchBatchBySeqNo returns the batch with the given seq number.
	FetchBatchBySeqNo(ctx context.Context, seqNum uint64) (*core.Batch, error)
	// FetchBatchHeaderBySeqNo returns the batch header with the given seq number.
	FetchBatchHeaderBySeqNo(ctx context.Context, seqNum uint64) (*common.BatchHeader, error)
	FetchHeadBatchHeader(ctx context.Context) (*common.BatchHeader, error)
	// FetchCurrentSequencerNo returns the sequencer number
	FetchCurrentSequencerNo(ctx context.Context) (*big.Int, error)
	// FetchBatchesByBlock returns all batches with the block hash as the L1 proof
	FetchBatchesByBlock(ctx context.Context, hash common.L1BlockHash) ([]*common.BatchHeader, error)
	// FetchNonCanonicalBatchesBetween - returns all reorged batches between the sequences
	FetchNonCanonicalBatchesBetween(ctx context.Context, startSeq uint64, endSeq uint64) ([]*common.BatchHeader, error)
	// FetchCanonicalBatchesBetween - returns all canon batches between the sequences
	FetchCanonicalBatchesBetween(ctx context.Context, startSeq uint64, endSeq uint64) ([]*common.BatchHeader, error)
	// IsBatchCanonical - true if the batch is canonical
	IsBatchCanonical(ctx context.Context, seq uint64) (bool, error)
	// FetchCanonicalUnexecutedBatches - return the list of the unexecuted batches that are canonical
	FetchCanonicalUnexecutedBatches(context.Context, *big.Int) ([]*common.BatchHeader, error)

	FetchConvertedHash(ctx context.Context, hash common.L2BatchHash) (gethcommon.Hash, error)

	// BatchWasExecuted - return true if the batch was executed
	BatchWasExecuted(ctx context.Context, hash common.L2BatchHash) (bool, error)

	// StoreBatch stores an un-executed batch.
	StoreBatch(ctx context.Context, batch *core.Batch, convertedHash gethcommon.Hash) error
	// StoreExecutedBatch - store the batch after it was executed
	StoreExecutedBatch(ctx context.Context, batch *core.Batch, results core.TxExecResults) error

	// StoreRollup
	StoreRollup(ctx context.Context, rollup *common.ExtRollup, header *common.CalldataRollupHeader) error
	FetchRollupMetadata(ctx context.Context, hash common.L2RollupHash) (*common.PublicRollupMetadata, error)
	FetchReorgedRollup(ctx context.Context, reorgedBlocks []common.L1BlockHash) (*common.L2BatchHash, error)
}

type GethStateDB interface {
	// CreateStateDB creates a database that can be used to execute transactions
	CreateStateDB(ctx context.Context, hash common.L2BatchHash) (*state.StateDB, error)
	// EmptyStateDB creates the original empty StateDB
	EmptyStateDB() (*state.StateDB, error)
}

type SharedSecretStorage interface {
	// FetchSecret returns the enclave's secret.
	FetchSecret(ctx context.Context) (*crypto.SharedEnclaveSecret, error)
	// StoreSecret stores a secret in the enclave
	StoreSecret(ctx context.Context, secret crypto.SharedEnclaveSecret) error
}

type TransactionStorage interface {
	// GetTransaction - returns the positional metadata of the tx by hash
	GetTransaction(ctx context.Context, txHash common.L2TxHash) (*types.Transaction, common.L2BatchHash, uint64, uint64, gethcommon.Address, error)
	// GetFilteredInternalReceipt - returns the receipt of a tx with event logs visible to the requester
	GetFilteredInternalReceipt(ctx context.Context, txHash common.L2TxHash, requester *gethcommon.Address, syntheticTx bool) (*core.InternalReceipt, error)
	ExistsTransactionReceipt(ctx context.Context, txHash common.L2TxHash) (bool, error)
}

type AttestationStorage interface {
	GetEnclavePubKey(ctx context.Context, enclaveId common.EnclaveID) (*AttestedEnclave, error)
	StoreNewEnclave(ctx context.Context, enclaveId common.EnclaveID, key *ecdsa.PublicKey) error
	StoreNodeType(ctx context.Context, enclaveId common.EnclaveID, nodeType common.NodeType) error
	GetSequencerEnclaveIDs(ctx context.Context) ([]common.EnclaveID, error)
}

type CrossChainMessagesStorage interface {
	StoreL1Messages(ctx context.Context, blockHash common.L1BlockHash, messages common.CrossChainMessages) error
	GetL1Messages(ctx context.Context, blockHash common.L1BlockHash) (common.CrossChainMessages, error)

	StoreValueTransfers(ctx context.Context, blockHash common.L1BlockHash, transfers common.ValueTransferEvents) error
	GetL1Transfers(ctx context.Context, blockHash common.L1BlockHash) (common.ValueTransferEvents, error)
}

type EnclaveKeyStorage interface {
	StoreEnclaveKey(ctx context.Context, enclaveKey []byte) error
	GetEnclaveKey(ctx context.Context) ([]byte, error)
}

type SystemContractAddressesStorage interface {
	StoreSystemContractAddresses(ctx context.Context, addresses common.SystemContractAddresses) error
	GetSystemContractAddresses(ctx context.Context) (common.SystemContractAddresses, error)
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
	SystemContractAddressesStorage
	io.Closer

	// HealthCheck returns whether the storage is deemed healthy or not
	HealthCheck(ctx context.Context) (bool, error)

	// FilterLogs - applies the properties the relevancy checks for the requestingAccount to all the stored log events
	// nil values will be ignored. Make sure to set all fields to the right values before calling this function
	// the blockHash should always be nil.
	FilterLogs(ctx context.Context, requestingAccount *gethcommon.Address, fromBlock, toBlock *big.Int, blockHash *common.L2BatchHash, addresses []gethcommon.Address, topics [][]gethcommon.Hash) ([]*types.Log, error)

	// DebugGetLogs returns logs for a given tx hash without any constraints - should only be used for debug purposes
	DebugGetLogs(ctx context.Context, from *big.Int, to *big.Int, address gethcommon.Address, eventSig gethcommon.Hash) ([]*common.DebugLogVisibility, error)

	// TrieDB - return the underlying trie database
	TrieDB() *triedb.Database

	// StateDB - return the underlying state database
	StateDB() state.Database

	ReadContract(ctx context.Context, address gethcommon.Address) (*enclavedb.Contract, error)
	ReadEventType(ctx context.Context, contractAddress gethcommon.Address, eventSignature gethcommon.Hash) (*enclavedb.EventType, error)
}

type ScanStorage interface {
	GetContractCount(ctx context.Context) (*big.Int, error)
	GetTransactionsPerAddress(ctx context.Context, address *gethcommon.Address, pagination *common.QueryPagination) ([]*core.InternalReceipt, error)

	CountTransactionsPerAddress(ctx context.Context, addr *gethcommon.Address) (uint64, error)
}
