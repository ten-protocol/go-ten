package rpc

import (
	"fmt"

	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
)

func GetCustomQueryValidate(reqParams []any, builder *CallBuilder[common.ListPrivateTransactionsQueryParams, common.PrivateQueryResponse], _ *EncryptionManager) error {
	// Parameters are [PrivateCustomQueryHeader, PrivateCustomQueryArgs, null]
	if len(reqParams) != 3 {
		builder.Err = fmt.Errorf("unexpected number of parameters (expected %d, got %d)", 3, len(reqParams))
		return nil
	}

	privateCustomQuery, err := gethencoding.ExtractPrivateCustomQuery(reqParams[0], reqParams[1])
	if err != nil {
		builder.Err = fmt.Errorf("unable to extract query - %w", err)
		return nil
	}
	builder.From = &privateCustomQuery.Address
	builder.Param = privateCustomQuery
	return nil
}

func GetCustomQueryExecute(builder *CallBuilder[common.ListPrivateTransactionsQueryParams, common.PrivateQueryResponse], rpc *EncryptionManager) error {
	err := authenticateFrom(builder.VK, builder.From)
	if err != nil {
		builder.Err = err
		return nil //nolint:nilerr
	}

	encryptReceipts, err := rpc.storage.GetTransactionsPerAddress(builder.ctx, &builder.Param.Address, &builder.Param.Pagination)
	if err != nil {
		return fmt.Errorf("GetTransactionsPerAddress - %w", err)
	}

	receiptsCount, err := rpc.storage.CountTransactionsPerAddress(builder.ctx, &builder.Param.Address)
	if err != nil {
		return fmt.Errorf("CountTransactionsPerAddress - %w", err)
	}

	builder.ReturnValue = &common.PrivateQueryResponse{
		Receipts: encryptReceipts,
		Total:    receiptsCount,
	}
	return nil
}
