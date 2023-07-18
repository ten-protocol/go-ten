package db

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/rlp"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/go-obscuro/go/enclave/db/orm"

	gethcore "github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/trie"

	"github.com/obscuronet/go-obscuro/go/common/syserr"

	"github.com/obscuronet/go-obscuro/go/enclave/db/sql"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/tracers"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/crypto"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	obscurorawdb "github.com/obscuronet/go-obscuro/go/enclave/db/rawdb"
)

// todo - this will require a dedicated table when updates are implemented
const master_seed_cfg = "MASTER_SEED"

// ErrNoRollups is returned if no rollups have been published yet in the history of the network
// Note: this is not just "not found", we cache at every L1 block what rollup we are up to so we also record that we haven't seen one yet
var ErrNoRollups = errors.New("no rollups have been published")

// todo (#1551) - consistency around whether we assert the secret is available or not

type storageImpl struct {
	db          *sql.EnclaveDB
	stateDB     state.Database
	chainConfig *params.ChainConfig
	logger      gethlog.Logger
}

func NewStorage(backingDB *sql.EnclaveDB, chainConfig *params.ChainConfig, logger gethlog.Logger) Storage {
	cacheConfig := &gethcore.CacheConfig{
		TrieCleanLimit: 256,
		TrieDirtyLimit: 256,
		TrieTimeLimit:  5 * time.Minute,
		SnapshotLimit:  256,
		SnapshotWait:   true,
	}

	return &storageImpl{
		db: backingDB,
		stateDB: state.NewDatabaseWithConfig(backingDB, &trie.Config{
			Cache:     cacheConfig.TrieCleanLimit,
			Journal:   cacheConfig.TrieCleanJournal,
			Preimages: cacheConfig.Preimages,
		}),
		chainConfig: chainConfig,
		logger:      logger,
	}
}

func (s *storageImpl) TrieDB() *trie.Database {
	return s.stateDB.TrieDB()
}

func (s *storageImpl) OpenBatch() *sql.Batch {
	return s.db.NewSQLBatch()
}

func (s *storageImpl) CommitBatch(dbBatch *sql.Batch) error {
	return dbBatch.Write()
}

func (s *storageImpl) Close() error {
	return s.db.GetSQLDB().Close()
}

func (s *storageImpl) FetchHeadBatch() (*core.Batch, error) {
	return orm.FetchHeadBatch(s.db.GetSQLDB())
}

func (s *storageImpl) FetchCurrentSequencerNo() (*big.Int, error) {
	return orm.FetchCurrentSequencerNo(s.db.GetSQLDB())
}

func (s *storageImpl) FetchBatch(hash common.L2BatchHash) (*core.Batch, error) {
	return orm.FetchBatchByHash(s.db.GetSQLDB(), hash)
}

func (s *storageImpl) FetchBatchHeader(hash common.L2BatchHash) (*common.BatchHeader, error) {
	return orm.ReadBatchHeader(s.db.GetSQLDB(), hash)
}

func (s *storageImpl) FetchBatchByHeight(height uint64) (*core.Batch, error) {
	return orm.FetchCanonicalBatchByHeight(s.db.GetSQLDB(), height)
}

func (s *storageImpl) StoreBlock(b *types.Block, canonical []common.L1BlockHash, nonCanonical []common.L1BlockHash) error {
	// todo - update canonical
	dbBatch := s.db.NewSQLBatch()
	if err := orm.WriteBlock(dbBatch, b.Header()); err != nil {
		return fmt.Errorf("could not store block %s. Cause: %w", b.Hash(), err)
	}
	orm.UpdateCanonicalBlocks(dbBatch, canonical, nonCanonical)

	if err := dbBatch.Write(); err != nil {
		return fmt.Errorf("could not store block %s. Cause: %w", b.Hash(), err)
	}
	return nil
}

func (s *storageImpl) FetchBlock(blockHash common.L1BlockHash) (*types.Block, error) {
	return orm.FetchBlock(s.db.GetSQLDB(), blockHash)
}

func (s *storageImpl) FetchHeadBlock() (*types.Block, error) {
	return orm.FetchHeadBlock(s.db.GetSQLDB())
}

func (s *storageImpl) StoreSecret(secret crypto.SharedEnclaveSecret) error {
	enc, err := rlp.EncodeToBytes(secret)
	if err != nil {
		return fmt.Errorf("could not encode shared secret. Cause: %w", err)
	}
	_, err = orm.WriteConfig(s.db.GetSQLDB(), master_seed_cfg, enc)
	if err != nil {
		return fmt.Errorf("could not shared secret in DB. Cause: %w", err)
	}
	return nil
}

func (s *storageImpl) FetchSecret() (*crypto.SharedEnclaveSecret, error) {
	var ss crypto.SharedEnclaveSecret

	cfg, err := orm.FetchConfig(s.db.GetSQLDB(), master_seed_cfg)
	if err != nil {
		return nil, err
	}
	if err := rlp.DecodeBytes(cfg, &ss); err != nil {
		return nil, fmt.Errorf("could not decode shared secret")
	}

	return &ss, nil
}

func (s *storageImpl) IsAncestor(block *types.Block, maybeAncestor *types.Block) bool {
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

func (s *storageImpl) IsBlockAncestor(block *types.Block, maybeAncestor common.L1BlockHash) bool {
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

func (s *storageImpl) FetchHeadBatchForBlock(blockHash common.L1BlockHash) (*core.Batch, error) {
	return orm.FetchHeadBatchForBlock(s.db.GetSQLDB(), blockHash)
}

func (s *storageImpl) FetchHeadRollupForBlock(blockHash *common.L1BlockHash) (*common.RollupHeader, error) {
	l2HeadBatch, err := obscurorawdb.ReadL2HeadRollup(s.db, blockHash)
	if err != nil {
		return nil, fmt.Errorf("could not read L2 head rollup for block. Cause: %w", err)
	}
	if *l2HeadBatch == (gethcommon.Hash{}) { // empty hash ==> no rollups yet up to this block
		return nil, ErrNoRollups
	}
	return obscurorawdb.ReadRollupHeader(s.db, *l2HeadBatch)
}

func (s *storageImpl) UpdateHeadRollup(l1Head *common.L1BlockHash, l2Head *common.L2BatchHash) error {
	dbBatch := s.db.NewBatch()
	if err := obscurorawdb.WriteL2HeadRollup(dbBatch, l1Head, l2Head); err != nil {
		return fmt.Errorf("could not write block state. Cause: %w", err)
	}
	if err := dbBatch.Write(); err != nil {
		return fmt.Errorf("could not save new head. Cause: %w", err)
	}
	return nil
}

func (s *storageImpl) CreateStateDB(hash common.L2BatchHash) (*state.StateDB, error) {
	batch, err := s.FetchBatch(hash)
	if err != nil {
		return nil, err
	}

	statedb, err := state.New(batch.Header.Root, s.stateDB, nil)
	if err != nil {
		return nil, syserr.NewInternalError(fmt.Errorf("could not create state DB. Cause: %w", err))
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
func (s *storageImpl) GetReceiptsByBatchHash(hash gethcommon.Hash) (types.Receipts, error) {
	return orm.ReadReceipts(s.db.GetSQLDB(), hash, s.chainConfig)
}

func (s *storageImpl) GetTransaction(txHash gethcommon.Hash) (*types.Transaction, gethcommon.Hash, uint64, uint64, error) {
	return orm.ReadTransaction(s.db.GetSQLDB(), txHash)
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
	return orm.ReadReceipt(s.db.GetSQLDB(), txHash, s.chainConfig)
}

func (s *storageImpl) FetchAttestedKey(address gethcommon.Address) (*ecdsa.PublicKey, error) {
	key, err := orm.FetchAttKey(s.db.GetSQLDB(), address)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve attestation key for address %s. Cause: %w", address, err)
	}

	publicKey, err := ethcrypto.DecompressPubkey(key)
	if err != nil {
		return nil, fmt.Errorf("could not parse key from db. Cause: %w", err)
	}

	return publicKey, nil
}

func (s *storageImpl) StoreAttestedKey(aggregator gethcommon.Address, key *ecdsa.PublicKey) error {
	_, err := orm.WriteAttKey(s.db.GetSQLDB(), aggregator, ethcrypto.CompressPubkey(key))
	return err
}

func (s *storageImpl) FetchBatchBySeqNo(seqNum uint64) (*core.Batch, error) {
	return orm.FindBatchBySeqNo(s.db.GetSQLDB(), seqNum)
}

func (s *storageImpl) StoreBatch(batch *core.Batch, receipts []*types.Receipt, dbBatch *sql.Batch) error {
	if dbBatch == nil {
		panic("StoreBatch called without an instance of sql.Batch")
	}

	if _, err := s.FetchBatchBySeqNo(batch.SeqNo().Uint64()); err == nil {
		return nil
		// return fmt.Errorf("batch with same sequence number already exists: %d", batch.SeqNo())
	}

	s.logger.Trace("write batch", "hash", batch.Hash(), "l1_proof", batch.Header.L1Proof)
	if err := orm.WriteBatch(dbBatch, batch); err != nil {
		return fmt.Errorf("could not write batch. Cause: %w", err)
	}

	for _, receipt := range receipts {
		s.logger.Trace("store receipt", "txHash", receipt.TxHash, "batch", receipt.BlockHash)
	}
	if err := orm.WriteReceipts(dbBatch, receipts); err != nil {
		return fmt.Errorf("could not write transaction receipts. Cause: %w", err)
	}

	if batch.Number().Int64() > 1 {
		stateDB, err := s.CreateStateDB(batch.Header.ParentHash)
		if err != nil {
			return fmt.Errorf("could not create state DB to filter logs. Cause: %w", err)
		}

		err2 := orm.StoreEventLogs(dbBatch, receipts, stateDB)
		if err2 != nil {
			return fmt.Errorf("could not save logs %w", err2)
		}
	}

	if err := obscurorawdb.WriteContractCreationTxs(dbBatch, receipts); err != nil {
		return fmt.Errorf("could not save contract creation transaction. Cause: %w", err)
	}
	// todo fix this as batches always stored even if not canonical
	if err := obscurorawdb.IncrementContractCreationCount(s.db, dbBatch, receipts); err != nil {
		return fmt.Errorf("unable to increment contract count")
	}
	return nil
}

func (s *storageImpl) StoreL1Messages(blockHash common.L1BlockHash, messages common.CrossChainMessages) error {
	return orm.WriteL1Messages(s.db.GetSQLDB(), blockHash, messages)
}

func (s *storageImpl) GetL1Messages(blockHash common.L1BlockHash) (common.CrossChainMessages, error) {
	return orm.FetchL1Messages(s.db.GetSQLDB(), blockHash)
}

func (s *storageImpl) StoreEnclaveKey(enclaveKey *ecdsa.PrivateKey) error {
	return obscurorawdb.StoreEnclaveKey(s.db, enclaveKey, s.logger)
}

func (s *storageImpl) GetEnclaveKey() (*ecdsa.PrivateKey, error) {
	return obscurorawdb.GetEnclaveKey(s.db, s.logger)
}

func (s *storageImpl) StoreRollup(rollup *common.ExtRollup) error {
	dbBatch := s.db.NewBatch()

	if err := obscurorawdb.WriteRollup(dbBatch, rollup); err != nil {
		return fmt.Errorf("could not write rollup. Cause: %w", err)
	}

	if err := dbBatch.Write(); err != nil {
		return fmt.Errorf("could not write rollup to storage. Cause: %w", err)
	}
	return nil
}

func (s *storageImpl) DebugGetLogs(txHash common.TxHash) ([]*tracers.DebugLogs, error) {
	return orm.DebugGetLogs(s.db.GetSQLDB(), txHash)
}

func (s *storageImpl) FilterLogs(
	requestingAccount *gethcommon.Address,
	fromBlock, toBlock *big.Int,
	blockHash *common.L2BatchHash,
	addresses []gethcommon.Address,
	topics [][]gethcommon.Hash,
) ([]*types.Log, error) {
	return orm.FilterLogs(s.db.GetSQLDB(), requestingAccount, fromBlock, toBlock, blockHash, addresses, topics)
}

func (s *storageImpl) GetContractCount() (*big.Int, error) {
	return obscurorawdb.ReadContractCreationCount(s.db)
}
