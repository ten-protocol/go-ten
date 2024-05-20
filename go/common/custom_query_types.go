package common

import "github.com/ethereum/go-ethereum/common"

// CustomQueries are Ten-specific queries that are not supported by the Ethereum RPC API but that we wish to support
// through the same interface.
//
// We currently use the eth_getStorageAt method to route these queries through the Ethereum RPC API, since it will not
// be supported by the Ten network.
//
// A custom query has a name (string), an address (if private request) and a params field (generic json object).
//
// NOTE: Private custom queries must include "address" as a top-level field in the params json object.

// CustomQuery methods
const (
	UserIDRequestMethodName           = "getUserID"
	ListPrivateTransactionsMethodName = "listPersonalTransactions"

	// DeprecatedUserIDRequestMethodName is still supported for backwards compatibility for User ID requests
	DeprecatedUserIDRequestMethodName = "0x0000000000000000000000000000000000000000"
)

type ListPrivateTransactionsQueryParams struct {
	Address    common.Address  `json:"address"`
	Pagination QueryPagination `json:"pagination"`
}
