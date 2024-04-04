package storage

import (
	"fmt"
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
	// Check if the Batch is already stored
	_, err := hostdb.GetBatchHeader(s.db.NewDBTransaction(), batch.Hash())
	if err == nil {
		return errutil.ErrAlreadyExists
	}

	dbTx := s.db.NewDBTransaction()
	if err := hostdb.AddBatch(dbTx, batch); err != nil {
		return fmt.Errorf("could not write batch. Cause: %w", err)
	}
	if err := dbTx.Write(); err != nil {
		return fmt.Errorf("could not commit batch %w", err)
	}
	return nil
}

func (s *storageImpl) AddRollup(rollup *common.ExtRollup, metadata *common.PublicRollupMetadata, block *common.L1Block) error {
	// Check if the Header is already stored
	_, err := hostdb.GetRollupHeader(s.db.NewDBTransaction(), rollup.Header.Hash())
	if err == nil {
		return errutil.ErrAlreadyExists
	}
	dbTx := s.db.NewDBTransaction()
	if err := hostdb.AddRollup(dbTx, rollup, metadata, block); err != nil {
		return fmt.Errorf("could not write batch. Cause: %w", err)
	}
	if err := dbTx.Write(); err != nil {
		return fmt.Errorf("could not commit batch %w", err)
	}
	return nil
}

func (s *storageImpl) AddBlock(b *types.Header, rollupHash common.L2RollupHash) error {
	dbTx := s.db.NewDBTransaction()
	if err := hostdb.AddBlock(dbTx, b, rollupHash); err != nil {
		return fmt.Errorf("could not write batch. Cause: %w", err)
	}
	if err := dbTx.Write(); err != nil {
		return fmt.Errorf("could not commit batch %w", err)
	}
	return nil
}

func (s *storageImpl) FetchBatchBySeqNo(seqNum uint64) (*common.ExtBatch, error) {
	return hostdb.GetBatchBySequenceNumber(s.db.NewDBTransaction(), seqNum)
}

func (s *storageImpl) FetchBatchHashByNumber(number *big.Int) (*gethcommon.Hash, error) {
	return hostdb.GetBatchHashByNumber(s.db.NewDBTransaction(), number)
}

func (s *storageImpl) FetchBatchHeaderByHash(hash gethcommon.Hash) (*common.BatchHeader, error) {
	return hostdb.GetBatchHeader(s.db.NewDBTransaction(), hash)
}

func (s *storageImpl) FetchHeadBatchHeader() (*common.BatchHeader, error) {
	return hostdb.GetHeadBatchHeader(s.db.GetSQLDB())
}

func (s *storageImpl) GetDB() hostdb.HostDB {
	return s.db
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
