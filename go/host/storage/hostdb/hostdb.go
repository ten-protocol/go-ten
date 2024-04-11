package hostdb

import (
	"database/sql"
	"fmt"
)

type HostDB interface {
	GetSQLDB() *sql.DB
	NewDBTransaction() (*dbTransaction, error)
	GetSQLStatement() *SQLStatements
}

type hostDB struct {
	sqldb      *sql.DB
	statements *SQLStatements
}

func (db *hostDB) GetSQLStatement() *SQLStatements {
	return db.statements
}

func NewHostDB(db *sql.DB, statements *SQLStatements) (HostDB, error) {
	return &hostDB{
		sqldb:      db,
		statements: statements,
	}, nil
}

func (db *hostDB) GetSQLDB() *sql.DB {
	return db.sqldb
}

func (db *hostDB) NewDBTransaction() (*dbTransaction, error) {
	tx, err := db.sqldb.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin host db transaction. Cause: %w", err)
	}

	return &dbTransaction{
		tx: tx,
	}, nil
}

func (db *hostDB) Close() error {
	if err := db.sqldb.Close(); err != nil {
		return fmt.Errorf("failed to close host sql db - %w", err)
	}
	return nil
}

type dbTransaction struct {
	tx *sql.Tx
}

func (b *dbTransaction) Write() error {
	if err := b.tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit host db transaction. Cause: %w", err)
	}
	return nil
}
