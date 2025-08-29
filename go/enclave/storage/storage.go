package storage

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ten-protocol/go-ten/go/common/storage"

	"github.com/jmoiron/sqlx"

	"github.com/ten-protocol/go-ten/go/common/errutil"
	enclaveconfig "github.com/ten-protocol/go-ten/go/enclave/config"

	"github.com/ethereum/go-ethereum/triedb"

	"github.com/ten-protocol/go-ten/go/common/measure"

	"github.com/ten-protocol/go-ten/go/enclave/storage/enclavedb"

	"github.com/ethereum/go-ethereum/rlp"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"github.com/ten-protocol/go-ten/go/enclave/crypto"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
)

// these are the keys from the config table
const (
	// todo - this will require a dedicated table when upgrades are implemented
	masterSeedCfg              = "MASTER_SEED"
	enclaveKeyCfg              = "ENCLAVE_KEY"
	systemContractAddressesCfg = "SYSTEM_CONTRACT_ADDRESSES"
)

type AttestedEnclave struct {
	PubKey    *ecdsa.PublicKey
	EnclaveID *common.EnclaveID
	Type      common.NodeType
}

// todo - this file needs splitting up based on concerns
type storageImpl struct {
	db                     enclavedb.EnclaveDB
	trieDB                 *triedb.Database
	preparedStatementCache *enclavedb.PreparedStatementCache
	cachingService         *CacheService
	eventsStorage          *eventsStorage

	stateCache  state.Database
	chainConfig *params.ChainConfig
	config      *enclaveconfig.EnclaveConfig
	logger      gethlog.Logger
}

func NewStorageFromConfig(config *enclaveconfig.EnclaveConfig, cachingService *CacheService, chainConfig *params.ChainConfig, logger gethlog.Logger) Storage {
	backingDB, err := CreateDBFromConfig(config, logger)
	if err != nil {
		logger.Crit("Failed to connect to backing database", log.ErrKey, err)
	}
	return NewStorage(backingDB, cachingService, config, chainConfig, logger)
}

func NewStorage(backingDB enclavedb.EnclaveDB, cachingService *CacheService, config *enclaveconfig.EnclaveConfig, chainConfig *params.ChainConfig, logger gethlog.Logger) Storage {
	// Open trie database with provided config
	trieDB := triedb.NewDatabase(backingDB, triedb.HashDefaults)

	// todo - figure out the snapshot tree
	stateDB := state.NewDatabase(trieDB, nil)

	prepStatementCache := enclavedb.NewStatementCache(backingDB.GetSQLDB(), logger)
	return &storageImpl{
		db:                     backingDB,
		trieDB:                 trieDB,
		stateCache:             stateDB,
		chainConfig:            chainConfig,
		config:                 config,
		cachingService:         cachingService,
		eventsStorage:          newEventsStorage(cachingService, backingDB, logger),
		preparedStatementCache: prepStatementCache,
		logger:                 logger,
	}
}

func (s *storageImpl) TrieDB() *triedb.Database {
	return s.stateCache.TrieDB()
}

func (s *storageImpl) StateDB() state.Database {
	return s.stateCache
}

func (s *storageImpl) Close() error {
	s.cachingService.Stop()
	s.preparedStatementCache.Clear()
	s.trieDB.Close()
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
	seqNo, err := s.cachingService.ReadBatchSeqByHash(ctx, hash, func() (*big.Int, error) {
		batch, err := enclavedb.ReadBatchHeaderByHash(ctx, s.db.GetSQLDB(), hash)
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

	convertedHash, err := s.cachingService.ReadConvertedHash(ctx, hash, func() (*gethcommon.Hash, error) {
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
	batch, err := s.cachingService.ReadBatch(ctx, seqNo, func() (*core.Batch, error) {
		batchHeader, err := s.FetchBatchHeaderBySeqNo(ctx, seqNo)
		if err != nil {
			return nil, err
		}
		txs, err := enclavedb.ReadBatchTransactions(ctx, s.db.GetSQLDB(), batchHeader.Number.Uint64())
		if err != nil {
			return nil, err
		}
		return &core.Batch{Header: batchHeader, Transactions: txs}, nil
	})
	if err != nil {
		return nil, err
	}
	return batch.Transactions, nil
}

func (s *storageImpl) FetchBatchByHeight(ctx context.Context, height uint64) (*core.Batch, error) {
	defer s.logDuration("FetchBatchByHeight", measure.NewStopwatch())
	seqNo, err := s.cachingService.ReadBatchSeqByHeight(ctx, height, func() (*big.Int, error) {
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

func (s *storageImpl) StoreBlock(ctx context.Context, block *types.Header, chainFork *common.ChainFork) error {
	defer s.logDuration("StoreBlock", measure.NewStopwatch())
	dbTx, err := s.db.NewDBTransaction(ctx)
	if err != nil {
		return fmt.Errorf("could not create DB transaction - %w", err)
	}
	defer dbTx.Rollback()

	// only insert the block if it doesn't exist already
	blockId, err := enclavedb.GetBlockId(ctx, dbTx, block.Hash())
	if errors.Is(err, sql.ErrNoRows) {
		if err := enclavedb.WriteBlock(ctx, dbTx, block); err != nil {
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

	var nonCanonical, canonical []common.L1BlockHash
	if chainFork != nil && chainFork.IsFork() {
		s.logger.Info(fmt.Sprintf("Update Fork. %s", chainFork))
		nonCanonical = chainFork.NonCanonicalPath
		canonical = chainFork.CanonicalPath
	} else {
		// handle the case when this block was canonical at some point, then reverted
		canonical = []common.L1BlockHash{block.Hash()}
	}

	err = enclavedb.UpdateCanonicalBlock(ctx, dbTx, false, nonCanonical)
	if err != nil {
		return err
	}
	err = enclavedb.UpdateCanonicalBlock(ctx, dbTx, true, canonical)
	if err != nil {
		return err
	}
	err = enclavedb.UpdateCanonicalBatch(ctx, dbTx, false, nonCanonical)
	if err != nil {
		return err
	}
	err = enclavedb.UpdateCanonicalBatch(ctx, dbTx, true, canonical)
	if err != nil {
		return err
	}

	// sanity check that there is always a single canonical batch or block per layer
	// called after forks, for the latest 50 blocks
	// uncomment when debugging a race
	//err = enclavedb.CheckCanonicalValidity(ctx, dbTx, blockId-50)
	//if err != nil {
	//	s.logger.Crit("Should not happen.", log.ErrKey, err)
	//}

	if err := dbTx.Commit(); err != nil {
		return fmt.Errorf("4. could not store block %s. Cause: %w", block.Hash(), err)
	}

	s.cachingService.CacheBlock(ctx, block)

	return nil
}

func (s *storageImpl) UpdateProcessed(ctx context.Context, block common.L1BlockHash) error {
	defer s.logDuration("UpdateProcessed", measure.NewStopwatch())
	dbTx, err := s.db.NewDBTransaction(ctx)
	if err != nil {
		return fmt.Errorf("could not create DB transaction - %w", err)
	}
	defer dbTx.Rollback()

	err = enclavedb.UpdateBlockProcessed(ctx, dbTx, block)
	if err != nil {
		return fmt.Errorf("could not update block processed - %w", err)
	}
	if err := dbTx.Commit(); err != nil {
		return fmt.Errorf(" could not update block processed. Cause: %w", err)
	}
	return nil
}

func (s *storageImpl) DeleteDirtyBlocks(ctx context.Context) error {
	defer s.logDuration("DeleteDirtyBlocks", measure.NewStopwatch())
	dbTx, err := s.db.NewDBTransaction(ctx)
	if err != nil {
		return fmt.Errorf("could not create DB transaction - %w", err)
	}
	defer dbTx.Rollback()
	dirtyBlocks, err := enclavedb.SelectUnprocessedBlocks(ctx, dbTx)
	if err != nil {
		return fmt.Errorf("could not select unprocessed blocks - %w", err)
	}
	if len(dirtyBlocks) == 0 {
		// nothing to do
		return nil
	}
	s.logger.Warn("Found unprocessed blocks", "count", len(dirtyBlocks))

	// delete rollups
	err = enclavedb.DeleteUnprocessedRollups(ctx, dbTx)
	if err != nil {
		return fmt.Errorf("could not delete unprocessed rollups. Cause: %w", err)
	}
	// delete cross chain messages
	err = enclavedb.DeleteUnprocessedL1Messages(ctx, dbTx)
	if err != nil {
		return fmt.Errorf("could not delete unprocessed cross chain messages. Cause: %w", err)
	}
	// delete block
	err = enclavedb.DeleteUnProcessedBlock(ctx, dbTx)
	if err != nil {
		return fmt.Errorf("could not delete unprocessed blocks. Cause: %w", err)
	}

	if err := dbTx.Commit(); err != nil {
		return fmt.Errorf(" could not remove dirty blocks. Cause: %w", err)
	}
	return nil
}

func (s *storageImpl) FetchBlock(ctx context.Context, blockHash common.L1BlockHash) (*types.Header, error) {
	defer s.logDuration("FetchBlockHeader", measure.NewStopwatch())
	return s.cachingService.ReadBlock(ctx, blockHash, func() (*types.Header, error) {
		return enclavedb.FetchBlockHeader(ctx, s.db.GetSQLDB(), blockHash)
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

func (s *storageImpl) FetchCanonicaBlockByHeight(ctx context.Context, height *big.Int) (*types.Header, error) {
	defer s.logDuration("FetchCanonicaBlockByHeight", measure.NewStopwatch())
	header, err := enclavedb.FetchBlockHeaderByHeight(ctx, s.db.GetSQLDB(), height)
	if err != nil {
		return nil, err
	}
	return s.FetchBlock(ctx, header.Hash())
}

func (s *storageImpl) FetchHeadBlock(ctx context.Context) (*types.Header, error) {
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
	_, err = storage.WriteConfig(ctx, dbTx, masterSeedCfg, enc)
	if err != nil {
		return fmt.Errorf("could not shared secret in DB. Cause: %w", err)
	}
	return dbTx.Commit()
}

// FetchSecret - this returns the most important secret, and should only be called during startup
func (s *storageImpl) FetchSecret(ctx context.Context) (*crypto.SharedEnclaveSecret, error) {
	defer s.logDuration("FetchSecret", measure.NewStopwatch())

	var ss crypto.SharedEnclaveSecret

	cfg, err := storage.FetchConfig(ctx, s.db.GetSQLDB(), masterSeedCfg)
	if err != nil {
		return nil, err
	}
	if err := rlp.DecodeBytes(cfg, &ss); err != nil {
		return nil, fmt.Errorf("could not decode shared secret")
	}

	return &ss, nil
}

func (s *storageImpl) IsAncestor(ctx context.Context, block *types.Header, maybeAncestor *types.Header) bool {
	defer s.logDuration("IsAncestor", measure.NewStopwatch())
	if bytes.Equal(maybeAncestor.Hash().Bytes(), block.Hash().Bytes()) {
		return true
	}

	if maybeAncestor.Number.Uint64() >= block.Number.Uint64() {
		return false
	}

	p, err := s.FetchBlock(ctx, block.ParentHash)
	if err != nil {
		s.logger.Debug("Could not find block with hash", log.BlockHashKey, block.ParentHash, log.ErrKey, err)
		return false
	}

	return s.IsAncestor(ctx, p, maybeAncestor)
}

func (s *storageImpl) HealthCheck(ctx context.Context) (bool, error) {
	defer s.logDuration("HealthCheck", measure.NewStopwatch())
	return s.db != nil, nil
	//seqNo, err := s.FetchCurrentSequencerNo(ctx)
	//if err != nil {
	//	return false, err
	//}
	//
	//if seqNo == nil {
	//	return false, fmt.Errorf("no batches are stored")
	//}
	//
	//return true, nil
}

func (s *storageImpl) CreateStateDB(ctx context.Context, batchHash common.L2BatchHash) (*state.StateDB, error) {
	defer s.logDuration("CreateStateDB", measure.NewStopwatch())
	batch, err := s.FetchBatchHeader(ctx, batchHash)
	if err != nil {
		return nil, err
	}

	statedb, err := state.New(batch.Root, s.stateCache)
	if err != nil {
		return nil, fmt.Errorf("could not create state DB for batch: %d. Cause: %w", batch.SequencerOrderNo, err)
	}
	return statedb, nil
}

func (s *storageImpl) EmptyStateDB() (*state.StateDB, error) {
	defer s.logDuration("EmptyStateDB", measure.NewStopwatch())
	statedb, err := state.New(types.EmptyRootHash, s.stateCache)
	if err != nil {
		return nil, fmt.Errorf("could not create state DB. Cause: %w", err)
	}
	return statedb, nil
}

func (s *storageImpl) GetTransaction(ctx context.Context, txHash common.L2TxHash) (*types.Transaction, common.L2BatchHash, uint64, uint64, gethcommon.Address, error) {
	defer s.logDuration("GetTransaction", measure.NewStopwatch())
	return enclavedb.ReadTransaction(ctx, s.db.GetSQLDB(), txHash)
}

func (s *storageImpl) GetFilteredInternalReceipt(ctx context.Context, txHash common.L2TxHash, requester *gethcommon.Address, syntheticTx bool) (*core.InternalReceipt, error) {
	defer s.logDuration("GetFilteredInternalReceipt", measure.NewStopwatch())
	if !syntheticTx && requester == nil {
		return nil, errors.New("requester address is required for non-synthetic transactions")
	}
	var requesterId *uint64
	var err error
	if requester != nil {
		requesterId, err = s.readOrWriteEOAWithTx(ctx, *requester)
		if err != nil {
			return nil, fmt.Errorf("could not read EOA for requester address %s. Cause: %w", requester.Hex(), err)
		}
	}
	return enclavedb.ReadReceipt(ctx, s.preparedStatementCache, txHash, requesterId)
}

func (s *storageImpl) ExistsTransactionReceipt(ctx context.Context, txHash common.L2TxHash) (bool, error) {
	defer s.logDuration("ExistsTransactionReceipt", measure.NewStopwatch())
	return enclavedb.ExistsReceipt(ctx, s.db.GetSQLDB(), txHash)
}

func (s *storageImpl) GetEnclavePubKey(ctx context.Context, enclaveId common.EnclaveID) (*AttestedEnclave, error) {
	defer s.logDuration("GetEnclavePubKey", measure.NewStopwatch())
	return s.cachingService.ReadEnclavePubKey(ctx, enclaveId, func() (*AttestedEnclave, error) {
		key, nodeType, err := enclavedb.FetchAttestation(ctx, s.db.GetSQLDB(), enclaveId)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve attestation key for enclave %s. Cause: %w", enclaveId, err)
		}

		publicKey, err := gethcrypto.DecompressPubkey(key)
		if err != nil {
			return nil, fmt.Errorf("could not parse key from db. Cause: %w", err)
		}

		return &AttestedEnclave{PubKey: publicKey, Type: nodeType, EnclaveID: &enclaveId}, nil
	})
}

func (s *storageImpl) StoreNodeType(ctx context.Context, enclaveId common.EnclaveID, nodeType common.NodeType) error {
	defer s.logDuration("StoreNodeType", measure.NewStopwatch())
	dbTx, err := s.db.NewDBTransaction(ctx)
	if err != nil {
		return fmt.Errorf("could not create DB transaction - %w", err)
	}
	defer dbTx.Rollback()
	res, err := enclavedb.UpdateAttestationType(ctx, dbTx, enclaveId, nodeType)
	if err != nil {
		return err
	}

	// we always expect to update exactly one row here, shouldn't be permissioning an enclave we don't recognize
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not get rows affected - %w", err)
	}
	if rowsAffected != 1 {
		return fmt.Errorf("expected to update nodetype for 1 attested enclave, but updated %d", rowsAffected)
	}
	err = dbTx.Commit()
	if err != nil {
		return fmt.Errorf("could not commit transaction - %w", err)
	}
	// set value in cache to ensure it is up to date
	s.cachingService.UpdateEnclaveNodeType(ctx, enclaveId, nodeType)
	// Fetch and update sequencer IDs cache
	sequencerIDs, err := enclavedb.FetchSequencerEnclaveIDs(ctx, s.db.GetSQLDB())
	if err != nil {
		return fmt.Errorf("could not fetch updated sequencer IDs - %w", err)
	}
	s.cachingService.CacheSequencerIDs(ctx, sequencerIDs)

	return nil
}

func (s *storageImpl) StoreNewEnclave(ctx context.Context, enclaveId common.EnclaveID, key *ecdsa.PublicKey) error {
	defer s.logDuration("StoreNewEnclave", measure.NewStopwatch())
	dbTx, err := s.db.NewDBTransaction(ctx)
	if err != nil {
		return fmt.Errorf("could not create DB transaction - %w", err)
	}
	defer dbTx.Rollback()

	alreadyExists, err := enclavedb.AttestationExists(ctx, dbTx, enclaveId)
	if err != nil {
		return fmt.Errorf("failed to check if attestation exists - %w", err)
	}

	compressedKey := gethcrypto.CompressPubkey(key)
	if alreadyExists {
		// this should be unusual, log it for visibility
		s.logger.Warn("Updating existing attestation key", "enclaveId", enclaveId)
		_, err = enclavedb.UpdateAttestationKey(ctx, dbTx, enclaveId, compressedKey)
	} else {
		_, err = enclavedb.WriteAttestation(ctx, dbTx, enclaveId, compressedKey, common.Validator)
	}
	if err != nil {
		return fmt.Errorf("failed to write/update attestation - %w", err)
	}

	return dbTx.Commit()
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
	return s.cachingService.ReadBatchHeader(ctx, seqNum, func() (*common.BatchHeader, error) {
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

	if err := enclavedb.WriteBatchHeader(ctx, dbTx, batch, convertedHash, isL1ProofCanonical); err != nil {
		return fmt.Errorf("could not write batch header. Cause: %w", err)
	}

	// only insert transactions if this is the first time a batch of this height is created
	if !existsHeight {
		transactionsWithSenders := make([]*core.TxWithSender, len(batch.Transactions))
		for i, tx := range batch.Transactions {
			sender, err := core.GetExternalTxSigner(tx)
			if err != nil {
				return fmt.Errorf("could not get tx sender. Cause: %w", err)
			}
			transactionsWithSenders[i] = &core.TxWithSender{Tx: tx, Sender: &sender}
		}

		// todo - there is an edge case where the contract is created in the same batch
		senderIds, toContractIds, toEoas, err := s.handleTxSendersAndReceivers(ctx, transactionsWithSenders, dbTx)
		if err != nil {
			return err
		}

		for _, txWithSender := range transactionsWithSenders {
			s.logger.Debug("Save TX", "tx", txWithSender.Tx.Hash(), "from", txWithSender.Sender.Hex(), "to", txWithSender.Tx.To())
		}

		if err := enclavedb.WriteTransactions(ctx, dbTx, transactionsWithSenders, batch.Header.Number.Uint64(), false, senderIds, toContractIds, toEoas, 0); err != nil {
			return fmt.Errorf("could not write transactions. Cause: %w", err)
		}
	}

	if err := dbTx.Commit(); err != nil {
		return fmt.Errorf("could not commit batch %w", err)
	}

	s.cachingService.CacheBatch(ctx, batch)
	return nil
}

func (s *storageImpl) handleTxSendersAndReceivers(ctx context.Context, transactionsWithSenders []*core.TxWithSender, dbTx *sqlx.Tx) ([]uint64, []*uint64, []*uint64, error) {
	senders := make([]uint64, len(transactionsWithSenders))
	toContracts := make([]*uint64, len(transactionsWithSenders))
	toEoas := make([]*uint64, len(transactionsWithSenders))
	// insert the tx signers as externally owned accounts
	for i, tx := range transactionsWithSenders {
		eoaID, err := s.readOrWriteEOA(ctx, dbTx, *tx.Sender)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("could not insert EOA. cause: %w", err)
		}
		s.logger.Trace("Tx sender", "tx", tx.Tx.Hash(), "sender", tx.Sender.Hex(), "eoaId", *eoaID)
		senders[i] = *eoaID

		to := tx.Tx.To()
		if to != nil {
			ctr, err := s.eventsStorage.ReadContract(ctx, *to)
			if err != nil && !errors.Is(err, errutil.ErrNotFound) {
				return nil, nil, nil, fmt.Errorf("could not read contract. cause: %w", err)
			}
			if ctr != nil {
				toContracts[i] = &ctr.Id
			} else {
				// if the "to" is not a contract, then it's an EOA
				toEoaID, err := s.readOrWriteEOA(ctx, dbTx, *to)
				if err != nil {
					return nil, nil, nil, fmt.Errorf("could not insert EOA. cause: %w", err)
				}
				toEoas[i] = toEoaID
			}
		}
	}
	return senders, toContracts, toEoas, nil
}

func (s *storageImpl) StoreExecutedBatch(ctx context.Context, batch *core.Batch, results core.TxExecResults) error {
	defer s.logDuration("StoreExecutedBatch", measure.NewStopwatch())
	executed, err := enclavedb.BatchWasExecuted(ctx, s.db.GetSQLDB(), batch.Hash())
	if err != nil {
		return err
	}
	if executed {
		s.logger.Debug("Batch was already executed", log.BatchHashKey, batch.Hash())
		return nil
	}

	s.logger.Trace("storing executed batch", log.BatchHashKey, batch.Hash(), log.BatchSeqNoKey, batch.Header.SequencerOrderNo, "receipts", len(results))

	dbTx, err := s.db.NewDBTransaction(ctx)
	if err != nil {
		return fmt.Errorf("could not create DB transaction - %w", err)
	}
	defer dbTx.Rollback()

	if err := enclavedb.MarkBatchExecuted(ctx, dbTx, batch.Header.SequencerOrderNo); err != nil {
		return fmt.Errorf("could not set the executed flag. Cause: %w", err)
	}

	// store the synthetic transactions
	syntheticTxs := results.SyntheticTransactions().ToTransactionsWithSenders()

	senders, toContracts, toEoas, err := s.handleTxSendersAndReceivers(ctx, syntheticTxs, dbTx)
	if err != nil {
		return fmt.Errorf("could not handle synthetic txs senders and receivers. Cause: %w", err)
	}

	if err := enclavedb.WriteTransactions(ctx, dbTx, syntheticTxs, batch.Header.Number.Uint64(), true, senders, toContracts, toEoas, len(batch.Transactions)); err != nil {
		return fmt.Errorf("could not write synthetic txs. Cause: %w", err)
	}

	if s.config.StoreExecutedTransactions {
		for _, txExecResult := range results {
			err = s.eventsStorage.storeReceiptAndEventLogs(ctx, dbTx, batch.Header, txExecResult)
			if err != nil {
				return fmt.Errorf("could not store receipt. Cause: %w", err)
			}
		}
	}

	if err = dbTx.Commit(); err != nil {
		return fmt.Errorf("could not commit batch %w", err)
	}

	// after a successful db commit, cache the receipts
	if s.config.StoreExecutedTransactions {
		s.cachingService.CacheReceipts(results)
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

func (s *storageImpl) StoreEnclaveKey(ctx context.Context, enclaveKey []byte) error {
	defer s.logDuration("StoreEnclaveKey", measure.NewStopwatch())
	if len(enclaveKey) == 0 {
		return errors.New("enclaveKey cannot be empty")
	}

	dbTx, err := s.db.NewDBTransaction(ctx)
	if err != nil {
		return fmt.Errorf("could not create DB transaction - %w", err)
	}
	defer dbTx.Rollback()
	_, err = storage.WriteConfig(ctx, dbTx, enclaveKeyCfg, enclaveKey)
	if err != nil {
		return err
	}
	return dbTx.Commit()
}

func (s *storageImpl) GetEnclaveKey(ctx context.Context) ([]byte, error) {
	defer s.logDuration("GetEnclaveKey", measure.NewStopwatch())
	return storage.FetchConfig(ctx, s.db.GetSQLDB(), enclaveKeyCfg)
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

func (s *storageImpl) DebugGetLogs(ctx context.Context, from *big.Int, to *big.Int, address gethcommon.Address, eventSig gethcommon.Hash) ([]*common.DebugLogVisibility, error) {
	defer s.logDuration("DebugGetLogs", measure.NewStopwatch())
	return enclavedb.DebugGetLogs(ctx, s.db.GetSQLDB(), from, to, address, eventSig)
}

const (
	MaxContractsPerFilter  = 5
	MaxEventTypesPerFilter = 15
)

var ErrInvalidFilter = errors.New("invalid filter")

func (s *storageImpl) FilterLogs(
	ctx context.Context,
	requestingAccount *gethcommon.Address,
	fromBlock, toBlock *big.Int,
	blockHash *common.L2BatchHash,
	contractAddresses []gethcommon.Address,
	topics [][]gethcommon.Hash,
) ([]*types.Log, error) {
	defer s.logDuration("FilterLogs", measure.NewStopwatch())

	var filterTopics enclavedb.FilterTopics = topics
	requestingId, err := s.readOrWriteEOAWithTx(ctx, *requestingAccount)
	if err != nil {
		return nil, err
	}

	noContracts := len(contractAddresses) == 0
	noEventTypes := !filterTopics.HasEventTypes()
	if noContracts && noEventTypes {
		return nil, fmt.Errorf("%w. The filter must have at least one contract address or one event type", ErrInvalidFilter)
	}

	if len(topics) > 4 {
		return nil, fmt.Errorf("%w. Too many topics", ErrInvalidFilter)
	}

	if len(contractAddresses) > MaxContractsPerFilter {
		return nil, fmt.Errorf("%w. Too many contract addresses. Max allowed: %d", ErrInvalidFilter, MaxContractsPerFilter)
	}

	if len(filterTopics.EventTypes()) > MaxEventTypesPerFilter {
		return nil, fmt.Errorf("%w. Too many event types. Max allowed: %d", ErrInvalidFilter, MaxEventTypesPerFilter)
	}

	// load the contracts from cache
	contracts := make([]*enclavedb.Contract, 0)
	for _, contractAddress := range contractAddresses {
		c, err := s.ReadContract(ctx, contractAddress)
		if err != nil && !errors.Is(err, errutil.ErrNotFound) {
			return nil, err
		}
		if c != nil {
			contracts = append(contracts, c)
		}
	}

	// there are no contracts for the requested addresses, so we can return early
	if !noContracts && len(contracts) == 0 {
		return nil, nil
	}

	eventTypes := make([]*enclavedb.EventType, 0)

	switch {
	case noContracts:
		// we load all contracts for which an event with the specified signature exists
		// we know there is at least one eventTypeHash from a previous check
		for _, eventTypeHash := range filterTopics.EventTypes() {
			// load all event types, and populate the contracts
			allEventTypes, err := s.readEventTypes(ctx, eventTypeHash)
			if err != nil && !errors.Is(err, errutil.ErrNotFound) {
				return nil, err
			}
			eventTypes = append(eventTypes, allEventTypes...)
		}

	case noEventTypes:
		if len(contracts) == 0 {
			s.logger.Crit("should not happen. there are no contracts for the requested addresses")
		}
		// if topics[0] were empty add all event Types found in the loaded contracts
		for _, contract := range contracts {
			if len(contract.EventTypes) > 0 {
				eventTypes = append(eventTypes, contract.EventTypeList()...)
			}
		}

	default:
		for _, contract := range contracts {
			for _, eventTypeHash := range filterTopics.EventTypes() {
				// load the event types for the specified contracts
				et := contract.EventType(eventTypeHash)
				if et != nil {
					eventTypes = append(eventTypes, et)
				}
			}
		}
	}

	// the specified contracts have no events, so we can return early
	if len(eventTypes) == 0 {
		return nil, nil
	}

	logs, err := enclavedb.FilterLogs(ctx, s.preparedStatementCache, requestingId, fromBlock, toBlock, blockHash, eventTypes, filterTopics)
	if err != nil {
		return nil, err
	}

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

func (s *storageImpl) MarkBatchAsUnexecuted(ctx context.Context, seqNo *big.Int) error {
	defer s.logDuration("MarkBatchAsUnexecuted", measure.NewStopwatch())
	dbTx, err := s.db.NewDBTransaction(ctx)
	if err != nil {
		return fmt.Errorf("could not create DB transaction - %w", err)
	}
	defer dbTx.Rollback()

	err = enclavedb.MarkBatchAsUnexecuted(ctx, dbTx, seqNo)
	if err != nil {
		return fmt.Errorf("could not mark batch as unexecuted. Cause: %w", err)
	}
	return dbTx.Commit()
}

func (s *storageImpl) GetTransactionsPerAddress(ctx context.Context, requester *gethcommon.Address, pagination *common.QueryPagination, showPublic bool, showSynthetic bool) ([]*core.InternalReceipt, error) {
	defer s.logDuration("GetTransactionsPerAddress", measure.NewStopwatch())
	requesterId, err := s.readOrWriteEOAWithTx(ctx, *requester)
	if err != nil {
		return nil, err
	}
	return enclavedb.GetTransactionsPerAddress(ctx, s.preparedStatementCache, requesterId, pagination, showPublic, showSynthetic)
}

func (s *storageImpl) CountTransactionsPerAddress(ctx context.Context, address *gethcommon.Address, showPublic bool, showSynthetic bool) (uint64, error) {
	defer s.logDuration("CountTransactionsPerAddress", measure.NewStopwatch())
	requesterId, err := s.readOrWriteEOAWithTx(ctx, *address)
	if err != nil {
		return 0, err
	}
	return enclavedb.CountTransactionsPerAddress(ctx, s.preparedStatementCache, requesterId, showPublic, showSynthetic)
}

func (s *storageImpl) readOrWriteEOAWithTx(ctx context.Context, addr gethcommon.Address) (*uint64, error) {
	dbTX, err := s.db.NewDBTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not create DB transaction - %w", err)
	}
	defer dbTX.Commit()

	return s.readOrWriteEOA(ctx, dbTX, addr)
}

func (s *storageImpl) readOrWriteEOA(ctx context.Context, dbTX *sqlx.Tx, addr gethcommon.Address) (*uint64, error) {
	defer s.logDuration("readOrWriteEOA", measure.NewStopwatch())
	return s.cachingService.ReadEOA(ctx, addr, func() (*uint64, error) {
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

func (s *storageImpl) ReadContract(ctx context.Context, address gethcommon.Address) (*enclavedb.Contract, error) {
	return s.eventsStorage.ReadContract(ctx, address)
}

func (s *storageImpl) readEventTypes(ctx context.Context, eventSignature gethcommon.Hash) ([]*enclavedb.EventType, error) {
	defer s.logDuration("readEventTypes", measure.NewStopwatch())
	dbTx, err := s.db.NewDBTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not create DB transaction - %w", err)
	}
	defer dbTx.Rollback()
	// todo - to use caching we have to update the cache when new types are added
	return enclavedb.ReadEventTypes(ctx, dbTx, eventSignature)
}

func (s *storageImpl) logDuration(method string, stopWatch *measure.Stopwatch) {
	core.LogMethodDuration(s.logger, stopWatch, fmt.Sprintf("Storage::%s completed", method))
}

func (s *storageImpl) StoreSystemContractAddresses(ctx context.Context, addresses common.SystemContractAddresses) error {
	defer s.logDuration("StoreSystemContractAddresses", measure.NewStopwatch())

	dbTx, err := s.db.NewDBTransaction(ctx)
	if err != nil {
		return fmt.Errorf("could not create DB transaction - %w", err)
	}
	defer dbTx.Rollback()

	addressesBytes, err := json.Marshal(addresses)
	if err != nil {
		return fmt.Errorf("could not marshal system contract addresses - %w", err)
	}
	_, err = storage.WriteConfig(ctx, dbTx, systemContractAddressesCfg, addressesBytes)
	if err != nil {
		return fmt.Errorf("could not write system contract addresses - %w", err)
	}

	if err := dbTx.Commit(); err != nil {
		return fmt.Errorf("could not commit system contract addresses - %w", err)
	}
	return nil
}

func (s *storageImpl) GetSystemContractAddresses(ctx context.Context) (common.SystemContractAddresses, error) {
	defer s.logDuration("GetSystemContractAddresses", measure.NewStopwatch())
	addressesBytes, err := storage.FetchConfig(ctx, s.db.GetSQLDB(), systemContractAddressesCfg)
	if err != nil {
		return nil, fmt.Errorf("could not fetch system contract addresses - %w", err)
	}
	var addresses common.SystemContractAddresses
	if err := json.Unmarshal(addressesBytes, &addresses); err != nil {
		return nil, fmt.Errorf("could not unmarshal system contract addresses - %w", err)
	}
	return addresses, nil
}

func (s *storageImpl) GetSequencerEnclaveIDs(ctx context.Context) ([]common.EnclaveID, error) {
	defer s.logDuration("GetSequencerEnclaveIDs", measure.NewStopwatch())

	ids, err := s.cachingService.ReadSequencerIDs(ctx, func() ([]common.EnclaveID, error) {
		sequencerIDs, err := enclavedb.FetchSequencerEnclaveIDs(ctx, s.db.GetSQLDB())
		if err != nil {
			return nil, fmt.Errorf("failed to fetch sequencer IDs from database. Cause: %w", err)
		}
		return sequencerIDs, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read sequencer IDs from cache. Cause: %w", err)
	}
	return ids, nil
}

// NetworkUpgradeStorage implementation

// Removed: StorePendingNetworkUpgrade, FinalizeNetworkUpgrade, GetPendingNetworkUpgrades, GetFinalizedNetworkUpgrades
// Upgrades are disabled; methods deleted.
