package enclavedb

import (
	"context"
	"database/sql"

	gethcommon "github.com/ethereum/go-ethereum/common"

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

// Contract - maps to the “contract“ table
type Contract struct {
	Id             uint64
	Address        gethcommon.Address
	AutoVisibility bool
	Transparent    *bool
}

// EventType - maps to the “event_type“ table
type EventType struct {
	Id                                          uint64
	ContractId                                  uint64
	EventSignature                              gethcommon.Hash
	AutoVisibility                              bool
	Public                                      bool
	Topic1CanView, Topic2CanView, Topic3CanView *bool
	SenderCanView                               *bool
}
