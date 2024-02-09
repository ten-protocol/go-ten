package rpc

import (
	"fmt"

	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
)

func ExtractGetCustomQueryRequest(reqParams []any, builder *CallBuilder[common.PrivateCustomQueryListTransactions, common.PrivateQueryResponse], _ *EncryptionManager) error {
	// Parameters are [PrivateCustomQueryHeader, PrivateCustomQueryArgs, null]
	if len(reqParams) != 3 {
		builder.Err = fmt.Errorf("unexpected number of parameters")
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

func ExecuteGetCustomQuery(rpcBuilder *CallBuilder[common.PrivateCustomQueryListTransactions, common.PrivateQueryResponse], rpc *EncryptionManager) error {
	// rpcBuilder are correct, fetch the receipts of the requested address
	encryptReceipts, err := rpc.storage.GetReceiptsPerAddress(&rpcBuilder.Param.Address, &rpcBuilder.Param.Pagination)
	if err != nil {
		return fmt.Errorf("GetReceiptsPerAddress - %w", err)
	}

	receiptsCount, err := rpc.storage.GetReceiptsPerAddressCount(&rpcBuilder.Param.Address)
	if err != nil {
		return fmt.Errorf("GetReceiptsPerAddressCount - %w", err)
	}

	rpcBuilder.ReturnValue = &common.PrivateQueryResponse{
		Receipts: encryptReceipts,
		Total:    receiptsCount,
	}
	return nil
}
