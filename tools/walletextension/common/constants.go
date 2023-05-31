package common

import (
	"fmt"

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
	PathRoot               = "/"
	PathReady              = "/ready/"
	PathViewingKeys        = "/viewingkeys/"
	PathGenerateViewingKey = "/generateviewingkey/"
	PathSubmitViewingKey   = "/submitviewingkey/"
	staticDir              = "static"
	WSProtocol             = "ws://"
	DefaultUser            = "defaultUser"

	SuccessMsg = "success"
)

var ErrSubscribeFailHTTP = fmt.Sprintf("received an %s request but the connection does not support subscriptions", rpc.Subscribe)
