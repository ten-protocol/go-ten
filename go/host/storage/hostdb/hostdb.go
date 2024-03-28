package hostdb

import (
	"database/sql"
	"fmt"
)

type HostDB interface {
	GetDB() *sql.DB
	NewDBTransaction() *dbTransaction
	BeginTx() (*sql.Tx, error)
}

type hostDB struct {
	sqldb *sql.DB
}

func NewHostDB(db *sql.DB) (HostDB, error) {
	return &hostDB{sqldb: db}, nil
}

func (db *hostDB) GetDB() *sql.DB {
	return db.sqldb
}

func (db *hostDB) BeginTx() (*sql.Tx, error) {
	return db.sqldb.Begin()
}

func (db *hostDB) NewDBTransaction() *dbTransaction {
	return &dbTransaction{
		db: db,
	}
}

func (db *hostDB) Close() error {
	if err := db.sqldb.Close(); err != nil {
		return fmt.Errorf("failed to close host sql db - %w", err)
	}
	return nil
}

type dbTransaction struct {
	db HostDB
}

func (b *dbTransaction) GetDB() *sql.DB {
	return b.db.GetDB()
}

func (b *dbTransaction) Write() error {
	tx, err := b.db.BeginTx()
	if err != nil {
		return fmt.Errorf("failed to create host db transaction - %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit host db transaction. Cause: %w", err)
	}
	return nil
}
