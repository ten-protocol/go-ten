package container

import (
	"fmt"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/metrics"
	"github.com/obscuronet/go-obscuro/go/config"
	"github.com/obscuronet/go-obscuro/go/ethadapter"
	"github.com/obscuronet/go-obscuro/go/ethadapter/mgmtcontractlib"
	"github.com/obscuronet/go-obscuro/go/host"
	"github.com/obscuronet/go-obscuro/go/host/p2p"
	"github.com/obscuronet/go-obscuro/go/host/rpc/clientapi"
	"github.com/obscuronet/go-obscuro/go/host/rpc/clientrpc"
	"github.com/obscuronet/go-obscuro/go/host/rpc/enclaverpc"
	"github.com/obscuronet/go-obscuro/go/wallet"

	gethlog "github.com/ethereum/go-ethereum/log"
	commonhost "github.com/obscuronet/go-obscuro/go/common/host"
)

const (
	APIVersion1             = "1.0"
	APINamespaceObscuro     = "obscuro"
	APINamespaceEth         = "eth"
	APINamespaceObscuroScan = "obscuroscan"
	APINamespaceNetwork     = "net"
	APINamespaceTest        = "test"
)

type HostContainer struct {
	host           commonhost.Host
	logger         gethlog.Logger
	metricsService *metrics.Service
	rpcServer      clientrpc.Server
}

func (h *HostContainer) Start() error {
	fmt.Println("Starting Obscuro host...")
	h.logger.Info("Starting Obscuro host...")
	h.metricsService.Start()
	// make sure the rpc server has a host to render requests
	h.host.Start()
	h.rpcServer.Start()
	return nil
}

func (h *HostContainer) Stop() error {
	h.metricsService.Stop()
	// make sure the rpc server does not request services from a stopped host
	h.rpcServer.Stop()
	h.host.Stop()
	return nil
}

// NewHostContainerFromConfig uses config to create all HostContainer dependencies and inject them into a new HostContainer
// (Note: it does not start the HostContainer process, `Start()` must be called on the container)
func NewHostContainerFromConfig(parsedConfig *config.HostInputConfig) *HostContainer {
	cfg := parsedConfig.ToHostConfig()

	logger := log.New(log.HostCmp, cfg.LogLevel, cfg.LogPath, log.NodeIDKey, cfg.ID)
	fmt.Printf("Building host container with config: %+v\n", cfg)
	logger.Info(fmt.Sprintf("Building host container with config: %+v", cfg))

	// set the Host ID as the Public Key Address
	ethWallet := wallet.NewInMemoryWalletFromConfig(cfg.PrivateKeyString, cfg.L1ChainID, log.New(log.HostCmp, cfg.LogLevel, cfg.LogPath))
	cfg.ID = ethWallet.Address()

	fmt.Println("Connecting to L1 network...")
	l1Client, err := ethadapter.NewEthClient(cfg.L1NodeHost, cfg.L1NodeWebsocketPort, cfg.L1RPCTimeout, cfg.ID, logger)
	if err != nil {
		logger.Crit("could not create Ethereum client.", log.ErrKey, err)
	}

	// update the wallet nonce
	nonce, err := l1Client.Nonce(ethWallet.Address())
	if err != nil {
		logger.Crit("could not retrieve Ethereum account nonce.", log.ErrKey, err)
	}
	ethWallet.SetNonce(nonce)

	// set the Host ID as the Public Key Address
	cfg.ID = ethWallet.Address()

	enclaveClient := enclaverpc.NewClient(cfg, logger)
	p2pLogger := logger.New(log.CmpKey, log.P2PCmp)
	metricsService := metrics.New(cfg.MetricsEnabled, cfg.MetricsHTTPPort, logger)
	aggP2P := p2p.NewSocketP2PLayer(cfg, p2pLogger, metricsService.Registry())
	rpcServer := clientrpc.NewServer(cfg, logger)

	return NewHostContainer(cfg, aggP2P, l1Client, enclaveClient, ethWallet, rpcServer, logger, metricsService)
}

// NewHostContainer builds a host container with dependency injection rather than from config.
// Useful for testing etc. (want to be able to pass in logger, and also have option to mock out dependencies)
func NewHostContainer(
	cfg *config.HostConfig, // provides various parameters that the host needs to function
	p2p commonhost.P2P, // provides the inbound and outbound p2p communication layer
	l1Client ethadapter.EthClient, // provides inbound and outbound L1 connectivity
	enclaveClient common.Enclave, // provides RPC connection to this host's Enclave
	hostWallet wallet.Wallet, // provides an L1 wallet for the host's transactions
	rpcServer clientrpc.Server, // For communication with Obscuro client applications
	logger gethlog.Logger, // provides logging with context
	metricsService *metrics.Service, // provides the metrics service for other packages to use
) *HostContainer {
	mgmtContractLib := mgmtcontractlib.NewMgmtContractLib(&cfg.RollupContractAddress, logger)

	h := host.NewHost(cfg, p2p, l1Client, enclaveClient, hostWallet, mgmtContractLib, logger, metricsService.Registry())

	hostContainer := &HostContainer{
		host:           h,
		logger:         logger,
		rpcServer:      rpcServer,
		metricsService: metricsService,
	}

	if cfg.HasClientRPCHTTP || cfg.HasClientRPCWebsockets {
		rpcServer.RegisterAPIs([]rpc.API{
			{
				Namespace: APINamespaceObscuro,
				Version:   APIVersion1,
				Service:   clientapi.NewObscuroAPI(h),
				Public:    true,
			},
			{
				Namespace: APINamespaceEth,
				Version:   APIVersion1,
				Service:   clientapi.NewEthereumAPI(h, logger),
				Public:    true,
			},
			{
				Namespace: APINamespaceObscuroScan,
				Version:   APIVersion1,
				Service:   clientapi.NewObscuroScanAPI(h),
				Public:    true,
			},
			{
				Namespace: APINamespaceNetwork,
				Version:   APIVersion1,
				Service:   clientapi.NewNetworkAPI(h),
				Public:    true,
			},
			{
				Namespace: APINamespaceTest,
				Version:   APIVersion1,
				Service:   clientapi.NewTestAPI(nil, hostContainer),
				Public:    true,
			},
			{
				Namespace: APINamespaceEth,
				Version:   APIVersion1,
				Service:   clientapi.NewFilterAPI(h, logger),
				Public:    true,
			},
		})
	}

	return hostContainer
}
