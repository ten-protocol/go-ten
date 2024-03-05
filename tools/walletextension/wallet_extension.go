package walletextension

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/ten-protocol/go-ten/tools/walletextension/cache"

	"github.com/ten-protocol/go-ten/tools/walletextension/accountmanager"

	"github.com/ten-protocol/go-ten/tools/walletextension/config"

	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/obsclient"

	"github.com/ten-protocol/go-ten/tools/walletextension/useraccountmanager"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/go-kit/kit/transport/http/jsonrpc"
	"github.com/ten-protocol/go-ten/go/common/stopcontrol"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"github.com/ten-protocol/go-ten/go/rpc"
	"github.com/ten-protocol/go-ten/tools/walletextension/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage"
	"github.com/ten-protocol/go-ten/tools/walletextension/userconn"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
)

// WalletExtension handles the management of viewing keys and the forwarding of Ethereum JSON-RPC requests.
type WalletExtension struct {
	hostAddrHTTP       string // The HTTP address on which the Ten host can be reached
	hostAddrWS         string // The WS address on which the Ten host can be reached
	userAccountManager *useraccountmanager.UserAccountManager
	unsignedVKs        map[gethcommon.Address]*viewingkey.ViewingKey // Map temporarily holding VKs that have been generated but not yet signed
	storage            storage.Storage
	logger             gethlog.Logger
	fileLogger         gethlog.Logger
	stopControl        *stopcontrol.StopControl
	version            string
	config             *config.Config
	tenClient          *obsclient.ObsClient
	cache              cache.Cache
}

func New(
	hostAddrHTTP string,
	hostAddrWS string,
	userAccountManager *useraccountmanager.UserAccountManager,
	storage storage.Storage,
	stopControl *stopcontrol.StopControl,
	version string,
	logger gethlog.Logger,
	config *config.Config,
) *WalletExtension {
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

	return &WalletExtension{
		hostAddrHTTP:       hostAddrHTTP,
		hostAddrWS:         hostAddrWS,
		userAccountManager: userAccountManager,
		unsignedVKs:        map[gethcommon.Address]*viewingkey.ViewingKey{},
		storage:            storage,
		logger:             logger,
		fileLogger:         newFileLogger,
		stopControl:        stopControl,
		version:            version,
		config:             config,
		tenClient:          newTenClient,
		cache:              newGatewayCache,
	}
}

// IsStopping returns whether the WE is stopping
func (w *WalletExtension) IsStopping() bool {
	return w.stopControl.IsStopping()
}

// Logger returns the WE set logger
func (w *WalletExtension) Logger() gethlog.Logger {
	return w.logger
}

// ProxyEthRequest proxys an incoming user request to the enclave
func (w *WalletExtension) ProxyEthRequest(request *common.RPCRequest, conn userconn.UserConn, hexUserID string) (map[string]interface{}, error) {
	response := map[string]interface{}{}
	// all responses must contain the request id. Both successful and unsuccessful.
	response[common.JSONKeyRPCVersion] = jsonrpc.Version
	response[common.JSONKeyID] = request.ID

	// start measuring time for request
	requestStartTime := time.Now()

	// Check if the request is in the cache
	isCacheable, key, ttl := cache.IsCacheable(request, hexUserID)

	// in case of cache hit return the response from the cache
	if isCacheable {
		if value, ok := w.cache.Get(key); ok {
			// do a shallow copy of the map to avoid concurrent map iteration and map write
			returnValue := make(map[string]interface{})
			for k, v := range value {
				returnValue[k] = v
			}

			requestEndTime := time.Now()
			duration := requestEndTime.Sub(requestStartTime)
			// adjust requestID
			returnValue[common.JSONKeyID] = request.ID
			w.fileLogger.Info(fmt.Sprintf("Request method: %s, request params: %s, encryptionToken of sender: %s, response: %s, duration: %d ", request.Method, request.Params, hexUserID, returnValue, duration.Milliseconds()))
			return returnValue, nil
		}
	}

	// proxyRequest will find the correct client to proxy the request (or try them all if appropriate)
	var rpcResp interface{}

	// wallet extension can override the GetStorageAt to retrieve the current userID
	if request.Method == rpc.GetStorageAt {
		if interceptedResponse := w.getStorageAtInterceptor(request, hexUserID); interceptedResponse != nil {
			w.logger.Info("interception successful for getStorageAt, returning userID response")
			requestEndTime := time.Now()
			duration := requestEndTime.Sub(requestStartTime)
			w.fileLogger.Info(fmt.Sprintf("Request method: %s, request params: %s, encryptionToken of sender: %s, response: %s, duration: %d ", request.Method, request.Params, hexUserID, interceptedResponse, duration.Milliseconds()))
			return interceptedResponse, nil
		}
	}

	// check if user is sending a new transaction and if we should store it in the database for debugging purposes
	if request.Method == rpc.SendRawTransaction && w.config.StoreIncomingTxs {
		userIDBytes, err := common.GetUserIDbyte(hexUserID)
		if err != nil {
			w.Logger().Error(fmt.Errorf("error decoding string (%s), %w", hexUserID[2:], err).Error())
			return nil, errors.New("error decoding userID. It should be in hex format")
		}
		err = w.storage.StoreTransaction(request.Params[0].(string), userIDBytes)
		if err != nil {
			w.Logger().Error(fmt.Errorf("error storing transaction in the database: %w", err).Error())
			return nil, err
		}
	}

	// get account manager for current user (if there is no users in the query parameters - use defaultUser for WE endpoints)
	selectedAccountManager, err := w.userAccountManager.GetUserAccountManager(hexUserID)
	if err != nil {
		w.Logger().Error(fmt.Errorf("error getting accountManager for user (%s), %w", hexUserID, err).Error())
		return nil, err
	}

	err = selectedAccountManager.ProxyRequest(request, &rpcResp, conn)
	if err != nil {
		if errors.Is(err, rpc.ErrNilResponse) {
			// if err was for a nil response then we will return an RPC result of null to the caller (this is a valid "not-found" response for some methods)
			response[common.JSONKeyResult] = nil
			requestEndTime := time.Now()
			duration := requestEndTime.Sub(requestStartTime)
			w.fileLogger.Info(fmt.Sprintf("Request method: %s, request params: %s, encryptionToken of sender: %s, response: %s, duration: %d ", request.Method, request.Params, hexUserID, response, duration.Milliseconds()))
			return response, nil
		}
		return nil, err
	}

	response[common.JSONKeyResult] = rpcResp

	// todo (@ziga) - fix this upstream on the decode
	// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-658.md
	adjustStateRoot(rpcResp, response)
	requestEndTime := time.Now()
	duration := requestEndTime.Sub(requestStartTime)
	w.fileLogger.Info(fmt.Sprintf("Request method: %s, request params: %s, encryptionToken of sender: %s, response: %s, duration: %d ", request.Method, request.Params, hexUserID, response, duration.Milliseconds()))

	// if the request is cacheable, store the response in the cache
	if isCacheable {
		w.cache.Set(key, response, ttl)
	}

	return response, nil
}

// GenerateViewingKey generates the user viewing key and waits for signature
func (w *WalletExtension) GenerateViewingKey(addr gethcommon.Address) (string, error) {
	w.fileLogger.Info(fmt.Sprintf("Requested to generate viewing key for address(old way): %s", addr.Hex()))
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
func (w *WalletExtension) SubmitViewingKey(address gethcommon.Address, signature []byte) error {
	w.fileLogger.Info(fmt.Sprintf("Requested to submit a viewing key (old way): %s", address.Hex()))
	vk, found := w.unsignedVKs[address]
	if !found {
		return fmt.Errorf(fmt.Sprintf("no viewing key found to sign for acc=%s, please call %s to generate key before sending signature", address, common.PathGenerateViewingKey))
	}

	// We transform the V from 27/28 to 0/1. This same change is made in Geth internals, for legacy reasons to be able
	// to recover the address: https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L452-L459
	signature[64] -= 27

	vk.SignatureWithAccountKey = signature

	err := w.storage.AddUser([]byte(common.DefaultUser), crypto.FromECDSA(vk.PrivateKey.ExportECDSA()))
	if err != nil {
		return fmt.Errorf("error saving user: %s", common.DefaultUser)
	}
	// create an encrypted RPC client with the signed VK and register it with the enclave
	// todo (@ziga) - Create the clients lazily, to reduce connections to the host.
	client, err := rpc.NewEncNetworkClient(w.hostAddrHTTP, vk, w.logger)
	if err != nil {
		return fmt.Errorf("failed to create encrypted RPC client for account %s - %w", address, err)
	}
	defaultAccountManager, err := w.userAccountManager.GetUserAccountManager(hex.EncodeToString([]byte(common.DefaultUser)))
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("error getting default user account manager: %s", err))
	}

	defaultAccountManager.AddClient(address, client)

	err = w.storage.AddAccount([]byte(common.DefaultUser), vk.Account.Bytes(), vk.SignatureWithAccountKey)
	if err != nil {
		return fmt.Errorf("error saving account %s for user %s", vk.Account.Hex(), common.DefaultUser)
	}

	if err != nil {
		return fmt.Errorf("error saving viewing key to the database: %w", err)
	}

	// finally we remove the VK from the pending 'unsigned VKs' map now the client has been created
	delete(w.unsignedVKs, address)

	return nil
}

// GenerateAndStoreNewUser generates new key-pair and userID, stores it in the database and returns hex encoded userID and error
func (w *WalletExtension) GenerateAndStoreNewUser() (string, error) {
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
	err = w.storage.AddUser(userID, crypto.FromECDSA(viewingPrivateKeyEcies.ExportECDSA()))
	if err != nil {
		w.Logger().Error(fmt.Sprintf("failed to save user to the database: %s", err))
		return "", err
	}

	hexUserID := hex.EncodeToString(userID)

	w.userAccountManager.AddAndReturnAccountManager(hexUserID)
	requestEndTime := time.Now()
	duration := requestEndTime.Sub(requestStartTime)
	w.fileLogger.Info(fmt.Sprintf("Storing new userID: %s, duration: %d ", hexUserID, duration.Milliseconds()))
	return hexUserID, nil
}

// AddAddressToUser checks if a message is in correct format and if signature is valid. If all checks pass we save address and signature against userID
func (w *WalletExtension) AddAddressToUser(hexUserID string, address string, signature []byte) error {
	requestStartTime := time.Now()
	addressFromMessage := gethcommon.HexToAddress(address)
	// check if a message was signed by the correct address and if the signature is valid
	sigAddrs, err := viewingkey.CheckSignature(hexUserID, signature, int64(w.config.TenChainID), address)
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
	err = w.storage.AddAccount(userIDBytes, addressFromMessage.Bytes(), signature)
	if err != nil {
		w.Logger().Error(fmt.Errorf("error while storing account (%s) for user (%s): %w", addressFromMessage.Hex(), hexUserID, err).Error())
		return err
	}

	// Get account manager for current userID (and create it if it doesn't exist)
	privateKeyBytes, err := w.storage.GetUserPrivateKey(userIDBytes)
	if err != nil {
		w.Logger().Error(fmt.Errorf("error getting private key for user: (%s), %w", hexUserID, err).Error())
	}

	accManager := w.userAccountManager.AddAndReturnAccountManager(hexUserID)

	encClient, err := common.CreateEncClient(w.hostAddrHTTP, addressFromMessage.Bytes(), privateKeyBytes, signature, w.Logger())
	if err != nil {
		w.Logger().Error(fmt.Errorf("error creating encrypted client for user: (%s), %w", hexUserID, err).Error())
		return fmt.Errorf("error creating encrypted client for user: (%s), %w", hexUserID, err)
	}

	accManager.AddClient(addressFromMessage, encClient)
	requestEndTime := time.Now()
	duration := requestEndTime.Sub(requestStartTime)
	w.fileLogger.Info(fmt.Sprintf("Storing new address for user: %s, address: %s, duration: %d ", hexUserID, address, duration.Milliseconds()))
	return nil
}

// UserHasAccount checks if provided account exist in the database for given userID
func (w *WalletExtension) UserHasAccount(hexUserID string, address string) (bool, error) {
	w.fileLogger.Info(fmt.Sprintf("Checkinf if user has account: %s, address: %s", hexUserID, address))
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
	accounts, err := w.storage.GetAccounts(userIDBytes)
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
func (w *WalletExtension) DeleteUser(hexUserID string) error {
	w.fileLogger.Info(fmt.Sprintf("Deleting user: %s", hexUserID))
	userIDBytes, err := common.GetUserIDbyte(hexUserID)
	if err != nil {
		w.Logger().Error(fmt.Errorf("error decoding string (%s), %w", hexUserID, err).Error())
		return err
	}

	err = w.storage.DeleteUser(userIDBytes)
	if err != nil {
		w.Logger().Error(fmt.Errorf("error deleting user (%s), %w", hexUserID, err).Error())
		return err
	}

	// Delete UserAccountManager for user that revoked userID
	err = w.userAccountManager.DeleteUserAccountManager(hexUserID)
	if err != nil {
		w.Logger().Error(fmt.Errorf("error deleting UserAccointManager for user (%s), %w", hexUserID, err).Error())
	}

	return nil
}

func (w *WalletExtension) UserExists(hexUserID string) bool {
	w.fileLogger.Info(fmt.Sprintf("Checking if user exists: %s", hexUserID))
	userIDBytes, err := common.GetUserIDbyte(hexUserID)
	if err != nil {
		w.Logger().Error(fmt.Errorf("error decoding string (%s), %w", hexUserID, err).Error())
		return false
	}

	// Check if user exists and don't log error if user doesn't exist, because we expect this to happen in case of
	// user revoking encryption token or using different testnet.
	// todo add a counter here in the future
	key, err := w.storage.GetUserPrivateKey(userIDBytes)
	if err != nil {
		return false
	}

	return len(key) > 0
}

func adjustStateRoot(rpcResp interface{}, respMap map[string]interface{}) {
	if resultMap, ok := rpcResp.(map[string]interface{}); ok {
		if val, foundRoot := resultMap[common.JSONKeyRoot]; foundRoot {
			if val == "0x" {
				respMap[common.JSONKeyResult].(map[string]interface{})[common.JSONKeyRoot] = nil
			}
		}
	}
}

// getStorageAtInterceptor checks if the parameters for getStorageAt are set to values that require interception
// and return response or nil if the gateway should forward the request to the node.
func (w *WalletExtension) getStorageAtInterceptor(request *common.RPCRequest, hexUserID string) map[string]interface{} {
	// check if parameters are correct, and we can intercept a request, otherwise return nil
	if w.checkParametersForInterceptedGetStorageAt(request.Params) {
		// check if userID in the parameters is also in our database
		userID, err := common.GetUserIDbyte(hexUserID)
		if err != nil {
			w.logger.Warn("GetStorageAt called with appropriate parameters to return userID, but not found in the database: ", "userId", hexUserID)
			return nil
		}

		// check if we have default user (we don't want to send userID of it out)
		if hexUserID == hex.EncodeToString([]byte(common.DefaultUser)) {
			response := map[string]interface{}{}
			response[common.JSONKeyRPCVersion] = jsonrpc.Version
			response[common.JSONKeyID] = request.ID
			response[common.JSONKeyResult] = fmt.Sprintf(accountmanager.ErrNoViewingKey, "eth_getStorageAt")
			return response
		}

		_, err = w.storage.GetUserPrivateKey(userID)
		if err != nil {
			w.logger.Info("Trying to get userID, but it is not present in our database: ", log.ErrKey, err)
			return nil
		}
		response := map[string]interface{}{}
		response[common.JSONKeyRPCVersion] = jsonrpc.Version
		response[common.JSONKeyID] = request.ID
		response[common.JSONKeyResult] = hexUserID
		return response
	}
	w.logger.Info(fmt.Sprintf("parameters used in the request do not match requited parameters for interception: %s", request.Params))

	return nil
}

// checkParametersForInterceptedGetStorageAt checks
// if parameters for getStorageAt are in the correct format to intercept the function
func (w *WalletExtension) checkParametersForInterceptedGetStorageAt(params []interface{}) bool {
	if len(params) != 3 {
		w.logger.Info(fmt.Sprintf("getStorageAt expects 3 parameters, but %d received", len(params)))
		return false
	}

	if methodName, ok := params[0].(string); ok {
		return methodName == common.GetStorageAtUserIDRequestMethodName
	}
	return false
}

func (w *WalletExtension) Version() string {
	return w.version
}

func (w *WalletExtension) GetTenNodeHealthStatus() (bool, error) {
	return w.tenClient.Health()
}
