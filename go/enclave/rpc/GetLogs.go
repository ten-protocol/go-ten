package rpc

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ten-protocol/go-ten/go/common/syserr"
)

func GetLogsValidate(reqParams []any, builder *CallBuilder[filters.FilterCriteria, []*types.Log], rpc *EncryptionManager) error {
	if !storeTxEnabled(rpc, builder) {
		return nil
	}
	// Parameters are [Filter]
	if len(reqParams) != 1 {
		builder.Err = fmt.Errorf("unexpected number of parameters")
		return nil
	}

	serialised, err := json.Marshal(reqParams[0])
	if err != nil {
		builder.Err = fmt.Errorf("invalid parameter %w", err)
		return nil
	}
	var crit common.FilterCriteriaJSON
	err = json.Unmarshal(serialised, &crit)
	if err != nil {
		builder.Err = fmt.Errorf("invalid parameter %w", err)
		return nil
	}
	filter := common.ToCriteria(crit)

	builder.Param = &filter
	return nil
}

func GetLogsExecute(builder *CallBuilder[filters.FilterCriteria, []*types.Log], rpc *EncryptionManager) error { //nolint:gocognit
	filter := builder.Param
	// can't have both from and blockhash
	if filter.BlockHash != nil && filter.FromBlock != nil {
		builder.Err = fmt.Errorf("invalid filter. Cannot have both blockhash and fromBlock")
		return nil
	}

	// from <=to
	from := filter.FromBlock
	if from != nil && from.Int64() < 0 {
		batch, err := rpc.storage.FetchBatchHeaderBySeqNo(builder.ctx, rpc.registry.HeadBatchSeq().Uint64())
		if err != nil {
			// system error
			return fmt.Errorf("could not retrieve head batch. Cause: %w", err)
		}
		from = batch.Number
	}

	// Set from to the height of the block hash
	if from == nil && filter.BlockHash != nil {
		batch, err := rpc.storage.FetchBatchHeader(builder.ctx, *filter.BlockHash)
		if err != nil {
			if errors.Is(err, errutil.ErrNotFound) {
				builder.Status = NotFound
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
		builder.Err = fmt.Errorf("invalid filter. from (%d) > to (%d)", from, to)
		return nil
	}

	// We retrieve the relevant logs that match the filter.
	filteredLogs, err := rpc.storage.FilterLogs(builder.ctx, builder.VK.AccountAddress, from, to, nil, filter.Addresses, filter.Topics)
	if err != nil {
		if errors.Is(err, syserr.InternalError{}) {
			return err
		}
		builder.Err = fmt.Errorf("could not retrieve logs matching the filter. Cause: %w", err)
		return nil
	}

	builder.ReturnValue = &filteredLogs
	return nil
}
