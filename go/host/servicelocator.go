package host

import (
	"github.com/ethereum/go-ethereum/log"
	hostcommon "github.com/ten-protocol/go-ten/go/common/host"
)

type ServicesRegistry struct {
	services map[string]hostcommon.Service
	logger   log.Logger
}

func NewServicesRegistry(logger log.Logger) *ServicesRegistry {
	return &ServicesRegistry{
		services: make(map[string]hostcommon.Service),
		logger:   logger,
	}
}

func (s *ServicesRegistry) All() map[string]hostcommon.Service {
	return s.services
}

func (s *ServicesRegistry) RegisterService(name string, service hostcommon.Service) {
	if _, ok := s.services[name]; ok {
		s.logger.Crit("service already registered", "name", name)
	}
	s.services[name] = service
}

func (s *ServicesRegistry) getService(name string) hostcommon.Service {
	service, ok := s.services[name]
	if !ok {
		s.logger.Crit("requested service not registered", "name", name)
	}
	return service
}

func (s *ServicesRegistry) P2P() hostcommon.P2P {
	return s.getService(hostcommon.P2PName).(hostcommon.P2P)
}

func (s *ServicesRegistry) L1Data() hostcommon.L1DataService {
	return s.getService(hostcommon.L1DataServiceName).(hostcommon.L1DataService)
}

func (s *ServicesRegistry) L1Publisher() hostcommon.L1Publisher {
	return s.getService(hostcommon.L1PublisherName).(hostcommon.L1Publisher)
}

func (s *ServicesRegistry) L2Repo() hostcommon.L2BatchRepository {
	return s.getService(hostcommon.L2BatchRepositoryName).(hostcommon.L2BatchRepository)
}

func (s *ServicesRegistry) Enclaves() hostcommon.EnclaveService {
	return s.getService(hostcommon.EnclaveServiceName).(hostcommon.EnclaveService)
}

func (s *ServicesRegistry) LogSubs() hostcommon.LogSubscriptionManager {
	return s.getService(hostcommon.LogSubscriptionServiceName).(hostcommon.LogSubscriptionManager)
}
