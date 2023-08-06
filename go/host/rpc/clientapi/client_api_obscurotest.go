package clientapi

// TestAPI implements JSON RPC operations required for testing.
type TestAPI struct {
	stopHost func() error
}

func NewTestAPI(stopHost func() error) *TestAPI {
	return &TestAPI{
		stopHost: stopHost,
	}
}

// StopHost gracefully stops the host.
func (api *TestAPI) StopHost() error {
	return api.stopHost()
}
