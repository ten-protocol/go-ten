package storage

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/triedb/hashdb"

	"github.com/ethereum/go-ethereum/triedb"

	"github.com/dgraph-io/ristretto"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/ten-protocol/go-ten/go/common/measure"

	ristretto_store "github.com/eko/gocache/store/ristretto/v4"

	"github.com/ten-protocol/go-ten/go/config"

	"github.com/ten-protocol/go-ten/go/enclave/storage/enclavedb"

	"github.com/ethereum/go-ethereum/rlp"

	gethcore "github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/tracers"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"github.com/ten-protocol/go-ten/go/enclave/crypto"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
)

// todo - this will require a dedicated table when upgrades are implemented
const (
	masterSeedCfg = "MASTER_SEED"
)

type storageImpl struct {
	db enclavedb.EnclaveDB

	// cache for the immutable blocks and batches.
	// this avoids a trip to the database.
	blockCache *cache.Cache[*types.Block]

	// stores batches using the sequence number as key
	batchCacheBySeqNo *cache.Cache[*core.Batch]

	// mapping between the hash and the sequence number
	// note:  to fetch a batch by hash will require 2 cache hits
	seqCacheByHash *cache.Cache[*big.Int]

	// mapping between the height and the sequence number
	// note: to fetch a batch by height will require 2 cache hits
	seqCacheByHeight *cache.Cache[*big.Int]

	cachedSharedSecret *crypto.SharedEnclaveSecret

	stateCache  state.Database
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

var defaultCacheConfig = &gethcore.CacheConfig{
	TrieCleanLimit: 256,
	TrieDirtyLimit: 256,
	TrieTimeLimit:  5 * time.Minute,
	SnapshotLimit:  256,
	SnapshotWait:   true,
	StateScheme:    rawdb.HashScheme,
}

var trieDBConfig = &triedb.Config{
	Preimages: defaultCacheConfig.Preimages,
	IsVerkle:  false,
	HashDB: &hashdb.Config{
		CleanCacheSize: defaultCacheConfig.TrieCleanLimit * 1024 * 1024,
	},
}

func NewStorage(backingDB enclavedb.EnclaveDB, chainConfig *params.ChainConfig, logger gethlog.Logger) Storage {
	// Open trie database with provided config
	triedb := triedb.NewDatabase(backingDB, trieDBConfig)

	stateDB := state.NewDatabaseWithNodeDB(backingDB, triedb)

	// todo (tudor) figure out the config
	ristrettoCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 20_000, // 10*MaxCost
		MaxCost:     2000,   // - how many items to cache
		BufferItems: 64,     // number of keys per Get buffer.
	})
	if err != nil {
		logger.Crit("Could not initialise ristretto cache", log.ErrKey, err)
	}
	ristrettoStore := ristretto_store.NewRistretto(ristrettoCache)
	return &storageImpl{
		db:                backingDB,
		stateCache:        stateDB,
		chainConfig:       chainConfig,
		blockCache:        cache.New[*types.Block](ristrettoStore),
		batchCacheBySeqNo: cache.New[*core.Batch](ristrettoStore),
		seqCacheByHash:    cache.New[*big.Int](ristrettoStore),
		seqCacheByHeight:  cache.New[*big.Int](ristrettoStore),
		logger:            logger,
	}
}

func (s *storageImpl) TrieDB() *triedb.Database {
	return s.stateCache.TrieDB()
}

func (s *storageImpl) StateDB() state.Database {
	return s.stateCache
}

func (s *storageImpl) Close() error {
	return s.db.GetSQLDB().Close()
}

func (s *storageImpl) FetchHeadBatch(ctx context.Context) (*core.Batch, error) {
	defer s.logDuration("FetchHeadBatch", measure.NewStopwatch())
	return enclavedb.ReadCurrentHeadBatch(ctx, s.db.GetSQLDB())
}

func (s *storageImpl) FetchCurrentSequencerNo(ctx context.Context) (*big.Int, error) {
	defer s.logDuration("FetchCurrentSequencerNo", measure.NewStopwatch())
	return enclavedb.ReadCurrentSequencerNo(ctx, s.db.GetSQLDB())
}

func (s *storageImpl) FetchBatch(ctx context.Context, hash common.L2BatchHash) (*core.Batch, error) {
	defer s.logDuration("FetchBatch", measure.NewStopwatch())
	seqNo, err := common.GetCachedValue(ctx, s.seqCacheByHash, s.logger, hash, func(v any) (*big.Int, error) {
		batch, err := enclavedb.ReadBatchByHash(ctx, s.db.GetSQLDB(), v.(common.L2BatchHash))
		if err != nil {
			return nil, err
		}
		return batch.SeqNo(), nil
	})
	if err != nil {
		return nil, err
	}
	return s.FetchBatchBySeqNo(ctx, seqNo.Uint64())
}

func (s *storageImpl) FetchConvertedHash(ctx context.Context, hash common.L2BatchHash) (gethcommon.Hash, error) {
	defer s.logDuration("FetchConvertedHash", measure.NewStopwatch())
	batch, err := s.FetchBatch(ctx, hash)
	if err != nil {
		return gethcommon.Hash{}, err
	}
	return enclavedb.FetchConvertedBatchHash(ctx, s.db.GetSQLDB(), batch.Header.SequencerOrderNo.Uint64())
}

func (s *storageImpl) FetchBatchHeader(ctx context.Context, hash common.L2BatchHash) (*common.BatchHeader, error) {
	defer s.logDuration("FetchBatchHeader", measure.NewStopwatch())
	b, err := s.FetchBatch(ctx, hash)
	if err != nil {
		return nil, err
	}
	return b.Header, nil
}

func (s *storageImpl) FetchBatchByHeight(ctx context.Context, height uint64) (*core.Batch, error) {
	defer s.logDuration("FetchBatchByHeight", measure.NewStopwatch())
	// the key is (height+1), because for some reason it doesn't like a key of 0
	seqNo, err := common.GetCachedValue(ctx, s.seqCacheByHeight, s.logger, height+1, func(h any) (*big.Int, error) {
		batch, err := enclavedb.ReadCanonicalBatchByHeight(ctx, s.db.GetSQLDB(), height)
		if err != nil {
			return nil, err
		}
		return batch.SeqNo(), nil
	})
	if err != nil {
		return nil, err
	}
	return s.FetchBatchBySeqNo(ctx, seqNo.Uint64())
}

func (s *storageImpl) FetchNonCanonicalBatchesBetween(ctx context.Context, startSeq uint64, endSeq uint64) ([]*core.Batch, error) {
	defer s.logDuration("FetchNonCanonicalBatchesBetween", measure.NewStopwatch())
	return enclavedb.ReadNonCanonicalBatches(ctx, s.db.GetSQLDB(), startSeq, endSeq)
}

func (s *storageImpl) StoreBlock(ctx context.Context, b *types.Block, chainFork *common.ChainFork) error {
	defer s.logDuration("StoreBlock", measure.NewStopwatch())
	dbTx, err := s.db.NewDBTransaction(ctx)
	if err != nil {
		return fmt.Errorf("could not create DB transaction - %w", err)
	}
	defer dbTx.Rollback()

	if err := enclavedb.WriteBlock(ctx, dbTx, b.Header()); err != nil {
		return fmt.Errorf("2. could not store block %s. Cause: %w", b.Hash(), err)
	}

	blockId, err := enclavedb.GetBlockId(ctx, dbTx, b.Hash())
	if err != nil {
		return fmt.Errorf("could not get block id - %w", err)
	}

	// In case there were any batches inserted before this block was received
	err = enclavedb.SetMissingBlockId(ctx, dbTx, blockId, b.Hash())
	if err != nil {
		return err
	}

	if chainFork != nil && chainFork.IsFork() {
		s.logger.Info(fmt.Sprintf("Update Fork. %s", chainFork))
		err = enclavedb.UpdateCanonicalBlocks(ctx, dbTx, chainFork.CanonicalPath, chainFork.NonCanonicalPath, s.logger)
		if err != nil {
			return err
		}
	}

	err = enclavedb.UpdateCanonicalBlocks(ctx, dbTx, []common.L1BlockHash{b.Hash()}, nil, s.logger)
	if err != nil {
		return err
	}

	if err := dbTx.Commit(); err != nil {
		return fmt.Errorf("3. could not store block %s. Cause: %w", b.Hash(), err)
	}

	common.CacheValue(ctx, s.blockCache, s.logger, b.Hash(), b)

	return nil
}

func (s *storageImpl) FetchBlock(ctx context.Context, blockHash common.L1BlockHash) (*types.Block, error) {
	defer s.logDuration("FetchBlock", measure.NewStopwatch())
	return common.GetCachedValue(ctx, s.blockCache, s.logger, blockHash, func(hash any) (*types.Block, error) {
		return enclavedb.FetchBlock(ctx, s.db.GetSQLDB(), hash.(common.L1BlockHash))
	})
}

func (s *storageImpl) FetchCanonicaBlockByHeight(ctx context.Context, height *big.Int) (*types.Block, error) {
	defer s.logDuration("FetchCanonicaBlockByHeight", measure.NewStopwatch())
	header, err := enclavedb.FetchBlockHeaderByHeight(ctx, s.db.GetSQLDB(), height)
	if err != nil {
		return nil, err
	}
	return s.FetchBlock(ctx, header.Hash())
}

func (s *storageImpl) FetchHeadBlock(ctx context.Context) (*types.Block, error) {
	defer s.logDuration("FetchHeadBlock", measure.NewStopwatch())
	return enclavedb.FetchHeadBlock(ctx, s.db.GetSQLDB())
}

func (s *storageImpl) StoreSecret(ctx context.Context, secret crypto.SharedEnclaveSecret) error {
	defer s.logDuration("StoreSecret", measure.NewStopwatch())
	enc, err := rlp.EncodeToBytes(secret)
	if err != nil {
		return fmt.Errorf("could not encode shared secret. Cause: %w", err)
	}
	dbTx, err := s.db.NewDBTransaction(ctx)
	if err != nil {
		return fmt.Errorf("could not create DB transaction - %w", err)
	}
	defer dbTx.Rollback()
	_, err = enclavedb.WriteConfig(ctx, dbTx, masterSeedCfg, enc)
	if err != nil {
		return fmt.Errorf("could not shared secret in DB. Cause: %w", err)
	}
	err = dbTx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (s *storageImpl) FetchSecret(ctx context.Context) (*crypto.SharedEnclaveSecret, error) {
	defer s.logDuration("FetchSecret", measure.NewStopwatch())

	if s.cachedSharedSecret != nil {
		return s.cachedSharedSecret, nil
	}

	var ss crypto.SharedEnclaveSecret

	cfg, err := enclavedb.FetchConfig(ctx, s.db.GetSQLDB(), masterSeedCfg)
	if err != nil {
		return nil, err
	}
	if err := rlp.DecodeBytes(cfg, &ss); err != nil {
		return nil, fmt.Errorf("could not decode shared secret")
	}

	s.cachedSharedSecret = &ss
	return s.cachedSharedSecret, nil
}

func (s *storageImpl) IsAncestor(ctx context.Context, block *types.Block, maybeAncestor *types.Block) bool {
	defer s.logDuration("IsAncestor", measure.NewStopwatch())
	if bytes.Equal(maybeAncestor.Hash().Bytes(), block.Hash().Bytes()) {
		return true
	}

	if maybeAncestor.NumberU64() >= block.NumberU64() {
		return false
	}

	p, err := s.FetchBlock(ctx, block.ParentHash())
	if err != nil {
		s.logger.Debug("Could not find block with hash", log.BlockHashKey, block.ParentHash(), log.ErrKey, err)
		return false
	}

	return s.IsAncestor(ctx, p, maybeAncestor)
}

func (s *storageImpl) IsBlockAncestor(ctx context.Context, block *types.Block, maybeAncestor common.L1BlockHash) bool {
	defer s.logDuration("IsBlockAncestor", measure.NewStopwatch())
	resolvedBlock, err := s.FetchBlock(ctx, maybeAncestor)
	if err != nil {
		return false
	}
	return s.IsAncestor(ctx, block, resolvedBlock)
}

func (s *storageImpl) HealthCheck(ctx context.Context) (bool, error) {
	defer s.logDuration("HealthCheck", measure.NewStopwatch())
	headBatch, err := s.FetchHeadBatch(ctx)
	if err != nil {
		return false, err
	}

	if headBatch == nil {
		return false, fmt.Errorf("head batch is nil")
	}

	return true, nil
}

func (s *storageImpl) CreateStateDB(ctx context.Context, batchHash common.L2BatchHash) (*state.StateDB, error) {
	defer s.logDuration("CreateStateDB", measure.NewStopwatch())
	batch, err := s.FetchBatch(ctx, batchHash)
	if err != nil {
		return nil, err
	}

	statedb, err := state.New(batch.Header.Root, s.stateCache, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create state DB for %s. Cause: %w", batch.Header.Root, err)
	}
	return statedb, nil
}

func (s *storageImpl) EmptyStateDB() (*state.StateDB, error) {
	defer s.logDuration("EmptyStateDB", measure.NewStopwatch())
	statedb, err := state.New(types.EmptyRootHash, s.stateCache, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create state DB. Cause: %w", err)
	}
	return statedb, nil
}

// GetReceiptsByBatchHash retrieves the receipts for all transactions in a given batch.
func (s *storageImpl) GetReceiptsByBatchHash(ctx context.Context, hash gethcommon.Hash) (types.Receipts, error) {
	defer s.logDuration("GetReceiptsByBatchHash", measure.NewStopwatch())
	return enclavedb.ReadReceiptsByBatchHash(ctx, s.db.GetSQLDB(), hash, s.chainConfig)
}

func (s *storageImpl) GetTransaction(ctx context.Context, txHash gethcommon.Hash) (*types.Transaction, common.L2BatchHash, uint64, uint64, error) {
	defer s.logDuration("GetTransaction", measure.NewStopwatch())
	return enclavedb.ReadTransaction(ctx, s.db.GetSQLDB(), txHash)
}

func (s *storageImpl) GetContractCreationTx(ctx context.Context, address gethcommon.Address) (*gethcommon.Hash, error) {
	defer s.logDuration("GetContractCreationTx", measure.NewStopwatch())
	return enclavedb.GetContractCreationTx(ctx, s.db.GetSQLDB(), address)
}

func (s *storageImpl) GetTransactionReceipt(ctx context.Context, txHash gethcommon.Hash) (*types.Receipt, error) {
	defer s.logDuration("GetTransactionReceipt", measure.NewStopwatch())
	return enclavedb.ReadReceipt(ctx, s.db.GetSQLDB(), txHash, s.chainConfig)
}

func (s *storageImpl) FetchAttestedKey(ctx context.Context, address gethcommon.Address) (*ecdsa.PublicKey, error) {
	defer s.logDuration("FetchAttestedKey", measure.NewStopwatch())
	key, err := enclavedb.FetchAttKey(ctx, s.db.GetSQLDB(), address)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve attestation key for address %s. Cause: %w", address, err)
	}

	publicKey, err := gethcrypto.DecompressPubkey(key)
	if err != nil {
		return nil, fmt.Errorf("could not parse key from db. Cause: %w", err)
	}

	return publicKey, nil
}

func (s *storageImpl) StoreAttestedKey(ctx context.Context, aggregator gethcommon.Address, key *ecdsa.PublicKey) error {
	defer s.logDuration("StoreAttestedKey", measure.NewStopwatch())
	dbTx, err := s.db.NewDBTransaction(ctx)
	if err != nil {
		return fmt.Errorf("could not create DB transaction - %w", err)
	}
	defer dbTx.Rollback()
	_, err = enclavedb.WriteAttKey(ctx, dbTx, aggregator, gethcrypto.CompressPubkey(key))
	if err != nil {
		return err
	}
	err = dbTx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (s *storageImpl) FetchBatchBySeqNo(ctx context.Context, seqNum uint64) (*core.Batch, error) {
	defer s.logDuration("FetchBatchBySeqNo", measure.NewStopwatch())
	b, err := common.GetCachedValue(ctx, s.batchCacheBySeqNo, s.logger, seqNum, func(seq any) (*core.Batch, error) {
		return enclavedb.ReadBatchBySeqNo(ctx, s.db.GetSQLDB(), seqNum)
	})
	if err == nil && b == nil {
		return nil, fmt.Errorf("not found")
	}
	return b, err
}

func (s *storageImpl) FetchBatchesByBlock(ctx context.Context, block common.L1BlockHash) ([]*core.Batch, error) {
	defer s.logDuration("FetchBatchesByBlock", measure.NewStopwatch())
	return enclavedb.ReadBatchesByBlock(ctx, s.db.GetSQLDB(), block)
}

func (s *storageImpl) StoreBatch(ctx context.Context, batch *core.Batch, convertedHash gethcommon.Hash) error {
	defer s.logDuration("StoreBatch", measure.NewStopwatch())
	// sanity check that this is not overlapping
	existingBatchWithSameSequence, _ := s.FetchBatchBySeqNo(ctx, batch.SeqNo().Uint64())
	if existingBatchWithSameSequence != nil && existingBatchWithSameSequence.Hash() != batch.Hash() {
		// todo - tudor - remove the Critical before production, and return a challenge
		s.logger.Crit(fmt.Sprintf("Conflicting batches for the same sequence %d: (previous) %+v != (incoming) %+v", batch.SeqNo(), existingBatchWithSameSequence.Header, batch.Header))
		return fmt.Errorf("a different batch with same sequence number already exists: %d", batch.SeqNo())
	}

	// already processed batch with this seq number and hash
	if existingBatchWithSameSequence != nil && existingBatchWithSameSequence.Hash() == batch.Hash() {
		return nil
	}

	dbTx, err := s.db.NewDBTransaction(ctx)
	if err != nil {
		return fmt.Errorf("could not create DB transaction - %w", err)
	}
	defer dbTx.Rollback()

	// it is possible that the block is not available if this is a validator
	blockId, err := enclavedb.GetBlockId(ctx, dbTx, batch.Header.L1Proof)
	if err != nil {
		s.logger.Warn("could not get block id from db", log.ErrKey, err)
	}
	s.logger.Trace("write batch", log.BatchHashKey, batch.Hash(), "l1Proof", batch.Header.L1Proof, log.BatchSeqNoKey, batch.SeqNo(), "block_id", blockId)

	if err := enclavedb.WriteBatchAndTransactions(ctx, dbTx, batch, convertedHash, blockId); err != nil {
		return fmt.Errorf("could not write batch. Cause: %w", err)
	}

	if err := dbTx.Commit(); err != nil {
		return fmt.Errorf("could not commit batch %w", err)
	}

	common.CacheValue(ctx, s.batchCacheBySeqNo, s.logger, batch.SeqNo().Uint64(), batch)
	common.CacheValue(ctx, s.seqCacheByHash, s.logger, batch.Hash(), batch.SeqNo())
	// note: the key is (height+1), because for some reason it doesn't like a key of 0
	// should always contain the canonical batch because the cache is overwritten by each new batch after a reorg
	common.CacheValue(ctx, s.seqCacheByHeight, s.logger, batch.NumberU64()+1, batch.SeqNo())
	return nil
}

func (s *storageImpl) StoreExecutedBatch(ctx context.Context, batch *core.Batch, receipts []*types.Receipt) error {
	defer s.logDuration("StoreExecutedBatch", measure.NewStopwatch())
	executed, err := enclavedb.BatchWasExecuted(ctx, s.db.GetSQLDB(), batch.Hash())
	if err != nil {
		return err
	}
	if executed {
		s.logger.Debug("Batch was already executed", log.BatchHashKey, batch.Hash())
		return nil
	}

	dbTx, err := s.db.NewDBTransaction(ctx)
	if err != nil {
		return fmt.Errorf("could not create DB transaction - %w", err)
	}
	defer dbTx.Rollback()
	if err := enclavedb.WriteBatchExecution(ctx, dbTx, batch.SeqNo(), receipts); err != nil {
		return fmt.Errorf("could not write transaction receipts. Cause: %w", err)
	}

	if batch.Number().Uint64() > common.L2GenesisSeqNo {
		stateDB, err := s.CreateStateDB(ctx, batch.Header.ParentHash)
		if err != nil {
			return fmt.Errorf("could not create state DB to filter logs. Cause: %w", err)
		}

		err = enclavedb.StoreEventLogs(ctx, dbTx, receipts, batch, stateDB)
		if err != nil {
			return fmt.Errorf("could not save logs %w", err)
		}
	}

	if err = dbTx.Commit(); err != nil {
		return fmt.Errorf("could not commit batch %w", err)
	}

	return nil
}

func (s *storageImpl) StoreValueTransfers(ctx context.Context, blockHash common.L1BlockHash, transfers common.ValueTransferEvents) error {
	dbtx, err := s.db.NewDBTransaction(ctx)
	if err != nil {
		return fmt.Errorf("could not create DB transaction - %w", err)
	}
	defer dbtx.Rollback()
	blockId, err := enclavedb.GetBlockId(ctx, dbtx, blockHash)
	if err != nil {
		return fmt.Errorf("could not get block id - %w", err)
	}
	err = enclavedb.WriteL1Messages(ctx, dbtx, blockId, transfers, true)
	if err != nil {
		return fmt.Errorf("could not write l1 messages - %w", err)
	}
	return dbtx.Commit()
}

func (s *storageImpl) StoreL1Messages(ctx context.Context, blockHash common.L1BlockHash, messages common.CrossChainMessages) error {
	defer s.logDuration("StoreL1Messages", measure.NewStopwatch())
	dbtx, err := s.db.NewDBTransaction(ctx)
	if err != nil {
		return fmt.Errorf("could not create DB transaction - %w", err)
	}
	defer dbtx.Rollback()
	blockId, err := enclavedb.GetBlockId(ctx, dbtx, blockHash)
	if err != nil {
		return fmt.Errorf("could not get block id - %w", err)
	}
	err = enclavedb.WriteL1Messages(ctx, dbtx, blockId, messages, false)
	if err != nil {
		return fmt.Errorf("could not write l1 messages - %w", err)
	}
	return dbtx.Commit()
}

func (s *storageImpl) GetL1Messages(ctx context.Context, blockHash common.L1BlockHash) (common.CrossChainMessages, error) {
	defer s.logDuration("GetL1Messages", measure.NewStopwatch())
	return enclavedb.FetchL1Messages[common.CrossChainMessage](ctx, s.db.GetSQLDB(), blockHash, false)
}

func (s *storageImpl) GetL1Transfers(ctx context.Context, blockHash common.L1BlockHash) (common.ValueTransferEvents, error) {
	return enclavedb.FetchL1Messages[common.ValueTransferEvent](ctx, s.db.GetSQLDB(), blockHash, true)
}

const enclaveKeyKey = "ek"

func (s *storageImpl) StoreEnclaveKey(ctx context.Context, enclaveKey *crypto.EnclaveKey) error {
	defer s.logDuration("StoreEnclaveKey", measure.NewStopwatch())
	if enclaveKey == nil {
		return errors.New("enclaveKey cannot be nil")
	}
	keyBytes := gethcrypto.FromECDSA(enclaveKey.PrivateKey())

	dbTx, err := s.db.NewDBTransaction(ctx)
	if err != nil {
		return fmt.Errorf("could not create DB transaction - %w", err)
	}
	defer dbTx.Rollback()
	_, err = enclavedb.WriteConfig(ctx, dbTx, enclaveKeyKey, keyBytes)
	if err != nil {
		return err
	}
	err = dbTx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (s *storageImpl) GetEnclaveKey(ctx context.Context) (*crypto.EnclaveKey, error) {
	defer s.logDuration("GetEnclaveKey", measure.NewStopwatch())
	keyBytes, err := enclavedb.FetchConfig(ctx, s.db.GetSQLDB(), enclaveKeyKey)
	if err != nil {
		return nil, err
	}
	ecdsaKey, err := gethcrypto.ToECDSA(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to construct ECDSA private key from enclave key bytes - %w", err)
	}
	return crypto.NewEnclaveKey(ecdsaKey), nil
}

func (s *storageImpl) StoreRollup(ctx context.Context, rollup *common.ExtRollup, internalHeader *common.CalldataRollupHeader) error {
	defer s.logDuration("StoreRollup", measure.NewStopwatch())

	dbTx, err := s.db.NewDBTransaction(ctx)
	if err != nil {
		return fmt.Errorf("could not create DB transaction - %w", err)
	}
	defer dbTx.Rollback()

	blockId, err := enclavedb.GetBlockId(ctx, dbTx, rollup.Header.CompressionL1Head)
	if err != nil {
		return fmt.Errorf("could not get block id - %w", err)
	}

	if err := enclavedb.WriteRollup(ctx, dbTx, rollup.Header, blockId, internalHeader); err != nil {
		return fmt.Errorf("could not write rollup. Cause: %w", err)
	}

	if err := dbTx.Commit(); err != nil {
		return fmt.Errorf("could not write rollup to storage. Cause: %w", err)
	}
	return nil
}

func (s *storageImpl) FetchReorgedRollup(ctx context.Context, reorgedBlocks []common.L1BlockHash) (*common.L2BatchHash, error) {
	return enclavedb.FetchReorgedRollup(ctx, s.db.GetSQLDB(), reorgedBlocks)
}

func (s *storageImpl) FetchRollupMetadata(ctx context.Context, hash common.L2RollupHash) (*common.PublicRollupMetadata, error) {
	return enclavedb.FetchRollupMetadata(ctx, s.db.GetSQLDB(), hash)
}

func (s *storageImpl) DebugGetLogs(ctx context.Context, txHash common.TxHash) ([]*tracers.DebugLogs, error) {
	defer s.logDuration("DebugGetLogs", measure.NewStopwatch())
	return enclavedb.DebugGetLogs(ctx, s.db.GetSQLDB(), txHash)
}

func (s *storageImpl) FilterLogs(
	ctx context.Context,
	requestingAccount *gethcommon.Address,
	fromBlock, toBlock *big.Int,
	blockHash *common.L2BatchHash,
	addresses []gethcommon.Address,
	topics [][]gethcommon.Hash,
) ([]*types.Log, error) {
	defer s.logDuration("FilterLogs", measure.NewStopwatch())
	return enclavedb.FilterLogs(ctx, s.db.GetSQLDB(), requestingAccount, fromBlock, toBlock, blockHash, addresses, topics)
}

func (s *storageImpl) GetContractCount(ctx context.Context) (*big.Int, error) {
	defer s.logDuration("GetContractCount", measure.NewStopwatch())
	return enclavedb.ReadContractCreationCount(ctx, s.db.GetSQLDB())
}

func (s *storageImpl) FetchCanonicalUnexecutedBatches(ctx context.Context, from *big.Int) ([]*core.Batch, error) {
	defer s.logDuration("FetchCanonicalUnexecutedBatches", measure.NewStopwatch())
	return enclavedb.ReadUnexecutedBatches(ctx, s.db.GetSQLDB(), from)
}

func (s *storageImpl) BatchWasExecuted(ctx context.Context, hash common.L2BatchHash) (bool, error) {
	defer s.logDuration("BatchWasExecuted", measure.NewStopwatch())
	return enclavedb.BatchWasExecuted(ctx, s.db.GetSQLDB(), hash)
}

func (s *storageImpl) GetReceiptsPerAddress(ctx context.Context, address *gethcommon.Address, pagination *common.QueryPagination) (types.Receipts, error) {
	defer s.logDuration("GetReceiptsPerAddress", measure.NewStopwatch())
	return enclavedb.GetReceiptsPerAddress(ctx, s.db.GetSQLDB(), s.chainConfig, address, pagination)
}

func (s *storageImpl) GetReceiptsPerAddressCount(ctx context.Context, address *gethcommon.Address) (uint64, error) {
	defer s.logDuration("GetReceiptsPerAddressCount", measure.NewStopwatch())
	return enclavedb.GetReceiptsPerAddressCount(ctx, s.db.GetSQLDB(), address)
}

func (s *storageImpl) GetPublicTransactionData(ctx context.Context, pagination *common.QueryPagination) ([]common.PublicTransaction, error) {
	defer s.logDuration("GetPublicTransactionData", measure.NewStopwatch())
	return enclavedb.GetPublicTransactionData(ctx, s.db.GetSQLDB(), pagination)
}

func (s *storageImpl) GetPublicTransactionCount(ctx context.Context) (uint64, error) {
	defer s.logDuration("GetPublicTransactionCount", measure.NewStopwatch())
	return enclavedb.GetPublicTransactionCount(ctx, s.db.GetSQLDB())
}

func (s *storageImpl) logDuration(method string, stopWatch *measure.Stopwatch) {
	core.LogMethodDuration(s.logger, stopWatch, fmt.Sprintf("Storage::%s completed", method))
}
