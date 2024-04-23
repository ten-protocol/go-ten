package enclavedb

import (
	"context"
	"database/sql"

	"github.com/ethereum/go-ethereum/ethdb"
)

// EnclaveDB - An abstraction that implements the `ethdb.Database` on top of SQL, and also exposes underling sql primitives.
// Note: This might not be the best approach.
// Todo - consider a few design alternatives:
// The EnclaveDB - can be a factory that returns an Sql implementation of the ethdb.Database
type EnclaveDB interface {
	ethdb.Database
	GetSQLDB() *sql.DB
	NewDBTransaction() *dbTransaction
	BeginTx(context.Context) (*sql.Tx, error)
}

// DBTransaction - represents a database transaction implemented unusually.
// Typically, databases have a "beginTransaction" command which is also exposed by the db drivers,
// and then the applications just sends commands on that connection.
// There are rules as to what data is returned when running selects.
// This implementation works by collecting all statements, and then writing them and committing in one go
// todo - does it need to be an ethdb.Batch?
// todo - can we use the typical
type DBTransaction interface {
	ethdb.Batch
	GetDB() *sql.DB
	ExecuteSQL(query string, args ...any)
}
