package walletextension

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/obscuronet/go-obscuro/tools/walletextension/useraccountmanager"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/go-kit/kit/transport/http/jsonrpc"
	"github.com/obscuronet/go-obscuro/go/common/stopcontrol"
	"github.com/obscuronet/go-obscuro/go/common/viewingkey"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/tools/walletextension/accountmanager"
	"github.com/obscuronet/go-obscuro/tools/walletextension/common"
	"github.com/obscuronet/go-obscuro/tools/walletextension/storage"
	"github.com/obscuronet/go-obscuro/tools/walletextension/userconn"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
)

var ErrSubscribeFailHTTP = fmt.Sprintf("received an %s request but the connection does not support subscriptions", rpc.Subscribe)

// WalletExtension handles the management of viewing keys and the forwarding of Ethereum JSON-RPC requests.
type WalletExtension struct {
	hostAddr           string // The address on which the Obscuro host can be reached.
	userAccountManager *useraccountmanager.UserAccountManager
	unsignedVKs        map[gethcommon.Address]*viewingkey.ViewingKey // Map temporarily holding VKs that have been generated but not yet signed
	storage            storage.Storage
	logger             gethlog.Logger
	stopControl        *stopcontrol.StopControl
}

func New(
	hostAddr string,
	userAccountManager *useraccountmanager.UserAccountManager,
	storage storage.Storage,
	stopControl *stopcontrol.StopControl,
	logger gethlog.Logger,
) *WalletExtension {
	return &WalletExtension{
		hostAddr:           hostAddr,
		userAccountManager: userAccountManager,
		unsignedVKs:        map[gethcommon.Address]*viewingkey.ViewingKey{},
		storage:            storage,
		logger:             logger,
		stopControl:        stopControl,
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
func (w *WalletExtension) ProxyEthRequest(request *accountmanager.RPCRequest, conn userconn.UserConn, hexUserID string) (map[string]interface{}, error) {
	response := map[string]interface{}{}
	// all responses must contain the request id. Both successful and unsuccessful.
	response[common.JSONKeyRPCVersion] = jsonrpc.Version
	response[common.JSONKeyID] = request.ID

	// proxyRequest will find the correct client to proxy the request (or try them all if appropriate)
	var rpcResp interface{}

	// wallet extension can override the GetStorageAt to retrieve the current userID
	if request.Method == rpc.GetStorageAt {
		if interceptedResponse := w.getStorageAtInterceptor(request, hexUserID); interceptedResponse != nil {
			w.logger.Info("interception successful for getStorageAt, returning userID response")
			return interceptedResponse, nil
		}
	}

	// get account manager for current user (if there is no users in the query parameters - use defaultUser for WE endpoints)
	selectedAccountManager, err := w.userAccountManager.GetUserAccountManager(hexUserID)
	if err != nil {
		w.Logger().Error(fmt.Errorf("error getting accountManager for user (%s), %w", hexUserID, err).Error())
		return nil, err
	}

	err = selectedAccountManager.ProxyRequest(request, &rpcResp, conn)

	if err != nil && !errors.Is(err, rpc.ErrNilResponse) {
		response = common.CraftErrorResponse(err)
	} else if errors.Is(err, rpc.ErrNilResponse) {
		// if err was for a nil response then we will return an RPC result of null to the caller (this is a valid "not-found" response for some methods)
		response[common.JSONKeyResult] = nil
	} else {
		response[common.JSONKeyResult] = rpcResp

		// todo (@ziga) - fix this upstream on the decode
		// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-658.md
		adjustStateRoot(rpcResp, response)
	}
	return response, nil
}

// GenerateViewingKey generates the user viewing key and waits for signature
func (w *WalletExtension) GenerateViewingKey(addr gethcommon.Address) (string, error) {
	viewingKeyPrivate, err := crypto.GenerateKey()
	if err != nil {
		return "", fmt.Errorf("unable to generate a new keypair - %w", err)
	}

	viewingPublicKeyBytes := crypto.CompressPubkey(&viewingKeyPrivate.PublicKey)
	viewingPrivateKeyEcies := ecies.ImportECDSA(viewingKeyPrivate)

	w.unsignedVKs[addr] = &viewingkey.ViewingKey{
		Account:    &addr,
		PrivateKey: viewingPrivateKeyEcies,
		PublicKey:  viewingPublicKeyBytes,
		Signature:  nil, // we await a signature from the user before we can set up the EncRPCClient
	}

	// compress the viewing key and convert it to hex string ( this is what Metamask signs)
	viewingKeyBytes := crypto.CompressPubkey(&viewingKeyPrivate.PublicKey)
	return hex.EncodeToString(viewingKeyBytes), nil
}

// SubmitViewingKey checks the signed viewing key and stores it
func (w *WalletExtension) SubmitViewingKey(address gethcommon.Address, signature []byte) error {
	vk, found := w.unsignedVKs[address]
	if !found {
		return fmt.Errorf(fmt.Sprintf("no viewing key found to sign for acc=%s, please call %s to generate key before sending signature", address, common.PathGenerateViewingKey))
	}

	// We transform the V from 27/28 to 0/1. This same change is made in Geth internals, for legacy reasons to be able
	// to recover the address: https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L452-L459
	signature[64] -= 27

	vk.Signature = signature
	// create an encrypted RPC client with the signed VK and register it with the enclave
	// todo (@ziga) - Create the clients lazily, to reduce connections to the host.
	client, err := rpc.NewEncNetworkClient(w.hostAddr, vk, w.logger)
	if err != nil {
		return fmt.Errorf("failed to create encrypted RPC client for account %s - %w", address, err)
	}
	defaultAccountManager, err := w.userAccountManager.GetUserAccountManager(hex.EncodeToString([]byte(common.DefaultUser)))
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("error getting default user account manager: %s", err))
	}

	defaultAccountManager.AddClient(address, client)

	err = w.storage.AddUser([]byte(common.DefaultUser), crypto.FromECDSA(vk.PrivateKey.ExportECDSA()))
	if err != nil {
		return fmt.Errorf("error saving user: %s", common.DefaultUser)
	}

	err = w.storage.AddAccount([]byte(common.DefaultUser), vk.Account.Bytes(), vk.Signature)
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
	// generate new key-pair
	viewingKeyPrivate, err := crypto.GenerateKey()
	viewingPrivateKeyEcies := ecies.ImportECDSA(viewingKeyPrivate)
	if err != nil {
		w.Logger().Error(fmt.Sprintf("could not generate new keypair: %s", err))
		return "", err
	}

	// create UserID and store it in the database with the private key
	userID := common.CalculateUserID(common.PrivateKeyToCompressedPubKey(viewingPrivateKeyEcies))
	err = w.storage.AddUser(userID, crypto.FromECDSA(viewingPrivateKeyEcies.ExportECDSA()))
	if err != nil {
		w.Logger().Error(fmt.Sprintf("failed to save user to the database: %s", err))
		return "", err
	}

	hexUserID := hex.EncodeToString(userID)

	w.userAccountManager.AddAndReturnAccountManager(hexUserID)

	return hexUserID, nil
}

// AddAddressToUser checks if message is in correct format and if signature is valid. If all checks pass we save address and signature against userID
func (w *WalletExtension) AddAddressToUser(hexUserID string, message string, signature []byte) error {
	// parse the message to get userID and account address
	messageUserID, messageAddressHex, err := common.GetUserIDAndAddressFromMessage(message)
	if err != nil {
		w.Logger().Error(fmt.Errorf("submitted message (%s) is not in the correct format", message).Error())
		return err
	}

	// check if userID corresponds to the one in the message and check if the length of hex encoded userID is correct
	if hexUserID != messageUserID || len(messageUserID) != common.MessageUserIDLen {
		w.Logger().Error(fmt.Errorf("submitted message (%s) is not in the correct format", message).Error())
		return errors.New("userID from message does not match userID from request")
	}

	addressFromMessage := gethcommon.HexToAddress(messageAddressHex)

	// check if message was signed by the correct address and if signature is valid
	valid, err := verifySignature(message, signature, addressFromMessage)
	if !valid && err != nil {
		w.Logger().Error(fmt.Errorf("error: signature is not valid: %s", string(signature)).Error())
		return err
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

	// Get account manager for current userID (and create it if it doesn't exist) accManager := w.userAccountManager.AddAndReturnAccountManager(messageUserID)
	privateKeyBytes, err := w.storage.GetUserPrivateKey(userIDBytes)
	if err != nil {
		w.Logger().Error(fmt.Errorf("error getting private key for user: (%s), %w", hexUserID, err).Error())
	}

	accManager := w.userAccountManager.AddAndReturnAccountManager(hexUserID)

	encClient, err := common.CreateEncClient(w.hostAddr, addressFromMessage.Bytes(), privateKeyBytes, signature)
	if err != nil {
		w.Logger().Error(fmt.Errorf("error creating encrypted client for user: (%s), %w", hexUserID, err).Error())
	}

	accManager.AddClient(addressFromMessage, encClient)

	return nil
}

// UserHasAccount checks if provided account exist in the database for given userID
func (w *WalletExtension) UserHasAccount(hexUserID string, address string) (bool, error) {
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
	userIDBytes, err := common.GetUserIDbyte(hexUserID)
	if err != nil {
		w.Logger().Error(fmt.Errorf("error decoding string (%s), %w", hexUserID, err).Error())
		return false
	}

	key, err := w.storage.GetUserPrivateKey(userIDBytes)
	if err != nil {
		w.Logger().Error(fmt.Errorf("error getting user's private key (%s), %w", hexUserID, err).Error())
		return false
	}

	return len(key) > 0
}

// verifySignature checks if a message was signed by the correct address and if signature is valid
func verifySignature(message string, signature []byte, address gethcommon.Address) (bool, error) {
	// prefix the message like in the personal_sign method
	prefixedMessage := fmt.Sprintf(common.PersonalSignMessagePrefix, len(message), message)
	messageHash := crypto.Keccak256([]byte(prefixedMessage))

	// check if the signature length is correct
	if len(signature) != common.SignatureLen {
		return false, errors.New("incorrect signature length")
	}

	// We transform the V from 27/28 to 0/1. This same change is made in Geth internals, for legacy reasons to be able
	// to recover the address: https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L452-L459
	signature[64] -= 27

	addressFromSignature, pubKeyFromSignature, err := common.GetAddressAndPubKeyFromSignature(messageHash, signature)
	if err != nil {
		return false, err
	}

	if !bytes.Equal(addressFromSignature.Bytes(), address.Bytes()) {
		return false, errors.New("address from signature not the same as expected")
	}

	// Split signature into r, s
	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:64])

	// Verify the signature
	isValid := ecdsa.Verify(pubKeyFromSignature, messageHash, r, s)

	if !isValid {
		return false, errors.New("signature is not valid")
	}

	return true, nil
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
func (w *WalletExtension) getStorageAtInterceptor(request *accountmanager.RPCRequest, hexUserID string) map[string]interface{} {
	// check if parameters are correct, and we can intercept a request, otherwise return nil
	if w.checkParametersForInterceptedGetStorageAt(request.Params) {
		// check if userID in the parameters is also in our database
		userID, err := common.GetUserIDbyte(hexUserID)
		if err != nil {
			w.logger.Warn("GetStorageAt called with appropriate parameters to return userID, but not found in the database: ", "userId", hexUserID)
			return nil
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
