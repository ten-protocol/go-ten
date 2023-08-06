package container

import (
	"fmt"

	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/metrics"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/host"
	"github.com/obscuronet/go-obscuro/go/host/enclave"
	"github.com/obscuronet/go-obscuro/go/host/events"
	"github.com/obscuronet/go-obscuro/go/host/l1"
	"github.com/obscuronet/go-obscuro/go/host/l2"
	"github.com/obscuronet/go-obscuro/go/host/p2p"
	"github.com/obscuronet/go-obscuro/go/host/rpc/clientrpc"
	"github.com/obscuronet/go-obscuro/go/wallet"

	gethlog "github.com/ethereum/go-ethereum/log"
	hostcommon "github.com/obscuronet/go-obscuro/go/common/host"
)

type HostContainer struct {
	// Host and Logger are temporarily public to allow fudge in the test wiring
	// todo (@matt) make these private again once we no longer need to wire the Stop method into the host factories
	Host   hostcommon.Host
	Logger gethlog.Logger
}

func (h *HostContainer) Start() error {
	// make sure the rpc server has a host to render requests
	err := h.Host.Start()
	if err != nil {
		return err
	}
	h.Logger.Info("Started Obscuro host...")
	fmt.Println("Started Obscuro host...")

	return nil
}

func (h *HostContainer) Stop() error {
	// host will not respond to further external requests
	h.Host.Stop()

	return nil
}

// NewHostContainerFromConfig uses config to create all HostContainer dependencies and inject them into a new HostContainer
// (Note: it does not start the HostContainer process, `Start()` must be called on the container)
func NewHostContainerFromConfig(parsedConfig *config.HostInputConfig, logger gethlog.Logger) *HostContainer {
	cfg := parsedConfig.ToHostConfig()

	addr, err := wallet.RetrieveAddress(parsedConfig.PrivateKeyString)
	if err != nil {
		panic("unable to retrieve the Node ID")
	}
	cfg.ID = *addr

	// create the logger if not set - used when the testlogger is injected
	if logger == nil {
		logger = log.New(log.HostCmp, cfg.LogLevel, cfg.LogPath, log.NodeIDKey, cfg.ID)
	}

	fmt.Printf("Building host container with config: %+v\n", cfg)
	logger.Info(fmt.Sprintf("Building host container with config: %+v", cfg))
	return NewHostContainer(cfg, logger)
}

// NewHostContainer builds a host container, passing in default service factories and the config so the host can create
// the services, wire them together and manage their lifecycles
func NewHostContainer(cfg *config.HostConfig, logger gethlog.Logger) *HostContainer {
	hostContainer := &HostContainer{
		Logger: logger,
	}
	h := host.NewHost(cfg, logger,
		p2p.ServiceFactory,
		l1.RepoFactory,
		l1.PublisherFactory,
		l2.RepoFactory,
		enclave.EnclavesFactory,
		events.LogEventFactory,
		metrics.ServiceFactory,
		// the rpc server currently needs a way to stop the host, so we wire through this callback (aim to remove this soon)
		clientrpc.CreateServerFactory(hostContainer.Stop),
		host.DBServiceFactory,
	)
	// the weird ordering on creating the HostContainer and then setting the host onto here is so that the .Stop method
	// can be passed in (see above), this is temporary
	hostContainer.Host = h

	return hostContainer
}
