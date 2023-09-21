package common

import (
	"fmt"
	"time"

	"github.com/obscuronet/go-obscuro/go/rpc"
)

const (
	Localhost = "127.0.0.1"

	JSONKeyAddress      = "address"
	JSONKeyData         = "data"
	JSONKeyErr          = "error"
	JSONKeyFrom         = "from"
	JSONKeyID           = "id"
	JSONKeyMethod       = "method"
	JSONKeyParams       = "params"
	JSONKeyResult       = "result"
	JSONKeyRoot         = "root"
	JSONKeyRPCVersion   = "jsonrpc"
	JSONKeySignature    = "signature"
	JSONKeySubscription = "subscription"
	JSONKeyCode         = "code"
	JSONKeyMessage      = "message"
)

const (
	PathRoot                            = "/"
	PathReady                           = "/ready/"
	PathViewingKeys                     = "/viewingkeys/"
	PathGenerateViewingKey              = "/generateviewingkey/"
	PathSubmitViewingKey                = "/submitviewingkey/"
	PathJoin                            = "/join/"
	PathAuthenticate                    = "/authenticate/"
	PathQuery                           = "/query/"
	PathRevoke                          = "/revoke/"
	PathObscuroGateway                  = "/"
	PathHealth                          = "/health/"
	WSProtocol                          = "ws://"
	DefaultUser                         = "defaultUser"
	UserQueryParameter                  = "u"
	AddressQueryParameter               = "a"
	MessageFormatRegex                  = `^Register\s(\w+)\sfor\s(\w+)$`
	MessageUserIDLen                    = 64
	SignatureLen                        = 65
	EthereumAddressLen                  = 42
	PersonalSignMessagePrefix           = "\x19Ethereum Signed Message:\n%d%s"
	GetStorageAtUserIDRequestMethodName = "getUserID"
	SuccessMsg                          = "success"
	APIVersion1                         = "/v1"
)

var (
	ErrSubscribeFailHTTP = fmt.Sprintf("received an %s request but the connection does not support subscriptions", rpc.Subscribe)
	ReaderHeadTimeout    = 10 * time.Second
)
