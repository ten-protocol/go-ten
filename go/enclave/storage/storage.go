package storage

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/ten-protocol/go-ten/go/common/errutil"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/triedb/hashdb"

	"github.com/ethereum/go-ethereum/triedb"

	"github.com/ten-protocol/go-ten/go/common/measure"

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

type EventType struct {
	id          uint64
	isLifecycle bool
}

// todo - this file needs splitting up based on concerns
type storageImpl struct {
	db                 enclavedb.EnclaveDB
	cachingService     *CacheService
	eventsStorage      *eventsStorage
	cachedSharedSecret *crypto.SharedEnclaveSecret

	stateCache  state.Database
	chainConfig *params.ChainConfig
	logger      gethlog.Logger
}

func NewStorageFromConfig(config *config.EnclaveConfig, cachingService *CacheService, chainConfig *params.ChainConfig, logger gethlog.Logger) Storage {
	backingDB, err := CreateDBFromConfig(config, logger)
	if err != nil {
		logger.Crit("Failed to connect to backing database", log.ErrKey, err)
	}
	return NewStorage(backingDB, cachingService, chainConfig, logger)
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

func NewStorage(backingDB enclavedb.EnclaveDB, cachingService *CacheService, chainConfig *params.ChainConfig, logger gethlog.Logger) Storage {
	// Open trie database with provided config
	triedb := triedb.NewDatabase(backingDB, trieDBConfig)

	stateDB := state.NewDatabaseWithNodeDB(backingDB, triedb)

	return &storageImpl{
		db:             backingDB,
		stateCache:     stateDB,
		chainConfig:    chainConfig,
		cachingService: cachingService,
		eventsStorage:  newEventsStorage(cachingService, logger),
		logger:         logger,
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

func (s *storageImpl) FetchHeadBatchHeader(ctx context.Context) (*common.BatchHeader, error) {
	defer s.logDuration("FetchHeadBatchHeader", measure.NewStopwatch())
	b, err := enclavedb.ReadCurrentHeadBatchHeader(ctx, s.db.GetSQLDB())
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (s *storageImpl) FetchCurrentSequencerNo(ctx context.Context) (*big.Int, error) {
	defer s.logDuration("FetchCurrentSequencerNo", measure.NewStopwatch())
	return enclavedb.ReadCurrentSequencerNo(ctx, s.db.GetSQLDB())
}

func (s *storageImpl) FetchBatch(ctx context.Context, hash common.L2BatchHash) (*core.Batch, error) {
	defer s.logDuration("FetchBatch", measure.NewStopwatch())
	seqNo, err := s.fetchSeqNoByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return s.FetchBatchBySeqNo(ctx, seqNo.Uint64())
}

func (s *storageImpl) fetchSeqNoByHash(ctx context.Context, hash common.L2BatchHash) (*big.Int, error) {
	seqNo, err := s.cachingService.ReadBatchSeqByHash(ctx, hash, func(v any) (*big.Int, error) {
		batch, err := enclavedb.ReadBatchHeaderByHash(ctx, s.db.GetSQLDB(), v.(common.L2BatchHash))
		if err != nil {
			return nil, err
		}
		return batch.SequencerOrderNo, nil
	})
	return seqNo, err
}

func (s *storageImpl) FetchConvertedHash(ctx context.Context, hash common.L2BatchHash) (gethcommon.Hash, error) {
	defer s.logDuration("FetchConvertedHash", measure.NewStopwatch())
	batch, err := s.FetchBatchHeader(ctx, hash)
	if err != nil {
		return gethcommon.Hash{}, err
	}

	convertedHash, err := s.cachingService.ReadConvertedHash(ctx, hash, func(v any) (*gethcommon.Hash, error) {
		ch, err := enclavedb.FetchConvertedBatchHash(ctx, s.db.GetSQLDB(), batch.SequencerOrderNo.Uint64())
		if err != nil {
			return nil, err
		}
		return &ch, nil
	})
	if err != nil {
		return gethcommon.Hash{}, err
	}
	return *convertedHash, nil
}

func (s *storageImpl) FetchBatchHeader(ctx context.Context, hash common.L2BatchHash) (*common.BatchHeader, error) {
	defer s.logDuration("FetchBatchHeader", measure.NewStopwatch())
	seqNo, err := s.fetchSeqNoByHash(ctx, hash)
	if err != nil {
		return nil, err
	}

	return s.FetchBatchHeaderBySeqNo(ctx, seqNo.Uint64())
}

func (s *storageImpl) FetchBatchTransactionsBySeq(ctx context.Context, seqNo uint64) ([]*common.L2Tx, error) {
	defer s.logDuration("FetchBatchTransactionsBySeq", measure.NewStopwatch())
	batch, err := s.FetchBatchHeaderBySeqNo(ctx, seqNo)
	if err != nil {
		return nil, err
	}
	return enclavedb.ReadBatchTransactions(ctx, s.db.GetSQLDB(), batch.Number.Uint64())
}

func (s *storageImpl) FetchBatchByHeight(ctx context.Context, height uint64) (*core.Batch, error) {
	defer s.logDuration("FetchBatchByHeight", measure.NewStopwatch())
	seqNo, err := s.cachingService.ReadBatchSeqByHeight(ctx, height, func(h any) (*big.Int, error) {
		batch, err := enclavedb.ReadCanonicalBatchHeaderByHeight(ctx, s.db.GetSQLDB(), height)
		if err != nil {
			return nil, err
		}
		return batch.SequencerOrderNo, nil
	})
	if err != nil {
		return nil, err
	}
	return s.FetchBatchBySeqNo(ctx, seqNo.Uint64())
}

func (s *storageImpl) FetchNonCanonicalBatchesBetween(ctx context.Context, startSeq uint64, endSeq uint64) ([]*common.BatchHeader, error) {
	defer s.logDuration("FetchNonCanonicalBatchesBetween", measure.NewStopwatch())
	return enclavedb.ReadNonCanonicalBatches(ctx, s.db.GetSQLDB(), startSeq, endSeq)
}

func (s *storageImpl) FetchCanonicalBatchesBetween(ctx context.Context, startSeq uint64, endSeq uint64) ([]*common.BatchHeader, error) {
	defer s.logDuration("FetchCanonicalBatchesBetween", measure.NewStopwatch())
	return enclavedb.ReadCanonicalBatches(ctx, s.db.GetSQLDB(), startSeq, endSeq)
}

func (s *storageImpl) IsBatchCanonical(ctx context.Context, seq uint64) (bool, error) {
	defer s.logDuration("IsBatchCanonical", measure.NewStopwatch())
	return enclavedb.IsCanonicalBatchSeq(ctx, s.db.GetSQLDB(), seq)
}

func (s *storageImpl) StoreBlock(ctx context.Context, block *types.Block, chainFork *common.ChainFork) error {
	defer s.logDuration("StoreBlock", measure.NewStopwatch())
	dbTx, err := s.db.NewDBTransaction(ctx)
	if err != nil {
		return fmt.Errorf("could not create DB transaction - %w", err)
	}
	defer dbTx.Rollback()

	// only insert the block if it doesn't exist already
	blockId, err := enclavedb.GetBlockId(ctx, dbTx, block.Hash())
	if errors.Is(err, sql.ErrNoRows) {
		if err := enclavedb.WriteBlock(ctx, dbTx, block.Header()); err != nil {
			return fmt.Errorf("2. could not store block %s. Cause: %w", block.Hash(), err)
		}

		blockId, err = enclavedb.GetBlockId(ctx, dbTx, block.Hash())
		if err != nil {
			return fmt.Errorf("3. could not get block id - %w", err)
		}
	}

	// In case there were any batches inserted before this block was received
	err = enclavedb.HandleBlockArrivedAfterBatches(ctx, dbTx, blockId, block.Hash())
	if err != nil {
		return err
	}

	if chainFork != nil && chainFork.IsFork() {
		s.logger.Info(fmt.Sprintf("Update Fork. %s", chainFork))
		err := enclavedb.UpdateCanonicalBlock(ctx, dbTx, false, chainFork.NonCanonicalPath)
		if err != nil {
			return err
		}
		err = enclavedb.UpdateCanonicalBlock(ctx, dbTx, true, chainFork.CanonicalPath)
		if err != nil {
			return err
		}
		err = enclavedb.UpdateCanonicalBatch(ctx, dbTx, false, chainFork.NonCanonicalPath)
		if err != nil {
			return err
		}
		err = enclavedb.UpdateCanonicalBatch(ctx, dbTx, true, chainFork.CanonicalPath)
		if err != nil {
			return err
		}
	}

	// double check that there is always a single canonical batch or block per layer
	// only for debugging
	//err = enclavedb.CheckCanonicalValidity(ctx, dbTx)
	//if err != nil {
	//	return err
	//}

	if err := dbTx.Commit(); err != nil {
		return fmt.Errorf("4. could not store block %s. Cause: %w", block.Hash(), err)
	}

	s.cachingService.CacheBlock(ctx, block)

	return nil
}

func (s *storageImpl) FetchBlock(ctx context.Context, blockHash common.L1BlockHash) (*types.Block, error) {
	defer s.logDuration("FetchBlock", measure.NewStopwatch())
	return s.cachingService.ReadBlock(ctx, blockHash, func(hash any) (*types.Block, error) {
		return enclavedb.FetchBlock(ctx, s.db.GetSQLDB(), hash.(common.L1BlockHash))
	})
}

func (s *storageImpl) IsBlockCanonical(ctx context.Context, blockHash common.L1BlockHash) (bool, error) {
	defer s.logDuration("IsBlockCanonical", measure.NewStopwatch())
	dbtx, err := s.db.NewDBTransaction(ctx)
	if err != nil {
		return false, err
	}
	defer dbtx.Rollback()
	return enclavedb.IsCanonicalBlock(ctx, dbtx, &blockHash)
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
	seqNo, err := s.FetchCurrentSequencerNo(ctx)
	if err != nil {
		return false, err
	}

	if seqNo == nil {
		return false, fmt.Errorf("no batches are stored")
	}

	return true, nil
}

func (s *storageImpl) CreateStateDB(ctx context.Context, batchHash common.L2BatchHash) (*state.StateDB, error) {
	defer s.logDuration("CreateStateDB", measure.NewStopwatch())
	batch, err := s.FetchBatchHeader(ctx, batchHash)
	if err != nil {
		return nil, err
	}

	statedb, err := state.New(batch.Root, s.stateCache, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create state DB for batch: %d. Cause: %w", batch.SequencerOrderNo, err)
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

func (s *storageImpl) GetTransaction(ctx context.Context, txHash gethcommon.Hash) (*types.Transaction, common.L2BatchHash, uint64, uint64, error) {
	defer s.logDuration("GetTransaction", measure.NewStopwatch())
	return enclavedb.ReadTransaction(ctx, s.db.GetSQLDB(), txHash)
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
	h, err := s.FetchBatchHeaderBySeqNo(ctx, seqNum)
	if err != nil {
		return nil, err
	}
	txs, err := s.FetchBatchTransactionsBySeq(ctx, seqNum)
	if err != nil {
		return nil, err
	}
	return &core.Batch{
		Header:       h,
		Transactions: txs,
	}, err
}

func (s *storageImpl) FetchBatchHeaderBySeqNo(ctx context.Context, seqNum uint64) (*common.BatchHeader, error) {
	defer s.logDuration("FetchBatchHeaderBySeqNo", measure.NewStopwatch())
	return s.cachingService.ReadBatch(ctx, seqNum, func(seq any) (*common.BatchHeader, error) {
		return enclavedb.ReadBatchHeaderBySeqNo(ctx, s.db.GetSQLDB(), seqNum)
	})
}

func (s *storageImpl) FetchBatchesByBlock(ctx context.Context, block common.L1BlockHash) ([]*common.BatchHeader, error) {
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

	// the batch is canonical only if the l1 proof is canonical
	isL1ProofCanonical, err := enclavedb.IsCanonicalBlock(ctx, dbTx, &batch.Header.L1Proof)
	if err != nil {
		return err
	}

	// sanity check: a batch can't be canonical if its parent is not
	parentIsCanon, err := enclavedb.IsCanonicalBatchHash(ctx, dbTx, &batch.Header.ParentHash)
	if err != nil {
		return err
	}
	parentIsCanon = parentIsCanon || batch.SeqNo().Uint64() <= common.L2GenesisSeqNo+2
	if isL1ProofCanonical && !parentIsCanon {
		s.logger.Crit("invalid chaining. Batch  is canonical. Parent  is not", log.BatchHashKey, batch.Hash(), "parentHash", batch.Header.ParentHash)
	}

	existsHeight, err := enclavedb.ExistsBatchAtHeight(ctx, dbTx, batch.Header.Number)
	if err != nil {
		return fmt.Errorf("could not read ExistsBatchAtHeight. Cause: %w", err)
	}

	if err := enclavedb.WriteBatchHeader(ctx, dbTx, batch, convertedHash, blockId, isL1ProofCanonical); err != nil {
		return fmt.Errorf("could not write batch header. Cause: %w", err)
	}

	// only insert transactions if this is the first time a batch of this height is created
	if !existsHeight {
		senders, err := s.handleTxSenders(ctx, batch, dbTx)
		if err != nil {
			return err
		}

		for _, tx := range batch.Transactions {
			fmt.Printf("Enclave storing tx: %s\n", tx.Hash().Hex())
		}
		if err := enclavedb.WriteTransactions(ctx, dbTx, batch, senders); err != nil {
			return fmt.Errorf("could not write transactions. Cause: %w", err)
		}
	}

	if err := dbTx.Commit(); err != nil {
		return fmt.Errorf("could not commit batch %w", err)
	}

	s.cachingService.CacheBatch(ctx, batch)
	return nil
}

func (s *storageImpl) handleTxSenders(ctx context.Context, batch *core.Batch, dbTx *sql.Tx) ([]*uint64, error) {
	senders := make([]*uint64, len(batch.Transactions))
	// insert the tx signers as externally owned accounts
	for i, tx := range batch.Transactions {
		sender, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
		if err != nil {
			return nil, fmt.Errorf("could not read tx sender. Cause: %w", err)
		}
		eoaID, err := s.readOrWriteEOA(ctx, dbTx, sender)
		if err != nil {
			return nil, fmt.Errorf("could not insert EOA. cause: %w", err)
		}
		senders[i] = eoaID
	}
	return senders, nil
}

func (s *storageImpl) StoreExecutedBatch(ctx context.Context, batch *common.BatchHeader, receipts []*types.Receipt, newContracts map[gethcommon.Hash][]*gethcommon.Address) error {
	defer s.logDuration("StoreExecutedBatch", measure.NewStopwatch())
	executed, err := enclavedb.BatchWasExecuted(ctx, s.db.GetSQLDB(), batch.Hash())
	if err != nil {
		return err
	}
	if executed {
		s.logger.Debug("Batch was already executed", log.BatchHashKey, batch.Hash())
		return nil
	}

	s.logger.Trace("storing executed batch", log.BatchHashKey, batch.Hash(), log.BatchSeqNoKey, batch.SequencerOrderNo, "receipts", len(receipts))

	dbTx, err := s.db.NewDBTransaction(ctx)
	if err != nil {
		return fmt.Errorf("could not create DB transaction - %w", err)
	}
	defer dbTx.Rollback()

	if err := enclavedb.MarkBatchExecuted(ctx, dbTx, batch.SequencerOrderNo); err != nil {
		return fmt.Errorf("could not set the executed flag. Cause: %w", err)
	}

	for _, receipt := range receipts {
		err = s.eventsStorage.storeReceiptAndEventLogs(ctx, dbTx, batch, receipt, newContracts[receipt.TxHash])
		if err != nil {
			return fmt.Errorf("could not store receipt. Cause: %w", err)
		}
	}
	if err = dbTx.Commit(); err != nil {
		return fmt.Errorf("could not commit batch %w", err)
	}

	return nil
}

func (s *storageImpl) StoreValueTransfers(ctx context.Context, blockHash common.L1BlockHash, transfers common.ValueTransferEvents) error {
	defer s.logDuration("StoreValueTransfers", measure.NewStopwatch())
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
	defer s.logDuration("GetL1Transfers", measure.NewStopwatch())
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
	defer s.logDuration("FetchReorgedRollup", measure.NewStopwatch())
	return enclavedb.FetchReorgedRollup(ctx, s.db.GetSQLDB(), reorgedBlocks)
}

func (s *storageImpl) FetchRollupMetadata(ctx context.Context, hash common.L2RollupHash) (*common.PublicRollupMetadata, error) {
	defer s.logDuration("FetchRollupMetadata", measure.NewStopwatch())
	return enclavedb.FetchRollupMetadata(ctx, s.db.GetSQLDB(), hash)
}

func (s *storageImpl) DebugGetLogs(ctx context.Context, txHash common.TxHash) ([]*tracers.DebugLogs, error) {
	defer s.logDuration("DebugGetLogs", measure.NewStopwatch())
	return enclavedb.DebugGetLogs(ctx, s.db.GetSQLDB(), txHash)
}

func (s *storageImpl) FilterLogsForReceipt(ctx context.Context, requestingAccount *gethcommon.Address, txHash gethcommon.Hash) ([]*types.Log, error) {
	defer s.logDuration("FilterLogs", measure.NewStopwatch())
	logs, err := enclavedb.FilterLogs(ctx, s.db.GetSQLDB(), requestingAccount, nil, nil, nil, nil, nil, &txHash)
	if err != nil {
		return nil, err
	}
	// the database returns an unsorted list of event logs.
	// we have to perform the sorting programatically
	sort.Slice(logs, func(i, j int) bool {
		if logs[i].BlockNumber == logs[j].BlockNumber {
			return logs[i].Index < logs[j].Index
		}
		return logs[i].BlockNumber < logs[j].BlockNumber
	})
	return logs, nil
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
	logs, err := enclavedb.FilterLogs(ctx, s.db.GetSQLDB(), requestingAccount, fromBlock, toBlock, blockHash, addresses, topics, nil)
	if err != nil {
		return nil, err
	}
	// the database returns an unsorted list of event logs.
	// we have to perform the sorting programatically
	sort.Slice(logs, func(i, j int) bool {
		if logs[i].BlockNumber == logs[j].BlockNumber {
			return logs[i].Index < logs[j].Index
		}
		return logs[i].BlockNumber < logs[j].BlockNumber
	})
	return logs, nil
}

func (s *storageImpl) GetContractCount(ctx context.Context) (*big.Int, error) {
	defer s.logDuration("GetContractCount", measure.NewStopwatch())
	return enclavedb.ReadContractCreationCount(ctx, s.db.GetSQLDB())
}

func (s *storageImpl) FetchCanonicalUnexecutedBatches(ctx context.Context, from *big.Int) ([]*common.BatchHeader, error) {
	defer s.logDuration("FetchCanonicalUnexecutedBatches", measure.NewStopwatch())
	return enclavedb.ReadUnexecutedBatches(ctx, s.db.GetSQLDB(), from)
}

func (s *storageImpl) BatchWasExecuted(ctx context.Context, hash common.L2BatchHash) (bool, error) {
	defer s.logDuration("BatchWasExecuted", measure.NewStopwatch())
	return enclavedb.BatchWasExecuted(ctx, s.db.GetSQLDB(), hash)
}

func (s *storageImpl) GetTransactionsPerAddress(ctx context.Context, address *gethcommon.Address, pagination *common.QueryPagination) (types.Receipts, error) {
	defer s.logDuration("GetTransactionsPerAddress", measure.NewStopwatch())
	return enclavedb.GetTransactionsPerAddress(ctx, s.db.GetSQLDB(), s.chainConfig, address, pagination)
}

func (s *storageImpl) CountTransactionsPerAddress(ctx context.Context, address *gethcommon.Address) (uint64, error) {
	defer s.logDuration("CountTransactionsPerAddress", measure.NewStopwatch())
	return enclavedb.CountTransactionsPerAddress(ctx, s.db.GetSQLDB(), address)
}

func (s *storageImpl) readOrWriteEOA(ctx context.Context, dbTX *sql.Tx, addr gethcommon.Address) (*uint64, error) {
	defer s.logDuration("readOrWriteEOA", measure.NewStopwatch())
	return s.cachingService.ReadEOA(ctx, addr, func(v any) (*uint64, error) {
		id, err := enclavedb.ReadEoa(ctx, dbTX, addr)
		if err != nil {
			if errors.Is(err, errutil.ErrNotFound) {
				wid, err := enclavedb.WriteEoa(ctx, dbTX, addr)
				if err != nil {
					return nil, fmt.Errorf("could not write the eoa. Cause: %w", err)
				}
				return &wid, nil
			}
			return nil, fmt.Errorf("count not read eoa. cause: %w", err)
		}
		return &id, nil
	})
}

func (s *storageImpl) ReadContractOwner(ctx context.Context, address gethcommon.Address) (*gethcommon.Address, error) {
	return enclavedb.ReadContractOwner(ctx, s.db.GetSQLDB(), address)
}

func (s *storageImpl) logDuration(method string, stopWatch *measure.Stopwatch) {
	core.LogMethodDuration(s.logger, stopWatch, fmt.Sprintf("Storage::%s completed", method))
}
