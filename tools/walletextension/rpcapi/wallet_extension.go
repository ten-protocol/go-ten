package rpcapi

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ten-protocol/go-ten/go/wallet"

	"github.com/ten-protocol/go-ten/tools/walletextension/cache"

	"github.com/ten-protocol/go-ten/go/obsclient"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/stopcontrol"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"github.com/ten-protocol/go-ten/go/rpc"
	"github.com/ten-protocol/go-ten/tools/walletextension/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage"
)

// Services handles the various business logic for the api endpoints
type Services struct {
	HostAddrHTTP string                                        // The HTTP address on which the Ten host can be reached
	HostAddrWS   string                                        // The WS address on which the Ten host can be reached
	unsignedVKs  map[gethcommon.Address]*viewingkey.ViewingKey // Map temporarily holding VKs that have been generated but not yet signed
	Storage      storage.Storage
	logger       gethlog.Logger
	FileLogger   gethlog.Logger
	stopControl  *stopcontrol.StopControl
	version      string
	tenClient    *obsclient.ObsClient
	Cache        cache.Cache
	Config       *common.Config
	DefaultUser  []byte
}

func NewServices(hostAddrHTTP string, hostAddrWS string, storage storage.Storage, stopControl *stopcontrol.StopControl, version string, logger gethlog.Logger, config *common.Config) *Services {
	rpcClient, err := rpc.NewNetworkClient(hostAddrHTTP)
	if err != nil {
		logger.Error(fmt.Errorf("could not create RPC client on %s. Cause: %w", hostAddrHTTP, err).Error())
		panic(err)
	}
	newTenClient := obsclient.NewObsClient(rpcClient)
	newFileLogger := common.NewFileLogger()
	newGatewayCache, err := cache.NewCache(logger)
	if err != nil {
		logger.Error(fmt.Errorf("could not create cache. Cause: %w", err).Error())
		panic(err)
	}

	userID := addDefaultUser(config, logger, storage)

	return &Services{
		HostAddrHTTP: hostAddrHTTP,
		HostAddrWS:   hostAddrWS,
		unsignedVKs:  map[gethcommon.Address]*viewingkey.ViewingKey{},
		Storage:      storage,
		logger:       logger,
		FileLogger:   newFileLogger,
		stopControl:  stopControl,
		version:      version,
		tenClient:    newTenClient,
		Cache:        newGatewayCache,
		Config:       config,
		DefaultUser:  userID,
	}
}

func addDefaultUser(config *common.Config, logger gethlog.Logger, storage storage.Storage) []byte {
	userAccountKey, err := crypto.GenerateKey()
	if err != nil {
		panic(fmt.Errorf("error generating default user key"))
	}

	wallet := wallet.NewInMemoryWalletFromPK(big.NewInt(int64(config.TenChainID)), userAccountKey, logger)
	vk, err := viewingkey.GenerateViewingKeyForWallet(wallet)
	if err != nil {
		panic(err)
	}

	// create UserID and store it in the database with the private key
	userID := viewingkey.CalculateUserID(common.PrivateKeyToCompressedPubKey(vk.PrivateKey))

	err = storage.AddUser(userID, crypto.FromECDSA(vk.PrivateKey.ExportECDSA()))
	if err != nil {
		panic(fmt.Errorf("error saving default user"))
	}

	err = storage.AddAccount(userID, vk.Account.Bytes(), vk.SignatureWithAccountKey)
	if err != nil {
		panic(fmt.Errorf("error saving account %s for default user %s", vk.Account.Hex(), userID))
	}
	return userID
}

// IsStopping returns whether the WE is stopping
func (w *Services) IsStopping() bool {
	return w.stopControl.IsStopping()
}

// Logger returns the WE set logger
func (w *Services) Logger() gethlog.Logger {
	return w.logger
}

// todo - once the logic in routes has been moved to RPC functions, these methods can be moved there

// GenerateViewingKey generates the user viewing key and waits for signature
func (w *Services) GenerateViewingKey(addr gethcommon.Address) (string, error) {
	// Requested to generate viewing key for address(old way): %s", addr.Hex()))
	viewingKeyPrivate, err := crypto.GenerateKey()
	if err != nil {
		return "", fmt.Errorf("unable to generate a new keypair - %w", err)
	}

	viewingPublicKeyBytes := crypto.CompressPubkey(&viewingKeyPrivate.PublicKey)
	viewingPrivateKeyEcies := ecies.ImportECDSA(viewingKeyPrivate)

	w.unsignedVKs[addr] = &viewingkey.ViewingKey{
		Account:                 &addr,
		PrivateKey:              viewingPrivateKeyEcies,
		PublicKey:               viewingPublicKeyBytes,
		SignatureWithAccountKey: nil, // we await a signature from the user before we can set up the EncRPCClient
	}

	// compress the viewing key and convert it to hex string ( this is what Metamask signs)
	viewingKeyBytes := crypto.CompressPubkey(&viewingKeyPrivate.PublicKey)
	return hex.EncodeToString(viewingKeyBytes), nil
}

// SubmitViewingKey checks the signed viewing key and stores it
func (w *Services) SubmitViewingKey(address gethcommon.Address, signature []byte) error {
	audit(w, "Requested to submit a viewing key (old way): %s", address.Hex())
	vk, found := w.unsignedVKs[address]
	if !found {
		return fmt.Errorf(fmt.Sprintf("no viewing key found to sign for acc=%s, please call %s to generate key before sending signature", address, common.PathGenerateViewingKey))
	}

	// We transform the V from 27/28 to 0/1. This same change is made in Geth internals, for legacy reasons to be able
	// to recover the address: https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L452-L459
	signature[64] -= 27

	vk.SignatureWithAccountKey = signature

	//err := w.Storage.AddUser(defaultUserId, crypto.FromECDSA(vk.PrivateKey.ExportECDSA()))
	//if err != nil {
	//	return fmt.Errorf("error saving user: %s", common.DefaultUser)
	//}
	//
	//err = w.Storage.AddAccount(defaultUserId, vk.Account.Bytes(), vk.SignatureWithAccountKey)
	//if err != nil {
	//	return fmt.Errorf("error saving account %s for user %s", vk.Account.Hex(), common.DefaultUser)
	//}

	// finally we remove the VK from the pending 'unsigned VKs' map now the client has been created
	delete(w.unsignedVKs, address)

	return nil
}

// GenerateAndStoreNewUser generates new key-pair and userID, stores it in the database and returns hex encoded userID and error
func (w *Services) GenerateAndStoreNewUser() (string, error) {
	requestStartTime := time.Now()
	// generate new key-pair
	viewingKeyPrivate, err := crypto.GenerateKey()
	viewingPrivateKeyEcies := ecies.ImportECDSA(viewingKeyPrivate)
	if err != nil {
		w.Logger().Error(fmt.Sprintf("could not generate new keypair: %s", err))
		return "", err
	}

	// create UserID and store it in the database with the private key
	userID := viewingkey.CalculateUserID(common.PrivateKeyToCompressedPubKey(viewingPrivateKeyEcies))
	err = w.Storage.AddUser(userID, crypto.FromECDSA(viewingPrivateKeyEcies.ExportECDSA()))
	if err != nil {
		w.Logger().Error(fmt.Sprintf("failed to save user to the database: %s", err))
		return "", err
	}

	hexUserID := hex.EncodeToString(userID)

	requestEndTime := time.Now()
	duration := requestEndTime.Sub(requestStartTime)
	audit(w, "Storing new userID: %s, duration: %d ", hexUserID, duration.Milliseconds())
	return hexUserID, nil
}

// AddAddressToUser checks if a message is in correct format and if signature is valid. If all checks pass we save address and signature against userID
func (w *Services) AddAddressToUser(hexUserID string, address string, signature []byte) error {
	requestStartTime := time.Now()
	addressFromMessage := gethcommon.HexToAddress(address)
	// check if a message was signed by the correct address and if the signature is valid
	sigAddrs, err := viewingkey.CheckEIP712Signature(hexUserID, signature, int64(w.Config.TenChainID))
	if err != nil {
		return fmt.Errorf("signature is not valid: %w", err)
	}

	if sigAddrs.Hex() != address {
		return fmt.Errorf("signature is not valid. Signature address %s!=%s ", sigAddrs, address)
	}

	// register the account for that viewing key
	userIDBytes, err := common.GetUserIDbyte(hexUserID)
	if err != nil {
		w.Logger().Error(fmt.Errorf("error decoding string (%s), %w", hexUserID[2:], err).Error())
		return errors.New("error decoding userID. It should be in hex format")
	}
	err = w.Storage.AddAccount(userIDBytes, addressFromMessage.Bytes(), signature)
	if err != nil {
		w.Logger().Error(fmt.Errorf("error while storing account (%s) for user (%s): %w", addressFromMessage.Hex(), hexUserID, err).Error())
		return err
	}

	requestEndTime := time.Now()
	duration := requestEndTime.Sub(requestStartTime)
	audit(w, "Storing new address for user: %s, address: %s, duration: %d ", hexUserID, address, duration.Milliseconds())
	return nil
}

// UserHasAccount checks if provided account exist in the database for given userID
func (w *Services) UserHasAccount(hexUserID string, address string) (bool, error) {
	audit(w, "Checkinf if user has account: %s, address: %s", hexUserID, address)
	userIDBytes, err := common.GetUserIDbyte(hexUserID)
	if err != nil {
		w.Logger().Error(fmt.Errorf("error decoding string (%s), %w", hexUserID[2:], err).Error())
		return false, err
	}

	addressBytes, err := hex.DecodeString(address[2:]) // remove 0x prefix from address
	if err != nil {
		w.Logger().Error(fmt.Errorf("error decoding string (%s), %w", address[2:], err).Error())
		return false, err
	}

	// todo - this can be optimised and done in the database if we will have users with large number of accounts
	// get all the accounts for the selected user
	accounts, err := w.Storage.GetAccounts(userIDBytes)
	if err != nil {
		w.Logger().Error(fmt.Errorf("error getting accounts for user (%s), %w", hexUserID, err).Error())
		return false, err
	}

	// check if any of the account matches given account
	found := false
	for _, account := range accounts {
		if bytes.Equal(account.AccountAddress, addressBytes) {
			found = true
		}
	}
	return found, nil
}

// DeleteUser deletes user and accounts associated with user from the database for given userID
func (w *Services) DeleteUser(hexUserID string) error {
	audit(w, "Deleting user: %s", hexUserID)
	userIDBytes, err := common.GetUserIDbyte(hexUserID)
	if err != nil {
		w.Logger().Error(fmt.Errorf("error decoding string (%s), %w", hexUserID, err).Error())
		return err
	}

	err = w.Storage.DeleteUser(userIDBytes)
	if err != nil {
		w.Logger().Error(fmt.Errorf("error deleting user (%s), %w", hexUserID, err).Error())
		return err
	}

	return nil
}

func (w *Services) UserExists(hexUserID string) bool {
	audit(w, "Checking if user exists: %s", hexUserID)
	userIDBytes, err := common.GetUserIDbyte(hexUserID)
	if err != nil {
		w.Logger().Error(fmt.Errorf("error decoding string (%s), %w", hexUserID, err).Error())
		return false
	}

	// Check if user exists and don't log error if user doesn't exist, because we expect this to happen in case of
	// user revoking encryption token or using different testnet.
	// todo add a counter here in the future
	key, err := w.Storage.GetUserPrivateKey(userIDBytes)
	if err != nil {
		return false
	}

	return len(key) > 0
}

func (w *Services) Version() string {
	return w.version
}

func (w *Services) GetTenNodeHealthStatus() (bool, error) {
	return w.tenClient.Health()
}

func (w *Services) UnauthenticatedClient() (rpc.Client, error) {
	return rpc.NewNetworkClient(w.HostAddrHTTP)
}
