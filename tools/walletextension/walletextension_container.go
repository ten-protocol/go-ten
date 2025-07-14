package walletextension

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ten-protocol/go-ten/tools/walletextension/metrics"
	"github.com/ten-protocol/go-ten/tools/walletextension/services"

	"github.com/ten-protocol/go-ten/go/common/subscription"

	"github.com/ten-protocol/go-ten/tools/walletextension/httpapi"

	"github.com/ten-protocol/go-ten/tools/walletextension/rpcapi"

	"github.com/ten-protocol/go-ten/lib/gethfork/node"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/stopcontrol"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"
	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/keymanager"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage"
	"golang.org/x/crypto/acme/autocert"
)

type Container struct {
	stopControl     *stopcontrol.StopControl
	logger          gethlog.Logger
	rpcServer       node.Server
	services        *services.Services
	newHeadsService *subscription.NewHeadsService
}

func NewContainerFromConfig(config wecommon.Config, logger gethlog.Logger) *Container {
	logger.Info("NewContainerFromConfig: Starting wallet extension container initialization")
	logger.Info("NewContainerFromConfig: Configuration", "dbType", config.DBType, "insideEnclave", config.InsideEnclave, "encryptionKeySource", config.EncryptionKeySource)
	
	// Log SGX environment details at startup
	logSGXStartupEnvironment(logger)
	
	// create the account manager with a single unauthenticated connection
	hostRPCBindAddrWS := wecommon.WSProtocol + config.NodeRPCWebsocketAddress
	hostRPCBindAddrHTTP := wecommon.HTTPProtocol + config.NodeRPCHTTPAddress
	logger.Info("NewContainerFromConfig: Node RPC addresses", "websocket", hostRPCBindAddrWS, "http", hostRPCBindAddrHTTP)

	// get the encryption key (method is determined by the config)
	logger.Info("NewContainerFromConfig: Attempting to get encryption key")
	encryptionKey, err := keymanager.GetEncryptionKey(config, logger)
	if err != nil {
		logger.Crit("unable to get encryption key", log.ErrKey, err)
		logger.Error("NewContainerFromConfig: Encryption key acquisition failed - this will prevent container startup")
		os.Exit(1)
	}
	logger.Info("NewContainerFromConfig: Successfully obtained encryption key")

	// Create metrics tracker
	var metricsTracker metrics.Metrics
	if config.DBType == "cosmosDB" {
		metricsStorage, err := storage.NewMetricsStorage(config.DBType, config.DBConnectionURL)
		if err != nil {
			logger.Crit("unable to create metrics storage", log.ErrKey, err)
			os.Exit(1)
		}
		metricsTracker = metrics.NewMetricsTracker(metricsStorage, logger)
	} else {
		metricsTracker = metrics.NewNoOpMetricsTracker(logger)
	}

	// start the database with the encryption key
	logger.Info("NewContainerFromConfig: Initializing database storage", "dbType", config.DBType, "connectionURL", config.DBConnectionURL, "pathOverride", config.DBPathOverride)
	userStorage, err := storage.New(config.DBType, config.DBConnectionURL, config.DBPathOverride, encryptionKey, logger)
	if err != nil {
		logger.Crit("unable to create database to store viewing keys ", log.ErrKey, err)
		logger.Error("NewContainerFromConfig: Database initialization failed")
		os.Exit(1)
	}
	logger.Info("NewContainerFromConfig: Database storage initialized successfully")

	// captures version in the env vars
	version := os.Getenv("OBSCURO_GATEWAY_VERSION")
	if version == "" {
		version = "dev"
	}

	stopControl := stopcontrol.New()
	logger.Info("NewContainerFromConfig: Creating wallet extension services", "version", version)
	walletExt := services.NewServices(hostRPCBindAddrHTTP, hostRPCBindAddrWS, userStorage, stopControl, version, logger, metricsTracker, &config)
	logger.Info("NewContainerFromConfig: Wallet extension services created successfully")
	cfg := &node.RPCConfig{
		EnableHTTP: true,
		HTTPPort:   config.WalletExtensionPortHTTP,
		EnableWs:   true,
		WsPort:     config.WalletExtensionPortWS,
		WsPath:     wecommon.APIVersion1 + "/",
		HTTPPath:   wecommon.APIVersion1 + "/",
		Host:       config.WalletExtensionHost,
	}

	// check if TLS is enabled
	if config.EnableTLS {
		// Create autocert manager for automatic certificate management
		// Generating a certificate consists of the following steps:
		// generating a new private key
		// domain ownership verification (HTTP-01 challenge since certManager.HTTPHandler(nil) is set)
		// Certificate Signing Request (CRS) is generated
		// CRS is sent to CA (Let's Encrypt) via ACME (automated certificate management environment) client
		// CA verifies CRS and issues a certificate
		// Store certificate and private key in certificate storage based on the database type
		certStorage, err := storage.NewCertStorage(config.DBType, config.DBConnectionURL, encryptionKey, config.EncryptingCertificateEnabled, logger)
		if err != nil {
			logger.Crit("unable to create certificate storage", log.ErrKey, err)
			os.Exit(1)
		}

		certManager := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(config.TLSDomain),
			Cache:      certStorage,
		}

		// Create HTTP-01 challenge handler
		httpServer := &http.Server{
			Addr:    ":http", // Port 80
			Handler: certManager.HTTPHandler(nil),
		}
		go httpServer.ListenAndServe() // Start HTTP server for ACME challenges

		tlsConfig := &tls.Config{
			GetCertificate: certManager.GetCertificate,
			MinVersion:     tls.VersionTLS12,
		}

		// Update RPC server config to use TLS
		cfg.TLSConfig = tlsConfig
	}

	rpcServer := node.NewServer(cfg, logger)

	rpcServer.RegisterRoutes(httpapi.NewHTTPRoutes(walletExt))

	// register all RPC endpoints exposed by a typical Geth node
	rpcServer.RegisterAPIs([]gethrpc.API{
		{
			Namespace: "eth",
			Service:   rpcapi.NewEthereumAPI(walletExt),
		}, {
			Namespace: "eth",
			Service:   rpcapi.NewBlockChainAPI(walletExt),
		}, {
			Namespace: "eth",
			Service:   rpcapi.NewTransactionAPI(walletExt),
		}, {
			Namespace: "txpool",
			Service:   rpcapi.NewTxPoolAPI(walletExt),
		}, {
			Namespace: "debug",
			Service:   rpcapi.NewDebugAPI(walletExt),
		}, {
			Namespace: "sessionkeys",
			Service:   rpcapi.NewSessionKeyAPI(walletExt),
		}, {
			Namespace: "eth",
			Service:   rpcapi.NewFilterAPI(walletExt),
		}, {
			Namespace: "net",
			Service:   rpcapi.NewNetAPI(walletExt),
		}, {
			Namespace: "web3",
			Service:   rpcapi.NewWeb3API(walletExt),
		}, {
			Namespace: "ten",
			Service:   rpcapi.NewTenAPI(walletExt),
		},
	})

	// Add metrics tracker to stop sequence
	stopControl.OnStop(func() {
		metricsTracker.Stop()
	})

	return &Container{
		stopControl:     stopControl,
		rpcServer:       rpcServer,
		newHeadsService: walletExt.NewHeadsService,
		services:        walletExt,
		logger:          logger,
	}
}

// Start starts the wallet extension container
func (w *Container) Start() error {
	err := w.newHeadsService.Start()
	if err != nil {
		return err
	}

	err = w.rpcServer.Start()
	if err != nil {
		return err
	}
	return nil
}

func (w *Container) Stop() error {
	w.stopControl.Stop()
	_ = w.newHeadsService.Stop()

	if w.rpcServer != nil {
		// rpc server cannot be stopped synchronously as it will kill current request
		go func() {
			// make sure it's not killing the connection before returning the response
			time.Sleep(time.Second) // todo review this sleep
			w.rpcServer.Stop()
		}()
	}

	w.services.Stop()
	return nil
}

// logSGXStartupEnvironment logs SGX environment information at container startup
func logSGXStartupEnvironment(logger gethlog.Logger) {
	logger.Info("SGX Startup Environment Check: Beginning comprehensive environment analysis")
	
	// Log critical environment variables
	envVars := []string{
		"OE_SIMULATION", "AESM_PATH", "PCCS_URL", "SGX_AESM_ADDR", 
		"SGX_SPID", "SGX_LINKABLE", "SGX_DEBUG", "SGX_MODE",
		"KUBERNETES_SERVICE_HOST", "KUBERNETES_SERVICE_PORT",
	}
	
	for _, envVar := range envVars {
		value := os.Getenv(envVar)
		if value != "" {
			logger.Info(fmt.Sprintf("SGX Startup Environment: %s=%s", envVar, value))
		} else {
			logger.Info(fmt.Sprintf("SGX Startup Environment: %s is not set", envVar))
		}
	}
	
	// Check container environment
	if _, err := os.Stat("/.dockerenv"); err == nil {
		logger.Info("SGX Startup Environment: Running inside Docker container")
	}
	
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		logger.Info("SGX Startup Environment: Running inside Kubernetes cluster")
		logger.Info("SGX Startup Environment: K8s namespace", "namespace", os.Getenv("KUBERNETES_NAMESPACE"))
		logger.Info("SGX Startup Environment: K8s pod name", "podName", os.Getenv("HOSTNAME"))
	}
	
	// Check SGX device availability
	sgxDevices := []string{"/dev/sgx_enclave", "/dev/sgx_provision", "/dev/sgx/enclave", "/dev/sgx/provision"}
	for _, device := range sgxDevices {
		if stat, err := os.Stat(device); err == nil {
			logger.Info(fmt.Sprintf("SGX Startup Environment: Device %s exists, mode: %v", device, stat.Mode()))
		} else {
			logger.Info(fmt.Sprintf("SGX Startup Environment: Device %s not accessible: %v", device, err))
		}
	}
	
	// Check AESM socket
	aesmSocket := os.Getenv("AESM_PATH")
	if aesmSocket == "" {
		aesmSocket = "/var/run/aesmd/aesm.socket"
	}
	if stat, err := os.Stat(aesmSocket); err == nil {
		logger.Info(fmt.Sprintf("SGX Startup Environment: AESM socket %s exists, mode: %v", aesmSocket, stat.Mode()))
	} else {
		logger.Info(fmt.Sprintf("SGX Startup Environment: AESM socket %s not accessible: %v", aesmSocket, err))
	}
	
	// Check data directory
	dataDir := "/data"
	if stat, err := os.Stat(dataDir); err == nil {
		logger.Info(fmt.Sprintf("SGX Startup Environment: Data directory %s exists, mode: %v", dataDir, stat.Mode()))
	} else {
		logger.Info(fmt.Sprintf("SGX Startup Environment: Data directory %s not accessible: %v", dataDir, err))
	}
	
	// Check memory limits (important for SGX)
	if memLimit, exists := os.LookupEnv("MEM_LIMIT"); exists {
		logger.Info(fmt.Sprintf("SGX Startup Environment: Memory limit set to %s", memLimit))
	}
	
	logger.Info("SGX Startup Environment Check: Environment analysis completed")
}
