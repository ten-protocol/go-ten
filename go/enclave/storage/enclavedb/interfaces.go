package enclavedb

import (
	"database/sql"

	"github.com/ethereum/go-ethereum/ethdb"
)

// EnclaveDB - An abstraction over
type EnclaveDB interface {
	ethdb.Database
	GetSQLDB() *sql.DB
	NewDBTransaction() *dbTransaction
	BeginTx() (*sql.Tx, error)
}

// DBTransaction - An abstraction over
type DBTransaction interface {
	ethdb.Batch
	GetDB() *sql.DB
	ExecuteSQL(query string, args ...any)
}
