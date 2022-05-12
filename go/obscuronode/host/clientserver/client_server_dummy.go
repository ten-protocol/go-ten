package clientserver

import (
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
)

// todo - joel - rename to in-mem (change file name too)
// A no-op implementation of `host.ClientServer`.
type clientServerDummyImpl struct{}

// NewClientServerDummy returns a no-op `host.ClientServer`.
func NewClientServerDummy() host.ClientServer {
	return clientServerDummyImpl{}
}

func (s clientServerDummyImpl) Start() {
}

func (s clientServerDummyImpl) Stop() {
}
