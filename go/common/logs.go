package common

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/rpc"
)

// LogSubscription is an authenticated subscription to logs.
type LogSubscription struct {
	// The account the events relate to.
	Account *common.Address
	// A signature over the account address using a private viewing key. Prevents attackers from subscribing to
	// (encrypted) logs for other accounts to see the pattern of logs.
	// TODO - This does not protect against replay attacks, where someone resends an intercepted subscription request.
	Signature *[]byte
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

// LogsByRollupByID is a double-map from subscription IDs to rollup numbers to logs generated in that rollup.
type LogsByRollupByID = map[rpc.ID]map[uint64][]*types.Log

// EncLogsByRollupByID is identical to LogsByRollupByID, but with the logs encrypted as bytes.
type EncLogsByRollupByID = map[rpc.ID]map[uint64][]byte

// FilterCriteriaJSON is a structure that JSON-serialises to the expected format for log filter criteria.
type FilterCriteriaJSON struct {
	BlockHash *common.Hash     `json:"blockHash"`
	FromBlock *rpc.BlockNumber `json:"fromBlock"`
	ToBlock   *rpc.BlockNumber `json:"toBlock"`
	Addresses interface{}      `json:"address"`
	Topics    []interface{}    `json:"topics"`
}
