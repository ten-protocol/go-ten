package walletextension

import (
	"os"
	"path/filepath"
	"time"

	"github.com/ten-protocol/go-ten/go/common/subscription"

	"github.com/ten-protocol/go-ten/tools/walletextension/httpapi"

	"github.com/ten-protocol/go-ten/tools/walletextension/rpcapi"

	"github.com/ten-protocol/go-ten/lib/gethfork/node"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/stopcontrol"
	"github.com/ten-protocol/go-ten/go/enclave/egoutils"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"
	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage"
)

var (
	encryptionKeyFilepath = filepath.Join("/data", "encryption_key.json")
)

type Container struct {
	stopControl     *stopcontrol.StopControl
	logger          gethlog.Logger
	rpcServer       node.Server
	services        *rpcapi.Services
	newHeadsService *subscription.NewHeadsService
}

func NewContainerFromConfig(config wecommon.Config, logger gethlog.Logger) *Container {
	// create the account manager with a single unauthenticated connection
	hostRPCBindAddrWS := wecommon.WSProtocol + config.NodeRPCWebsocketAddress
	hostRPCBindAddrHTTP := wecommon.HTTPProtocol + config.NodeRPCHTTPAddress

	// TODO fix passing this random key to next enclave and not generating it here every time we start the wallet extension
	randomKey, err := wecommon.GenerateRandomKey()
	if err != nil {
		logger.Crit("unable to generate random encryption key", log.ErrKey, err)
		os.Exit(1)
	}

	// seal the encryption key
	err = egoutils.SealAndPersist(string(randomKey), encryptionKeyFilepath, true)
	if err != nil {
		logger.Error("unable to seal and persist encryption key", log.ErrKey, err)
		// os.Exit(1)
	}

	// start the database
	databaseStorage, err := storage.New(config.DBType, config.DBConnectionURL, config.DBPathOverride, randomKey)
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
		WsPath:     wecommon.APIVersion1 + "/",
		HTTPPath:   wecommon.APIVersion1 + "/",
		Host:       config.WalletExtensionHost,
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
			Namespace: "eth",
			Service:   rpcapi.NewFilterAPI(walletExt),
		}, {
			Namespace: "net",
			Service:   rpcapi.NewNetAPI(walletExt),
		}, {
			Namespace: "web3",
			Service:   rpcapi.NewWeb3API(walletExt),
		},
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
