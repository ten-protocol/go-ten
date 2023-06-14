package api

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"

	"github.com/obscuronet/go-obscuro/go/common/httputil"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/tools/walletextension"
	"github.com/obscuronet/go-obscuro/tools/walletextension/common"
	"github.com/obscuronet/go-obscuro/tools/walletextension/userconn"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// Route defines the path plus handler for a given path
type Route struct {
	Name string
	Func func(resp http.ResponseWriter, req *http.Request)
}

// NewHTTPRoutes returns the http specific routes
func NewHTTPRoutes(walletExt *walletextension.WalletExtension) []Route {
	return []Route{
		{
			Name: common.PathRoot,
			Func: httpHandler(walletExt, ethRequestHandler),
		},
		{
			Name: common.PathReady,
			Func: httpHandler(walletExt, readyRequestHandler),
		},
		{
			Name: common.PathGenerateViewingKey,
			Func: httpHandler(walletExt, generateViewingKeyRequestHandler),
		},

		{
			Name: common.PathSubmitViewingKey,
			Func: httpHandler(walletExt, submitViewingKeyRequestHandler),
		},
		{
			Name: common.PathJoin,
			Func: httpHandler(walletExt, joinRequestHandler),
		},
		{
			Name: common.PathAuthenticate,
			Func: httpHandler(walletExt, authenticateRequestHandler),
		},
		{
			Name: common.PathQuery,
			Func: httpHandler(walletExt, queryRequestHandler),
		},
		{
			Name: common.PathRevoke,
			Func: httpHandler(walletExt, revokeRequestHandler),
		},
	}
}

func httpHandler(
	walletExt *walletextension.WalletExtension,
	fun func(walletExt *walletextension.WalletExtension, conn userconn.UserConn),
) func(resp http.ResponseWriter, req *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		httpRequestHandler(walletExt, resp, req, fun)
	}
}

// Overall request handler for http requests
func httpRequestHandler(walletExt *walletextension.WalletExtension, resp http.ResponseWriter, req *http.Request, fun func(walletExt *walletextension.WalletExtension, conn userconn.UserConn)) {
	if walletExt.IsStopping() {
		return
	}
	if httputil.EnableCORS(resp, req) {
		return
	}
	userConn := userconn.NewUserConnHTTP(resp, req, walletExt.Logger())
	fun(walletExt, userConn)
}

// NewWSRoutes returns the WS specific routes
func NewWSRoutes(walletExt *walletextension.WalletExtension) []Route {
	return []Route{
		{
			Name: common.PathRoot,
			Func: wsHandler(walletExt, ethRequestHandler),
		},
		{
			Name: common.PathReady,
			Func: wsHandler(walletExt, readyRequestHandler),
		},
		{
			Name: common.PathGenerateViewingKey,
			Func: wsHandler(walletExt, generateViewingKeyRequestHandler),
		},

		{
			Name: common.PathSubmitViewingKey,
			Func: wsHandler(walletExt, submitViewingKeyRequestHandler),
		},
	}
}

func wsHandler(
	walletExt *walletextension.WalletExtension,
	fun func(walletExt *walletextension.WalletExtension, conn userconn.UserConn),
) func(resp http.ResponseWriter, req *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		wsRequestHandler(walletExt, resp, req, fun)
	}
}

// Overall request handler for WS requests
func wsRequestHandler(walletExt *walletextension.WalletExtension, resp http.ResponseWriter, req *http.Request, fun func(walletExt *walletextension.WalletExtension, conn userconn.UserConn)) {
	if walletExt.IsStopping() {
		return
	}

	userConn, err := userconn.NewUserConnWS(resp, req, walletExt.Logger())
	if err != nil {
		return
	}
	// We handle requests in a loop until the connection is closed on the client side.
	for !userConn.IsClosed() {
		fun(walletExt, userConn)
	}
}

// ethRequestHandler parses the user eth request, passes it on to the WE to proxy it and processes the response
func ethRequestHandler(walletExt *walletextension.WalletExtension, conn userconn.UserConn) {
	body, err := conn.ReadRequest()
	if err != nil {
		conn.HandleError("Error: bad request")
		walletExt.Logger().Error(fmt.Errorf("error reading request: %w", err).Error())
		return
	}

	request, err := parseRequest(body)
	if err != nil {
		conn.HandleError(err.Error())
		return
	}
	walletExt.Logger().Debug("REQUEST", "method", request.Method, "body", string(body))

	if request.Method == rpc.Subscribe && !conn.SupportsSubscriptions() {
		conn.HandleError(common.ErrSubscribeFailHTTP)
		return
	}

	// todo (@pedro) remove this conn dependency
	response, err := walletExt.ProxyEthRequest(request, conn)
	if err != nil {
		walletExt.Logger().Error("error while proxying request", log.ErrKey, err)
		response = common.CraftErrorResponse(err)
	}

	rpcResponse, err := json.Marshal(response)
	if err != nil {
		conn.HandleError(fmt.Sprintf("failed to remarshal RPC response to return to caller: %s", err))
		return
	}
	walletExt.Logger().Info(fmt.Sprintf("Forwarding %s response from Obscuro node: %s", request.Method, rpcResponse))

	err = conn.WriteResponse(rpcResponse)
	if err != nil {
		return
	}
}

// readyRequestHandler is used to check whether the server is ready
func readyRequestHandler(_ *walletextension.WalletExtension, _ userconn.UserConn) {}

// generateViewingKeyRequestHandler parses the gen vk request
func generateViewingKeyRequestHandler(walletExt *walletextension.WalletExtension, conn userconn.UserConn) {
	body, err := conn.ReadRequest()
	if err != nil {
		conn.HandleError("Error: bad request")
		walletExt.Logger().Error(fmt.Errorf("error reading request: %w", err).Error())
		return
	}

	var reqJSONMap map[string]string
	err = json.Unmarshal(body, &reqJSONMap)
	if err != nil {
		conn.HandleError(fmt.Sprintf("could not unmarshal address request - %s", err))
		return
	}

	address := gethcommon.HexToAddress(reqJSONMap[common.JSONKeyAddress])

	pubViewingKey, err := walletExt.GenerateViewingKey(address)
	if err != nil {
		conn.HandleError(fmt.Sprintf("unable to generate vieweing key: %s", err))
		return
	}

	err = conn.WriteResponse([]byte(pubViewingKey))
	if err != nil {
		return
	}
}

// submitViewingKeyRequestHandler submits the viewing key and signed bytes to the WE
func submitViewingKeyRequestHandler(walletExt *walletextension.WalletExtension, userConn userconn.UserConn) {
	body, err := userConn.ReadRequest()
	if err != nil {
		userConn.HandleError("Error: bad request")
		walletExt.Logger().Error(fmt.Errorf("error reading request: %w", err).Error())
		return
	}

	var reqJSONMap map[string]string
	err = json.Unmarshal(body, &reqJSONMap)
	if err != nil {
		userConn.HandleError(fmt.Sprintf("could not unmarshal address and signature from client to JSON: %s", err))
		return
	}
	accAddress := gethcommon.HexToAddress(reqJSONMap[common.JSONKeyAddress])

	signature, err := hex.DecodeString(reqJSONMap[common.JSONKeySignature][2:])
	if err != nil {
		userConn.HandleError(fmt.Sprintf("could not decode signature from client to hex: %s", err))
		return
	}

	err = walletExt.SubmitViewingKey(accAddress, signature)
	if err != nil {
		userConn.HandleError(fmt.Sprintf("could not submit viewing key - %s", err))
		return
	}

	err = userConn.WriteResponse([]byte(common.SuccessMsg))
	if err != nil {
		walletExt.Logger().Error(fmt.Errorf("error writing success response, %w", err).Error())
		return
	}
}

// This function handles request to /join endpoint.
// It generates new key-pair and userID, stores it in the database and returns userID back to the user.
func joinRequestHandler(walletExt *walletextension.WalletExtension, userConn userconn.UserConn) {
	// todo (@ziga) add protection against DDOS attacks
	_, err := userConn.ReadRequest()
	if err != nil {
		userConn.HandleError("Error: bad request")
		walletExt.Logger().Error(fmt.Errorf("error reading request: %w", err).Error())
		return
	}

	// generate new key-pair
	viewingKeyPrivate, err := crypto.GenerateKey()
	viewingPrivateKeyEcies := ecies.ImportECDSA(viewingKeyPrivate)
	if err != nil {
		userConn.HandleError("Internal error")
		walletExt.Logger().Error(fmt.Sprintf("could not generate new keypair: %s", err))
		return
	}

	// create UserID and store it in the database with the private key
	userID := calculateUserID(viewingKeyPrivate)
	err = walletExt.Storage.AddUser(userID, crypto.FromECDSA(viewingPrivateKeyEcies.ExportECDSA()))
	if err != nil {
		userConn.HandleError("Internal error")
		walletExt.Logger().Error(fmt.Sprintf("failed to save user to the database: %s", err))
		return
	}

	// write hex encoded userID in the response
	err = userConn.WriteResponse([]byte(hex.EncodeToString(userID)))

	if err != nil {
		walletExt.Logger().Error(fmt.Errorf("error writing success response, %w", err).Error())
		return
	}
}

// This function handles request to /authenticate endpoint.
// In the request we receive message, signature and address in JSON as request body and userID and address as query parameters
// We then check if message is in correct format and if signature is valid. If all checks pass we save address and signature against userID
func authenticateRequestHandler(walletExt *walletextension.WalletExtension, userConn userconn.UserConn) {
	// read the request
	body, err := userConn.ReadRequest()
	if err != nil {
		userConn.HandleError("Error: bad request")
		walletExt.Logger().Error(fmt.Errorf("error reading request: %w", err).Error())
		return
	}

	// get the text that was signed and signature
	var reqJSONMap map[string]string
	err = json.Unmarshal(body, &reqJSONMap)
	if err != nil {
		userConn.HandleError("Internal error")
		walletExt.Logger().Error(fmt.Errorf("error unmarshaling request to authentcate: %w", err).Error())
		return
	}

	// get signature from the request and remove leading two bytes (0x)
	signature, err := hex.DecodeString(reqJSONMap[common.JSONKeySignature][2:])
	if err != nil {
		userConn.HandleError("Error: unable to decode signature")
		walletExt.Logger().Error(fmt.Errorf("could not find or decode signature from client to hex: %w", err).Error())
		return
	}

	// get message from the request
	message, ok := reqJSONMap[common.JSONKeyMessage]
	if !ok || message == "" {
		userConn.HandleError("Error: unable to read message field from the request")
		walletExt.Logger().Error(fmt.Errorf("could not find message in the request: %w", err).Error())
		return
	}

	// read userID from query params
	hexUserID, err := getQueryParameter(userConn.ReadRequestParams(), common.UserQueryParameter)
	if err != nil {
		userConn.HandleError("Malformed query: 'u' required - representing userID")
		walletExt.Logger().Error(fmt.Errorf("user not found in the query params: %w", err).Error())
		return
	}

	// parse the message to get userID and account address
	messageUserID, messageAddressHex, err := getUserIDAndAddressFromMessage(message)
	if err != nil {
		userConn.HandleError("Internal error")
		walletExt.Logger().Error(fmt.Errorf("submitted message (%s) is not in the correct format", message).Error())
		return
	}

	// check if userID corresponds to the one in the message and check if the length of hex encoded userID is correct
	// todo: do we need userID in query param, because we get it already from the message
	if hexUserID != messageUserID || len(messageUserID) != common.MessageUserIDLen {
		userConn.HandleError(fmt.Sprintf("User in submitted message (%s) does not match user provided in the request (%s) of is wrong size", messageUserID, hexUserID))
		return
	}

	// Check if the signature is valid
	// prefix the message like in the personal_sign method
	prefixedMessage := fmt.Sprintf(common.PersonalSignMessagePrefix, len(message), message)
	messageHash := crypto.Keccak256([]byte(prefixedMessage))

	// check if the signature length is correct
	if len(signature) != common.SignatureLen {
		userConn.HandleError("Error: signature must be 64 bytes long")
		walletExt.Logger().Error(fmt.Errorf("signature must be 64 bytes long, but %d bytes long signature received", len(signature)).Error())
		return
	}

	// We transform the V from 27/28 to 0/1. This same change is made in Geth internals, for legacy reasons to be able
	// to recover the address: https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L452-L459
	signature[64] -= 27

	// get addresses from signature and message and compare if they are the same
	addressFromSignature, err := getAddressFromSignature(messageHash, signature)
	if err != nil {
		userConn.HandleError("Internal error")
		walletExt.Logger().Error(fmt.Errorf("error getting address from signature: %w", err).Error())
	}
	addressFromMessage := gethcommon.HexToAddress(messageAddressHex)

	// verify that message was signed by the same address as in the message
	if addressFromSignature != addressFromMessage {
		userConn.HandleError("Message not signed by the same address as provided in message")
		walletExt.Logger().Error(fmt.Errorf("address from signature (%s) is not the same as address from message (%s)", addressFromSignature, addressFromSignature).Error())
		return
	}

	// register the account for that viewing key
	userIDBytes, err := getUserIDbyte(hexUserID)
	if err != nil {
		userConn.HandleError("Error decoding userID. It should be in hex format")
		walletExt.Logger().Error(fmt.Errorf("error decoding string (%s), %w", hexUserID[2:], err).Error())
		return
	}
	err = walletExt.Storage.AddAccount(userIDBytes, addressFromMessage.Bytes(), signature)
	if err != nil {
		userConn.HandleError("Internal error")
		walletExt.Logger().Error(fmt.Errorf("error while storing account (%s) for user (%s): %w", addressFromMessage.Hex(), hexUserID, err).Error())
		return
	}

	err = userConn.WriteResponse([]byte("success!"))
	if err != nil {
		walletExt.Logger().Error(fmt.Errorf("error writing success response, %w", err).Error())
		return
	}
}

// This function handles request to /query endpoint.
// In the query parameters address and userID are required. We check if provided address is registered for given userID
// and return true/false in json response
func queryRequestHandler(walletExt *walletextension.WalletExtension, userConn userconn.UserConn) {
	// read the request
	_, err := userConn.ReadRequest()
	if err != nil {
		userConn.HandleError("Error: bad request")
		walletExt.Logger().Error(fmt.Errorf("error reading request: %w", err).Error())
		return
	}

	hexUserID, err := getQueryParameter(userConn.ReadRequestParams(), common.UserQueryParameter)
	if err != nil {
		userConn.HandleError("user ('u') not found in query parameters")
		walletExt.Logger().Error(fmt.Errorf("user not found in the query params: %w", err).Error())
		return
	}
	userIDBytes, err := getUserIDbyte(hexUserID)
	if err != nil {
		userConn.HandleError("error decoding userID. It should be in hex format")
		walletExt.Logger().Error(fmt.Errorf("error decoding string (%s), %w", hexUserID[2:], err).Error())
		return
	}

	address, err := getQueryParameter(userConn.ReadRequestParams(), common.AddressQueryParameter)
	if err != nil {
		userConn.HandleError("address ('a') not found in query parameters")
		walletExt.Logger().Error(fmt.Errorf("address not found in the query params: %w", err).Error())
		return
	}
	addressBytes, err := hex.DecodeString(address[2:]) // remove 0x prefix from address
	if err != nil {
		userConn.HandleError("Internal error")
		walletExt.Logger().Error(fmt.Errorf("error decoding string (%s), %w", address[2:], err).Error())
		return
	}

	// todo - this can be optimised and done in the database if we will have users with large number of accounts
	// get all the accounts for the selected user
	accounts, err := walletExt.Storage.GetAccounts(userIDBytes)
	if err != nil {
		userConn.HandleError("Internal error")
		walletExt.Logger().Error(fmt.Errorf("error getting accounts for user (%s), %w", hexUserID, err).Error())
		return
	}

	found := false
	for _, account := range accounts {
		if bytes.Equal(account.AccountAddress, addressBytes) {
			found = true
		}
	}

	// create and write the response
	res := struct {
		Status bool `json:"status"`
	}{Status: found}

	msg, err := json.Marshal(res)
	if err != nil {
		userConn.HandleError("Internal error")
		walletExt.Logger().Error(fmt.Errorf("error marshalling: %w", err).Error())
		return
	}

	err = userConn.WriteResponse(msg)
	if err != nil {
		walletExt.Logger().Error(fmt.Errorf("error writing success response, %w", err).Error())
		return
	}
}

// This function handles request to /revoke endpoint.
// It requires userID as query parameter and deletes given user and all associated viewing keys
func revokeRequestHandler(walletExt *walletextension.WalletExtension, userConn userconn.UserConn) {
	// read the request
	_, err := userConn.ReadRequest()
	if err != nil {
		walletExt.Logger().Error(fmt.Errorf("error reading request: %w", err).Error())
		return
	}

	hexUserID, err := getQueryParameter(userConn.ReadRequestParams(), common.UserQueryParameter)
	if err != nil {
		userConn.HandleError("user ('u') not found in query parameters")
		walletExt.Logger().Error(fmt.Errorf("user not found in the query params: %w", err).Error())
		return
	}
	userIDBytes, err := getUserIDbyte(hexUserID)
	if err != nil {
		userConn.HandleError("Internal error")
		walletExt.Logger().Error(fmt.Errorf("error decoding string (%s), %w", hexUserID, err).Error())
		return
	}

	err = walletExt.Storage.DeleteUser(userIDBytes)
	if err != nil {
		userConn.HandleError("Internal error")
		walletExt.Logger().Error(fmt.Errorf("error deleting user (%s), %w", hexUserID, err).Error())
		return
	}

	err = userConn.WriteResponse([]byte("success!"))
	if err != nil {
		walletExt.Logger().Error(fmt.Errorf("error writing success response, %w", err).Error())
	}
}

// calculate userID from public key
func calculateUserID(pk *ecdsa.PrivateKey) []byte {
	viewingPublicKeyBytes := crypto.CompressPubkey(&pk.PublicKey)
	return crypto.Keccak256Hash(viewingPublicKeyBytes).Bytes()
}

// check if message is in correct format and extracts userID and address from it
func getUserIDAndAddressFromMessage(message string) (string, string, error) {
	regex := regexp.MustCompile(common.MessageFormatRegex)
	if regex.MatchString(message) {
		params := regex.FindStringSubmatch(message)
		return params[1], params[2], nil
	}
	return "", "", errors.New("invalid message format")
}

// get an address that was used to sign given signature
func getAddressFromSignature(messageHash []byte, signature []byte) (gethcommon.Address, error) {
	pubKey, err := crypto.SigToPub(messageHash, signature)
	if err != nil {
		return gethcommon.Address{}, err
	}

	return crypto.PubkeyToAddress(*pubKey), nil
}

// convert userID from string to correct byte format
func getUserIDbyte(userID string) ([]byte, error) {
	return hex.DecodeString(userID[2:]) // remove 0x prefix from userID
}
