package sql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/ethdb"
)

const (
	getQry    = `select keyvalue.val from keyvalue where keyvalue.ky = ?;`
	putQry    = `insert or replace into keyvalue values(?, ?);`
	delQry    = `delete from keyvalue where keyvalue.ky = ?;`
	searchQry = `select * from keyvalue where substring(keyvalue.ky, 1, ?) = ? and keyvalue.ky >= ? order by keyvalue.ky asc`
)

// sqlEthDatabase implements ethdb.Database
type sqlEthDatabase struct {
	db *sql.DB
}

func CreateSQLEthDatabase(db *sql.DB) (ethdb.Database, error) {
	return &sqlEthDatabase{db: db}, nil
}

func (m *sqlEthDatabase) Has(key []byte) (bool, error) {
	err := m.db.QueryRow(getQry, key).Scan()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (m *sqlEthDatabase) Get(key []byte) ([]byte, error) {
	var res []byte

	err := m.db.QueryRow(getQry, key).Scan(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *sqlEthDatabase) Put(key []byte, value []byte) error {
	_, err := m.db.Exec(putQry, key, value)
	return err
}

func (m *sqlEthDatabase) Delete(key []byte) error {
	_, err := m.db.Exec(delQry, key)
	return err
}

func (m *sqlEthDatabase) Close() error {
	if err := m.db.Close(); err != nil {
		return fmt.Errorf("failed to close sql db - %w", err)
	}
	return nil
}

func (m *sqlEthDatabase) NewBatch() ethdb.Batch {
	return &sqlBatch{
		db: m,
	}
}

func (m *sqlEthDatabase) NewIterator(prefix []byte, start []byte) ethdb.Iterator {
	pr := prefix
	st := append(prefix, start...)
	// iterator clean-up handles closing this rows iterator
	rows, err := m.db.Query(searchQry, len(pr), pr, st) //nolint:sqlclosecheck
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

func (m *sqlEthDatabase) Stat(property string) (string, error) {
	// TODO implement me
	panic("implement me")
}

func (m *sqlEthDatabase) Compact(start []byte, limit []byte) error {
	// TODO implement me
	panic("implement me")
}

// no-freeze! Copied from the geth in-memory db implementation these ancient method implementations disable the 'freezer'

// errNotSupported is returned if the database doesn't support the required operation.
var errNotSupported = errors.New("this operation is not supported")

// HasAncient returns an error as we don't have a backing chain freezer.
func (m *sqlEthDatabase) HasAncient(kind string, number uint64) (bool, error) {
	return false, errNotSupported
}

// Ancient returns an error as we don't have a backing chain freezer.
func (m *sqlEthDatabase) Ancient(kind string, number uint64) ([]byte, error) {
	return nil, errNotSupported
}

// AncientRange returns an error as we don't have a backing chain freezer.
func (m *sqlEthDatabase) AncientRange(kind string, start, max, maxByteSize uint64) ([][]byte, error) {
	return nil, errNotSupported
}

// Ancients returns an error as we don't have a backing chain freezer.
func (m *sqlEthDatabase) Ancients() (uint64, error) {
	return 0, errNotSupported
}

// AncientSize returns an error as we don't have a backing chain freezer.
func (m *sqlEthDatabase) AncientSize(kind string) (uint64, error) {
	return 0, errNotSupported
}

// ModifyAncients is not supported.
func (m *sqlEthDatabase) ModifyAncients(func(ethdb.AncientWriteOp) error) (int64, error) {
	return 0, errNotSupported
}

// TruncateAncients returns an error as we don't have a backing chain freezer.
func (m *sqlEthDatabase) TruncateAncients(items uint64) error {
	return errNotSupported
}

// Sync returns an error as we don't have a backing chain freezer.
func (m *sqlEthDatabase) Sync() error {
	return errNotSupported
}

func (m *sqlEthDatabase) ReadAncients(fn func(reader ethdb.AncientReader) error) (err error) {
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
	return fn(m)
}
