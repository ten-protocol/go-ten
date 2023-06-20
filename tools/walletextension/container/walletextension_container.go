package container

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/common/stopcontrol"
	"github.com/obscuronet/go-obscuro/go/common/viewingkey"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/tools/walletextension"
	"github.com/obscuronet/go-obscuro/tools/walletextension/api"
	"github.com/obscuronet/go-obscuro/tools/walletextension/config"
	"github.com/obscuronet/go-obscuro/tools/walletextension/storage"
	"github.com/obscuronet/go-obscuro/tools/walletextension/useraccountmanager"

	gethcrypto "github.com/ethereum/go-ethereum/crypto"
	gethlog "github.com/ethereum/go-ethereum/log"
	wecommon "github.com/obscuronet/go-obscuro/tools/walletextension/common"
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
	hostRPCBindAddr := wecommon.WSProtocol + config.NodeRPCWebsocketAddress
	unAuthedClient, err := rpc.NewNetworkClient(hostRPCBindAddr)
	if err != nil {
		logger.Crit("unable to create temporary client for request ", log.ErrKey, err)
	}

	userAccountManager := useraccountmanager.NewUserAccountManager(unAuthedClient, logger)

	// TODO (@ziga) - change logic here and get VKs for all users and just add them to the correct accountManager in userAccountManager
	// This line needs to be replaced in future and is here only to enable defaultUser
	userAccountManager.AddUserAccountManager(wecommon.DefaultUser)

	// todo (@ziga) - remove this code below and generalize it for all users
	defaultAccountManager, err := userAccountManager.GetUserAccountManager(wecommon.DefaultUser)
	if err != nil {
		logger.Crit("Error getting account manager for user:", wecommon.DefaultUser)
	}

	// start the database
	databaseStorage, err := storage.New(config.DBPathOverride)
	if err != nil {
		logger.Crit("unable to create database to store viewing keys ", log.ErrKey, err)
	}

	// We reload data from the database
	// todo (@ziga) - change this code once we can handle multiple users
	defaultPrivateKeyBytes, err := databaseStorage.GetUserPrivateKey([]byte(wecommon.DefaultUser))
	if err != nil {
		logger.Crit("Error getting private key for user: ", wecommon.DefaultUser)
	}

	// if we have account load the data and add clients
	if !bytes.Equal(defaultPrivateKeyBytes, []byte{}) {
		userAccounts, err := databaseStorage.GetAccounts([]byte(wecommon.DefaultUser))
		if err != nil {
			logger.Crit("Error getting addresses for user: ", wecommon.DefaultUser)
		}

		privKeyECDSA, err := gethcrypto.ToECDSA(defaultPrivateKeyBytes)
		if err != nil {
			logger.Crit("Unable to covert private key to ECDSA")
		}
		privKey := ecies.ImportECDSA(privKeyECDSA)
		for _, acc := range userAccounts {
			a := common.BytesToAddress(acc.AccountAddress)
			viewingKey := &viewingkey.ViewingKey{
				Account:    &a,
				PrivateKey: privKey,
				PublicKey:  gethcrypto.CompressPubkey(&privKeyECDSA.PublicKey),
				Signature:  acc.Signature,
			}

			client, err := rpc.NewEncNetworkClient(hostRPCBindAddr, viewingKey, logger)
			if err != nil {
				logger.Error(fmt.Sprintf("failed to create encrypted RPC client for persisted account %s", viewingKey.Account.Hex()), log.ErrKey, err)
				continue
			}
			defaultAccountManager.AddClient(*viewingKey.Account, client)
		}
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
