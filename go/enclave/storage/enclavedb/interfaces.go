package enclavedb

import (
	"context"
	"database/sql"
	"fmt"

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

func (contract Contract) IsTransparent() bool {
	return contract.Transparent != nil && *contract.Transparent
}

// EventType - maps to the “event_type“ table
type EventType struct {
	Id                                          uint64
	Contract                                    *Contract
	EventSignature                              gethcommon.Hash
	AutoVisibility                              bool
	Public                                      bool
	Topic1CanView, Topic2CanView, Topic3CanView *bool
	SenderCanView                               *bool
}

func (et EventType) IsPublic() bool {
	return (et.Contract.Transparent != nil && *et.Contract.Transparent) || et.Public
}

func (et EventType) IsTopicRelevant(topicNo int) bool {
	switch topicNo {
	case 1:
		return et.Topic1CanView != nil && *et.Topic1CanView
	case 2:
		return et.Topic2CanView != nil && *et.Topic2CanView
	case 3:
		return et.Topic3CanView != nil && *et.Topic3CanView
	}
	// this should not happen under any circumstance
	panic(fmt.Sprintf("unknown topic no: %d", topicNo))
}
