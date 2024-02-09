package rpc

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ten-protocol/go-ten/go/common/errutil"

	"github.com/ethereum/go-ethereum/core/types"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ten-protocol/go-ten/go/common/syserr"
)

func ExtractGetLogsRequest(reqParams []any, builder *CallBuilder[filters.FilterCriteria, []*types.Log], _ *EncryptionManager) error {
	// Parameters are [Filter, Address]
	if len(reqParams) != 2 {
		builder.Err = fmt.Errorf("unexpected number of parameters")
		return nil
	}
	// We extract the arguments from the param bytes.
	filter, forAddress, err := extractGetLogsParams(reqParams)
	if err != nil {
		builder.Err = err
		return nil //nolint:nilerr
	}
	builder.From = forAddress
	builder.Param = filter
	return nil
}

func ExecuteGetLogs(rpcBuilder *CallBuilder[filters.FilterCriteria, []*types.Log], rpc *EncryptionManager) error { //nolint:gocognit
	filter := rpcBuilder.Param
	// todo logic to check that the filter is valid
	// can't have both from and blockhash
	// from <=to
	// todo (@stefan) - return user error
	if filter.BlockHash != nil && filter.FromBlock != nil {
		rpcBuilder.Err = fmt.Errorf("invalid filter. Cannot have both blockhash and fromBlock")
		return nil
	}

	from := filter.FromBlock
	if from != nil && from.Int64() < 0 {
		batch, err := rpc.storage.FetchBatchBySeqNo(rpc.registry.HeadBatchSeq().Uint64())
		if err != nil {
			// system error
			return fmt.Errorf("could not retrieve head batch. Cause: %w", err)
		}
		from = batch.Number()
	}

	// Set from to the height of the block hash
	if from == nil && filter.BlockHash != nil {
		batch, err := rpc.storage.FetchBatchHeader(*filter.BlockHash)
		if err != nil {
			if errors.Is(err, errutil.ErrNotFound) {
				rpcBuilder.Status = NotFound
				return nil
			}
			return err
		}
		from = batch.Number
	}

	to := filter.ToBlock
	// when to=="latest", don't filter on it
	if to != nil && to.Int64() < 0 {
		to = nil
	}

	if from != nil && to != nil && from.Cmp(to) > 0 {
		rpcBuilder.Err = fmt.Errorf("invalid filter. from (%d) > to (%d)", from, to)
		return nil
	}

	// We retrieve the relevant logs that match the filter.
	filteredLogs, err := rpc.storage.FilterLogs(rpcBuilder.From, from, to, nil, filter.Addresses, filter.Topics)
	if err != nil {
		if errors.Is(err, syserr.InternalError{}) {
			return err
		}
		rpcBuilder.Err = fmt.Errorf("could not retrieve logs matching the filter. Cause: %w", err)
		return nil
	}

	rpcBuilder.ReturnValue = &filteredLogs
	return nil
}

// Returns the params extracted from an eth_getLogs request.
func extractGetLogsParams(paramList []interface{}) (*filters.FilterCriteria, *gethcommon.Address, error) {
	// We extract the first param, the filter for the logs.
	// We marshal the filter criteria from a map to JSON, then back from JSON into a FilterCriteria. This is
	// because the filter criteria arrives as a map, and there is no way to convert it to a map directly into a
	// FilterCriteria.
	filterJSON, err := json.Marshal(paramList[0])
	if err != nil {
		return nil, nil, fmt.Errorf("could not marshal filter criteria to JSON. Cause: %w", err)
	}
	filter := filters.FilterCriteria{}
	err = filter.UnmarshalJSON(filterJSON)
	if err != nil {
		return nil, nil, fmt.Errorf("could not unmarshal filter criteria from JSON. Cause: %w", err)
	}

	// We extract the second param, the address the logs are for.
	forAddressHex, ok := paramList[1].(string)
	if !ok {
		return nil, nil, fmt.Errorf("expected second argument in GetLogs request to be of type string, but got %T", paramList[0])
	}
	forAddress := gethcommon.HexToAddress(forAddressHex)
	return &filter, &forAddress, nil
}
