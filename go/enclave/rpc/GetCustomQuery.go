package rpc

import (
	"fmt"

	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
	"github.com/ten-protocol/go-ten/go/responses"
)

func (rpc *EncryptionManager) GetCustomQuery(encryptedParams common.EncryptedParamsGetStorageAt) (*responses.PrivateQueryResponse, common.SystemError) {
	return withVKEncryption1[common.PrivateCustomQueryListTransactions, common.PrivateQueryResponse](
		rpc,
		rpc.config.ObscuroChainID,
		encryptedParams,
		// extract sender
		func(reqParams []any) (*UserRPCRequest1[common.PrivateCustomQueryListTransactions], error) {
			// Parameters are [PrivateCustomQueryHeader, PrivateCustomQueryArgs, null]
			if len(reqParams) != 3 {
				return nil, fmt.Errorf("unexpected number of parameters")
			}

			privateCustomQuery, err := gethencoding.ExtractPrivateCustomQuery(reqParams[0], reqParams[1])
			if err != nil {
				return nil, fmt.Errorf("unable to extract query - %w", err)
			}

			return &UserRPCRequest1[common.PrivateCustomQueryListTransactions]{&privateCustomQuery.Address, privateCustomQuery}, nil
		},
		// execute
		func(params *UserRPCRequest1[common.PrivateCustomQueryListTransactions]) (*UserResponse[common.PrivateQueryResponse], error) {
			// params are correct, fetch the receipts of the requested address
			encryptReceipts, err := rpc.storage.GetReceiptsPerAddress(&params.Param1.Address, &params.Param1.Pagination)
			if err != nil {
				return nil, fmt.Errorf("unable to get storage - %w", err)
			}

			receiptsCount, err := rpc.storage.GetReceiptsPerAddressCount(&params.Param1.Address)
			if err != nil {
				return nil, fmt.Errorf("unable to get storage - %w", err)
			}

			return &UserResponse[common.PrivateQueryResponse]{&common.PrivateQueryResponse{
				Receipts: encryptReceipts,
				Total:    receiptsCount,
			}, nil}, nil
		})
}
