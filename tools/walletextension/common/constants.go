package common

import (
	"time"
)

const (
	Localhost = "127.0.0.1"

	JSONKeyAddress    = "address"
	JSONKeyID         = "id"
	JSONKeyMethod     = "method"
	JSONKeyParams     = "params"
	JSONKeyResult     = "result"
	JSONKeyRPCVersion = "jsonrpc"
	JSONKeySignature  = "signature"
)

const (
	PathReady                           = "/ready/"
	PathViewingKeys                     = "/viewingkeys/"
	PathGenerateViewingKey              = "/generateviewingkey/"
	PathSubmitViewingKey                = "/submitviewingkey/"
	PathJoin                            = "/join/"
	PathAuthenticate                    = "/authenticate/"
	PathQuery                           = "/query/"
	PathRevoke                          = "/revoke/"
	PathHealth                          = "/health/"
	PathNetworkHealth                   = "/network-health/"
	WSProtocol                          = "ws://"
	HTTPProtocol                        = "http://"
	UserQueryParameter                  = "u"
	EncryptedTokenQueryParameter        = "token"
	AddressQueryParameter               = "a"
	MessageUserIDLen                    = 40
	EthereumAddressLen                  = 42
	GetStorageAtUserIDRequestMethodName = "0x0000000000000000000000000000000000000000"
	SuccessMsg                          = "success"
	APIVersion1                         = "/v1"
	PathVersion                         = "/version/"
	DeduplicationBufferSize             = 20
	ErrNoViewingKey                     = "authentication failed %s"
)

var ReaderHeadTimeout = 10 * time.Second
