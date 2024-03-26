package walletextension

import (
	"os"
	"time"

	"github.com/ten-protocol/go-ten/tools/walletextension/httpapi"

	"github.com/ten-protocol/go-ten/tools/walletextension/rpcapi"

	"github.com/ten-protocol/go-ten/lib/gethfork/node"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/stopcontrol"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"
	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage"
)

type Container struct {
	stopControl *stopcontrol.StopControl
	logger      gethlog.Logger
	rpcServer   node.Server
}

func NewContainerFromConfig(config wecommon.Config, logger gethlog.Logger) *Container {
	// create the account manager with a single unauthenticated connection
	hostRPCBindAddrWS := wecommon.WSProtocol + config.NodeRPCWebsocketAddress
	hostRPCBindAddrHTTP := wecommon.HTTPProtocol + config.NodeRPCHTTPAddress
	// start the database
	databaseStorage, err := storage.New(config.DBType, config.DBConnectionURL, config.DBPathOverride)
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
	walletExt := rpcapi.NewServices(hostRPCBindAddrHTTP, hostRPCBindAddrWS, databaseStorage, stopControl, version, logger, &config)
	cfg := &node.RPCConfig{
		EnableHTTP: true,
		HTTPPort:   config.WalletExtensionPortHTTP,
		EnableWs:   true,
		WsPort:     config.WalletExtensionPortWS,
		WsPath:     "/v1/",
		HTTPPath:   "/v1/",
		Host:       config.WalletExtensionHost,
	}
	rpcServer := node.NewServer(cfg, logger)

	rpcServer.RegisterRoutes(httpapi.NewHTTPRoutes(walletExt))

	// register the static files
	//staticHandler := api.StaticFilesHandler()
	//rpcServer.RegisterRoutes([]node.Route{{
	//	Name: "/",
	//	Func: func(resp http.ResponseWriter, req *http.Request) {
	//		staticHandler.ServeHTTP(resp, req)
	//	},
	//}})

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
		},
	})

	return NewWalletExtensionContainer(
		stopControl,
		rpcServer,
		logger,
	)
}

func NewWalletExtensionContainer(
	stopControl *stopcontrol.StopControl,
	rpcServer node.Server,
	logger gethlog.Logger,
) *Container {
	return &Container{
		stopControl: stopControl,
		rpcServer:   rpcServer,
		logger:      logger,
	}
}

// Start starts the wallet extension container
func (w *Container) Start() error {
	err := w.rpcServer.Start()
	if err != nil {
		return err
	}
	return nil
}

func (w *Container) Stop() error {
	w.stopControl.Stop()

	if w.rpcServer != nil {
		// rpc server cannot be stopped synchronously as it will kill current request
		go func() {
			// make sure it's not killing the connection before returning the response
			time.Sleep(time.Second) // todo review this sleep
			w.rpcServer.Stop()
		}()
	}

	return nil
}
