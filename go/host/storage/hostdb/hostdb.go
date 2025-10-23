package hostdb

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/ten-protocol/go-ten/go/common/errutil"

	gethlog "github.com/ethereum/go-ethereum/log"
)

type HostDB interface {
	GetSQLDB() *sqlx.DB
	NewDBTransaction() (*dbTransaction, error)
	Logger() gethlog.Logger
}

type hostDB struct {
	sqldb  *sqlx.DB
	logger gethlog.Logger
}

func NewHostDB(db *sqlx.DB, logger gethlog.Logger) (HostDB, error) {
	return &hostDB{
		sqldb:  db,
		logger: logger,
	}, nil
}

func (db *hostDB) GetSQLDB() *sqlx.DB {
	return db.sqldb
}

func (db *hostDB) Logger() gethlog.Logger {
	return db.logger
}

func (db *hostDB) NewDBTransaction() (*dbTransaction, error) {
	tx, err := db.sqldb.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin host db transaction. Cause: %w", err)
	}

	return &dbTransaction{
		Tx: tx,
	}, nil
}

func (db *hostDB) Close() error {
	if err := db.sqldb.Close(); err != nil {
		return fmt.Errorf("failed to close host sql db - %w", err)
	}
	return nil
}

type dbTransaction struct {
	Tx *sql.Tx
}

func (b *dbTransaction) Write() error {
	if err := b.Tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit host db transaction. Cause: %w", err)
	}
	return nil
}

func (b *dbTransaction) Rollback() error {
	if err := b.Tx.Rollback(); err != nil {
		return fmt.Errorf("failed to rollback host transaction. Cause: %w", err)
	}
	return nil
}

func GetMetadata(db HostDB, key string) (uint64, error) {
	var bytea []byte
	query := db.GetSQLDB().Rebind("SELECT val FROM config WHERE ky = ?")
	if err := db.GetSQLDB().Get(&bytea, query, key); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errutil.ErrNotFound
		}
		return 0, fmt.Errorf("failed to get metadata: %w", err)
	}
	// we can't cast to integer on postgres so have to convert the raw bytes outside the query
	s := strings.TrimSpace(string(bytea))
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid metadata value %q: %w", s, err)
	}
	return v, nil
}

func SetMetadata(db HostDB, key string, value uint64) error {
	query := "INSERT OR REPLACE INTO config (ky, val) VALUES (?, ?)"
	reboundQuery := db.GetSQLDB().Rebind(query)
	_, err := db.GetSQLDB().Exec(reboundQuery, key, value)
	if err != nil {
		return fmt.Errorf("failed to set metadata: %w", err)
	}
	return nil
}
