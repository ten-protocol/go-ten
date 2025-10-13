package walletextension

import (
	"crypto/tls"
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
	// secureHeaders is a small middleware that sets defensive HTTP headers.
	secureHeaders := func(next http.Handler, enableHSTS bool) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// MIME sniffing protection
			w.Header().Set("X-Content-Type-Options", "nosniff")
			// Legacy anti-framing
			w.Header().Set("X-Frame-Options", "DENY")
			// Modern anti-framing (CSP)
			w.Header().Set("Content-Security-Policy", "frame-ancestors 'none'")
			// Referrer policy
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			// Feature gating
			w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=(), payment=(), usb=(), accelerometer=(), gyroscope=(), magnetometer=(), fullscreen=(), autoplay=()")

			// HSTS: only on HTTPS responses
			if enableHSTS && r.TLS != nil {
				w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			}

			next.ServeHTTP(w, r)
		})
	}

	// create the account manager with a single unauthenticated connection
	hostRPCBindAddrWS := wecommon.WSProtocol + config.NodeRPCWebsocketAddress
	hostRPCBindAddrHTTP := wecommon.HTTPProtocol + config.NodeRPCHTTPAddress

	// get the encryption key (method is determined by the config)
	encryptionKey, err := keymanager.GetEncryptionKey(config, logger)
	if err != nil {
		logger.Crit("unable to get encryption key", log.ErrKey, err)
		os.Exit(1)
	}

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
	userStorage, err := storage.New(config.DBType, config.DBConnectionURL, config.DBPathOverride, encryptionKey, logger)
	if err != nil {
		logger.Crit("unable to create database to store viewing keys ", log.ErrKey, err)
		os.Exit(1)
	}

	// captures version in the env vars
	version := os.Getenv("OBSCURO_GATEWAY_VERSION")
	if version == "" {
		version = "dev"
	}

	stopControl := stopcontrol.New()
	walletExt := services.NewServices(hostRPCBindAddrHTTP, hostRPCBindAddrWS, userStorage, stopControl, version, logger, metricsTracker, &config)
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

		// Create HTTP-01 challenge handler + global HTTPS redirect on port 80
		httpMux := http.NewServeMux()
		// Serve only ACME challenges
		httpMux.HandleFunc("/.well-known/acme-challenge/", certManager.HTTPHandler(nil).ServeHTTP)
		// Redirect everything else to HTTPS
		httpMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			target := "https://" + r.Host + r.URL.RequestURI()
			http.Redirect(w, r, target, http.StatusMovedPermanently)
		})

		httpServer := &http.Server{
			Addr:    ":http", // Port 80
			Handler: httpMux,
		}
		go func() {
			if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Error("http redirect server error", "err", err)
			}
		}()

		tlsConfig := &tls.Config{
			GetCertificate: certManager.GetCertificate,
			MinVersion:     tls.VersionTLS12,
			MaxVersion:     tls.VersionTLS13,
			// Prefer strong cipher suites and explicitly exclude legacy/3DES suites (SWEET32)
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				// TLS 1.2 modern suites (TLS 1.3 suites are not configurable and always enabled in Go)
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			},
			CurvePreferences: []tls.CurveID{
				tls.X25519,
				tls.CurveP256,
			},
		}

		// Update RPC server config to use TLS
		cfg.TLSConfig = tlsConfig
	}

	rpcServer := node.NewServer(cfg, logger)

	// Build routes and wrap each with secure headers
	routes := httpapi.NewHTTPRoutes(walletExt)
	for i := range routes {
		// Wrap route function with middleware
		handler := http.HandlerFunc(routes[i].Func)
		secured := secureHeaders(handler, config.EnableTLS)
		// Adapt back to func(resp, req)
		routes[i].Func = func(w http.ResponseWriter, r *http.Request) { secured.ServeHTTP(w, r) }
	}
	rpcServer.RegisterRoutes(routes)

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
