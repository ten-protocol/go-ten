package hostdb

import (
	"database/sql"
	"fmt"
)

type HostDB interface {
	GetSQLDB() *sql.DB
	NewDBTransaction() *dbTransaction
	BeginTx() *sql.Tx
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

func (db *hostDB) BeginTx() *sql.Tx {
	tx, _ := db.sqldb.Begin()
	return tx
}

func (db *hostDB) NewDBTransaction() *dbTransaction {
	return &dbTransaction{
		tx: db.BeginTx(),
	}
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
	err := b.tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit host db transaction. Cause: %w", err)
	}
	return nil
}
