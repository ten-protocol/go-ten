package sql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/ethdb"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
)

const (
	getQry = `select keyvalue.val from keyvalue where keyvalue.ky = ?;`
	// `replace` will perform insert or replace if existing and this syntax works for both sqlite and edgeless db
	putQry    = `replace into keyvalue values(?, ?);`
	delQry    = `delete from keyvalue where keyvalue.ky = ?;`
	searchQry = `select * from keyvalue where substring(keyvalue.ky, 1, ?) = ? and keyvalue.ky >= ? order by keyvalue.ky asc`
)

// EnclaveDB - Implements the key-value ethdb.Database and also exposes the underlying sql database
type EnclaveDB struct {
	db     *sql.DB
	logger gethlog.Logger
}

func CreateSQLEthDatabase(db *sql.DB, logger gethlog.Logger) (*EnclaveDB, error) {
	return &EnclaveDB{db: db, logger: logger}, nil
}

func (sqlDB *EnclaveDB) GetSQLDB() *sql.DB {
	return sqlDB.db
}

func (sqlDB *EnclaveDB) Has(key []byte) (bool, error) {
	err := sqlDB.db.QueryRow(getQry, key).Scan()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (sqlDB *EnclaveDB) Get(key []byte) ([]byte, error) {
	var res []byte

	err := sqlDB.db.QueryRow(getQry, key).Scan(&res)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// make sure the error is converted to obscuro-wide not found error
			return nil, errutil.ErrNotFound
		}
		return nil, err
	}
	return res, nil
}

func (sqlDB *EnclaveDB) Put(key []byte, value []byte) error {
	_, err := sqlDB.db.Exec(putQry, key, value)
	return err
}

func (sqlDB *EnclaveDB) Delete(key []byte) error {
	_, err := sqlDB.db.Exec(delQry, key)
	return err
}

func (sqlDB *EnclaveDB) Close() error {
	if err := sqlDB.db.Close(); err != nil {
		return fmt.Errorf("failed to close sql db - %w", err)
	}
	return nil
}

func (sqlDB *EnclaveDB) NewBatch() ethdb.Batch {
	return &sqlBatch{
		db: sqlDB,
	}
}

func (sqlDB *EnclaveDB) NewIterator(prefix []byte, start []byte) ethdb.Iterator {
	pr := prefix
	st := append(prefix, start...)
	// iterator clean-up handles closing this rows iterator
	rows, err := sqlDB.db.Query(searchQry, len(pr), pr, st)
	if err != nil {
		return &iterator{
			err: fmt.Errorf("failed to get rows, iter will be empty, %w", err),
		}
	}
	if err = rows.Err(); err != nil {
		return &iterator{
			err: fmt.Errorf("failed to get rows, iter will be empty, %w", err),
		}
	}
	return &iterator{
		rows: rows,
	}
}

func (sqlDB *EnclaveDB) Stat(property string) (string, error) {
	// todo - implement me
	sqlDB.logger.Crit("implement me")
	return "", nil
}

func (sqlDB *EnclaveDB) Compact(start []byte, limit []byte) error {
	// todo - implement me
	sqlDB.logger.Crit("implement me")
	return nil
}

// no-freeze! Copied from the geth in-memory db implementation these ancient method implementations disable the 'freezer'

// errNotSupported is returned if the database doesn't support the required operation.
var errNotSupported = errors.New("this operation is not supported")

// HasAncient returns an error as we don't have a backing chain freezer.
func (sqlDB *EnclaveDB) HasAncient(kind string, number uint64) (bool, error) {
	return false, errNotSupported
}

// Ancient returns an error as we don't have a backing chain freezer.
func (sqlDB *EnclaveDB) Ancient(kind string, number uint64) ([]byte, error) {
	return nil, errNotSupported
}

// AncientRange returns an error as we don't have a backing chain freezer.
func (sqlDB *EnclaveDB) AncientRange(kind string, start, max, maxByteSize uint64) ([][]byte, error) {
	return nil, errNotSupported
}

// Ancients returns an error as we don't have a backing chain freezer.
func (sqlDB *EnclaveDB) Ancients() (uint64, error) {
	return 0, errNotSupported
}

// AncientSize returns an error as we don't have a backing chain freezer.
func (sqlDB *EnclaveDB) AncientSize(kind string) (uint64, error) {
	return 0, errNotSupported
}

// ModifyAncients is not supported.
func (sqlDB *EnclaveDB) ModifyAncients(func(ethdb.AncientWriteOp) error) (int64, error) {
	return 0, errNotSupported
}

// TruncateAncients returns an error as we don't have a backing chain freezer.
func (sqlDB *EnclaveDB) TruncateAncients(items uint64) error {
	return errNotSupported
}

// Sync returns an error as we don't have a backing chain freezer.
func (sqlDB *EnclaveDB) Sync() error {
	return errNotSupported
}

func (sqlDB *EnclaveDB) ReadAncients(fn func(reader ethdb.AncientReader) error) (err error) {
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
