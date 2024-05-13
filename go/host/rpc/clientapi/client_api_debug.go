package clientapi

import (
	"context"

	"github.com/ten-protocol/go-ten/go/common/host"
	"github.com/ten-protocol/go-ten/go/common/tracers"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// NetworkDebug implements a subset of the Ethereum network JSON RPC operations.
type NetworkDebug struct {
	host host.Host
}

func NewNetworkDebug(host host.Host) *NetworkDebug {
	return &NetworkDebug{
		host: host,
	}
}

// TraceTransaction returns the structured logs created during the execution of EVM
// and returns them as a JSON object.
func (api *NetworkDebug) TraceTransaction(ctx context.Context, hash gethcommon.Hash, config *tracers.TraceConfig) (interface{}, error) {
	response, err := api.host.EnclaveClient().DebugTraceTransaction(ctx, hash, config)
	if err != nil {
		return "", err
	}
	return response, nil
}

// EventLogRelevancy returns the events for a given transactions and the revelancy params
func (api *NetworkDebug) EventLogRelevancy(ctx context.Context, hash gethcommon.Hash) (interface{}, error) {
	response, err := api.host.EnclaveClient().DebugEventLogRelevancy(ctx, hash)
	if err != nil {
		return "", err
	}
	return response, nil
}
