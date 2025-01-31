package common

const (
	Localhost = "127.0.0.1"

	JSONKeyAddress         = "address"
	JSONKeyID              = "id"
	JSONKeyMethod          = "method"
	JSONKeyParams          = "params"
	JSONKeyRPCVersion      = "jsonrpc"
	JSONKeySignature       = "signature"
	JSONKeyType            = "type"
	JSONKeyEncryptionToken = "encryptionToken"
	JSONKeyFormats         = "formats"
)

const (
	PathReady                     = "/ready/"
	PathJoin                      = "/join/"
	PathGetMessage                = "/getmessage/"
	PathAuthenticate              = "/authenticate/"
	PathQuery                     = "/query/"
	PathRevoke                    = "/revoke/"
	PathHealth                    = "/health/"
	PathSessionKeys               = "/session-key/"
	PathNetworkHealth             = "/network-health/"
	PathNetworkConfig             = "/network-config/"
	PathKeyExchange               = "/key-exchange/"
	WSProtocol                    = "ws://"
	HTTPProtocol                  = "http://"
	EncryptedTokenQueryParameter  = "token"
	AddressQueryParameter         = "a"
	MessageUserIDLen              = 40
	MessageUserIDLenWithPrefix    = 42
	EthereumAddressLen            = 42
	SuccessMsg                    = "success"
	APIVersion1                   = "/v1"
	PathVersion                   = "/version/"
	DeduplicationBufferSize       = 20
	DefaultGatewayAuthMessageType = "EIP712"
)
