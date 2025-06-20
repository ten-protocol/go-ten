package enclavedb

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/ethdb"
)

// EnclaveDB - An abstraction that implements the `ethdb.Database` on top of SQL, and also exposes underling sql primitives.
// Note: This might not be the best approach.
// Todo - consider a few design alternatives:
// The EnclaveDB - can be a factory that returns an Sql implementation of the ethdb.Database
type EnclaveDB interface {
	ethdb.Database
	GetSQLDB() *sqlx.DB
	NewDBTransaction(ctx context.Context) (*sqlx.Tx, error)
}

// Contract - maps to the “contract“ table
type Contract struct {
	Id             uint64
	Address        gethcommon.Address
	Creator        gethcommon.Address
	AutoVisibility bool
	Transparent    *bool
	EventTypes     map[gethcommon.Hash]*EventType
}

func (contract *Contract) EventType(eventSignature gethcommon.Hash) *EventType {
	if contract.EventTypes == nil {
		return nil
	}
	return contract.EventTypes[eventSignature]
}

func (contract *Contract) EventTypeList() []*EventType {
	if contract.EventTypes == nil {
		return nil
	}
	result := make([]*EventType, 0)
	for _, eventType := range contract.EventTypes {
		result = append(result, eventType)
	}
	return result
}

func (contract *Contract) IsTransparent() bool {
	return contract.Transparent != nil && *contract.Transparent
}

func (contract *Contract) SetEventTypes(ets []*EventType) {
	contract.EventTypes = make(map[gethcommon.Hash]*EventType)
	for _, et := range ets {
		contract.EventTypes[et.EventSignature] = et
	}
}

// EventType - maps to the “event_type“ table
type EventType struct {
	Id                                          uint64
	Contract                                    *Contract
	EventSignature                              gethcommon.Hash
	AutoVisibility                              bool
	AutoPublic                                  *bool // true -when the event is autodetected as public
	ConfigPublic                                bool
	Topic1CanView, Topic2CanView, Topic3CanView *bool
	SenderCanView                               *bool
}

func (et EventType) Validate() error {
	if !et.IsPublic() && !et.AutoVisibility {
		noneRelevant := true
		for i := 1; i <= 3; i++ {
			if et.IsTopicRelevant(i) {
				noneRelevant = false
			}
		}
		if noneRelevant {
			return fmt.Errorf("event type %s is not public and has no relevant topics", et.EventSignature.Hex())
		}
	}
	return nil
}

func (et EventType) IsPublicConfig() bool {
	return (et.Contract.Transparent != nil && *et.Contract.Transparent) || et.ConfigPublic
}

func (et EventType) IsPublic() bool {
	return et.IsPublicConfig() || (et.AutoPublic != nil && *et.AutoPublic)
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

// EventTopic - maps to the "event_topic" table
type EventTopic struct {
	Id                uint64
	RelevantAddressId *uint64
}
