package clientserver

import (
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
)

// A no-op implementation of `host.ClientServer`.
type clientServerDummyImpl struct{}

// NewClientServerDummy returns a no-op `host.ClientServer`.
func NewClientServerDummy() host.ClientServer {
	return clientServerDummyImpl{}
}

func (server clientServerDummyImpl) Start() {
}

func (server clientServerDummyImpl) Stop() {
}
