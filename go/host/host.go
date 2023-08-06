package host

import (
	"encoding/json"
	"fmt"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/naoina/toml"
	"github.com/obscuronet/go-obscuro/go/common"
	hostcommon "github.com/obscuronet/go-obscuro/go/common/host"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/profiler"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/host/db"
)

// The host has very little logic of its own.
// It is responsible for creating, starting, stopping and monitoring a variety of interacting services.
type host struct {
	config *config.HostConfig

	// services
	allServices map[string]hostcommon.Service
	p2p         hostcommon.P2P                    // provides inbound and outbound communication with other obscuro hosts
	l1Repo      hostcommon.L1BlockRepository      // provides L1 block data and subscriptions for new blocks
	l1Publisher hostcommon.L1Publisher            // provides access to management contract functions on the L1
	l2Repo      hostcommon.L2BatchRepository      // provides L2 batch data and subscriptions for new batches
	enclaves    hostcommon.EnclaveService         // provides access to the node enclave(s) to submit and request data
	logSubs     hostcommon.LogSubscriptionManager // manages event log subscriptions
	metrics     hostcommon.Metrics                // record and expose metrics about the host
	rpcServer   hostcommon.RPCServer              // provides RPC interfaces for the host
	db          *db.DB                            // Stores the host's publicly-available data

	logger gethlog.Logger
}

// NewHost creates a new host instance.
// It takes config, a logger, and a range of service factories.
func NewHost(
	config *config.HostConfig,
	logger gethlog.Logger,
	p2pFactory ServiceFactory[P2PService],
	l1RepoFactory ServiceFactory[L1BlockRepositoryService],
	l1PublisherFactory ServiceFactory[L1PublisherService],
	l2RepoFactory ServiceFactory[L2BatchRepositoryService],
	enclavesFactory ServiceFactory[EnclaveHostService],
	logSubsFactory ServiceFactory[LogSubscriptionManagerService],
	metricsFactory ServiceFactory[MetricsService],
	rpcServerFactory ServiceFactory[RPCServerService],
	dbFactory ServiceFactory[*db.DB],
) hostcommon.Host {
	h := &host{allServices: make(map[string]hostcommon.Service), config: config, logger: logger}

	h.p2p = setupService("p2p", h, p2pFactory)
	h.l1Repo = setupService("l1-repo", h, l1RepoFactory)
	h.l1Publisher = setupService("l1-pub", h, l1PublisherFactory)
	h.l2Repo = setupService("l2-repo", h, l2RepoFactory)
	h.enclaves = setupService("enclaves", h, enclavesFactory)
	h.logSubs = setupService("log-subs", h, logSubsFactory)
	h.metrics = setupService("metrics", h, metricsFactory)
	h.rpcServer = setupService("rpc-server", h, rpcServerFactory)
	h.db = setupService("db", h, dbFactory)

	var prof *profiler.Profiler
	if config.ProfilerEnabled {
		prof = profiler.NewProfiler(profiler.DefaultHostPort, logger)
		err := prof.Start()
		if err != nil {
			logger.Crit("unable to start the profiler: %s", log.ErrKey, err)
		}
	}

	jsonConfig, _ := json.MarshalIndent(config, "", "  ")
	logger.Info("Host service created with following config:", log.CfgKey, string(jsonConfig))

	return h
}

func (h *host) P2P() hostcommon.P2P {
	return h.p2p
}

func (h *host) L1Repo() hostcommon.L1BlockRepository {
	return h.l1Repo
}

func (h *host) L1Publisher() hostcommon.L1Publisher {
	return h.l1Publisher
}

func (h *host) L2Repo() hostcommon.L2BatchRepository {
	return h.l2Repo
}

func (h *host) Enclave() hostcommon.EnclaveService {
	return h.enclaves
}

func (h *host) LogSubs() hostcommon.LogSubscriptionManager {
	return h.logSubs
}

func (h *host) Metrics() hostcommon.Metrics {
	return h.metrics
}

func (h *host) DB() *db.DB {
	return h.db
}

func (h *host) HostControls() hostcommon.HostControls {
	return h
}

// Start validates the host config and starts the Host in a go routine - immediately returns after
func (h *host) Start() error {
	h.validateConfig()

	// start all registered services
	for name, service := range h.allServices {
		err := service.Start()
		if err != nil {
			return fmt.Errorf("could not start service=%s: %w", name, err)
		}
	}

	tomlConfig, err := toml.Marshal(h.config)
	if err != nil {
		return fmt.Errorf("could not print host config - %w", err)
	}
	h.logger.Info("Host started with following config", log.CfgKey, string(tomlConfig))

	return nil
}

func (h *host) Stop() {
	h.logger.Info("Host received a stop command. Attempting shutdown...")

	// stop all registered services
	for _, service := range h.allServices {
		service.Stop()
	}

	h.logger.Info("Host shut down complete.")
}

// HealthCheck returns whether the host, enclave and DB are healthy
func (h *host) HealthCheck() (*hostcommon.HealthCheck, error) {
	healthErrors := make([]string, 0)

	// loop through all registered services and collect their health statuses
	for name, service := range h.allServices {
		status := service.HealthStatus()
		if !status.OK() {
			healthErrors = append(healthErrors, fmt.Sprintf("[%s] not healthy - %s", name, status.Message()))
		}
	}

	return &hostcommon.HealthCheck{
		OverallHealth: len(healthErrors) == 0,
		Errors:        healthErrors,
	}, nil
}

// Checks the host config is valid.
func (h *host) validateConfig() {
	if h.config.IsGenesis && h.config.NodeType != common.Sequencer {
		h.logger.Crit("genesis node must be the sequencer")
	}
	if !h.config.IsGenesis && h.config.NodeType == common.Sequencer {
		h.logger.Crit("only the genesis node can be a sequencer")
	}

	if h.config.P2PPublicAddress == "" {
		h.logger.Crit("the host must specify a public P2P address")
	}
}

func setupService[S hostcommon.Service](serviceName string, h *host, factory ServiceFactory[S]) S {
	logger := h.logger.New(log.ServiceKey, serviceName)
	service, err := factory(h.config, h, logger)
	if err != nil {
		h.logger.Crit("could not create service", log.ServiceKey, serviceName, log.ErrKey, err)
	}
	h.allServices[serviceName] = service
	return service
}
