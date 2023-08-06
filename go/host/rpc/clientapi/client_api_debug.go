package clientapi

import (
	"context"

	"github.com/obscuronet/go-obscuro/go/host"

	"github.com/obscuronet/go-obscuro/go/common/tracers"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

type networkAPIServiceLocator interface {
	host.EnclaveLocator
}

// NetworkDebug implements a subset of the Ethereum network JSON RPC operations.
type NetworkDebug struct {
	sl networkAPIServiceLocator
}

func NewNetworkDebug(serviceLocator networkAPIServiceLocator) *NetworkDebug {
	return &NetworkDebug{
		sl: serviceLocator,
	}
}

// TraceTransaction returns the structured logs created during the execution of EVM
// and returns them as a JSON object.
func (api *NetworkDebug) TraceTransaction(_ context.Context, hash gethcommon.Hash, config *tracers.TraceConfig) (interface{}, error) {
	response, err := api.sl.Enclave().GetEnclaveClient().DebugTraceTransaction(hash, config)
	if err != nil {
		return "", err
	}
	return response, nil
}

// EventLogRelevancy returns the events for a given transactions and the revelancy params
func (api *NetworkDebug) EventLogRelevancy(_ context.Context, hash gethcommon.Hash) (interface{}, error) {
	response, err := api.sl.Enclave().GetEnclaveClient().DebugEventLogRelevancy(hash)
	if err != nil {
		return "", err
	}
	return response, nil
}
