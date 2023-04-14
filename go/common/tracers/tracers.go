// Package tracers: This file was copied/adapted from geth - go-ethereum/eth/tracers
//
//nolint
package tracers

import (
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/vm"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	gethlogger "github.com/ethereum/go-ethereum/eth/tracers/logger"
)

// DebugLogs are the logs returned when using the DebugGetLogs endpoint
type DebugLogs struct {
	RelAddress1    gethcommon.Hash `json:"relAddress1"`
	RelAddress2    gethcommon.Hash `json:"relAddress2"`
	RelAddress3    gethcommon.Hash `json:"relAddress3"`
	RelAddress4    gethcommon.Hash `json:"relAddress4"`
	LifecycleEvent bool            `json:"lifecycleEvent"`

	gethtypes.Log
}

// MarshalJSON marshals as JSON.
// this holds a copy of the gethtypes.Log log marshaller
func (l DebugLogs) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		RelAddress1    gethcommon.Hash    `json:"relAddress1"`
		RelAddress2    gethcommon.Hash    `json:"relAddress2"`
		RelAddress3    gethcommon.Hash    `json:"relAddress3"`
		RelAddress4    gethcommon.Hash    `json:"relAddress4"`
		LifecycleEvent bool               `json:"lifecycleEvent"`
		Address        gethcommon.Address `json:"address" gencodec:"required"`
		Topics         []gethcommon.Hash  `json:"topics" gencodec:"required"`
		Data           hexutil.Bytes      `json:"data" gencodec:"required"`
		BlockNumber    hexutil.Uint64     `json:"blockNumber"`
		TxHash         gethcommon.Hash    `json:"transactionHash" gencodec:"required"`
		TxIndex        hexutil.Uint       `json:"transactionIndex"`
		BlockHash      gethcommon.Hash    `json:"blockHash"`
		Index          hexutil.Uint       `json:"logIndex"`
		Removed        bool               `json:"removed"`
	}{
		l.RelAddress1,
		l.RelAddress2,
		l.RelAddress3,
		l.RelAddress4,
		l.LifecycleEvent,
		l.Address,
		l.Topics,
		l.Data,
		hexutil.Uint64(l.BlockNumber),
		l.TxHash,
		hexutil.Uint(l.TxIndex),
		l.BlockHash,
		hexutil.Uint(l.Index),
		l.Removed,
	})
}

// TraceConfig holds extra parameters to trace functions.
type TraceConfig struct {
	*gethlogger.Config
	Tracer  *string
	Timeout *string
	Reexec  *uint64
	// Config specific to given tracer. Note struct logger
	// config are historically embedded in main object.
	TracerConfig json.RawMessage
}

// Context contains some contextual infos for a transaction execution that is not
// available from within the EVM object.
type Context struct {
	BlockHash gethcommon.Hash // Hash of the block the tx is contained within (zero if dangling tx or call)
	TxIndex   int             // Index of the transaction within a block (zero if dangling tx or call)
	TxHash    gethcommon.Hash // Hash of the transaction being traced (zero if dangling call)
}

// Tracer interface extends vm.EVMLogger and additionally
// allows collecting the tracing result.
type Tracer interface {
	vm.EVMLogger
	GetResult() (json.RawMessage, error)
	// Stop terminates execution of the tracer at the first opportune moment.
	Stop(err error)
}

type lookupFunc func(string, *Context) (Tracer, error)

var lookups []lookupFunc

// RegisterLookup registers a method as a lookup for tracers, meaning that
// users can invoke a named tracer through that lookup. If 'wildcard' is true,
// then the lookup will be placed last. This is typically meant for interpreted
// engines (js) which can evaluate dynamic user-supplied code.
func RegisterLookup(wildcard bool, lookup lookupFunc) {
	if wildcard {
		lookups = append(lookups, lookup)
	} else {
		lookups = append([]lookupFunc{lookup}, lookups...)
	}
}

// New returns a new instance of a tracer, by iterating through the
// registered lookups.
func New(code string, ctx *Context) (Tracer, error) {
	for _, lookup := range lookups {
		if tracer, err := lookup(code, ctx); err == nil {
			return tracer, nil
		}
	}
	return nil, errors.New("tracer not found")
}
