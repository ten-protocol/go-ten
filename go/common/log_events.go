package common

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

// LogSubscription is an authenticated subscription to logs.
type LogSubscription struct {
	// ViewingKey - links this subscription request to an externally owed account
	ViewingKey *viewingkey.RPCSignedViewingKey

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
	Addresses []common.Address `json:"addresses"`
	Topics    [][]common.Hash  `json:"topics"`
}

func FromCriteria(crit filters.FilterCriteria) FilterCriteriaJSON {
	var from *rpc.BlockNumber
	if crit.FromBlock != nil {
		f := (rpc.BlockNumber)(crit.FromBlock.Int64())
		from = &f
	}

	var to *rpc.BlockNumber
	if crit.ToBlock != nil {
		t := (rpc.BlockNumber)(crit.ToBlock.Int64())
		to = &t
	}

	return FilterCriteriaJSON{
		BlockHash: crit.BlockHash,
		FromBlock: from,
		ToBlock:   to,
		Addresses: crit.Addresses,
		Topics:    crit.Topics,
	}
}

func ToCriteria(jsonCriteria FilterCriteriaJSON) filters.FilterCriteria {
	var from *big.Int
	if jsonCriteria.FromBlock != nil {
		from = big.NewInt(jsonCriteria.FromBlock.Int64())
	}
	var to *big.Int
	if jsonCriteria.ToBlock != nil {
		to = big.NewInt(jsonCriteria.ToBlock.Int64())
	}

	return filters.FilterCriteria{
		BlockHash: jsonCriteria.BlockHash,
		FromBlock: from,
		ToBlock:   to,
		Addresses: jsonCriteria.Addresses,
		Topics:    jsonCriteria.Topics,
	}
}
