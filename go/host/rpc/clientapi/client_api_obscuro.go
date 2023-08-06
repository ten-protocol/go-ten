package clientapi

import (
	hostcommon "github.com/obscuronet/go-obscuro/go/common/host"
	"github.com/obscuronet/go-obscuro/go/host"
)

type obscuroAPIServiceLocator interface {
	host.HostControlsLocator
}

// ObscuroAPI implements Obscuro-specific JSON RPC operations.
type ObscuroAPI struct {
	sl obscuroAPIServiceLocator
}

func NewObscuroAPI(serviceLocator obscuroAPIServiceLocator) *ObscuroAPI {
	return &ObscuroAPI{
		sl: serviceLocator,
	}
}

// Health returns the health status of obscuro host + enclave + db
func (api *ObscuroAPI) Health() (*hostcommon.HealthCheck, error) {
	return api.sl.HostControls().HealthCheck()
}
