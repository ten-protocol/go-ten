package storage

import (
	"errors"
	"fmt"
	"io"
	"math/big"
	"strings"

	smt "github.com/FantasyJony/openzeppelin-merkle-tree-go/standard_merkle_tree"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/log"
	hostconfig "github.com/ten-protocol/go-ten/go/host/config"
	"github.com/ten-protocol/go-ten/go/host/storage/hostdb"
)

type storageImpl struct {
	db     hostdb.HostDB
	logger gethlog.Logger
	io.Closer
}

func (s *storageImpl) AddBatch(batch *common.ExtBatch) error {
	_, err := hostdb.GetBatchHeader(s.db, batch.Hash())
	if err == nil {
		// batch already exists don't error
		return nil
	}

	dbtx, err := s.db.NewDBTransaction()
	if err != nil {
		return err
	}
	defer dbtx.Rollback()

	if err := hostdb.AddBatch(dbtx, s.db, batch); err != nil {
		if errors.Is(err, errutil.ErrAlreadyExists) {
			return nil
		}
		return fmt.Errorf("could not add batch to host. Cause: %w", err)
	}

	if err := dbtx.Write(); err != nil {
		if hostdb.IsRowExistsError(err) {
			return nil
		}
		return fmt.Errorf("could not commit batch tx. Cause: %w", err)
	}
	return nil
}

func (s *storageImpl) AddRollup(rollup *common.ExtRollup, extMetadata *common.ExtRollupMetadata, metadata *common.PublicRollupMetadata, block *types.Header) error {
	_, err := hostdb.GetRollupHeader(s.db, rollup.Header.Hash())
	if err == nil {
		// rollup already exists don't error
		return nil
	}

	dbtx, err := s.db.NewDBTransaction()
	if err != nil {
		return err
	}
	defer dbtx.Rollback()

	if err := hostdb.AddRollup(dbtx, s.db, rollup, extMetadata, metadata, block); err != nil {
		if errors.Is(err, errutil.ErrAlreadyExists) {
			return nil
		}
		return fmt.Errorf("could not add rollup to host. Cause: %w", err)
	}

	if err := dbtx.Write(); err != nil {
		if hostdb.IsRowExistsError(err) {
			return nil
		}
		return fmt.Errorf("could not commit rollup tx. Cause %w", err)
	}
	return nil
}

func (s *storageImpl) ReadBlock(blockHash *gethcommon.Hash) (*types.Header, error) {
	return hostdb.GetBlock(s.db, blockHash)
}

func (s *storageImpl) AddBlock(b *types.Header) error {
	dbtx, err := s.db.NewDBTransaction()
	if err != nil {
		return fmt.Errorf("could not create DB transaction - %w", err)
	}
	defer dbtx.Rollback()

	_, err = hostdb.GetBlockId(dbtx.Tx, s.db.GetSQLDB(), b.Hash())
	switch {
	case err == nil:
		// Block already exists
		s.logger.Debug("Block already exists", "hash", b.Hash().Hex())
		return nil
	case !errors.Is(err, errutil.ErrNotFound):
		return fmt.Errorf("error checking block existence: %w", err)
	}

	if err := hostdb.AddBlock(dbtx.Tx, s.db.GetSQLDB(), b); err != nil {
		if IsConstraintError(err) {
			s.logger.Debug("Block already exists",
				"hash", b.Hash().Hex(),
				"error", err)
			return nil
		}
		return fmt.Errorf("could not add block to host: %w", err)
	}

	if err := dbtx.Write(); err != nil {
		return fmt.Errorf("could not commit block tx. Cause %w", err)
	}
	return nil
}

func (s *storageImpl) FetchCrossChainProof(messageType string, crossChainMessage gethcommon.Hash) ([][]byte, gethcommon.Hash, error) {
	tree, err := hostdb.GetCrossChainMessagesTree(s.db, crossChainMessage)
	if err != nil {
		return nil, gethcommon.Hash{}, err
	}

	for k, value := range tree {
		tree[k][1] = gethcommon.BytesToHash(value[1].([]byte))
	}

	merkleTree, err := smt.Of(tree, []string{smt.SOL_STRING, smt.SOL_BYTES32})
	if err != nil {
		return nil, gethcommon.Hash{}, err
	}
	proof, err := merkleTree.GetProof([]interface{}{messageType, crossChainMessage})
	if err != nil {
		return nil, gethcommon.Hash{}, err
	}
	return proof, gethcommon.Hash(merkleTree.GetRoot()), nil
}

func (s *storageImpl) FetchBatchBySeqNo(seqNum uint64) (*common.ExtBatch, error) {
	return hostdb.GetBatchBySequenceNumber(s.db, seqNum)
}

func (s *storageImpl) FetchBatchHashByHeight(number *big.Int) (*gethcommon.Hash, error) {
	return hostdb.GetBatchHashByNumber(s.db, number)
}

func (s *storageImpl) FetchBatchHeaderByHash(hash gethcommon.Hash) (*common.BatchHeader, error) {
	return hostdb.GetBatchHeader(s.db, hash)
}

func (s *storageImpl) FetchHeadBatchHeader() (*common.BatchHeader, error) {
	return hostdb.GetHeadBatchHeader(s.db)
}

func (s *storageImpl) FetchPublicBatchByHash(batchHash common.L2BatchHash) (*common.PublicBatch, error) {
	return hostdb.GetPublicBatch(s.db, batchHash)
}

func (s *storageImpl) FetchBatch(batchHash gethcommon.Hash) (*common.PublicBatch, error) {
	return hostdb.GetBatchByHash(s.db, batchHash)
}

func (s *storageImpl) FetchBatchByTx(txHash gethcommon.Hash) (*common.PublicBatch, error) {
	return hostdb.GetBatchByTx(s.db, txHash)
}

func (s *storageImpl) FetchLatestBatch() (*common.BatchHeader, error) {
	return hostdb.GetHeadBatchHeader(s.db)
}

func (s *storageImpl) FetchBatchHeaderByHeight(height *big.Int) (*common.BatchHeader, error) {
	return hostdb.GetBatchHeaderByHeight(s.db, height)
}

func (s *storageImpl) FetchPublicBatchByHeight(height *big.Int) (*common.PublicBatch, error) {
	return hostdb.GetBatchByHeight(s.db, height)
}

func (s *storageImpl) FetchPublicBatchBySeqNo(seqNum *big.Int) (*common.PublicBatch, error) {
	return hostdb.GetPublicBatchBySequenceNumber(s.db, seqNum.Uint64())
}

func (s *storageImpl) FetchBatchListing(pagination *common.QueryPagination) (*common.BatchListingResponse, error) {
	return hostdb.GetBatchListing(s.db, pagination)
}

func (s *storageImpl) FetchLatestRollupHeader() (*common.RollupHeader, error) {
	return hostdb.GetLatestRollup(s.db)
}

func (s *storageImpl) FetchRollupListing(pagination *common.QueryPagination) (*common.RollupListingResponse, error) {
	return hostdb.GetRollupListing(s.db, pagination)
}

func (s *storageImpl) FetchBlockListing(pagination *common.QueryPagination) (*common.BlockListingResponse, error) {
	return hostdb.GetBlockListing(s.db, pagination)
}

func (s *storageImpl) FetchTotalTxCount() (*big.Int, error) {
	return hostdb.GetTotalTxCount(s.db)
}

func (s *storageImpl) FetchTotalTxsQuery() (*big.Int, error) {
	return hostdb.GetTotalTxsQuery(s.db)
}

func (s *storageImpl) FetchHistoricalTransactionCount() (*big.Int, error) {
	return hostdb.GetHistoricalTransactionCount(s.db)
}

func (s *storageImpl) FetchTransaction(hash gethcommon.Hash) (*common.PublicTransaction, error) {
	return hostdb.GetTransaction(s.db, hash)
}

func (s *storageImpl) FetchRollupByHash(rollupHash gethcommon.Hash) (*common.PublicRollup, error) {
	return hostdb.GetRollupByHash(s.db, rollupHash)
}

func (s *storageImpl) FetchRollupBySeqNo(seqNo uint64) (*common.PublicRollup, error) {
	return hostdb.GetRollupBySeqNo(s.db, seqNo)
}

func (s *storageImpl) FetchRollupBatches(rollupHash gethcommon.Hash, pagination *common.QueryPagination) (*common.BatchListingResponse, error) {
	return hostdb.GetRollupBatches(s.db, rollupHash, pagination)
}

func (s *storageImpl) FetchBatchTransactions(batchHash gethcommon.Hash, pagination *common.QueryPagination) (*common.TransactionListingResponse, error) {
	return hostdb.GetBatchTransactions(s.db, batchHash, pagination)
}

func (s *storageImpl) FetchTransactionListing(pagination *common.QueryPagination) (*common.TransactionListingResponse, error) {
	return hostdb.GetTransactionListing(s.db, pagination)
}

func (s *storageImpl) EstimateRollupSize(fromSeqNo *big.Int) (uint64, error) {
	return hostdb.EstimateRollupSize(s.db, fromSeqNo)
}

func (s *storageImpl) Search(query string) (*common.SearchResponse, error) {
	return hostdb.Search(s.db, query)
}

func (s *storageImpl) Close() error {
	return s.db.GetSQLDB().Close()
}

func NewHostStorageFromConfig(config *hostconfig.HostConfig, logger gethlog.Logger) Storage {
	backingDB, err := CreateDBFromConfig(config, logger)
	if err != nil {
		logger.Crit("Failed to connect to backing database", log.ErrKey, err)
	}
	return NewStorage(backingDB, logger)
}

func NewStorage(backingDB hostdb.HostDB, logger gethlog.Logger) Storage {
	return &storageImpl{
		db:     backingDB,
		logger: logger,
	}
}

// SQLite constraint error messages
const (
	ErrUniqueBlockHash = "UNIQUE constraint failed: block_host.hash"
	ErrForeignKey      = "FOREIGN KEY constraint failed"
)

// IsConstraintError returns true if the error is a known constraint error
func IsConstraintError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return strings.Contains(errMsg, ErrUniqueBlockHash) ||
		strings.Contains(errMsg, ErrForeignKey)
}
