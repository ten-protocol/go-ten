package host

import (
	hostcommon "github.com/obscuronet/go-obscuro/go/common/host"
	"github.com/obscuronet/go-obscuro/go/host/enclave"
	"github.com/obscuronet/go-obscuro/go/host/l1"
	"github.com/obscuronet/go-obscuro/go/host/l2"
)

type P2PHostService interface {
	hostcommon.Service
	hostcommon.P2P
}

// Services are components that have their own lifecycle of Start/Stop/Health
type Services struct {
	P2P            P2PHostService
	L1Repo         *l1.Repository
	L2Repo         *l2.Repository
	EnclaveService *enclave.Service
}

// All returns the serv
func (s *Services) All() []hostcommon.Service {
	return []hostcommon.Service{
		s.EnclaveService,
		s.L1Repo,
		s.L2Repo,
		s.P2P,
	}
}
