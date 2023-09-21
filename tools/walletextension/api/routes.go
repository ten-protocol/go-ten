package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

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
			Name: common.APIVersion1 + common.PathRoot,
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
			Name: common.APIVersion1 + common.PathJoin,
			Func: httpHandler(walletExt, joinRequestHandler),
		},
		{
			Name: common.APIVersion1 + common.PathAuthenticate,
			Func: httpHandler(walletExt, authenticateRequestHandler),
		},
		{
			Name: common.APIVersion1 + common.PathQuery,
			Func: httpHandler(walletExt, queryRequestHandler),
		},
		{
			Name: common.APIVersion1 + common.PathRevoke,
			Func: httpHandler(walletExt, revokeRequestHandler),
		},
		{
			Name: common.PathHealth,
			Func: httpHandler(walletExt, healthRequestHandler),
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
		err = fmt.Errorf("error reading request: %w", err)
		conn.HandleError(err.Error())
		walletExt.Logger().Error(err.Error())
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

	// Get userID and check if user exists (if not - use default user)
	hexUserID, err := getQueryParameter(conn.ReadRequestParams(), common.UserQueryParameter)
	if err != nil || !walletExt.UserExists(hexUserID) {
		walletExt.Logger().Error(fmt.Errorf("user not found in the query params: %w. Using the default user", err).Error())
		hexUserID = hex.EncodeToString([]byte(common.DefaultUser)) // todo (@ziga) - this can be removed once old WE endpoints are removed
	}

	// todo (@pedro) remove this conn dependency
	response, err := walletExt.ProxyEthRequest(request, conn, hexUserID)
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
		err = fmt.Errorf("error reading request: %w", err)
		conn.HandleError(err.Error())
		walletExt.Logger().Error(err.Error())
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

// This function handles request to /join endpoint. It is responsible to create new user (new key-pair) and store it to the db
func joinRequestHandler(walletExt *walletextension.WalletExtension, userConn userconn.UserConn) {
	// todo (@ziga) add protection against DDOS attacks
	_, err := userConn.ReadRequest()
	if err != nil {
		userConn.HandleError("Error: bad request")
		walletExt.Logger().Error(fmt.Errorf("error reading request: %w", err).Error())
		return
	}

	// generate new key-pair and store it in the database
	hexUserID, err := walletExt.GenerateAndStoreNewUser()
	if err != nil {
		userConn.HandleError("Internal error")
		walletExt.Logger().Error(fmt.Errorf("error creating new user, %w", err).Error())
	}

	// write hex encoded userID in the response
	err = userConn.WriteResponse([]byte(hexUserID))

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

	// check signature and add address and signature for that user
	err = walletExt.AddAddressToUser(hexUserID, message, signature)
	if err != nil {
		userConn.HandleError("Internal error")
		walletExt.Logger().Error(fmt.Errorf("error adding address to user with message: %s, %w", message, err).Error())
		return
	}
	err = userConn.WriteResponse([]byte(common.SuccessMsg))
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
	address, err := getQueryParameter(userConn.ReadRequestParams(), common.AddressQueryParameter)
	if err != nil {
		userConn.HandleError("address ('a') not found in query parameters")
		walletExt.Logger().Error(fmt.Errorf("address not found in the query params: %w", err).Error())
		return
	}

	// check if this account is registered with given user
	found, err := walletExt.UserHasAccount(hexUserID, address)
	if err != nil {
		userConn.HandleError("Internal error")
		walletExt.Logger().Error(fmt.Errorf("error during checking if account exists for user %s: %w", hexUserID, err).Error())
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

	// delete user and accounts associated with it from the databse
	err = walletExt.DeleteUser(hexUserID)
	if err != nil {
		userConn.HandleError("Internal error")
		walletExt.Logger().Error(fmt.Errorf("unable to delete user %s: %w", hexUserID, err).Error())
		return
	}

	err = userConn.WriteResponse([]byte(common.SuccessMsg))
	if err != nil {
		walletExt.Logger().Error(fmt.Errorf("error writing success response, %w", err).Error())
	}
}

// Handles request to /health endpoint.
func healthRequestHandler(walletExt *walletextension.WalletExtension, userConn userconn.UserConn) {
	// read the request
	_, err := userConn.ReadRequest()
	if err != nil {
		walletExt.Logger().Error(fmt.Errorf("error reading request: %w", err).Error())
		return
	}

	err = userConn.WriteResponse([]byte(common.SuccessMsg))
	if err != nil {
		walletExt.Logger().Error(fmt.Errorf("error writing success response, %w", err).Error())
	}
}
