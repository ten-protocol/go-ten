package rpc

import (
	"fmt"

	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
)

func GetPersonalTransactionsValidate(reqParams []any, builder *CallBuilder[common.ListPrivateTransactionsQueryParams, common.PrivateTransactionsQueryResponse], _ *EncryptionManager) error {
	// Parameters are [PrivateCustomQueryHeader, PrivateCustomQueryArgs, null]
	if len(reqParams) != 3 {
		builder.Err = fmt.Errorf("unexpected number of parameters (expected %d, got %d)", 3, len(reqParams))
		return nil
	}

	privateCustomQuery, err := gethencoding.ExtractPrivateTransactionsQuery(reqParams[0], reqParams[1])
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
		return nil //nolint:nilerr
	}
	addr := builder.Param.Address
	encryptReceipts, err := rpc.storage.GetTransactionsPerAddress(builder.ctx, &addr, &builder.Param.Pagination)
	if err != nil {
		return fmt.Errorf("GetTransactionsPerAddress - %w", err)
	}

	receiptsCount, err := rpc.storage.CountTransactionsPerAddress(builder.ctx, &addr)
	if err != nil {
		return fmt.Errorf("CountTransactionsPerAddress - %w", err)
	}

	builder.ReturnValue = &common.PrivateTransactionsQueryResponse{
		Receipts: encryptReceipts,
		Total:    receiptsCount,
	}
	return nil
}
