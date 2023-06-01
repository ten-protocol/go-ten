package container

import (
	"errors"
	"fmt"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/stopcontrol"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/tools/walletextension"
	"github.com/obscuronet/go-obscuro/tools/walletextension/api"
	"github.com/obscuronet/go-obscuro/tools/walletextension/common"
	"github.com/obscuronet/go-obscuro/tools/walletextension/config"
	"github.com/obscuronet/go-obscuro/tools/walletextension/storage"
	"github.com/obscuronet/go-obscuro/tools/walletextension/useraccountmanager"
	"net/http"

	gethlog "github.com/ethereum/go-ethereum/log"
)

type WalletExtensionContainer struct {
	hostAddr           string
	userAccountManager *useraccountmanager.UserAccountManager
	storage            *storage.Storage
	stopControl        *stopcontrol.StopControl
	logger             gethlog.Logger
	walletExt          *walletextension.WalletExtension
	httpServer         *api.Server
	wsServer           *api.Server
}

func NewWalletExtensionContainerFromConfig(config config.Config, logger gethlog.Logger) *WalletExtensionContainer {
	// create the account manager with a single unauthed connection
	hostRPCBindAddr := common.WSProtocol + config.NodeRPCWebsocketAddress
	unAuthedClient, err := rpc.NewNetworkClient(hostRPCBindAddr)
	if err != nil {
		logger.Crit("unable to create temporary client for request ", log.ErrKey, err)
	}

	userAccountManager := useraccountmanager.NewUserAccountManager(unAuthedClient, logger)

	// TODO (@ziga) - change logic here and get VKs for all users and just add them to the correct accountManager in userAccountManager
	// This line needs to be replaced in future and is here only to enable defaultUser
	userAccountManager.AddUserAccountManager(common.DefaultUser)

	// start the database
	databaseStorage, err := storage.New(config.DBPathOverride)
	if err != nil {
		logger.Crit("unable to create database to store viewing keys ", log.ErrKey, err)
	}

	// We reload the existing viewing keys from the database.
	viewingKeys, err := databaseStorage.GetUserVKs(common.DefaultUser)
	if err != nil {
		logger.Crit("Error getting viewing keys for user:", common.DefaultUser)
	}

	// todo (@ziga) - remove this code below and generalize it for all users
	defaultAccountManager, err := userAccountManager.GetUserAccountManager(common.DefaultUser)
	if err != nil {
		logger.Crit("Error getting account manager for user:", common.DefaultUser)
	}

	for accountAddr, viewingKey := range viewingKeys {
		// create an encrypted RPC client with the signed VK and register it with the enclave
		// todo(@ziga) - Create the clients lazily, to reduce connections to the host.
		client, err := rpc.NewEncNetworkClient(hostRPCBindAddr, viewingKey, logger)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to create encrypted RPC client for persisted account %s", accountAddr), log.ErrKey, err)
			continue
		}
		defaultAccountManager.AddClient(accountAddr, client)
	}

	stopControl := stopcontrol.New()
	walletExt := walletextension.New(hostRPCBindAddr, &userAccountManager, databaseStorage, stopControl, logger)
	httpRoutes := api.NewHTTPRoutes(walletExt)
	httpServer := api.NewHTTPServer(fmt.Sprintf("%s:%d", config.WalletExtensionHost, config.WalletExtensionPortHTTP), httpRoutes)

	wsRoutes := api.NewWSRoutes(walletExt)
	wsServer := api.NewWSServer(fmt.Sprintf("%s:%d", config.WalletExtensionHost, config.WalletExtensionPortWS), wsRoutes)

	return NewWalletExtensionContainer(
		hostRPCBindAddr,
		walletExt,
		&userAccountManager,
		databaseStorage,
		stopControl,
		httpServer,
		wsServer,
		logger,
	)
}

func NewWalletExtensionContainer(
	hostAddr string,
	walletExt *walletextension.WalletExtension,
	userAccountManager *useraccountmanager.UserAccountManager,
	storage *storage.Storage,
	stopControl *stopcontrol.StopControl,
	httpServer *api.Server,
	wsServer *api.Server,
	logger gethlog.Logger,
) *WalletExtensionContainer {
	return &WalletExtensionContainer{
		hostAddr:           hostAddr,
		walletExt:          walletExt,
		userAccountManager: userAccountManager,
		storage:            storage,
		stopControl:        stopControl,
		httpServer:         httpServer,
		wsServer:           wsServer,
		logger:             logger,
	}
}

func (w *WalletExtensionContainer) Start() error {
	httpErrChan := w.httpServer.Start()
	wsErrChan := w.wsServer.Start()

	select {
	case err := <-httpErrChan:
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	case err := <-wsErrChan:
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}
	return nil
}

func (w *WalletExtensionContainer) Stop() error {
	w.stopControl.Stop()

	err := w.httpServer.Stop()
	if err != nil {
		w.logger.Warn("could not shut down wallet extension", log.ErrKey, err)
	}

	err = w.wsServer.Stop()
	if err != nil {
		w.logger.Warn("could not shut down wallet extension", log.ErrKey, err)
	}

	// todo (@pedro) correctly surface shutdown errors
	return nil
}
