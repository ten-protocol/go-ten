package storage

import (
	"database/sql"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/config"
	"github.com/ten-protocol/go-ten/go/host/storage/hostdb"
	"io"
)

type storageImpl struct {
	db     hostdb.HostDB
	logger gethlog.Logger
	io.Closer
}

func (s *storageImpl) AddBatch(batch *common.ExtBatch) error {
	// Check if the Batch is already stored
	//_, err := hostdb.GetBatchHeader(s.db.GetDB(), batch.Hash())
	_, err := hostdb.GetBatchBySequenceNumber(s.db.GetDB(), batch.SeqNo().Uint64())
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
	_, err := hostdb.GetRollupHeader(s.db.GetDB(), rollup.Header.Hash())
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
	return hostdb.GetBatchBySequenceNumber(s.db.GetDB(), seqNum)
}

func (s *storageImpl) GetDB() *sql.DB {
	return s.db.GetDB()
}

func (s *storageImpl) Close() error {
	return s.db.GetDB().Close()
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
