package host

import (
	"database/sql"
	"encoding/json"
	"fmt"
	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ten-protocol/go-ten/go/host/l2"

	"github.com/ten-protocol/go-ten/go/host/enclave"
	"github.com/ten-protocol/go-ten/go/host/l1"
	"github.com/ten-protocol/go-ten/go/host/storage"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/naoina/toml"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/profiler"
	"github.com/ten-protocol/go-ten/go/common/stopcontrol"
	"github.com/ten-protocol/go-ten/go/config"
	"github.com/ten-protocol/go-ten/go/ethadapter"
	"github.com/ten-protocol/go-ten/go/ethadapter/mgmtcontractlib"
	"github.com/ten-protocol/go-ten/go/host/events"
	"github.com/ten-protocol/go-ten/go/responses"
	"github.com/ten-protocol/go-ten/go/wallet"

	gethlog "github.com/ethereum/go-ethereum/log"
	gethmetrics "github.com/ethereum/go-ethereum/metrics"
	hostcommon "github.com/ten-protocol/go-ten/go/common/host"
)

// Implementation of host.Host.
type host struct {
	config   *config.HostConfig
	shortID  uint64
	services *ServicesRegistry // registry of services that the host manages and makes available

	// ignore incoming requests
	stopControl *stopcontrol.StopControl

	db *sql.DB // Stores the host's publicly-available data

	logger gethlog.Logger

	metricRegistry gethmetrics.Registry

	// l2MessageBusAddress is fetched from the enclave but cache it here because it never changes
	l2MessageBusAddress *gethcommon.Address
}

func NewHost(config *config.HostConfig, hostServices *ServicesRegistry, p2p hostcommon.P2PHostService, ethClient ethadapter.EthClient, l1Repo hostcommon.L1RepoService, enclaveClient common.Enclave, ethWallet wallet.Wallet, mgmtContractLib mgmtcontractlib.MgmtContractLib, logger gethlog.Logger, regMetrics gethmetrics.Registry) hostcommon.Host {
	database, err := storage.CreateDBFromConfig(config, logger)
	if err != nil {
		logger.Crit("unable to create database for host", log.ErrKey, err)
	}
	hostIdentity := hostcommon.NewIdentity(config)
	host := &host{
		// config
		config:  config,
		shortID: common.ShortAddress(config.ID),

		// services
		services: hostServices,

		// Initialize the host DB
		db: database,

		logger:         logger,
		metricRegistry: regMetrics,

		stopControl: stopcontrol.New(),
	}

	enclGuardian := enclave.NewGuardian(config, hostIdentity, hostServices, enclaveClient, database, host.stopControl, logger)
	enclService := enclave.NewService(hostIdentity, hostServices, enclGuardian, logger)
	l2Repo := l2.NewBatchRepository(config, hostServices, database, logger)
	subsService := events.NewLogEventManager(hostServices, logger)

	hostServices.RegisterService(hostcommon.P2PName, p2p)
	hostServices.RegisterService(hostcommon.L1BlockRepositoryName, l1Repo)
	maxWaitForL1Receipt := 6 * config.L1BlockTime   // wait ~10 blocks to see if tx gets published before retrying
	retryIntervalForL1Receipt := config.L1BlockTime // retry ~every block
	l1Publisher := l1.NewL1Publisher(hostIdentity, ethWallet, ethClient, mgmtContractLib, l1Repo, host.stopControl, logger, maxWaitForL1Receipt, retryIntervalForL1Receipt)
	hostServices.RegisterService(hostcommon.L1PublisherName, l1Publisher)
	hostServices.RegisterService(hostcommon.L2BatchRepositoryName, l2Repo)
	hostServices.RegisterService(hostcommon.EnclaveServiceName, enclService)
	hostServices.RegisterService(hostcommon.LogSubscriptionServiceName, subsService)

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

	return host
}

// Start validates the host config and starts the Host in a go routine - immediately returns after
func (h *host) Start() error {
	if h.stopControl.IsStopping() {
		return responses.ToInternalError(fmt.Errorf("requested Start with the host stopping"))
	}

	h.validateConfig()

	// start all registered services
	for name, service := range h.services.All() {
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

func (h *host) Config() *config.HostConfig {
	return h.config
}

func (h *host) DB() *sql.DB {
	return h.db
}

func (h *host) EnclaveClient() common.Enclave {
	return h.services.Enclaves().GetEnclaveClient()
}

func (h *host) SubmitAndBroadcastTx(encryptedParams common.EncryptedParamsSendRawTx) (*responses.RawTx, error) {
	if h.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested SubmitAndBroadcastTx with the host stopping"))
	}
	return h.services.Enclaves().SubmitAndBroadcastTx(encryptedParams)
}

func (h *host) Subscribe(id rpc.ID, encryptedLogSubscription common.EncryptedParamsLogSubscription, matchedLogsCh chan []byte) error {
	if h.stopControl.IsStopping() {
		return responses.ToInternalError(fmt.Errorf("requested Subscribe with the host stopping"))
	}
	return h.services.LogSubs().Subscribe(id, encryptedLogSubscription, matchedLogsCh)
}

func (h *host) Unsubscribe(id rpc.ID) {
	if h.stopControl.IsStopping() {
		h.logger.Debug("requested Subscribe with the host stopping")
	}
	h.services.LogSubs().Unsubscribe(id)
}

func (h *host) Stop() error {
	// block all incoming requests
	h.stopControl.Stop()

	h.logger.Info("Host received a stop command. Attempting shutdown...")

	// stop all registered services
	for name, service := range h.services.All() {
		if err := service.Stop(); err != nil {
			h.logger.Error("Failed to stop service", "service", name, log.ErrKey, err)
		}
	}

	if err := h.db.Stop(); err != nil {
		h.logger.Error("Failed to stop DB", log.ErrKey, err)
	}

	h.logger.Info("Host shut down complete.")
	return nil
}

// HealthCheck returns whether the host, enclave and DB are healthy
func (h *host) HealthCheck() (*hostcommon.HealthCheck, error) {
	if h.stopControl.IsStopping() {
		return nil, responses.ToInternalError(fmt.Errorf("requested HealthCheck with the host stopping"))
	}

	healthErrors := make([]string, 0)

	// loop through all registered services and collect their health statuses
	for name, service := range h.services.All() {
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

// ObscuroConfig returns info on the Obscuro network
func (h *host) ObscuroConfig() (*common.ObscuroNetworkInfo, error) {
	if h.l2MessageBusAddress == nil {
		publicCfg, err := h.EnclaveClient().EnclavePublicConfig()
		if err != nil {
			return nil, responses.ToInternalError(fmt.Errorf("unable to get L2 message bus address - %w", err))
		}
		h.l2MessageBusAddress = &publicCfg.L2MessageBusAddress
	}
	return &common.ObscuroNetworkInfo{
		ManagementContractAddress: h.config.ManagementContractAddress,
		L1StartHash:               h.config.L1StartHash,

		SequencerID:         h.config.SequencerID,
		MessageBusAddress:   h.config.MessageBusAddress,
		L2MessageBusAddress: *h.l2MessageBusAddress,
		ImportantContracts:  h.services.L1Publisher().GetImportantContracts(),
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

	if h.config.L1BlockTime == 0 {
		h.logger.Crit("the host must specify an L1 block time")
	}
}
