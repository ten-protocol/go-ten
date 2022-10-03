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

// LogsByID is a map of subscription IDs to logs.
type LogsByID = map[rpc.ID][]*types.Log

// EncLogsByID is a map of subscription IDs to encrypted logs.
type EncLogsByID = map[rpc.ID]EncryptedLogs
