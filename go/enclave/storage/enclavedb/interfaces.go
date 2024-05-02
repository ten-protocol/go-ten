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
	NewDBTransaction(ctx context.Context) (*sql.Tx, error)
}
