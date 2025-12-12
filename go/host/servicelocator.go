package host

import (
	"github.com/ethereum/go-ethereum/log"
	hostcommon "github.com/ten-protocol/go-ten/go/common/host"
)

type serviceEntry struct {
	name    string
	service hostcommon.Service
	order   int
}

type ServicesRegistry struct {
	orderedServices []serviceEntry
	serviceIndex    map[string]hostcommon.Service
	logger          log.Logger
}

func NewServicesRegistry(logger log.Logger) *ServicesRegistry {
	return &ServicesRegistry{
		orderedServices: make([]serviceEntry, 0),
		serviceIndex:    make(map[string]hostcommon.Service),
		logger:          logger,
	}
}

func (s *ServicesRegistry) All() map[string]hostcommon.Service {
	return s.serviceIndex
}

func (s *ServicesRegistry) RegisterService(name string, service hostcommon.Service, order int) {
	if _, ok := s.serviceIndex[name]; ok {
		s.logger.Crit("service already registered", "name", name)
	}

	entry := serviceEntry{name: name, service: service, order: order}

	// Insert in order
	inserted := false
	for i, existing := range s.orderedServices {
		if order < existing.order {
			s.orderedServices = append(s.orderedServices[:i], append([]serviceEntry{entry}, s.orderedServices[i:]...)...)
			inserted = true
			break
		}
	}
	if !inserted {
		s.orderedServices = append(s.orderedServices, entry)
	}

	s.serviceIndex[name] = service
}

func (s *ServicesRegistry) getService(name string) hostcommon.Service {
	service, ok := s.serviceIndex[name]
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

func (s *ServicesRegistry) Start() error {
	for _, entry := range s.orderedServices {
		s.logger.Info("Starting service", "name", entry.name)
		if err := entry.service.Start(); err != nil {
			s.logger.Error("Failed to start service", "name", entry.name, "error", err)
			return err
		}
	}
	return nil
}

func (s *ServicesRegistry) Stop() error {
	for i := len(s.orderedServices) - 1; i >= 0; i-- {
		entry := s.orderedServices[i]
		s.logger.Info("Stopping service", "name", entry.name)
		if err := entry.service.Stop(); err != nil {
			s.logger.Error("Failed to stop service", "name", entry.name, "error", err)
			return err
		}
	}
	return nil
}
