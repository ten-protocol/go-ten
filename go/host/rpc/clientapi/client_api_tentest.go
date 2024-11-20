package clientapi

import (
	"github.com/ten-protocol/go-ten/go/common/container"
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
func (api *TestAPI) StopHost() error {
	if api.container != nil {
		err := api.container.Stop()
		if err != nil {
			return err
		}
	}
	return nil
}
