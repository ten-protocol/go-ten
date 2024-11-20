package rpc

import "github.com/ten-protocol/go-ten/go/common/viewingkey"

// encryptedMethods for which the RPC requests and responses should be encrypted
// these are names of virtual RPC methods exposed by a TEN node
// they all get routed through "ten_encryptedRPC"
const (
	ERPCCall                    = "ten_call"
	ERPCGetBalance              = "ten_getBalance"
	ERPCGetTransactionByHash    = "ten_getTransactionByHash"
	ERPCGetRawTransactionByHash = "ten_getRawTransactionByHash"
	ERPCGetTransactionCount     = "ten_getTransactionCount"
	ERPCGetTransactionReceipt   = "ten_getTransactionReceipt"
	ERPCSendRawTransaction      = "ten_sendRawTransaction"
	ERPCResend                  = "ten_resend"
	ERPCEstimateGas             = "ten_estimateGas"
	ERPCGetLogs                 = "ten_getLogs"
	ERPCGetStorageAt            = "ten_getStorageAt"
	ERPCDebugLogs               = "debug_eventLogRelevancy"
	ERPCGetPersonalTransactions = "scan_getPersonalTransactions"
)

var encryptedMethods = []string{
	ERPCCall,
	ERPCGetBalance,
	ERPCGetTransactionByHash,
	ERPCGetRawTransactionByHash,
	ERPCGetTransactionCount,
	ERPCGetTransactionReceipt,
	ERPCSendRawTransaction,
	ERPCResend,
	ERPCEstimateGas,
	ERPCGetLogs,
	ERPCGetStorageAt,
	ERPCDebugLogs,
	ERPCGetPersonalTransactions,
}

// IsEncryptedMethod indicates whether the RPC method's requests and responses should be encrypted.
func IsEncryptedMethod(method string) bool {
	for _, m := range encryptedMethods {
		if m == method {
			return true
		}
	}
	return false
}

// RequestWithVk - wraps the eth parameters with a viewing key
type RequestWithVk struct {
	VK     *viewingkey.RPCSignedViewingKey
	Method string // can be only one of the encrypted methods above
	Params []any
}
