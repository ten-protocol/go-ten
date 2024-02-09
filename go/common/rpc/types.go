package rpc

import "github.com/ten-protocol/go-ten/go/common/viewingkey"

// RequestWithVk - wraps the eth parameters with a viewing key
type RequestWithVk struct {
	VK     *viewingkey.RPCSignedViewingKey
	Params []any
}
