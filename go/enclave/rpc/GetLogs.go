package rpc

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/syserr"
	"github.com/ten-protocol/go-ten/go/responses"
)

func (rpc *EncryptionManager) GetLogs(encryptedParams common.EncryptedParamsGetLogs) (*responses.Logs, common.SystemError) { //nolint
	return withVKEncryption1[filters.FilterCriteria, []*types.Log](
		rpc,
		rpc.config.ObscuroChainID,
		encryptedParams,
		// extract sender
		func(reqParams []any) (*UserRPCRequest1[filters.FilterCriteria], error) {
			// Parameters are [Filter, Address]
			if len(reqParams) != 2 {
				return nil, fmt.Errorf("unexpected number of parameters")
			}
			// We extract the arguments from the param bytes.
			filter, forAddress, err := extractGetLogsParams(reqParams)
			if err != nil {
				return nil, err
			}

			return &UserRPCRequest1[filters.FilterCriteria]{forAddress, filter}, nil
		},
		// execute
		func(decodedParams *UserRPCRequest1[filters.FilterCriteria]) (*UserResponse[[]*types.Log], error) {
			filter := decodedParams.Param1
			// todo logic to check that the filter is valid
			// can't have both from and blockhash
			// from <=to
			// todo (@stefan) - return user error
			if filter.BlockHash != nil && filter.FromBlock != nil {
				return &UserResponse[[]*types.Log]{nil, fmt.Errorf("invalid filter. Cannot have both blockhash and fromBlock")}, nil
			}

			from := filter.FromBlock
			if from != nil && from.Int64() < 0 {
				batch, err := rpc.storage.FetchBatchBySeqNo(rpc.registry.HeadBatchSeq().Uint64())
				if err != nil {
					return &UserResponse[[]*types.Log]{nil, fmt.Errorf("could not retrieve head batch. Cause: %w", err)}, nil
				}
				from = batch.Number()
			}

			// Set from to the height of the block hash
			if from == nil && filter.BlockHash != nil {
				batch, err := rpc.storage.FetchBatchHeader(*filter.BlockHash)
				if err != nil {
					return nil, err
				}
				from = batch.Number
			}

			to := filter.ToBlock
			// when to=="latest", don't filter on it
			if to != nil && to.Int64() < 0 {
				to = nil
			}

			if from != nil && to != nil && from.Cmp(to) > 0 {
				return &UserResponse[[]*types.Log]{nil, fmt.Errorf("invalid filter. from (%d) > to (%d)", from, to)}, nil
			}

			// We retrieve the relevant logs that match the filter.
			filteredLogs, err := rpc.storage.FilterLogs(decodedParams.Sender, from, to, nil, filter.Addresses, filter.Topics)
			if err != nil {
				if errors.Is(err, syserr.InternalError{}) {
					return nil, err
				}
				err = fmt.Errorf("could not retrieve logs matching the filter. Cause: %w", err)
				return &UserResponse[[]*types.Log]{nil, err}, nil
			}
			return &UserResponse[[]*types.Log]{&filteredLogs, nil}, nil
		})
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
