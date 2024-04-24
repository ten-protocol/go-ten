package storage

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/ten-protocol/go-ten/go/common/measure"

	ristretto_store "github.com/eko/gocache/store/ristretto/v4"

	"github.com/ten-protocol/go-ten/go/config"

	"github.com/ten-protocol/go-ten/go/enclave/storage/enclavedb"

	"github.com/ethereum/go-ethereum/rlp"

	gethcore "github.com/ethereum/go-ethereum/core"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/trie"

	"github.com/ten-protocol/go-ten/go/common/syserr"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
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
	// these are the twice the default configs from geth
	// todo - consider tweaking these independently on the validator and on the sequencer
	// the validator probably need higher values on this cache?
	cacheConfig := &gethcore.CacheConfig{
		TrieCleanLimit: 256 * 2,
		TrieDirtyLimit: 256 * 2,
		TrieTimeLimit:  5 * time.Minute,
		SnapshotLimit:  256 * 2,
		SnapshotWait:   true,
	}

	stateDB := state.NewDatabaseWithConfig(backingDB, &trie.Config{
		Preimages: cacheConfig.Preimages,
	})

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
		stateDB:           stateDB,
		chainConfig:       chainConfig,
		blockCache:        cache.New[*types.Block](ristrettoStore),
		batchCacheBySeqNo: cache.New[*core.Batch](ristrettoStore),
		seqCacheByHash:    cache.New[*big.Int](ristrettoStore),
		seqCacheByHeight:  cache.New[*big.Int](ristrettoStore),
		logger:            logger,
	}
}

func (s *storageImpl) TrieDB() *trie.Database {
	return s.stateDB.TrieDB()
}

func (s *storageImpl) StateDB() state.Database {
	return s.stateDB
}

func (s *storageImpl) Close() error {
	return s.db.GetSQLDB().Close()
}

func (s *storageImpl) FetchHeadBatch() (*core.Batch, error) {
	defer s.logDuration("FetchHeadBatch", measure.NewStopwatch())
	return enclavedb.ReadCurrentHeadBatch(s.db.GetSQLDB())
}

func (s *storageImpl) FetchCurrentSequencerNo() (*big.Int, error) {
	defer s.logDuration("FetchCurrentSequencerNo", measure.NewStopwatch())
	return enclavedb.ReadCurrentSequencerNo(s.db.GetSQLDB())
}

func (s *storageImpl) FetchBatch(hash common.L2BatchHash) (*core.Batch, error) {
	defer s.logDuration("FetchBatch", measure.NewStopwatch())
	seqNo, err := common.GetCachedValue(s.seqCacheByHash, s.logger, hash, func(v any) (*big.Int, error) {
		batch, err := enclavedb.ReadBatchByHash(s.db.GetSQLDB(), v.(common.L2BatchHash))
		if err != nil {
			return nil, err
		}
		return batch.SeqNo(), nil
	})
	if err != nil {
		return nil, err
	}
	return s.FetchBatchBySeqNo(seqNo.Uint64())
}

func (s *storageImpl) FetchConvertedHash(hash common.L2BatchHash) (gethcommon.Hash, error) {
	defer s.logDuration("FetchConvertedHash", measure.NewStopwatch())
	batch, err := s.FetchBatch(hash)
	if err != nil {
		return gethcommon.Hash{}, err
	}
	return enclavedb.FetchConvertedBatchHash(s.db.GetSQLDB(), batch.Header.SequencerOrderNo.Uint64())
}

func (s *storageImpl) FetchBatchHeader(hash common.L2BatchHash) (*common.BatchHeader, error) {
	defer s.logDuration("FetchBatchHeader", measure.NewStopwatch())
	b, err := s.FetchBatch(hash)
	if err != nil {
		return nil, err
	}
	return b.Header, nil
}

func (s *storageImpl) FetchBatchByHeight(height uint64) (*core.Batch, error) {
	defer s.logDuration("FetchBatchByHeight", measure.NewStopwatch())
	// the key is (height+1), because for some reason it doesn't like a key of 0
	seqNo, err := common.GetCachedValue(s.seqCacheByHeight, s.logger, height+1, func(h any) (*big.Int, error) {
		batch, err := enclavedb.ReadCanonicalBatchByHeight(s.db.GetSQLDB(), height)
		if err != nil {
			return nil, err
		}
		return batch.SeqNo(), nil
	})
	if err != nil {
		return nil, err
	}
	return s.FetchBatchBySeqNo(seqNo.Uint64())
}

func (s *storageImpl) FetchNonCanonicalBatchesBetween(startSeq uint64, endSeq uint64) ([]*core.Batch, error) {
	defer s.logDuration("FetchNonCanonicalBatchesBetween", measure.NewStopwatch())
	return enclavedb.ReadNonCanonicalBatches(s.db.GetSQLDB(), startSeq, endSeq)
}

func (s *storageImpl) StoreBlock(b *types.Block, chainFork *common.ChainFork) error {
	defer s.logDuration("StoreBlock", measure.NewStopwatch())
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

	common.CacheValue(s.blockCache, s.logger, b.Hash(), b)

	return nil
}

func (s *storageImpl) FetchBlock(blockHash common.L1BlockHash) (*types.Block, error) {
	defer s.logDuration("FetchBlock", measure.NewStopwatch())
	return common.GetCachedValue(s.blockCache, s.logger, blockHash, func(hash any) (*types.Block, error) {
		return enclavedb.FetchBlock(s.db.GetSQLDB(), hash.(common.L1BlockHash))
	})
}

func (s *storageImpl) FetchCanonicaBlockByHeight(height *big.Int) (*types.Block, error) {
	defer s.logDuration("FetchCanonicaBlockByHeight", measure.NewStopwatch())
	header, err := enclavedb.FetchBlockHeaderByHeight(s.db.GetSQLDB(), height)
	if err != nil {
		return nil, err
	}
	return s.FetchBlock(header.Hash())
}

func (s *storageImpl) FetchCanonicalBlocksBetween(start *big.Int, end *big.Int) ([]*types.Header, error) {
	defer s.logDuration("FetchCanonicalBlocksBetween", measure.NewStopwatch())
	return enclavedb.FetchBlockHeadersBetween(s.db.GetSQLDB(), start, end)
}

func (s *storageImpl) FetchHeadBlock() (*types.Block, error) {
	defer s.logDuration("FetchHeadBlock", measure.NewStopwatch())
	return enclavedb.FetchHeadBlock(s.db.GetSQLDB())
}

func (s *storageImpl) StoreSecret(secret crypto.SharedEnclaveSecret) error {
	defer s.logDuration("StoreSecret", measure.NewStopwatch())
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
	defer s.logDuration("FetchSecret", measure.NewStopwatch())

	if s.cachedSharedSecret != nil {
		return s.cachedSharedSecret, nil
	}

	var ss crypto.SharedEnclaveSecret

	cfg, err := enclavedb.FetchConfig(s.db.GetSQLDB(), masterSeedCfg)
	if err != nil {
		return nil, err
	}
	if err := rlp.DecodeBytes(cfg, &ss); err != nil {
		return nil, fmt.Errorf("could not decode shared secret")
	}

	s.cachedSharedSecret = &ss
	return s.cachedSharedSecret, nil
}

func (s *storageImpl) IsAncestor(block *types.Block, maybeAncestor *types.Block) bool {
	defer s.logDuration("IsAncestor", measure.NewStopwatch())
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
	defer s.logDuration("IsBlockAncestor", measure.NewStopwatch())
	resolvedBlock, err := s.FetchBlock(maybeAncestor)
	if err != nil {
		return false
	}
	return s.IsAncestor(block, resolvedBlock)
}

func (s *storageImpl) HealthCheck() (bool, error) {
	defer s.logDuration("HealthCheck", measure.NewStopwatch())
	headBatch, err := s.FetchHeadBatch()
	if err != nil {
		return false, err
	}

	if headBatch == nil {
		return false, fmt.Errorf("head batch is nil")
	}

	return true, nil
}

func (s *storageImpl) FetchHeadBatchForBlock(blockHash common.L1BlockHash) (*core.Batch, error) {
	defer s.logDuration("FetchHeadBatchForBlock", measure.NewStopwatch())
	return enclavedb.ReadHeadBatchForBlock(s.db.GetSQLDB(), blockHash)
}

func (s *storageImpl) CreateStateDB(batchHash common.L2BatchHash) (*state.StateDB, error) {
	defer s.logDuration("CreateStateDB", measure.NewStopwatch())
	batch, err := s.FetchBatch(batchHash)
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
	defer s.logDuration("EmptyStateDB", measure.NewStopwatch())
	statedb, err := state.New(types.EmptyRootHash, s.stateDB, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create state DB. Cause: %w", err)
	}
	return statedb, nil
}

// GetReceiptsByBatchHash retrieves the receipts for all transactions in a given batch.
func (s *storageImpl) GetReceiptsByBatchHash(hash gethcommon.Hash) (types.Receipts, error) {
	defer s.logDuration("GetReceiptsByBatchHash", measure.NewStopwatch())
	return enclavedb.ReadReceiptsByBatchHash(s.db.GetSQLDB(), hash, s.chainConfig)
}

func (s *storageImpl) GetTransaction(txHash gethcommon.Hash) (*types.Transaction, common.L2BatchHash, uint64, uint64, error) {
	defer s.logDuration("GetTransaction", measure.NewStopwatch())
	return enclavedb.ReadTransaction(s.db.GetSQLDB(), txHash)
}

func (s *storageImpl) GetContractCreationTx(address gethcommon.Address) (*gethcommon.Hash, error) {
	defer s.logDuration("GetContractCreationTx", measure.NewStopwatch())
	return enclavedb.GetContractCreationTx(s.db.GetSQLDB(), address)
}

func (s *storageImpl) GetTransactionReceipt(txHash gethcommon.Hash) (*types.Receipt, error) {
	defer s.logDuration("GetTransactionReceipt", measure.NewStopwatch())
	return enclavedb.ReadReceipt(s.db.GetSQLDB(), txHash, s.chainConfig)
}

func (s *storageImpl) FetchAttestedKey(address gethcommon.Address) (*ecdsa.PublicKey, error) {
	defer s.logDuration("FetchAttestedKey", measure.NewStopwatch())
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
	defer s.logDuration("StoreAttestedKey", measure.NewStopwatch())
	_, err := enclavedb.WriteAttKey(s.db.GetSQLDB(), aggregator, gethcrypto.CompressPubkey(key))
	return err
}

func (s *storageImpl) FetchBatchBySeqNo(seqNum uint64) (*core.Batch, error) {
	defer s.logDuration("FetchBatchBySeqNo", measure.NewStopwatch())
	b, err := common.GetCachedValue(s.batchCacheBySeqNo, s.logger, seqNum, func(seq any) (*core.Batch, error) {
		return enclavedb.ReadBatchBySeqNo(s.db.GetSQLDB(), seqNum)
	})
	if err == nil && b == nil {
		return nil, fmt.Errorf("not found")
	}
	return b, err
}

func (s *storageImpl) FetchBatchesByBlock(block common.L1BlockHash) ([]*core.Batch, error) {
	defer s.logDuration("FetchBatchesByBlock", measure.NewStopwatch())
	return enclavedb.ReadBatchesByBlock(s.db.GetSQLDB(), block)
}

func (s *storageImpl) StoreBatch(batch *core.Batch, convertedHash gethcommon.Hash) error {
	defer s.logDuration("StoreBatch", measure.NewStopwatch())
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

	if err := enclavedb.WriteBatchAndTransactions(dbTx, batch, convertedHash); err != nil {
		return fmt.Errorf("could not write batch. Cause: %w", err)
	}

	if err := dbTx.Write(); err != nil {
		return fmt.Errorf("could not commit batch %w", err)
	}

	common.CacheValue(s.batchCacheBySeqNo, s.logger, batch.SeqNo().Uint64(), batch)
	common.CacheValue(s.seqCacheByHash, s.logger, batch.Hash(), batch.SeqNo())
	// note: the key is (height+1), because for some reason it doesn't like a key of 0
	// should always contain the canonical batch because the cache is overwritten by each new batch after a reorg
	common.CacheValue(s.seqCacheByHeight, s.logger, batch.NumberU64()+1, batch.SeqNo())
	return nil
}

func (s *storageImpl) StoreExecutedBatch(batch *core.Batch, receipts []*types.Receipt) error {
	defer s.logDuration("StoreExecutedBatch", measure.NewStopwatch())
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
	defer s.logDuration("StoreL1Messages", measure.NewStopwatch())
	return enclavedb.WriteL1Messages(s.db.GetSQLDB(), blockHash, messages, false)
}

func (s *storageImpl) GetL1Messages(blockHash common.L1BlockHash) (common.CrossChainMessages, error) {
	defer s.logDuration("GetL1Messages", measure.NewStopwatch())
	return enclavedb.FetchL1Messages[common.CrossChainMessage](s.db.GetSQLDB(), blockHash, false)
}

func (s *storageImpl) GetL1Transfers(blockHash common.L1BlockHash) (common.ValueTransferEvents, error) {
	return enclavedb.FetchL1Messages[common.ValueTransferEvent](s.db.GetSQLDB(), blockHash, true)
}

const enclaveKeyKey = "ek"

func (s *storageImpl) StoreEnclaveKey(enclaveKey *crypto.EnclaveKey) error {
	defer s.logDuration("StoreEnclaveKey", measure.NewStopwatch())
	if enclaveKey == nil {
		return errors.New("enclaveKey cannot be nil")
	}
	keyBytes := gethcrypto.FromECDSA(enclaveKey.PrivateKey())

	_, err := enclavedb.WriteConfig(s.db.GetSQLDB(), enclaveKeyKey, keyBytes)
	return err
}

func (s *storageImpl) GetEnclaveKey() (*crypto.EnclaveKey, error) {
	defer s.logDuration("GetEnclaveKey", measure.NewStopwatch())
	keyBytes, err := enclavedb.FetchConfig(s.db.GetSQLDB(), enclaveKeyKey)
	if err != nil {
		return nil, err
	}
	ecdsaKey, err := gethcrypto.ToECDSA(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to construct ECDSA private key from enclave key bytes - %w", err)
	}
	return crypto.NewEnclaveKey(ecdsaKey), nil
}

func (s *storageImpl) StoreRollup(rollup *common.ExtRollup, internalHeader *common.CalldataRollupHeader) error {
	defer s.logDuration("StoreRollup", measure.NewStopwatch())
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

func (s *storageImpl) FetchRollupMetadata(hash common.L2RollupHash) (*common.PublicRollupMetadata, error) {
	return enclavedb.FetchRollupMetadata(s.db.GetSQLDB(), hash)
}

func (s *storageImpl) DebugGetLogs(txHash common.TxHash) ([]*tracers.DebugLogs, error) {
	defer s.logDuration("DebugGetLogs", measure.NewStopwatch())
	return enclavedb.DebugGetLogs(s.db.GetSQLDB(), txHash)
}

func (s *storageImpl) FilterLogs(
	requestingAccount *gethcommon.Address,
	fromBlock, toBlock *big.Int,
	blockHash *common.L2BatchHash,
	addresses []gethcommon.Address,
	topics [][]gethcommon.Hash,
) ([]*types.Log, error) {
	defer s.logDuration("FilterLogs", measure.NewStopwatch())
	return enclavedb.FilterLogs(s.db.GetSQLDB(), requestingAccount, fromBlock, toBlock, blockHash, addresses, topics)
}

func (s *storageImpl) GetContractCount() (*big.Int, error) {
	defer s.logDuration("GetContractCount", measure.NewStopwatch())
	return enclavedb.ReadContractCreationCount(s.db.GetSQLDB())
}

func (s *storageImpl) FetchCanonicalUnexecutedBatches(from *big.Int) ([]*core.Batch, error) {
	defer s.logDuration("FetchCanonicalUnexecutedBatches", measure.NewStopwatch())
	return enclavedb.ReadUnexecutedBatches(s.db.GetSQLDB(), from)
}

func (s *storageImpl) BatchWasExecuted(hash common.L2BatchHash) (bool, error) {
	defer s.logDuration("BatchWasExecuted", measure.NewStopwatch())
	return enclavedb.BatchWasExecuted(s.db.GetSQLDB(), hash)
}

func (s *storageImpl) GetReceiptsPerAddress(address *gethcommon.Address, pagination *common.QueryPagination) (types.Receipts, error) {
	defer s.logDuration("GetReceiptsPerAddress", measure.NewStopwatch())
	return enclavedb.GetReceiptsPerAddress(s.db.GetSQLDB(), s.chainConfig, address, pagination)
}

func (s *storageImpl) GetReceiptsPerAddressCount(address *gethcommon.Address) (uint64, error) {
	defer s.logDuration("GetReceiptsPerAddressCount", measure.NewStopwatch())
	return enclavedb.GetReceiptsPerAddressCount(s.db.GetSQLDB(), address)
}

func (s *storageImpl) GetPublicTransactionData(pagination *common.QueryPagination) ([]common.PublicTransaction, error) {
	defer s.logDuration("GetPublicTransactionData", measure.NewStopwatch())
	return enclavedb.GetPublicTransactionData(s.db.GetSQLDB(), pagination)
}

func (s *storageImpl) GetPublicTransactionCount() (uint64, error) {
	defer s.logDuration("GetPublicTransactionCount", measure.NewStopwatch())
	return enclavedb.GetPublicTransactionCount(s.db.GetSQLDB())
}

func (s *storageImpl) logDuration(method string, stopWatch *measure.Stopwatch) {
	core.LogMethodDuration(s.logger, stopWatch, fmt.Sprintf("Storage::%s completed", method))
}
