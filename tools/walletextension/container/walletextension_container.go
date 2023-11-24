package container

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/stopcontrol"
	"github.com/ten-protocol/go-ten/go/rpc"
	"github.com/ten-protocol/go-ten/tools/walletextension"
	"github.com/ten-protocol/go-ten/tools/walletextension/api"
	"github.com/ten-protocol/go-ten/tools/walletextension/config"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage"
	"github.com/ten-protocol/go-ten/tools/walletextension/useraccountmanager"

	gethlog "github.com/ethereum/go-ethereum/log"
	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"
)

type WalletExtensionContainer struct {
	hostAddr           string
	userAccountManager *useraccountmanager.UserAccountManager
	storage            storage.Storage
	stopControl        *stopcontrol.StopControl
	logger             gethlog.Logger
	walletExt          *walletextension.WalletExtension
	httpServer         *api.Server
	wsServer           *api.Server
}

func NewWalletExtensionContainerFromConfig(config config.Config, logger gethlog.Logger) *WalletExtensionContainer {
	// create the account manager with a single unauthenticated connection
	hostRPCBindAddr := wecommon.WSProtocol + config.NodeRPCWebsocketAddress
	unAuthedClient, err := rpc.NewNetworkClient(hostRPCBindAddr)
	if err != nil {
		logger.Crit("unable to create temporary client for request ", log.ErrKey, err)
		os.Exit(1)
	}

	// start the database
	databaseStorage, err := storage.New(config.DBType, config.DBConnectionURL, config.DBPathOverride)
	if err != nil {
		logger.Crit("unable to create database to store viewing keys ", log.ErrKey, err)
		os.Exit(1)
	}
	userAccountManager := useraccountmanager.NewUserAccountManager(unAuthedClient, logger, databaseStorage, hostRPCBindAddr)

	// add default user (when no UserID is provided in the query parameter - for WE endpoints)
	userAccountManager.AddAndReturnAccountManager(hex.EncodeToString([]byte(wecommon.DefaultUser)))

	// add default user to the database (temporary fix before removing wallet extension endpoints)
	accountPrivateKey, err := crypto.GenerateKey()
	if err != nil {
		logger.Error("Unable to generate hey pair for default user", log.ErrKey, err)
	}
	err = databaseStorage.AddUser([]byte(wecommon.DefaultUser), crypto.FromECDSA(accountPrivateKey))
	if err != nil {
		logger.Error("Unable to save default user to the database", log.ErrKey, err)
	}

	// get all users and their private keys from the database
	allUsers, err := databaseStorage.GetAllUsers()
	if err != nil {
		logger.Error(fmt.Errorf("error getting all users from database, %w", err).Error())
	}

	// iterate over users create accountManagers and add all accounts to them per user
	for _, user := range allUsers {
		userAccountManager.AddAndReturnAccountManager(hex.EncodeToString(user.UserID))
		logger.Info(fmt.Sprintf("account manager added for user: %s", hex.EncodeToString(user.UserID)))
	}

	// captures version in the env vars
	version := os.Getenv("OBSCURO_GATEWAY_VERSION")
	if version == "" {
		version = "dev"
	}

	stopControl := stopcontrol.New()
	walletExt := walletextension.New(hostRPCBindAddr, &userAccountManager, databaseStorage, stopControl, version, logger, &config)
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
	storage storage.Storage,
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

// TODO Start should not be a locking process
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
