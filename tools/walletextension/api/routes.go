package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ten-protocol/go-ten/go/common/log"

	"github.com/ten-protocol/go-ten/go/common/httputil"
	"github.com/ten-protocol/go-ten/go/rpc"
	"github.com/ten-protocol/go-ten/tools/walletextension"
	"github.com/ten-protocol/go-ten/tools/walletextension/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/userconn"

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
		{
			Name: common.PathNetworkHealth,
			Func: httpHandler(walletExt, networkHealthRequestHandler),
		},
		{
			Name: common.PathVersion,
			Func: httpHandler(walletExt, versionRequestHandler),
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
		handleEthError(nil, conn, walletExt.Logger(), fmt.Errorf("error reading request - %w", err))
		return
	}

	request, err := parseRequest(body)
	if err != nil {
		handleError(conn, walletExt.Logger(), err)
		return
	}
	walletExt.Logger().Debug("REQUEST", "method", request.Method, "body", string(body))

	if request.Method == rpc.Subscribe && !conn.SupportsSubscriptions() {
		handleError(conn, walletExt.Logger(), fmt.Errorf("received an %s request but the connection does not support subscriptions", rpc.Subscribe))
		return
	}

	// Get userID
	// TODO: @ziga - after removing old wallet extension endpoints we should prevent users doing anything without valid encryption token
	hexUserID, err := getUserID(conn, 1)
	if err != nil || !walletExt.UserExists(hexUserID) {
		walletExt.Logger().Info("user not found in the query params: %w. Using the default user", log.ErrKey, err)
		hexUserID = hex.EncodeToString([]byte(common.DefaultUser)) // todo (@ziga) - this can be removed once old WE endpoints are removed
	}

	if len(hexUserID) < 3 {
		handleError(conn, walletExt.Logger(), fmt.Errorf("encryption token length is incorrect"))
		return
	}

	// todo (@pedro) remove this conn dependency
	response, err := walletExt.ProxyEthRequest(request, conn, hexUserID)
	if err != nil {
		handleEthError(request, conn, walletExt.Logger(), err)
		return
	}

	rpcResponse, err := json.Marshal(response)
	if err != nil {
		handleEthError(request, conn, walletExt.Logger(), err)
		return
	}

	walletExt.Logger().Info(fmt.Sprintf("Forwarding %s response from Obscuro node: %s", request.Method, rpcResponse))
	if err = conn.WriteResponse(rpcResponse); err != nil {
		walletExt.Logger().Error("error writing success response", log.ErrKey, err)
	}
}

// readyRequestHandler is used to check whether the server is ready
func readyRequestHandler(_ *walletextension.WalletExtension, _ userconn.UserConn) {}

// generateViewingKeyRequestHandler parses the gen vk request
func generateViewingKeyRequestHandler(walletExt *walletextension.WalletExtension, conn userconn.UserConn) {
	body, err := conn.ReadRequest()
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("error reading request: %w", err))
		return
	}

	var reqJSONMap map[string]string
	err = json.Unmarshal(body, &reqJSONMap)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("could not unmarshal address request - %w", err))
		return
	}

	address := gethcommon.HexToAddress(reqJSONMap[common.JSONKeyAddress])

	pubViewingKey, err := walletExt.GenerateViewingKey(address)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("unable to generate vieweing key - %w", err))
		return
	}

	err = conn.WriteResponse([]byte(pubViewingKey))
	if err != nil {
		walletExt.Logger().Error("error writing success response", log.ErrKey, err)
	}
}

// submitViewingKeyRequestHandler submits the viewing key and signed bytes to the WE
func submitViewingKeyRequestHandler(walletExt *walletextension.WalletExtension, conn userconn.UserConn) {
	body, err := conn.ReadRequest()
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("error reading request: %w", err))
		return
	}

	var reqJSONMap map[string]string
	err = json.Unmarshal(body, &reqJSONMap)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("could not unmarshal address request - %w", err))
		return
	}
	accAddress := gethcommon.HexToAddress(reqJSONMap[common.JSONKeyAddress])

	signature, err := hex.DecodeString(reqJSONMap[common.JSONKeySignature][2:])
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("could not decode signature from client to hex - %w", err))
		return
	}

	err = walletExt.SubmitViewingKey(accAddress, signature)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("could not submit viewing key - %w", err))
		return
	}

	err = conn.WriteResponse([]byte(common.SuccessMsg))
	if err != nil {
		walletExt.Logger().Error("error writing success response", log.ErrKey, err)
		return
	}
}

// This function handles request to /join endpoint. It is responsible to create new user (new key-pair) and store it to the db
func joinRequestHandler(walletExt *walletextension.WalletExtension, conn userconn.UserConn) {
	// todo (@ziga) add protection against DDOS attacks
	_, err := conn.ReadRequest()
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("error reading request: %w", err))
		return
	}

	// generate new key-pair and store it in the database
	hexUserID, err := walletExt.GenerateAndStoreNewUser()
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("internal Error"))
		walletExt.Logger().Error("error creating new user", log.ErrKey, err)
	}

	// write hex encoded userID in the response
	err = conn.WriteResponse([]byte(hexUserID))
	if err != nil {
		walletExt.Logger().Error("error writing success response", log.ErrKey, err)
	}
}

// This function handles request to /authenticate endpoint.
// In the request we receive message, signature and address in JSON as request body and userID and address as query parameters
// We then check if message is in correct format and if signature is valid. If all checks pass we save address and signature against userID
func authenticateRequestHandler(walletExt *walletextension.WalletExtension, conn userconn.UserConn) {
	// read the request
	body, err := conn.ReadRequest()
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("error reading request: %w", err))
		return
	}

	// get the text that was signed and signature
	var reqJSONMap map[string]string
	err = json.Unmarshal(body, &reqJSONMap)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("could not unmarshal address request - %w", err))
		return
	}

	// get signature from the request and remove leading two bytes (0x)
	signature, err := hex.DecodeString(reqJSONMap[common.JSONKeySignature][2:])
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("unable to decode signature - %w", err))
		return
	}

	// get address from the request
	address, ok := reqJSONMap[common.JSONKeyAddress]
	if !ok || address == "" {
		handleError(conn, walletExt.Logger(), fmt.Errorf("unable to read address field from the request"))
		return
	}

	// get optional type of the message that was signed
	messageTypeValue := common.DefaultGatewayAuthMessageType
	if typeFromRequest, ok := reqJSONMap[common.JSONKeyType]; ok && typeFromRequest != "" {
		messageTypeValue = typeFromRequest
	}

	// check if message type is valid
	messageType, ok := common.TypeMap[messageTypeValue]
	if !ok {
		handleError(conn, walletExt.Logger(), fmt.Errorf("invalid message type + %s", messageTypeValue))
	}

	// read userID from query params
	hexUserID, err := getUserID(conn, 2)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("malformed query: 'u' required - representing encryption token - %w", err))
		return
	}

	// check signature and add address and signature for that user
	err = walletExt.AddAddressToUser(hexUserID, address, signature, messageType)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("internal error"))
		walletExt.Logger().Error(fmt.Sprintf("error adding address: %s to user: %s with signature: %s", address, hexUserID, signature))
		return
	}
	err = conn.WriteResponse([]byte(common.SuccessMsg))
	if err != nil {
		walletExt.Logger().Error("error writing success response", log.ErrKey, err)
	}
}

// This function handles request to /query endpoint.
// In the query parameters address and userID are required. We check if provided address is registered for given userID
// and return true/false in json response
func queryRequestHandler(walletExt *walletextension.WalletExtension, conn userconn.UserConn) {
	// read the request
	_, err := conn.ReadRequest()
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("error reading request: %w", err))
		return
	}

	hexUserID, err := getUserID(conn, 2)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("user ('u') not found in query parameters"))
		walletExt.Logger().Info("user not found in the query params", log.ErrKey, err)
		return
	}
	address, err := getQueryParameter(conn.ReadRequestParams(), common.AddressQueryParameter)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("address ('a') not found in query parameters"))
		walletExt.Logger().Error("address ('a') not found in query parameters", log.ErrKey, err)
		return
	}
	// check if address length is correct
	if len(address) != common.EthereumAddressLen {
		handleError(conn, walletExt.Logger(), fmt.Errorf("provided address length is %d, expected: %d", len(address), common.EthereumAddressLen))
		return
	}

	// check if this account is registered with given user
	found, err := walletExt.UserHasAccount(hexUserID, address)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("internal error"))
		walletExt.Logger().Error("error during checking if account exists for user", "hexUserID", hexUserID, log.ErrKey, err)
	}

	// create and write the response
	res := struct {
		Status bool `json:"status"`
	}{Status: found}

	msg, err := json.Marshal(res)
	if err != nil {
		handleError(conn, walletExt.Logger(), err)
		return
	}

	err = conn.WriteResponse(msg)
	if err != nil {
		walletExt.Logger().Error("error writing success response", log.ErrKey, err)
	}
}

// This function handles request to /revoke endpoint.
// It requires userID as query parameter and deletes given user and all associated viewing keys
func revokeRequestHandler(walletExt *walletextension.WalletExtension, conn userconn.UserConn) {
	// read the request
	_, err := conn.ReadRequest()
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("error reading request: %w", err))
		return
	}

	hexUserID, err := getUserID(conn, 2)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("user ('u') not found in query parameters"))
		walletExt.Logger().Info("user not found in the query params", log.ErrKey, err)
		return
	}

	// delete user and accounts associated with it from the database
	err = walletExt.DeleteUser(hexUserID)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("internal error"))
		walletExt.Logger().Error("unable to delete user", "hexUserID", hexUserID, log.ErrKey, err)
		return
	}

	err = conn.WriteResponse([]byte(common.SuccessMsg))
	if err != nil {
		walletExt.Logger().Error("error writing success response", log.ErrKey, err)
	}
}

// Handles request to /health endpoint.
func healthRequestHandler(walletExt *walletextension.WalletExtension, conn userconn.UserConn) {
	// read the request
	_, err := conn.ReadRequest()
	if err != nil {
		walletExt.Logger().Error("error reading request", log.ErrKey, err)
		return
	}

	// TODO: connect to database and check if it is healthy
	err = conn.WriteResponse([]byte(common.SuccessMsg))
	if err != nil {
		walletExt.Logger().Error("error writing success response", log.ErrKey, err)
	}
}

// Handles request to /network-health endpoint.
func networkHealthRequestHandler(walletExt *walletextension.WalletExtension, userConn userconn.UserConn) {
	// read the request
	_, err := userConn.ReadRequest()
	if err != nil {
		walletExt.Logger().Error("error reading request", log.ErrKey, err)
		return
	}

	healthStatus, err := walletExt.GetTenNodeHealthStatus()

	data, err := json.Marshal(map[string]interface{}{
		"result": healthStatus,
		"error":  err,
	})
	if err != nil {
		walletExt.Logger().Error("error marshaling response", log.ErrKey, err)
		return
	}

	err = userConn.WriteResponse(data)

	if err != nil {
		walletExt.Logger().Error("error writing success response", log.ErrKey, err)
	}
}

// Handles request to /version endpoint.
func versionRequestHandler(walletExt *walletextension.WalletExtension, userConn userconn.UserConn) {
	// read the request
	_, err := userConn.ReadRequest()
	if err != nil {
		walletExt.Logger().Error("error reading request", log.ErrKey, err)
		return
	}

	err = userConn.WriteResponse([]byte(walletExt.Version()))
	if err != nil {
		walletExt.Logger().Error("error writing success response", log.ErrKey, err)
	}
}
