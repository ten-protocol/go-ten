package storage

import (
	"fmt"
	"github.com/ten-protocol/go-ten/go/common/measure"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"io"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/config"
	"github.com/ten-protocol/go-ten/go/host/storage/hostdb"
)

type storageImpl struct {
	db     hostdb.HostDB
	logger gethlog.Logger
	io.Closer
}

func (s *storageImpl) AddBatch(batch *common.ExtBatch) error {
	defer core.LogMethodDuration(s.logger, measure.NewStopwatch(), "Add batch", log.BatchHashKey, batch.Hash())
	// Check if the Batch is already stored
	_, err := hostdb.GetBatchHeader(s.db, batch.Hash())
	if err == nil {
		return errutil.ErrAlreadyExists
	}

	dbtx, err := s.db.NewDBTransaction()
	if err != nil {
		return err
	}

	if err := hostdb.AddBatch(dbtx, s.db.GetSQLStatement(), batch); err != nil {
		if err := dbtx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("could not add batch to host. Cause: %w", err)
	}

	if err := dbtx.Write(); err != nil {
		return fmt.Errorf("could not commit batch tx. Cause: %w", err)
	}
	return nil
}

func (s *storageImpl) AddRollup(rollup *common.ExtRollup, metadata *common.PublicRollupMetadata, block *common.L1Block) error {
	defer core.LogMethodDuration(s.logger, measure.NewStopwatch(), "Add rollup", log.RollupHashKey, rollup.Hash())
	// Check if the Header is already stored
	_, err := hostdb.GetRollupHeader(s.db, rollup.Header.Hash())
	if err == nil {
		return errutil.ErrAlreadyExists
	}

	dbtx, err := s.db.NewDBTransaction()
	if err != nil {
		return err
	}

	if err := hostdb.AddRollup(dbtx, s.db.GetSQLStatement(), rollup, metadata, block); err != nil {
		if err := dbtx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("could not add rollup to host. Cause: %w", err)
	}

	if err := dbtx.Write(); err != nil {
		return fmt.Errorf("could not commit rollup tx. Cause %w", err)
	}
	return nil
}

func (s *storageImpl) AddBlock(b *types.Header, rollupHash common.L2RollupHash) error {
	defer core.LogMethodDuration(s.logger, measure.NewStopwatch(), "Add block", log.BlockHashKey, b.Hash())
	dbtx, err := s.db.NewDBTransaction()
	if err != nil {
		return err
	}

	if err := hostdb.AddBlock(dbtx, s.db.GetSQLStatement(), b, rollupHash); err != nil {
		if err := dbtx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("could not add block to host. Cause: %w", err)
	}

	if err := dbtx.Write(); err != nil {
		return fmt.Errorf("could not commit block tx. Cause %w", err)
	}
	return nil
}

func (s *storageImpl) FetchBatchBySeqNo(seqNum uint64) (*common.ExtBatch, error) {
	defer core.LogMethodDuration(s.logger, measure.NewStopwatch(), "Fetch batch by seq no", log.BatchSeqNoKey, seqNum)
	return hostdb.GetBatchBySequenceNumber(s.db, seqNum)
}

func (s *storageImpl) FetchBatchHashByHeight(number *big.Int) (*gethcommon.Hash, error) {
	defer core.LogMethodDuration(s.logger, measure.NewStopwatch(), "Fetch batch by height", log.BatchHeightKey, number.Uint64())
	return hostdb.GetBatchHashByNumber(s.db, number)
}

func (s *storageImpl) FetchBatchHeaderByHash(hash gethcommon.Hash) (*common.BatchHeader, error) {
	defer core.LogMethodDuration(s.logger, measure.NewStopwatch(), "Fetch batch by hash", log.BatchHashKey, hash)
	return hostdb.GetBatchHeader(s.db, hash)
}

func (s *storageImpl) FetchHeadBatchHeader() (*common.BatchHeader, error) {
	defer core.LogMethodDuration(s.logger, measure.NewStopwatch(), "Fetch head batch header")
	return hostdb.GetHeadBatchHeader(s.db)
}

func (s *storageImpl) FetchPublicBatchByHash(batchHash common.L2BatchHash) (*common.PublicBatch, error) {
	defer core.LogMethodDuration(s.logger, measure.NewStopwatch(), "Fetch public batch by hash", log.BatchHashKey, batchHash)
	return hostdb.GetPublicBatch(s.db, batchHash)
}

func (s *storageImpl) FetchBatch(batchHash gethcommon.Hash) (*common.ExtBatch, error) {
	defer core.LogMethodDuration(s.logger, measure.NewStopwatch(), "Fetch batch by hash", log.BatchHashKey, batchHash)
	return hostdb.GetBatchByHash(s.db, batchHash)
}

func (s *storageImpl) FetchBatchByTx(txHash gethcommon.Hash) (*common.ExtBatch, error) {
	defer core.LogMethodDuration(s.logger, measure.NewStopwatch(), "Fetch batch by transaction", log.TxKey, txHash)
	return hostdb.GetBatchByTx(s.db, txHash)
}

func (s *storageImpl) FetchLatestBatch() (*common.BatchHeader, error) {
	defer core.LogMethodDuration(s.logger, measure.NewStopwatch(), "Fetch latest batch")
	return hostdb.GetLatestBatch(s.db)
}

func (s *storageImpl) FetchBatchHeaderByHeight(height *big.Int) (*common.BatchHeader, error) {
	defer core.LogMethodDuration(s.logger, measure.NewStopwatch(), "Fetch batch header by height", log.BatchHeightKey, height.Uint64())
	return hostdb.GetBatchHeaderByHeight(s.db, height)
}

func (s *storageImpl) FetchBatchByHeight(height *big.Int) (*common.PublicBatch, error) {
	defer core.LogMethodDuration(s.logger, measure.NewStopwatch(), "Fetch batch by height", log.BatchHeightKey, height.Uint64())
	return hostdb.GetBatchByHeight(s.db, height)
}

func (s *storageImpl) FetchBatchListing(pagination *common.QueryPagination) (*common.BatchListingResponse, error) {
	defer core.LogMethodDuration(s.logger, measure.NewStopwatch(), "Fetch batch listing")
	return hostdb.GetBatchListing(s.db, pagination)
}

func (s *storageImpl) FetchBatchListingDeprecated(pagination *common.QueryPagination) (*common.BatchListingResponseDeprecated, error) {
	defer core.LogMethodDuration(s.logger, measure.NewStopwatch(), "Fetch batch listing (deprecated)")
	return hostdb.GetBatchListingDeprecated(s.db, pagination)
}

func (s *storageImpl) FetchLatestRollupHeader() (*common.RollupHeader, error) {
	defer core.LogMethodDuration(s.logger, measure.NewStopwatch(), "Fetch latest rollup header")
	return hostdb.GetLatestRollup(s.db)
}

func (s *storageImpl) FetchRollupListing(pagination *common.QueryPagination) (*common.RollupListingResponse, error) {
	defer core.LogMethodDuration(s.logger, measure.NewStopwatch(), "Fetch rollup listing")
	return hostdb.GetRollupListing(s.db, pagination)
}

func (s *storageImpl) FetchBlockListing(pagination *common.QueryPagination) (*common.BlockListingResponse, error) {
	defer core.LogMethodDuration(s.logger, measure.NewStopwatch(), "Fetch block listing")
	return hostdb.GetBlockListing(s.db, pagination)
}

func (s *storageImpl) FetchTotalTxCount() (*big.Int, error) {
	defer core.LogMethodDuration(s.logger, measure.NewStopwatch(), "Fetch total transaction count")
	return hostdb.GetTotalTxCount(s.db)
}

func (s *storageImpl) FetchTransaction(hash gethcommon.Hash) (*common.PublicTransaction, error) {
	defer core.LogMethodDuration(s.logger, measure.NewStopwatch(), "Fetch transaction by hash", log.TxKey, hash)
	return hostdb.GetTransaction(s.db, hash)
}

func (s *storageImpl) FetchRollupByHash(rollupHash gethcommon.Hash) (*common.PublicRollup, error) {
	defer core.LogMethodDuration(s.logger, measure.NewStopwatch(), "Fetch rollup by hash", log.RollupHashKey, rollupHash)
	return hostdb.GetRollupByHash(s.db, rollupHash)
}

func (s *storageImpl) FetchRollupBySeqNo(seqNo uint64) (*common.PublicRollup, error) {
	defer core.LogMethodDuration(s.logger, measure.NewStopwatch(), "Fetch rollup by sequence number", log.RollupHashKey, seqNo)
	return hostdb.GetRollupBySeqNo(s.db, seqNo)
}

func (s *storageImpl) FetchRollupBatches(rollupHash gethcommon.Hash) (*common.BatchListingResponse, error) {
	defer core.LogMethodDuration(s.logger, measure.NewStopwatch(), "Fetch rollup batches", log.RollupHashKey, rollupHash)
	return hostdb.GetRollupBatches(s.db, rollupHash)
}

func (s *storageImpl) FetchBatchTransactions(batchHash gethcommon.Hash) (*common.TransactionListingResponse, error) {
	defer core.LogMethodDuration(s.logger, measure.NewStopwatch(), "Fetch batch transactions", log.BatchHashKey, batchHash)
	return hostdb.GetBatchTransactions(s.db, batchHash)
}

func (s *storageImpl) FetchTransactionListing(pagination *common.QueryPagination) (*common.TransactionListingResponse, error) {
	defer core.LogMethodDuration(s.logger, measure.NewStopwatch(), "Fetch transaction listing")
	return hostdb.GetTransactionListing(s.db, pagination)
}

func (s *storageImpl) Close() error {
	return s.db.GetSQLDB().Close()
}

func NewHostStorageFromConfig(config *config.HostConfig, logger gethlog.Logger) Storage {
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
