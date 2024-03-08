package container

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/common"

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
	hostRPCBindAddrWS := wecommon.WSProtocol + config.NodeRPCWebsocketAddress
	hostRPCBindAddrHTTP := wecommon.HTTPProtocol + config.NodeRPCHTTPAddress
	unAuthedClient, err := rpc.NewNetworkClient(hostRPCBindAddrHTTP)
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
	userAccountManager := useraccountmanager.NewUserAccountManager(unAuthedClient, logger, databaseStorage, hostRPCBindAddrHTTP, hostRPCBindAddrWS)

	// add default user (when no UserID is provided in the query parameter - for WE endpoints)
	defaultUserAccountManager := userAccountManager.AddAndReturnAccountManager(hex.EncodeToString([]byte(wecommon.DefaultUser)))

	// add default user to the database (temporary fix before removing wallet extension endpoints)
	accountPrivateKey, err := crypto.GenerateKey()
	if err != nil {
		logger.Error("Unable to generate key pair for default user", log.ErrKey, err)
		os.Exit(1)
	}

	// get all users and their private keys from the database
	allUsers, err := databaseStorage.GetAllUsers()
	if err != nil {
		logger.Error(fmt.Errorf("error getting all users from database, %w", err).Error())
		os.Exit(1)
	}

	// iterate over users create accountManagers and add all defaultUserAccounts to them per user
	for _, user := range allUsers {
		userAccountManager.AddAndReturnAccountManager(hex.EncodeToString(user.UserID))
		logger.Info(fmt.Sprintf("account manager added for user: %s", hex.EncodeToString(user.UserID)))

		// to ensure backwards compatibility we want to load clients for the default user
		// TODO @ziga - this code needs to be removed when removing old wallet extension endpoints
		if bytes.Equal(user.UserID, []byte(wecommon.DefaultUser)) {
			accounts, err := databaseStorage.GetAccounts(user.UserID)
			if err != nil {
				logger.Error(fmt.Errorf("error getting accounts for user: %s, %w", hex.EncodeToString(user.UserID), err).Error())
				os.Exit(1)
			}
			for _, account := range accounts {
				encClient, err := wecommon.CreateEncClient(hostRPCBindAddrWS, account.AccountAddress, user.PrivateKey, account.Signature, account.SignatureType, logger)
				if err != nil {
					logger.Error(fmt.Errorf("error creating new client, %w", err).Error())
					os.Exit(1)
				}

				// add a client to default user
				defaultUserAccountManager.AddClient(common.BytesToAddress(account.AccountAddress), encClient)
			}
		}
	}
	// TODO @ziga - remove this when removing wallet extension endpoints
	err = databaseStorage.AddUser([]byte(wecommon.DefaultUser), crypto.FromECDSA(accountPrivateKey))
	if err != nil {
		logger.Error("Unable to save default user to the database", log.ErrKey, err)
		os.Exit(1)
	}

	// captures version in the env vars
	version := os.Getenv("OBSCURO_GATEWAY_VERSION")
	if version == "" {
		version = "dev"
	}

	stopControl := stopcontrol.New()
	walletExt := walletextension.New(hostRPCBindAddrHTTP, hostRPCBindAddrWS, &userAccountManager, databaseStorage, stopControl, version, logger, &config)
	httpRoutes := api.NewHTTPRoutes(walletExt)
	httpServer := api.NewHTTPServer(fmt.Sprintf("%s:%d", config.WalletExtensionHost, config.WalletExtensionPortHTTP), httpRoutes)

	wsRoutes := api.NewWSRoutes(walletExt)
	wsServer := api.NewWSServer(fmt.Sprintf("%s:%d", config.WalletExtensionHost, config.WalletExtensionPortWS), wsRoutes)
	return NewWalletExtensionContainer(
		hostRPCBindAddrWS,
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

// Start starts the wallet extension container
func (w *WalletExtensionContainer) Start() error {
	httpErrChan := w.httpServer.Start()
	wsErrChan := w.wsServer.Start()

	// Start a goroutine for handling HTTP and WS server errors
	go func() {
		for {
			select {
			case err := <-httpErrChan:
				if errors.Is(err, http.ErrServerClosed) {
					err = w.Stop() // Stop the container when the HTTP server is closed
					if err != nil {
						fmt.Printf("failed to stop gracefully - %s\n", err)
						os.Exit(1)
					}
				} else {
					// for other errors, we just log them
					w.logger.Error("HTTP server error: %v", err)
				}
			case err := <-wsErrChan:
				if errors.Is(err, http.ErrServerClosed) {
					err = w.Stop() // Stop the container when the WS server is closed
					if err != nil {
						fmt.Printf("failed to stop gracefully - %s\n", err)
						os.Exit(1)
					}
				} else {
					// for other errors, we just log them
					w.logger.Error("HTTP server error: %v", err)
				}
			case <-w.stopControl.Done():
				return // Exit the goroutine when stop signal is received
			}
		}
	}()
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
