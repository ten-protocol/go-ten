package enclavedb

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/ethereum/go-ethereum/ethdb"
	gethlog "github.com/ethereum/go-ethereum/log"
	enclaveconfig "github.com/ten-protocol/go-ten/go/enclave/config"
)

// enclaveDB - Implements the key-value ethdb.Database and also exposes the underlying sql database
// should not be used directly outside the db package
type enclaveDB struct {
	sqldb   *sqlx.DB
	rwSqldb *sqlx.DB // required only by sqlite. For a normal db, it will be the same instance as sqldb
	config  *enclaveconfig.EnclaveConfig
	logger  gethlog.Logger
}

func (sqlDB *enclaveDB) SyncKeyValue() error {
	// do nothing because we use db transactions
	return nil
}

func (sqlDB *enclaveDB) SyncAncient() error {
	// TODO implement me
	panic("implement me3")
}

func (sqlDB *enclaveDB) DeleteRange(start, end []byte) error {
	// TODO implement me
	panic("implement me4")
}

func (sqlDB *enclaveDB) Tail() (uint64, error) {
	// TODO implement me
	panic("implement me5")
}

func (sqlDB *enclaveDB) TruncateHead(uint64) (uint64, error) {
	// TODO implement me
	panic("implement me6")
}

func (sqlDB *enclaveDB) TruncateTail(uint64) (uint64, error) {
	// TODO implement me
	panic("implement me7")
}

func (sqlDB *enclaveDB) MigrateTable(string, func([]byte) ([]byte, error)) error {
	// TODO implement me
	panic("implement me8")
}

func (sqlDB *enclaveDB) NewBatchWithSize(int) ethdb.Batch {
	return &dbTxBatch{
		timeout: sqlDB.config.RPCTimeout,
		db:      sqlDB,
	}
}

func (sqlDB *enclaveDB) AncientDatadir() (string, error) {
	return "", fmt.Errorf("not implemented")
}

func NewEnclaveDB(db *sqlx.DB, rwdb *sqlx.DB, config *enclaveconfig.EnclaveConfig, logger gethlog.Logger) (EnclaveDB, error) {
	return &enclaveDB{sqldb: db, rwSqldb: rwdb, config: config, logger: logger}, nil
}

func (sqlDB *enclaveDB) GetSQLDB() *sqlx.DB {
	return sqlDB.sqldb
}

func (sqlDB *enclaveDB) Has(key []byte) (bool, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), sqlDB.config.RPCTimeout)
	defer cancelCtx()
	return Has(ctx, sqlDB.sqldb, key)
}

func (sqlDB *enclaveDB) Get(key []byte) ([]byte, error) {
	// ctx, cancelCtx := context.WithTimeout(context.Background(), sqlDB.config.RPCTimeout)
	// defer cancelCtx()
	val, err := Get(context.Background(), sqlDB.sqldb, key)

	trieJournalKey := []byte("vTrieJournal")
	if bytes.Equal(key, trieJournalKey) {
		sqlDB.logger.Debug("TrieJournal GET", "key", key, "err", err, " len_val", len(val))
	}

	return val, err
}

func (sqlDB *enclaveDB) Put(key []byte, value []byte) error {
	if key == nil {
		return errors.New("key cannot be nil")
	}
	if value == nil {
		return fmt.Errorf("value cannot be nil. key: %x", key)
	}
	// ctx, cancelCtx := context.WithTimeout(context.Background(), sqlDB.config.RPCTimeout)
	// defer cancelCtx()
	err := Put(context.Background(), sqlDB.rwSqldb, key, value)
	trieJournalKey := []byte("vTrieJournal")
	if bytes.Equal(key, trieJournalKey) {
		sqlDB.logger.Debug("TrieJournal PUT", "key", key, "err", err, "len_val", len(value))
		_, err := sqlDB.Get(trieJournalKey)
		if err != nil {
			sqlDB.logger.Crit("TrieJournal GET failed", "key", key, "err", err)
		}
	}
	return err
}

func (sqlDB *enclaveDB) Delete(key []byte) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), sqlDB.config.RPCTimeout)
	defer cancelCtx()
	return Delete(ctx, sqlDB.rwSqldb, key)
}

func (sqlDB *enclaveDB) Close() error {
	if err := sqlDB.sqldb.Close(); err != nil {
		return fmt.Errorf("failed to close sql db - %w", err)
	}
	return nil
}

func (sqlDB *enclaveDB) NewDBTransaction(ctx context.Context) (*sqlx.Tx, error) {
	tx, err := sqlDB.rwSqldb.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to create db transaction - %w", err)
	}
	return tx, nil
}

func (sqlDB *enclaveDB) NewBatch() ethdb.Batch {
	return &dbTxBatch{
		timeout: sqlDB.config.RPCTimeout,
		db:      sqlDB,
	}
}

func (sqlDB *enclaveDB) NewIterator(prefix []byte, start []byte) ethdb.Iterator {
	// we can't use a timeout context here, because the cleanup function must be called
	return NewIterator(context.Background(), sqlDB.sqldb, prefix, start)
}

func (sqlDB *enclaveDB) Stat() (string, error) {
	// todo - implement me
	sqlDB.logger.Crit("implement me")
	return "", nil
}

func (sqlDB *enclaveDB) Compact(_ []byte, _ []byte) error {
	// todo - implement me
	sqlDB.logger.Crit("implement me")
	return nil
}

// no-freeze! Copied from the geth in-memory sqldb implementation these ancient method implementations disable the 'freezer'

// errNotSupported is returned if the database doesn't support the required operation.
var errNotSupported = errors.New("this operation is not supported")

// HasAncient returns an error as we don't have a backing chain freezer.
func (sqlDB *enclaveDB) HasAncient(_ string, _ uint64) (bool, error) {
	return false, errNotSupported
}

// Ancient returns an error as we don't have a backing chain freezer.
func (sqlDB *enclaveDB) Ancient(_ string, _ uint64) ([]byte, error) {
	return nil, errNotSupported
}

// AncientRange returns an error as we don't have a backing chain freezer.
func (sqlDB *enclaveDB) AncientRange(_ string, _, _, _ uint64) ([][]byte, error) {
	return nil, errNotSupported
}

// Ancients returns an error as we don't have a backing chain freezer.
func (sqlDB *enclaveDB) Ancients() (uint64, error) {
	return 0, errNotSupported
}

// AncientSize returns an error as we don't have a backing chain freezer.
func (sqlDB *enclaveDB) AncientSize(_ string) (uint64, error) {
	return 0, errNotSupported
}

// ModifyAncients is not supported.
func (sqlDB *enclaveDB) ModifyAncients(func(ethdb.AncientWriteOp) error) (int64, error) {
	return 0, errNotSupported
}

// TruncateAncients returns an error as we don't have a backing chain freezer.
func (sqlDB *enclaveDB) TruncateAncients(_ uint64) error {
	return errNotSupported
}

// Sync returns an error as we don't have a backing chain freezer.
func (sqlDB *enclaveDB) Sync() error {
	return errNotSupported
}

func (sqlDB *enclaveDB) ReadAncients(fn func(reader ethdb.AncientReaderOp) error) (err error) {
	// Unlike other ancient-related methods, this method does not return
	// errNotSupported when invoked.
	// The reason for this is that the caller might want to do several things:
	// 1. Check if something is in freezer,
	// 2. If not, check leveldb.
	//
	// This will work, since the ancient-checks inside 'fn' will return errors,
	// and the leveldb work will continue.
	//
	// If we instead were to return errNotSupported here, then the caller would
	// have to explicitly check for that, having an extra clause to do the
	// non-ancient operations.
	return fn(sqlDB)
}

func (sqlDB *enclaveDB) AncientBytes(kind string, id, offset, length uint64) ([]byte, error) {
	// TODO implement me
	panic("implement me")
}
