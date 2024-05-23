package common

import "github.com/ethereum/go-ethereum/common"

// CustomQueries are Ten-specific queries that are not supported by the Ethereum RPC API but that we wish to support
// through the same interface.
//
// We currently use the eth_getStorageAt method to route these queries through the Ethereum RPC API.
//
// In order to match the eth_getStorageAt method signature, we require that all custom queries use an incrementing "address"
// to specify the method we are calling (e.g. 0x000...001 is getUserID, 0x000...002 is listPrivateTransactions).
//
// The signature is: eth_getStorageAt(method, params, nil) where:
// - method is the address of the custom query as an address (e.g. 0x000...001)
// - params is a JSON string with the parameters for the query (this complies with the eth_getStorageAt method signature since position gets encoded as a hex string)
//
// NOTE: Private custom queries must also include "address" as a top-level field in the params json object to indicate
// the account the query is being made for.

// CustomQuery methods
const (
	UserIDRequestCQMethod           = "0x0000000000000000000000000000000000000001"
	ListPrivateTransactionsCQMethod = "0x0000000000000000000000000000000000000002"
)

type ListPrivateTransactionsQueryParams struct {
	Address    common.Address  `json:"address"`
	Pagination QueryPagination `json:"pagination"`
}
