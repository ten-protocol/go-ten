package common

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
)

// LogSubscription is an authenticated subscription to logs.
type LogSubscription struct {
	// ViewingKey - links this subscription request to an externally owed account
	ViewingKey viewingkey.RPCSignedViewingKey

	// A subscriber-defined filter to apply to the stream of logs.
	Filter *filters.FilterCriteria
}

// IDAndEncLog pairs an encrypted log with the ID of the subscription that generated it.
type IDAndEncLog struct {
	SubID  rpc.ID
	EncLog []byte
}

// IDAndLog pairs a log with the ID of the subscription that generated it.
type IDAndLog struct {
	SubID rpc.ID
	Log   *types.Log
}

// FilterCriteriaJSON is a structure that JSON-serialises to a format that can be successfully deserialised into a
// filters.FilterCriteria object (round-tripping a filters.FilterCriteria to JSON and back doesn't work, due to a
// custom serialiser implemented by filters.FilterCriteria).
type FilterCriteriaJSON struct {
	BlockHash *common.Hash     `json:"blockHash"`
	FromBlock *rpc.BlockNumber `json:"fromBlock"`
	ToBlock   *rpc.BlockNumber `json:"toBlock"`
	Addresses interface{}      `json:"address"`
	Topics    []interface{}    `json:"topics"`
}
