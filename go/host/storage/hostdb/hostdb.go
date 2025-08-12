package hostdb

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

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
