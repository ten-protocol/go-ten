package db

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/crypto"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	obscurorawdb "github.com/obscuronet/go-obscuro/go/enclave/db/rawdb"
)

// ErrNoRollups is returned if no rollups have been published yet in the history of the network
// Note: this is not just "not found", we cache at every L1 block what rollup we are up to so we also record that we haven't seen one yet
var ErrNoRollups = errors.New("no rollups have been published")

// TODO - Consistency around whether we assert the secret is available or not.

type storageImpl struct {
	db          ethdb.Database
	stateDB     state.Database
	chainConfig *params.ChainConfig
	logger      gethlog.Logger
}

func NewStorage(backingDB ethdb.Database, chainConfig *params.ChainConfig, logger gethlog.Logger) Storage {
	return &storageImpl{
		db:          backingDB,
		stateDB:     state.NewDatabase(backingDB),
		chainConfig: chainConfig,
		logger:      logger,
	}
}

func (s *storageImpl) FetchHeadBatch() (*core.Batch, error) {
	headHash, err := obscurorawdb.ReadL2HeadBatch(s.db)
	if err != nil {
		return nil, err
	}
	if (bytes.Equal(headHash.Bytes(), gethcommon.Hash{}.Bytes())) {
		return nil, errutil.ErrNotFound
	}
	return s.FetchBatch(*headHash)
}

func (s *storageImpl) FetchBatch(hash common.L2RootHash) (*core.Batch, error) {
	s.assertSecretAvailable()
	batch, err := obscurorawdb.ReadBatch(s.db, hash)
	if err != nil {
		return nil, err
	}
	return batch, nil
}

func (s *storageImpl) FetchBatchByHeight(height uint64) (*core.Batch, error) {
	hash, err := obscurorawdb.ReadCanonicalBatchHash(s.db, height)
	if err != nil {
		return nil, err
	}
	return s.FetchBatch(*hash)
}

func (s *storageImpl) StoreBlock(b *types.Block) {
	s.assertSecretAvailable()
	rawdb.WriteBlock(s.db, b)
}

func (s *storageImpl) FetchBlock(blockHash common.L1RootHash) (*types.Block, error) {
	s.assertSecretAvailable()
	height := rawdb.ReadHeaderNumber(s.db, blockHash)
	if height == nil {
		return nil, errutil.ErrNotFound
	}
	b := rawdb.ReadBlock(s.db, blockHash, *height)
	if b == nil {
		return nil, errutil.ErrNotFound
	}
	return b, nil
}

func (s *storageImpl) FetchHeadBlock() (*types.Block, error) {
	s.assertSecretAvailable()
	block, err := s.FetchBlock(rawdb.ReadHeadHeaderHash(s.db))
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (s *storageImpl) StoreSecret(secret crypto.SharedEnclaveSecret) error {
	return obscurorawdb.WriteSharedSecret(s.db, secret)
}

func (s *storageImpl) FetchSecret() (*crypto.SharedEnclaveSecret, error) {
	return obscurorawdb.ReadSharedSecret(s.db)
}

func (s *storageImpl) IsAncestor(block *types.Block, maybeAncestor *types.Block) bool {
	s.assertSecretAvailable()
	if bytes.Equal(maybeAncestor.Hash().Bytes(), block.Hash().Bytes()) {
		return true
	}

	if maybeAncestor.NumberU64() >= block.NumberU64() {
		return false
	}

	p, err := s.FetchBlock(block.ParentHash())
	if err != nil {
		return false
	}

	return s.IsAncestor(p, maybeAncestor)
}

func (s *storageImpl) IsBlockAncestor(block *types.Block, maybeAncestor common.L1RootHash) bool {
	resolvedBlock, err := s.FetchBlock(maybeAncestor)
	if err != nil {
		return false
	}
	return s.IsAncestor(block, resolvedBlock)
}

func (s *storageImpl) HealthCheck() (bool, error) {
	headBatch, err := s.FetchHeadBatch()
	if err != nil {
		s.logger.Error("unable to HealthCheck storage", log.ErrKey, err)
		return false, err
	}
	return headBatch != nil, nil
}

func (s *storageImpl) assertSecretAvailable() {
	// TODO uncomment this
	//if s.FetchSecret() == nil {
	//	panic("Enclave not initialized")
	//}
}

func (s *storageImpl) FetchHeadBatchForBlock(blockHash common.L1RootHash) (*core.Batch, error) {
	l2HeadBatch, err := obscurorawdb.ReadL2HeadBatchForBlock(s.db, blockHash)
	if err != nil {
		return nil, fmt.Errorf("could not read L2 head batch for block. Cause: %w", err)
	}
	return obscurorawdb.ReadBatch(s.db, *l2HeadBatch)
}

func (s *storageImpl) FetchHeadRollupForBlock(blockHash *common.L1RootHash) (*core.Rollup, error) {
	l2HeadBatch, err := obscurorawdb.ReadL2HeadRollup(s.db, blockHash)
	if err != nil {
		return nil, fmt.Errorf("could not read L2 head rollup for block. Cause: %w", err)
	}
	if *l2HeadBatch == (gethcommon.Hash{}) { // empty hash ==> no rollups yet up to this block
		return nil, ErrNoRollups
	}
	return obscurorawdb.ReadRollup(s.db, *l2HeadBatch)
}

func (s *storageImpl) FetchLogs(blockHash common.L1RootHash) ([]*types.Log, error) {
	logs, err := obscurorawdb.ReadBlockLogs(s.db, blockHash)
	if err != nil {
		// TODO - Return the error itself, once we move from `errutil.ErrNotFound` to `ethereum.NotFound`
		return nil, errutil.ErrNotFound
	}
	return logs, nil
}

func (s *storageImpl) UpdateHeadBatch(l1Head common.L1RootHash, l2Head *core.Batch, receipts []*types.Receipt) error {
	dbBatch := s.db.NewBatch()

	if err := obscurorawdb.SetL2HeadBatch(dbBatch, *l2Head.Hash()); err != nil {
		return fmt.Errorf("could not write block state. Cause: %w", err)
	}
	if err := obscurorawdb.WriteL1ToL2BatchMapping(dbBatch, l1Head, *l2Head.Hash()); err != nil {
		return fmt.Errorf("could not write block state. Cause: %w", err)
	}

	// We update the canonical hash of the batch at this height.
	if err := obscurorawdb.WriteCanonicalHash(dbBatch, l2Head); err != nil {
		return fmt.Errorf("could not write canonical hash. Cause: %w", err)
	}

	// We update the block's logs, based on the batch's logs.
	var logs []*types.Log
	for _, receipt := range receipts {
		logs = append(logs, receipt.Logs...)
	}
	if err := obscurorawdb.WriteBlockLogs(dbBatch, l1Head, logs); err != nil {
		return fmt.Errorf("could not write block logs. Cause: %w", err)
	}

	if err := dbBatch.Write(); err != nil {
		return fmt.Errorf("could not save new head. Cause: %w", err)
	}
	return nil
}

func (s *storageImpl) SetHeadBatchPointer(l2Head *core.Batch) error {
	dbBatch := s.db.NewBatch()

	// We update the canonical hash of the batch at this height.
	if err := obscurorawdb.SetL2HeadBatch(dbBatch, *l2Head.Hash()); err != nil {
		return fmt.Errorf("could not write canonical hash. Cause: %w", err)
	}
	if err := dbBatch.Write(); err != nil {
		return fmt.Errorf("could not save new head. Cause: %w", err)
	}
	return nil
}

func (s *storageImpl) UpdateHeadRollup(l1Head *common.L1RootHash, l2Head *common.L2RootHash) error {
	dbBatch := s.db.NewBatch()
	if err := obscurorawdb.WriteL2HeadRollup(dbBatch, l1Head, l2Head); err != nil {
		return fmt.Errorf("could not write block state. Cause: %w", err)
	}
	if err := dbBatch.Write(); err != nil {
		return fmt.Errorf("could not save new head. Cause: %w", err)
	}
	return nil
}

func (s *storageImpl) UpdateL1Head(l1Head common.L1RootHash) error {
	dbBatch := s.db.NewBatch()
	rawdb.WriteHeadHeaderHash(dbBatch, l1Head)
	if err := dbBatch.Write(); err != nil {
		return fmt.Errorf("could not save new L1 head. Cause: %w", err)
	}
	return nil
}

func (s *storageImpl) CreateStateDB(hash common.L2RootHash) (*state.StateDB, error) {
	batch, err := s.FetchBatch(hash)
	if err != nil {
		return nil, err
	}

	// todo - snapshots?
	statedb, err := state.New(batch.Header.Root, s.stateDB, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create state DB. Cause: %w", err)
	}

	return statedb, nil
}

func (s *storageImpl) EmptyStateDB() (*state.StateDB, error) {
	statedb, err := state.New(gethcommon.BigToHash(big.NewInt(0)), s.stateDB, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create state DB. Cause: %w", err)
	}
	return statedb, nil
}

// GetReceiptsByHash retrieves the receipts for all transactions in a given batch.
func (s *storageImpl) GetReceiptsByHash(hash gethcommon.Hash) (types.Receipts, error) {
	number, err := obscurorawdb.ReadBatchNumber(s.db, hash)
	if err != nil {
		return nil, err
	}
	return obscurorawdb.ReadReceipts(s.db, hash, *number, s.chainConfig)
}

func (s *storageImpl) GetTransaction(txHash gethcommon.Hash) (*types.Transaction, gethcommon.Hash, uint64, uint64, error) {
	tx, blockHash, blockNumber, index, err := obscurorawdb.ReadTransaction(s.db, txHash)
	if err != nil {
		return nil, gethcommon.Hash{}, 0, 0, err
	}
	return tx, blockHash, blockNumber, index, nil
}

func (s *storageImpl) GetSender(txHash gethcommon.Hash) (gethcommon.Address, error) {
	tx, _, _, _, err := s.GetTransaction(txHash) //nolint:dogsled
	if err != nil {
		return gethcommon.Address{}, err
	}
	// todo - make the signer a field of the rollup chain
	msg, err := tx.AsMessage(types.NewLondonSigner(tx.ChainId()), nil)
	if err != nil {
		return gethcommon.Address{}, fmt.Errorf("could not convert transaction to message to retrieve sender address in eth_getTransactionReceipt request. Cause: %w", err)
	}
	return msg.From(), nil
}

func (s *storageImpl) GetContractCreationTx(address gethcommon.Address) (*gethcommon.Hash, error) {
	return obscurorawdb.ReadContractTransaction(s.db, address)
}

func (s *storageImpl) GetTransactionReceipt(txHash gethcommon.Hash) (*types.Receipt, error) {
	_, blockHash, _, index, err := s.GetTransaction(txHash)
	if err != nil {
		return nil, err
	}

	receipts, err := s.GetReceiptsByHash(blockHash)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve receipts for transaction. Cause: %w", err)
	}

	if len(receipts) <= int(index) {
		return nil, fmt.Errorf("receipt index not matching the transactions in block: %s", blockHash.Hex())
	}
	receipt := receipts[index]

	return receipt, nil
}

func (s *storageImpl) FetchAttestedKey(aggregator gethcommon.Address) (*ecdsa.PublicKey, error) {
	return obscurorawdb.ReadAttestationKey(s.db, aggregator)
}

func (s *storageImpl) StoreAttestedKey(aggregator gethcommon.Address, key *ecdsa.PublicKey) error {
	return obscurorawdb.WriteAttestationKey(s.db, aggregator, key)
}

func (s *storageImpl) StoreBatch(batch *core.Batch, receipts []*types.Receipt) error {
	dbBatch := s.db.NewBatch()

	if err := obscurorawdb.WriteBatch(dbBatch, batch); err != nil {
		return fmt.Errorf("could not write batch. Cause: %w", err)
	}
	if err := obscurorawdb.WriteTxLookupEntriesByBatch(dbBatch, batch); err != nil {
		return fmt.Errorf("could not write transaction lookup entries by batch. Cause: %w", err)
	}
	if err := obscurorawdb.WriteReceipts(dbBatch, *batch.Hash(), receipts); err != nil {
		return fmt.Errorf("could not write transaction receipts. Cause: %w", err)
	}
	if err := obscurorawdb.WriteContractCreationTxs(dbBatch, receipts); err != nil {
		return fmt.Errorf("could not save contract creation transaction. Cause: %w", err)
	}

	if err := dbBatch.Write(); err != nil {
		return fmt.Errorf("could not write batch to storage. Cause: %w", err)
	}
	return nil
}

func (s *storageImpl) StoreL1Messages(blockHash common.L1RootHash, messages common.CrossChainMessages) error {
	return obscurorawdb.StoreL1Messages(s.db, blockHash, messages, s.logger)
}

func (s *storageImpl) GetL1Messages(blockHash common.L1RootHash) (common.CrossChainMessages, error) {
	return obscurorawdb.GetL1Messages(s.db, blockHash, s.logger)
}

func (s *storageImpl) StoreRollup(rollup *core.Rollup) error {
	dbBatch := s.db.NewBatch()

	if err := obscurorawdb.WriteRollup(dbBatch, rollup); err != nil {
		return fmt.Errorf("could not write rollup. Cause: %w", err)
	}

	if err := dbBatch.Write(); err != nil {
		return fmt.Errorf("could not write rollup to storage. Cause: %w", err)
	}
	return nil
}
