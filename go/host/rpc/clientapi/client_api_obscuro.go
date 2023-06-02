package clientapi

import (
	"github.com/obscuronet/go-obscuro/go/common/host"
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
