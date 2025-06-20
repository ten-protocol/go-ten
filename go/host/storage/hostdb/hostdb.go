package hostdb

import (
	"database/sql"
	"fmt"

	gethlog "github.com/ethereum/go-ethereum/log"
)

type HostDB interface {
	GetSQLDB() *sql.DB
	NewDBTransaction() (*dbTransaction, error)
	GetSQLStatement() *SQLStatements
	Logger() gethlog.Logger
}

type hostDB struct {
	sqldb      *sql.DB
	statements *SQLStatements
	logger     gethlog.Logger
}

func (db *hostDB) GetSQLStatement() *SQLStatements {
	return db.statements
}

func NewHostDB(db *sql.DB, statements *SQLStatements, logger gethlog.Logger) (HostDB, error) {
	return &hostDB{
		sqldb:      db,
		statements: statements,
		logger:     logger,
	}, nil
}

func (db *hostDB) GetSQLDB() *sql.DB {
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
