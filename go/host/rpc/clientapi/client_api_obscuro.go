package clientapi

import (
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/host"
)

// ObscuroAPI implements Obscuro-specific JSON RPC operations.
type ObscuroAPI struct {
	host host.Host
}

func NewObscuroAPI(host host.Host) *ObscuroAPI {
	return &ObscuroAPI{
		host: host,
	}
}

// Health returns the health status of obscuro host + enclave + db
func (api *ObscuroAPI) Health() (*host.HealthCheck, error) {
	return api.host.HealthCheck()
}

// Config returns the config status of obscuro host + enclave + db
func (api *ObscuroAPI) Config() (*common.ObscuroNetworkInfo, error) {
	return api.host.ObscuroConfig()
}
