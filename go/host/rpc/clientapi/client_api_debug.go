package clientapi

import (
	"context"
	"github.com/ten-protocol/go-ten/go/responses"

	"github.com/ten-protocol/go-ten/go/common"

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

func (api *NetworkDebug) EventLogRelevancy(ctx context.Context, encryptedParams common.EncryptedParamsDebugLogRelevancy) (responses.DebugLogs, error) {
	enclaveResponse, sysError := api.host.EnclaveClient().DebugEventLogRelevancy(ctx, encryptedParams)
	if sysError != nil {
		return responses.EnclaveResponse{
			Err: &responses.InternalErrMsg,
		}, nil
	}
	return *enclaveResponse, nil
}
