package storage

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/eko/gocache/lib/v4/cache"
	bigcache_store "github.com/eko/gocache/store/bigcache/v4"

	"github.com/obscuronet/go-obscuro/go/config"

	"github.com/obscuronet/go-obscuro/go/enclave/storage/enclavedb"

	"github.com/ethereum/go-ethereum/rlp"

	gethcore "github.com/ethereum/go-ethereum/core"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/trie"

	"github.com/obscuronet/go-obscuro/go/common/syserr"

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
)

// todo - this will require a dedicated table when updates are implemented
const (
	masterSeedCfg            = "MASTER_SEED"
	_slowCallThresholdMillis = 50 // requests that take longer than this will be logged
)

type storageImpl struct {
	db enclavedb.EnclaveDB

	// cache for the immutable batches and blocks.
	// this avoids a trip to the database.
	batchCache *cache.Cache[[]byte]
	blockCache *cache.Cache[[]byte]

	stateDB     state.Database
	chainConfig *params.ChainConfig
	logger      gethlog.Logger
}

func NewStorageFromConfig(config *config.EnclaveConfig, chainConfig *params.ChainConfig, logger gethlog.Logger) Storage {
	backingDB, err := CreateDBFromConfig(config, logger)
	if err != nil {
		logger.Crit("Failed to connect to backing database", log.ErrKey, err)
	}
	return NewStorage(backingDB, chainConfig, logger)
}

func NewStorage(backingDB enclavedb.EnclaveDB, chainConfig *params.ChainConfig, logger gethlog.Logger) Storage {
	cacheConfig := &gethcore.CacheConfig{
		TrieCleanLimit: 256,
		TrieDirtyLimit: 256,
		TrieTimeLimit:  5 * time.Minute,
		SnapshotLimit:  256,
		SnapshotWait:   true,
	}

	// todo (tudor) figure out the context and the config
	bigcacheClient, err := bigcache.New(context.Background(), bigcache.DefaultConfig(5*time.Minute))
	if err != nil {
		logger.Crit("Could not initialise bigcache", log.ErrKey, err)
	}

	bigcacheStore := bigcache_store.NewBigcache(bigcacheClient)

	return &storageImpl{
		db: backingDB,
		stateDB: state.NewDatabaseWithConfig(backingDB, &trie.Config{
			Cache:     cacheConfig.TrieCleanLimit,
			Preimages: cacheConfig.Preimages,
		}),
		chainConfig: chainConfig,
		batchCache:  cache.New[[]byte](bigcacheStore),
		blockCache:  cache.New[[]byte](bigcacheStore),
		logger:      logger,
	}
}

func (s *storageImpl) TrieDB() *trie.Database {
	return s.stateDB.TrieDB()
}

func (s *storageImpl) Close() error {
	return s.db.GetSQLDB().Close()
}

func (s *storageImpl) FetchHeadBatch() (*core.Batch, error) {
	callStart := time.Now()
	defer s.logDuration("FetchHeadBatch", callStart)
	return enclavedb.ReadCurrentHeadBatch(s.db.GetSQLDB())
}

func (s *storageImpl) FetchCurrentSequencerNo() (*big.Int, error) {
	callStart := time.Now()
	defer s.logDuration("FetchCurrentSequencerNo", callStart)
	return enclavedb.ReadCurrentSequencerNo(s.db.GetSQLDB())
}

func (s *storageImpl) FetchBatch(hash common.L2BatchHash) (*core.Batch, error) {
	callStart := time.Now()
	defer s.logDuration("FetchBatch", callStart)
	return getCachedValue(s.batchCache, s.logger, hash, func(v any) (*core.Batch, error) {
		return enclavedb.ReadBatchByHash(s.db.GetSQLDB(), v.(common.L2BatchHash))
	})
}

func (s *storageImpl) FetchBatchHeader(hash common.L2BatchHash) (*common.BatchHeader, error) {
	callStart := time.Now()
	defer s.logDuration("FetchBatchHeader", callStart)
	b, err := s.FetchBatch(hash)
	if err != nil {
		return nil, err
	}
	return b.Header, nil
}

func (s *storageImpl) FetchBatchByHeight(height uint64) (*core.Batch, error) {
	callStart := time.Now()
	defer s.logDuration("FetchBatchByHeight", callStart)
	return enclavedb.ReadCanonicalBatchByHeight(s.db.GetSQLDB(), height)
}

func (s *storageImpl) StoreBlock(b *types.Block, chainFork *common.ChainFork) error {
	callStart := time.Now()
	defer s.logDuration("StoreBlock", callStart)
	dbTransaction := s.db.NewDBTransaction()
	if chainFork != nil && chainFork.IsFork() {
		s.logger.Info(fmt.Sprintf("Fork. %s", chainFork))
		enclavedb.UpdateCanonicalBlocks(dbTransaction, chainFork.CanonicalPath, chainFork.NonCanonicalPath)
	}

	// In case there were any batches inserted before this block was received
	enclavedb.UpdateCanonicalBlocks(dbTransaction, []common.L1BlockHash{b.Hash()}, nil)

	if err := enclavedb.WriteBlock(dbTransaction, b.Header()); err != nil {
		return fmt.Errorf("2. could not store block %s. Cause: %w", b.Hash(), err)
	}

	if err := dbTransaction.Write(); err != nil {
		return fmt.Errorf("3. could not store block %s. Cause: %w", b.Hash(), err)
	}

	cacheValue(s.blockCache, s.logger, b.Hash(), b)

	return nil
}

func (s *storageImpl) FetchBlock(blockHash common.L1BlockHash) (*types.Block, error) {
	callStart := time.Now()
	defer s.logDuration("FetchBlock", callStart)
	return getCachedValue(s.blockCache, s.logger, blockHash, func(hash any) (*types.Block, error) {
		return enclavedb.FetchBlock(s.db.GetSQLDB(), hash.(common.L1BlockHash))
	})
}

func (s *storageImpl) FetchCanonicaBlockByHeight(height *big.Int) (*types.Block, error) {
	callStart := time.Now()
	defer s.logDuration("FetchCanonicaBlockByHeight", callStart)
	header, err := enclavedb.FetchBlockHeaderByHeight(s.db.GetSQLDB(), height)
	if err != nil {
		return nil, err
	}
	blockHash := header.Hash()
	return getCachedValue(s.blockCache, s.logger, blockHash, func(hash any) (*types.Block, error) {
		return enclavedb.FetchBlock(s.db.GetSQLDB(), hash.(common.L2BatchHash))
	})
}

func (s *storageImpl) FetchHeadBlock() (*types.Block, error) {
	callStart := time.Now()
	defer s.logDuration("FetchHeadBlock", callStart)
	return enclavedb.FetchHeadBlock(s.db.GetSQLDB())
}

func (s *storageImpl) StoreSecret(secret crypto.SharedEnclaveSecret) error {
	callStart := time.Now()
	defer s.logDuration("StoreSecret", callStart)
	enc, err := rlp.EncodeToBytes(secret)
	if err != nil {
		return fmt.Errorf("could not encode shared secret. Cause: %w", err)
	}
	_, err = enclavedb.WriteConfig(s.db.GetSQLDB(), masterSeedCfg, enc)
	if err != nil {
		return fmt.Errorf("could not shared secret in DB. Cause: %w", err)
	}
	return nil
}

func (s *storageImpl) FetchSecret() (*crypto.SharedEnclaveSecret, error) {
	callStart := time.Now()
	defer s.logDuration("FetchSecret", callStart)
	var ss crypto.SharedEnclaveSecret

	cfg, err := enclavedb.FetchConfig(s.db.GetSQLDB(), masterSeedCfg)
	if err != nil {
		return nil, err
	}
	if err := rlp.DecodeBytes(cfg, &ss); err != nil {
		return nil, fmt.Errorf("could not decode shared secret")
	}

	return &ss, nil
}

func (s *storageImpl) IsAncestor(block *types.Block, maybeAncestor *types.Block) bool {
	callStart := time.Now()
	defer s.logDuration("IsAncestor", callStart)
	if bytes.Equal(maybeAncestor.Hash().Bytes(), block.Hash().Bytes()) {
		return true
	}

	if maybeAncestor.NumberU64() >= block.NumberU64() {
		return false
	}

	p, err := s.FetchBlock(block.ParentHash())
	if err != nil {
		s.logger.Debug("Could not find block with hash", log.BlockHashKey, block.ParentHash(), log.ErrKey, err)
		return false
	}

	return s.IsAncestor(p, maybeAncestor)
}

func (s *storageImpl) IsBlockAncestor(block *types.Block, maybeAncestor common.L1BlockHash) bool {
	callStart := time.Now()
	defer s.logDuration("IsBlockAncestor", callStart)
	resolvedBlock, err := s.FetchBlock(maybeAncestor)
	if err != nil {
		return false
	}
	return s.IsAncestor(block, resolvedBlock)
}

func (s *storageImpl) HealthCheck() (bool, error) {
	callStart := time.Now()
	defer s.logDuration("HealthCheck", callStart)
	headBatch, err := s.FetchHeadBatch()
	if err != nil {
		s.logger.Info("HealthCheck failed for enclave storage", log.ErrKey, err)
		return false, err
	}
	return headBatch != nil, nil
}

func (s *storageImpl) FetchHeadBatchForBlock(blockHash common.L1BlockHash) (*core.Batch, error) {
	callStart := time.Now()
	defer s.logDuration("FetchHeadBatchForBlock", callStart)
	return enclavedb.ReadHeadBatchForBlock(s.db.GetSQLDB(), blockHash)
}

func (s *storageImpl) CreateStateDB(hash common.L2BatchHash) (*state.StateDB, error) {
	callStart := time.Now()
	defer s.logDuration("CreateStateDB", callStart)
	batch, err := s.FetchBatch(hash)
	if err != nil {
		return nil, err
	}

	statedb, err := state.New(batch.Header.Root, s.stateDB, nil)
	if err != nil {
		return nil, syserr.NewInternalError(fmt.Errorf("could not create state DB for %s. Cause: %w", batch.Header.Root, err))
	}

	return statedb, nil
}

func (s *storageImpl) EmptyStateDB() (*state.StateDB, error) {
	callStart := time.Now()
	defer s.logDuration("EmptyStateDB", callStart)
	statedb, err := state.New(types.EmptyRootHash, s.stateDB, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create state DB. Cause: %w", err)
	}
	return statedb, nil
}

// GetReceiptsByBatchHash retrieves the receipts for all transactions in a given batch.
func (s *storageImpl) GetReceiptsByBatchHash(hash gethcommon.Hash) (types.Receipts, error) {
	callStart := time.Now()
	defer s.logDuration("GetReceiptsByBatchHash", callStart)
	return enclavedb.ReadReceiptsByBatchHash(s.db.GetSQLDB(), hash, s.chainConfig)
}

func (s *storageImpl) GetTransaction(txHash gethcommon.Hash) (*types.Transaction, gethcommon.Hash, uint64, uint64, error) {
	callStart := time.Now()
	defer s.logDuration("GetTransaction", callStart)
	return enclavedb.ReadTransaction(s.db.GetSQLDB(), txHash)
}

func (s *storageImpl) GetContractCreationTx(address gethcommon.Address) (*gethcommon.Hash, error) {
	callStart := time.Now()
	defer s.logDuration("GetContractCreationTx", callStart)
	return enclavedb.GetContractCreationTx(s.db.GetSQLDB(), address)
}

func (s *storageImpl) GetTransactionReceipt(txHash gethcommon.Hash) (*types.Receipt, error) {
	callStart := time.Now()
	defer s.logDuration("GetTransactionReceipt", callStart)
	return enclavedb.ReadReceipt(s.db.GetSQLDB(), txHash, s.chainConfig)
}

func (s *storageImpl) FetchAttestedKey(address gethcommon.Address) (*ecdsa.PublicKey, error) {
	callStart := time.Now()
	defer s.logDuration("FetchAttestedKey", callStart)
	key, err := enclavedb.FetchAttKey(s.db.GetSQLDB(), address)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve attestation key for address %s. Cause: %w", address, err)
	}

	publicKey, err := gethcrypto.DecompressPubkey(key)
	if err != nil {
		return nil, fmt.Errorf("could not parse key from db. Cause: %w", err)
	}

	return publicKey, nil
}

func (s *storageImpl) StoreAttestedKey(aggregator gethcommon.Address, key *ecdsa.PublicKey) error {
	callStart := time.Now()
	defer s.logDuration("StoreAttestedKey", callStart)
	_, err := enclavedb.WriteAttKey(s.db.GetSQLDB(), aggregator, gethcrypto.CompressPubkey(key))
	return err
}

func (s *storageImpl) FetchBatchBySeqNo(seqNum uint64) (*core.Batch, error) {
	callStart := time.Now()
	defer s.logDuration("FetchBatchBySeqNo", callStart)
	return getCachedValue(s.batchCache, s.logger, seqNum, func(seq any) (*core.Batch, error) {
		return enclavedb.ReadBatchBySeqNo(s.db.GetSQLDB(), seq.(uint64))
	})
}

func (s *storageImpl) FetchBatchesByBlock(block common.L1BlockHash) ([]*core.Batch, error) {
	callStart := time.Now()
	defer s.logDuration("FetchBatchesByBlock", callStart)
	return enclavedb.ReadBatchesByBlock(s.db.GetSQLDB(), block)
}

func (s *storageImpl) StoreBatch(batch *core.Batch) error {
	callStart := time.Now()
	defer s.logDuration("StoreBatch", callStart)
	// sanity check that this is not overlapping
	existingBatchWithSameSequence, _ := s.FetchBatchBySeqNo(batch.SeqNo().Uint64())
	if existingBatchWithSameSequence != nil && existingBatchWithSameSequence.Hash() != batch.Hash() {
		// todo - tudor - remove the Critical before production, and return a challenge
		s.logger.Crit(fmt.Sprintf("Conflicting batches for the same sequence %d: (previous) %+v != (incoming) %+v", batch.SeqNo(), existingBatchWithSameSequence.Header, batch.Header))
		return fmt.Errorf("a different batch with same sequence number already exists: %d", batch.SeqNo())
	}

	// already processed batch with this seq number and hash
	if existingBatchWithSameSequence != nil && existingBatchWithSameSequence.Hash() == batch.Hash() {
		return nil
	}

	dbTx := s.db.NewDBTransaction()
	s.logger.Trace("write batch", log.BatchHashKey, batch.Hash(), "l1Proof", batch.Header.L1Proof, log.BatchSeqNoKey, batch.SeqNo())
	if err := enclavedb.WriteBatchAndTransactions(dbTx, batch); err != nil {
		return fmt.Errorf("could not write batch. Cause: %w", err)
	}

	if err := dbTx.Write(); err != nil {
		return fmt.Errorf("could not commit batch %w", err)
	}

	cacheValue(s.batchCache, s.logger, batch.Hash(), batch)
	return nil
}

func (s *storageImpl) StoreExecutedBatch(batch *core.Batch, receipts []*types.Receipt) error {
	callStart := time.Now()
	defer s.logDuration("StoreExecutedBatch", callStart)
	executed, err := enclavedb.BatchWasExecuted(s.db.GetSQLDB(), batch.Hash())
	if err != nil {
		return err
	}
	if executed {
		s.logger.Debug("Batch was already executed", log.BatchHashKey, batch.Hash())
		return nil
	}

	dbTx := s.db.NewDBTransaction()
	if err := enclavedb.WriteBatchExecution(dbTx, batch.SeqNo(), receipts); err != nil {
		return fmt.Errorf("could not write transaction receipts. Cause: %w", err)
	}

	if batch.Number().Int64() > 1 {
		stateDB, err := s.CreateStateDB(batch.Header.ParentHash)
		if err != nil {
			return fmt.Errorf("could not create state DB to filter logs. Cause: %w", err)
		}

		err = enclavedb.StoreEventLogs(dbTx, receipts, stateDB)
		if err != nil {
			return fmt.Errorf("could not save logs %w", err)
		}
	}

	if err = dbTx.Write(); err != nil {
		return fmt.Errorf("could not commit batch %w", err)
	}

	return nil
}

func (s *storageImpl) StoreValueTransfers(blockHash common.L1BlockHash, transfers common.ValueTransferEvents) error {
	return enclavedb.WriteL1Messages(s.db.GetSQLDB(), blockHash, transfers, true)
}

func (s *storageImpl) StoreL1Messages(blockHash common.L1BlockHash, messages common.CrossChainMessages) error {
	callStart := time.Now()
	defer s.logDuration("StoreL1Messages", callStart)
	return enclavedb.WriteL1Messages(s.db.GetSQLDB(), blockHash, messages, false)
}

func (s *storageImpl) GetL1Messages(blockHash common.L1BlockHash) (common.CrossChainMessages, error) {
	callStart := time.Now()
	defer s.logDuration("GetL1Messages", callStart)
	return enclavedb.FetchL1Messages[common.CrossChainMessage](s.db.GetSQLDB(), blockHash, false)
}

func (s *storageImpl) GetL1Transfers(blockHash common.L1BlockHash) (common.ValueTransferEvents, error) {
	return enclavedb.FetchL1Messages[common.ValueTransferEvent](s.db.GetSQLDB(), blockHash, true)
}

const enclaveKeyKey = "ek"

func (s *storageImpl) StoreEnclaveKey(enclaveKey *ecdsa.PrivateKey) error {
	callStart := time.Now()
	defer s.logDuration("StoreEnclaveKey", callStart)
	if enclaveKey == nil {
		return errors.New("enclaveKey cannot be nil")
	}
	keyBytes := gethcrypto.FromECDSA(enclaveKey)

	_, err := enclavedb.WriteConfig(s.db.GetSQLDB(), enclaveKeyKey, keyBytes)
	return err
}

func (s *storageImpl) GetEnclaveKey() (*ecdsa.PrivateKey, error) {
	callStart := time.Now()
	defer s.logDuration("GetEnclaveKey", callStart)
	keyBytes, err := enclavedb.FetchConfig(s.db.GetSQLDB(), enclaveKeyKey)
	if err != nil {
		return nil, err
	}
	enclaveKey, err := gethcrypto.ToECDSA(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to construct ECDSA private key from enclave key bytes - %w", err)
	}
	return enclaveKey, nil
}

func (s *storageImpl) StoreRollup(rollup *common.ExtRollup, internalHeader *common.CalldataRollupHeader) error {
	callStart := time.Now()
	defer s.logDuration("StoreRollup", callStart)
	dbBatch := s.db.NewDBTransaction()

	if err := enclavedb.WriteRollup(dbBatch, rollup.Header, internalHeader); err != nil {
		return fmt.Errorf("could not write rollup. Cause: %w", err)
	}

	if err := dbBatch.Write(); err != nil {
		return fmt.Errorf("could not write rollup to storage. Cause: %w", err)
	}
	return nil
}

func (s *storageImpl) FetchReorgedRollup(reorgedBlocks []common.L1BlockHash) (*common.L2BatchHash, error) {
	return enclavedb.FetchReorgedRollup(s.db.GetSQLDB(), reorgedBlocks)
}

func (s *storageImpl) DebugGetLogs(txHash common.TxHash) ([]*tracers.DebugLogs, error) {
	callStart := time.Now()
	defer s.logDuration("DebugGetLogs", callStart)
	return enclavedb.DebugGetLogs(s.db.GetSQLDB(), txHash)
}

func (s *storageImpl) FilterLogs(
	requestingAccount *gethcommon.Address,
	fromBlock, toBlock *big.Int,
	blockHash *common.L2BatchHash,
	addresses []gethcommon.Address,
	topics [][]gethcommon.Hash,
) ([]*types.Log, error) {
	callStart := time.Now()
	defer s.logDuration("FilterLogs", callStart)
	return enclavedb.FilterLogs(s.db.GetSQLDB(), requestingAccount, fromBlock, toBlock, blockHash, addresses, topics)
}

func (s *storageImpl) GetContractCount() (*big.Int, error) {
	callStart := time.Now()
	defer s.logDuration("GetContractCount", callStart)
	return enclavedb.ReadContractCreationCount(s.db.GetSQLDB())
}

func (s *storageImpl) FetchCanonicalUnexecutedBatches(from *big.Int) ([]*core.Batch, error) {
	callStart := time.Now()
	defer s.logDuration("FetchCanonicalUnexecutedBatches", callStart)
	return enclavedb.ReadUnexecutedBatches(s.db.GetSQLDB(), from)
}

func (s *storageImpl) BatchWasExecuted(hash common.L2BatchHash) (bool, error) {
	callStart := time.Now()
	defer s.logDuration("BatchWasExecuted", callStart)
	return enclavedb.BatchWasExecuted(s.db.GetSQLDB(), hash)
}

func (s *storageImpl) GetReceiptsPerAddress(address *gethcommon.Address, pagination *common.QueryPagination) (types.Receipts, error) {
	callStart := time.Now()
	defer s.logDuration("GetReceiptsPerAddress", callStart)
	return enclavedb.GetReceiptsPerAddress(s.db.GetSQLDB(), s.chainConfig, address, pagination)
}

func (s *storageImpl) GetReceiptsPerAddressCount(address *gethcommon.Address) (uint64, error) {
	callStart := time.Now()
	defer s.logDuration("GetReceiptsPerAddressCount", callStart)
	return enclavedb.GetReceiptsPerAddressCount(s.db.GetSQLDB(), address)
}

func (s *storageImpl) GetPublicTransactionData(pagination *common.QueryPagination) ([]common.PublicTransaction, error) {
	callStart := time.Now()
	defer s.logDuration("GetPublicTransactionData", callStart)
	return enclavedb.GetPublicTransactionData(s.db.GetSQLDB(), pagination)
}

func (s *storageImpl) GetPublicTransactionCount() (uint64, error) {
	callStart := time.Now()
	defer s.logDuration("GetPublicTransactionCount", callStart)
	return enclavedb.GetPublicTransactionCount(s.db.GetSQLDB())
}

func (s *storageImpl) logDuration(method string, callStart time.Time) {
	durationMillis := time.Since(callStart).Milliseconds()
	// we only log 'slow' calls to reduce noise
	if durationMillis > _slowCallThresholdMillis {
		s.logger.Info(fmt.Sprintf("Storage::%s completed", method), log.DurationMilliKey, durationMillis)
	}
}
