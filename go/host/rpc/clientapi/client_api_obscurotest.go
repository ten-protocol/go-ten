package clientapi

import (
	"github.com/obscuronet/go-obscuro/go/common/container"
)

// TestAPI implements JSON RPC operations required for testing.
type TestAPI struct {
	container container.Container
}

func NewTestAPI(container container.Container) *TestAPI {
	return &TestAPI{
		container: container,
	}
}

// StopHost gracefully stops the host.
func (api *TestAPI) StopHost() {
	if api.container != nil {
		_ = api.container.Stop()
	}
}
