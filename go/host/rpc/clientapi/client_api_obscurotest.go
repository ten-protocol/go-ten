package clientapi

import (
	"github.com/obscuronet/go-obscuro/go/common/container"
	"github.com/obscuronet/go-obscuro/go/common/host"
)

// TODO - Disable these methods on a non-test node.

// TestAPI implements JSON RPC operations required for testing.
type TestAPI struct {
	host      host.Host
	container container.Container
}

func NewTestAPI(host host.Host, container container.Container) *TestAPI {
	return &TestAPI{
		host:      host,
		container: container,
	}
}

// StopHost gracefully stops the host.
// TODO - Investigate how to authenticate this and other sensitive methods in production (Geth uses JWT).
// TODO - Change inmemory tests to use the host container instead of the host. Remove the `api.host` after. - https://github.com/obscuronet/obscuro-internal/issues/1314
func (api *TestAPI) StopHost() {
	if api.host != nil {
		api.host.Stop()
	}

	if api.container != nil {
		_ = api.container.Stop()
	}
}
