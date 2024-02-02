package txpool

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/tracers"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"github.com/ten-protocol/go-ten/go/enclave/crypto"
	"github.com/ten-protocol/go-ten/go/enclave/limiters"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

type mockBatchRegistry struct {
	currentBatch *core.Batch
}

func (m *mockBatchRegistry) BatchesAfter(_ uint64, _ uint64, _ limiters.RollupLimiter) ([]*core.Batch, []*types.Block, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockBatchRegistry) GetBatchStateAtHeight(_ *rpc.BlockNumber) (*state.StateDB, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockBatchRegistry) GetBatchAtHeight(_ rpc.BlockNumber) (*core.Batch, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockBatchRegistry) SubscribeForExecutedBatches(_ func(*core.Batch, types.Receipts)) {
	// TODO implement me
	panic("implement me")
}

func (m *mockBatchRegistry) UnsubscribeFromBatches() {
	// TODO implement me
	panic("implement me")
}

func (m *mockBatchRegistry) OnBatchExecuted(batch *core.Batch, _ types.Receipts) {
	m.currentBatch = batch
}

func (m *mockBatchRegistry) HasGenesisBatch() (bool, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockBatchRegistry) HeadBatchSeq() *big.Int {
	return m.currentBatch.SeqNo()
}

func newMockBatchRegistry() *mockBatchRegistry {
	return &mockBatchRegistry{}
}

type mockStorage struct {
	currentBatch  *core.Batch
	batchesSeqNo  map[uint64]*core.Batch
	batchesHeight map[uint64]*core.Batch
	batchesHash   map[gethcommon.Hash]*core.Batch
	stateDB       state.Database
}

func newMockStorage() *mockStorage {
	db := state.NewDatabaseWithConfig(rawdb.NewMemoryDatabase(), &trie.Config{})
	stateDB, err := state.New(types.EmptyRootHash, db, nil)
	if err != nil {
		panic(err)
	}

	_, err = stateDB.Commit(0, true)
	if err != nil {
		panic(err)
	}

	return &mockStorage{
		batchesSeqNo:  map[uint64]*core.Batch{},
		batchesHeight: map[uint64]*core.Batch{},
		batchesHash:   map[gethcommon.Hash]*core.Batch{},
		stateDB:       db,
	}
}

func (m *mockStorage) FetchBlock(_ common.L1BlockHash) (*types.Block, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) FetchCanonicaBlockByHeight(_ *big.Int) (*types.Block, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) FetchHeadBlock() (*types.Block, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) StoreBlock(_ *types.Block, _ *common.ChainFork) error {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) IsAncestor(_ *types.Block, _ *types.Block) bool {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) IsBlockAncestor(_ *types.Block, _ common.L1BlockHash) bool {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) FetchBatch(_ common.L2BatchHash) (*core.Batch, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) FetchBatchHeader(_ common.L2BatchHash) (*common.BatchHeader, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) FetchBatchByHeight(height uint64) (*core.Batch, error) {
	batch, found := m.batchesHeight[height]
	if !found {
		return nil, errutil.ErrNotFound
	}
	return batch, nil
}

func (m *mockStorage) FetchBatchBySeqNo(seqNum uint64) (*core.Batch, error) {
	batch, found := m.batchesSeqNo[seqNum]
	if !found {
		return nil, errutil.ErrNotFound
	}
	return batch, nil
}

func (m *mockStorage) FetchHeadBatch() (*core.Batch, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) FetchCurrentSequencerNo() (*big.Int, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) FetchBatchesByBlock(_ common.L1BlockHash) ([]*core.Batch, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) FetchNonCanonicalBatchesBetween(_ uint64, _ uint64) ([]*core.Batch, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) FetchCanonicalUnexecutedBatches(_ *big.Int) ([]*core.Batch, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) BatchWasExecuted(_ common.L2BatchHash) (bool, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) FetchHeadBatchForBlock(_ common.L1BlockHash) (*core.Batch, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) StoreBatch(_ *core.Batch) error {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) StoreExecutedBatch(batch *core.Batch, _ []*types.Receipt) error {
	m.currentBatch = batch
	m.batchesSeqNo[batch.SeqNo().Uint64()] = batch
	m.batchesHeight[batch.Number().Uint64()] = batch
	m.batchesHash[batch.Hash()] = batch
	return nil
}

func (m *mockStorage) StoreRollup(_ *common.ExtRollup, _ *common.CalldataRollupHeader) error {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) FetchReorgedRollup(_ []common.L1BlockHash) (*common.L2BatchHash, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) CreateStateDB(hash common.L2BatchHash) (*state.StateDB, error) {
	batch, found := m.batchesHash[hash]
	if !found {
		return nil, errutil.ErrNotFound
	}
	return state.New(batch.Header.Root, m.stateDB, nil)
}

func (m *mockStorage) EmptyStateDB() (*state.StateDB, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) FetchSecret() (*crypto.SharedEnclaveSecret, error) {
	return &crypto.SharedEnclaveSecret{}, nil
}

func (m *mockStorage) StoreSecret(_ crypto.SharedEnclaveSecret) error {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) GetTransaction(_ common.L2TxHash) (*types.Transaction, gethcommon.Hash, uint64, uint64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) GetTransactionReceipt(_ common.L2TxHash) (*types.Receipt, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) GetReceiptsByBatchHash(_ common.L2BatchHash) (types.Receipts, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) GetContractCreationTx(_ gethcommon.Address) (*gethcommon.Hash, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) FetchAttestedKey(_ gethcommon.Address) (*ecdsa.PublicKey, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) StoreAttestedKey(_ gethcommon.Address, _ *ecdsa.PublicKey) error {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) StoreL1Messages(_ common.L1BlockHash, _ common.CrossChainMessages) error {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) GetL1Messages(_ common.L1BlockHash) (common.CrossChainMessages, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) StoreValueTransfers(_ common.L1BlockHash, _ common.ValueTransferEvents) error {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) GetL1Transfers(_ common.L1BlockHash) (common.ValueTransferEvents, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) StoreEnclaveKey(_ *ecdsa.PrivateKey) error {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) GetEnclaveKey() (*ecdsa.PrivateKey, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) GetContractCount() (*big.Int, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) GetReceiptsPerAddress(_ *gethcommon.Address, _ *common.QueryPagination) (types.Receipts, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) GetPublicTransactionData(_ *common.QueryPagination) ([]common.PublicTransaction, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) GetPublicTransactionCount() (uint64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) GetReceiptsPerAddressCount(_ *gethcommon.Address) (uint64, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) Close() error {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) HealthCheck() (bool, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) FilterLogs(_ *gethcommon.Address, _, _ *big.Int, _ *common.L2BatchHash, _ []gethcommon.Address, _ [][]gethcommon.Hash) ([]*types.Log, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) DebugGetLogs(_ common.TxHash) ([]*tracers.DebugLogs, error) {
	// TODO implement me
	panic("implement me")
}

func (m *mockStorage) TrieDB() *trie.Database {
	// TODO implement me
	panic("implement me")
}
