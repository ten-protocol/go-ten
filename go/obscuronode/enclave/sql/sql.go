package sql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/obscuronet/obscuro-playground/go/log"

	"github.com/ethereum/go-ethereum/ethdb"
)

// sqlEthDatabase implements ethdb.Database
type sqlEthDatabase struct {
	db *sql.DB
}

func CreateSQLEthDatabase(db *sql.DB) (ethdb.Database, error) {
	s := &sqlEthDatabase{db: db}
	if err := s.Initialise(); err != nil {
		return nil, err
	}
	return s, nil
}

func (m *sqlEthDatabase) Initialise() error {
	stmt := `create table if not exists kv (key text primary key, value blob); delete from kv;`
	if _, err := m.db.Exec(stmt); err != nil {
		return fmt.Errorf("failed to initialise sql eth db - %w", err)
	}

	return nil
}

func (m *sqlEthDatabase) Has(key []byte) (bool, error) {
	pFind, err := m.getFindPrepStmt()
	if err != nil {
		return false, err
	}
	defer pFind.Close()

	err = pFind.QueryRow(key).Scan()
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
	pFind, err := m.getFindPrepStmt()
	if err != nil {
		return []byte{}, err
	}
	defer pFind.Close()

	err = pFind.QueryRow(key).Scan(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *sqlEthDatabase) Put(key []byte, value []byte) error {
	pIns, err := m.getInsertPrepStmt()
	if err != nil {
		return err
	}
	defer pIns.Close()

	_, err = pIns.Exec(key, value)
	return err
}

func (m *sqlEthDatabase) Delete(key []byte) error {
	pDel, err := m.getDeletePrepStmt()
	if err != nil {
		return err
	}
	defer pDel.Close()

	_, err = pDel.Exec(key)
	return err
}

func (m *sqlEthDatabase) Close() error {
	if err := m.db.Close(); err != nil {
		return fmt.Errorf("failed to close sql db - %w", err)
	}
	return nil
}

func (m *sqlEthDatabase) NewBatch() ethdb.Batch {
	log.Trace("SQL :: New batch")
	return &sqlBatch{
		db: m,
	}
}

func (m *sqlEthDatabase) NewIterator(prefix []byte, start []byte) ethdb.Iterator {
	pFindLike, err := m.getIteratorPrepStmt()
	defer func() { _ = pFindLike.Close() }()
	if err != nil {
		return &iterator{
			err: fmt.Errorf("failed to get prepared SQL stmt, iter will be empty, %w", err),
		}
	}
	// todo: make sure we're stringifying that prefix correctly
	pr := string(prefix)
	st := string(append(prefix, start...))
	rows, err := pFindLike.Query(pr+"%", st) //nolint:sqlclosecheck
	defer func() { _ = pFindLike.Close() }()
	if err != nil {
		return &iterator{
			err: fmt.Errorf("failed to get rows, iter will be empty, %w", err),
		}
	}
	if err = rows.Err(); rows.Err() != nil {
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

func (m *sqlEthDatabase) getPrepStmt(query string) (*sql.Stmt, error) {
	prep, err := m.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement `%s` - %w", query, err)
	}
	return prep, nil
}

func (m *sqlEthDatabase) getFindPrepStmt() (*sql.Stmt, error) {
	return m.getPrepStmt(`select kv.value from kv where kv.key = ?;`)
}

func (m *sqlEthDatabase) getInsertPrepStmt() (*sql.Stmt, error) {
	return m.getPrepStmt(`insert or replace into kv values(?, ?);`)
}

func (m *sqlEthDatabase) getDeletePrepStmt() (*sql.Stmt, error) {
	return m.getPrepStmt(`delete from kv where kv.key = ?;`)
}

func (m *sqlEthDatabase) getIteratorPrepStmt() (*sql.Stmt, error) {
	return m.getPrepStmt(`select * from kv where kv.key like ? and kv.key > ?`)
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
