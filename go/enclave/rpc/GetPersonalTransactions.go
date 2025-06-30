package rpc

import (
	"errors"
	"fmt"

	"github.com/ten-protocol/go-ten/go/common/errutil"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
)

func GetPersonalTransactionsValidate(reqParams []any, builder *CallBuilder[common.ListPrivateTransactionsQueryParams, common.PrivateTransactionsQueryResponse], rpc *EncryptionManager) error {
	if !storeTxEnabled(rpc, builder) {
		return nil
	}

	// Parameters are [PrivateTransactionListParams]
	if len(reqParams) != 1 {
		builder.Err = fmt.Errorf("unexpected number of parameters (expected %d, got %d)", 1, len(reqParams))
		return nil
	}

	privateCustomQuery, err := gethencoding.ExtractPrivateTransactionsQuery(reqParams[0])
	if err != nil {
		builder.Err = fmt.Errorf("unable to extract query - %w", err)
		return nil
	}
	addr := privateCustomQuery.Address
	builder.From = &addr
	builder.Param = privateCustomQuery
	return nil
}

func GetPersonalTransactionsExecute(builder *CallBuilder[common.ListPrivateTransactionsQueryParams, common.PrivateTransactionsQueryResponse], rpc *EncryptionManager) error {
	err := authenticateFrom(builder.VK, builder.From)
	if err != nil {
		builder.Err = err
		return nil
	}
	addr := builder.Param.Address

	receiptsCount, err := rpc.storage.CountTransactionsPerAddress(builder.ctx, &addr, builder.Param.ShowAllPublicTxs, builder.Param.ShowSyntheticTxs)
	if err != nil {
		return fmt.Errorf("CountTransactionsPerAddress - %w", err)
	}

	if receiptsCount == 0 {
		builder.ReturnValue = &common.PrivateTransactionsQueryResponse{
			Receipts: types.Receipts{},
			Total:    0,
		}
		return nil
	}

	internalReceipts, err := rpc.storage.GetTransactionsPerAddress(builder.ctx, &addr, &builder.Param.Pagination, builder.Param.ShowAllPublicTxs, builder.Param.ShowSyntheticTxs)
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			builder.ReturnValue = &common.PrivateTransactionsQueryResponse{
				Receipts: types.Receipts{},
				Total:    0,
			}
			return nil
		}
		return fmt.Errorf("GetTransactionsPerAddress - %w", err)
	}

	var receipts types.Receipts
	for _, receipt := range internalReceipts {
		receipts = append(receipts, receipt.ToReceipt())
	}

	builder.ReturnValue = &common.PrivateTransactionsQueryResponse{
		Receipts: receipts,
		Total:    receiptsCount,
	}
	return nil
}
