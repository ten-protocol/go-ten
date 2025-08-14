package httpapi

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	tencommon "github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/cache"
	"github.com/ten-protocol/go-ten/tools/walletextension/keymanager"
	"github.com/ten-protocol/go-ten/tools/walletextension/services"

	"github.com/status-im/keycard-go/hexutils"

	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/viewingkey"
	"github.com/ten-protocol/go-ten/lib/gethfork/node"

	"github.com/ten-protocol/go-ten/go/common/httputil"
	"github.com/ten-protocol/go-ten/tools/walletextension/common"
)

// NewHTTPRoutes returns the http specific routes
// todo - move these to the rpc framework.
func NewHTTPRoutes(walletExt *services.Services) []node.Route {
	return []node.Route{
		{
			Name: common.APIVersion1 + common.PathReady,
			Func: httpHandler(walletExt, readyRequestHandler),
		},
		{
			Name: common.APIVersion1 + common.PathJoin,
			Func: httpHandler(walletExt, joinRequestHandler),
		},
		{
			Name: common.APIVersion1 + common.PathGetToken,
			Func: httpHandler(walletExt, getTokenRequestHandler),
		},
		{
			Name: common.APIVersion1 + common.PathSetToken,
			Func: httpHandler(walletExt, setTokenRequestHandler),
		},
		{
			Name: common.APIVersion1 + common.PathGetMessage,
			Func: httpHandler(walletExt, getMessageRequestHandler),
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
			Name: common.APIVersion1 + common.PathHealth,
			Func: httpHandler(walletExt, healthRequestHandler),
		},
		{
			Name: common.APIVersion1 + common.PathNetworkHealth,
			Func: httpHandler(walletExt, networkHealthRequestHandler),
		},
		{
			Name: common.APIVersion1 + common.PathVersion,
			Func: httpHandler(walletExt, versionRequestHandler),
		},
		{
			Name: common.APIVersion1 + common.PathNetworkConfig,
			Func: httpHandler(walletExt, networkConfigRequestHandler),
		},
		{
			Name: common.APIVersion1 + common.PathKeyExchange,
			Func: httpHandler(walletExt, keyExchangeRequestHandler),
		},
		{
			Name: common.APIVersion1 + common.PathSessionKeys + "create",
			Func: httpHandler(walletExt, createSKRequestHandler),
		},
		{
			Name: common.APIVersion1 + common.PathSessionKeys + "activate",
			Func: httpHandler(walletExt, activateSKRequestHandler),
		},
		{
			Name: common.APIVersion1 + common.PathSessionKeys + "deactivate",
			Func: httpHandler(walletExt, deactivateSKRequestHandler),
		},
		{
			Name: common.APIVersion1 + common.PathSessionKeys + "delete",
			Func: httpHandler(walletExt, deleteSKRequestHandler),
		},
		{
			Name: common.APIVersion1 + common.PathSessionKeys + "list",
			Func: httpHandler(walletExt, listSKRequestHandler),
		},
	}
}

func httpHandler(
	walletExt *services.Services,
	fun func(walletExt *services.Services, conn UserConn),
) func(resp http.ResponseWriter, req *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		httpRequestHandler(walletExt, resp, req, fun)
	}
}

// Overall request handler for http requests
func httpRequestHandler(walletExt *services.Services, resp http.ResponseWriter, req *http.Request, fun func(walletExt *services.Services, conn UserConn)) {
	if walletExt.IsStopping() {
		return
	}
	if httputil.EnableCORS(resp, req) {
		return
	}
	userConn := NewUserConnHTTP(resp, req, walletExt.Logger())
	fun(walletExt, userConn)
}

// readyRequestHandler is used to check whether the server is ready
func readyRequestHandler(_ *services.Services, _ UserConn) {}

// This function handles request to /join endpoint. It is responsible to create new user (new key-pair) and store it to the db
func joinRequestHandler(walletExt *services.Services, conn UserConn) {
	// todo (@ziga) add protection against DDOS attacks
	_, err := conn.ReadRequest()
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("error reading request: %w", err))
		return
	}

	// generate new key-pair and store it in the database
	userID, err := walletExt.GenerateAndStoreNewUser()
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("internal Error"))
		walletExt.Logger().Error("error creating new user", log.ErrKey, err)
	}

	// set secure HTTP-only cookie with userID
	cookie := &http.Cookie{
		Name:     "gateway_token",
		Value:    hexutils.BytesToHex(userID),
		Path:     "/",
		Domain:   ".ten.xyz",              // Share across all .ten.xyz subdomains
		HttpOnly: true,                    // Prevents XSS
		Secure:   true,                    // HTTPS only
		SameSite: http.SameSiteNoneMode,   // Required for cross-origin AJAX requests
		MaxAge:   365 * 24 * 60 * 60 * 10, // 10 years (effectively permanent)
	}

	err = conn.SetCookie(cookie)
	if err != nil {
		walletExt.Logger().Error("error setting cookie", log.ErrKey, err)
	}

	// write hex encoded userID in the response
	err = conn.WriteResponse([]byte(hexutils.BytesToHex(userID)))
	if err != nil {
		walletExt.Logger().Error("error writing success response", log.ErrKey, err)
	}
}

// This function handles request to /get-token endpoint. It reads the session key from the cookie and returns it to the user.
func getTokenRequestHandler(walletExt *services.Services, conn UserConn) {
	// Get the HTTP request to access cookies
	req := conn.GetHTTPRequest()
	if req == nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("could not access request"))
		return
	}

	// Find the gateway_token cookie
	var userID string
	for _, cookie := range req.Cookies() {
		if cookie.Name == "gateway_token" {
			userID = cookie.Value
			break
		}
	}

	if userID == "" {
		handleError(conn, walletExt.Logger(), fmt.Errorf("gateway_token cookie not found"))
		return
	}

	// Validate the token format (should be hex)
	userIDBytes, err := hex.DecodeString(userID)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("invalid token format: %w", err))
		return
	}

	if len(userIDBytes) == 0 {
		handleError(conn, walletExt.Logger(), fmt.Errorf("token cannot be empty"))
		return
	}

	// Verify the user exists in the database
	_, err = walletExt.Storage.GetUser(userIDBytes)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("user not found in database"))
		return
	}

	// Return the userID from the cookie (same format as /join)
	err = conn.WriteResponse([]byte(userID))
	if err != nil {
		walletExt.Logger().Error("error writing token response", log.ErrKey, err)
	}
}

// This function handles request to /set-token endpoint. It receives a token in JSON format and sets it as a session cookie.
func setTokenRequestHandler(walletExt *services.Services, conn UserConn) {
	// Read the request body to get the token
	requestBody, err := conn.ReadRequest()
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("error reading request: %w", err))
		return
	}

	// Parse the JSON request
	type SetTokenRequest struct {
		Token string `json:"token"`
	}

	var req SetTokenRequest
	err = json.Unmarshal(requestBody, &req)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("invalid JSON format: %w", err))
		return
	}

	if req.Token == "" {
		handleError(conn, walletExt.Logger(), fmt.Errorf("token is required"))
		return
	}

	// Validate the token format (should be hex)
	userIDBytes, err := hex.DecodeString(req.Token)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("invalid token format: %w", err))
		return
	}

	if len(userIDBytes) == 0 {
		handleError(conn, walletExt.Logger(), fmt.Errorf("token cannot be empty"))
		return
	}

	// Verify the user exists in the database
	_, err = walletExt.Storage.GetUser(userIDBytes)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("user not found in database"))
		return
	}

	// Set the cookie with the provided token
	cookie := &http.Cookie{
		Name:     "gateway_token",
		Value:    req.Token,
		Path:     "/",
		Domain:   ".ten.xyz",              // Share across all .ten.xyz subdomains
		HttpOnly: true,                    // Prevents XSS
		Secure:   true,                    // HTTPS only
		SameSite: http.SameSiteNoneMode,   // Required for cross-origin AJAX requests
		MaxAge:   365 * 24 * 60 * 60 * 10, // 10 years (effectively permanent)
	}

	err = conn.SetCookie(cookie)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("error setting cookie: %w", err))
		return
	}

	// Return success response
	successResponse := map[string]string{"status": "success", "message": "Token cookie set successfully"}
	responseBytes, err := json.Marshal(successResponse)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("error marshaling response: %w", err))
		return
	}

	err = conn.WriteResponse(responseBytes)
	if err != nil {
		walletExt.Logger().Error("error writing success response", log.ErrKey, err)
	}
}

// This function handles request to /authenticate endpoint.
// In the request we receive message, signature and address in JSON as request body and userID and address as query parameters
// We then check if message is in correct format and if signature is valid. If all checks pass we save address and signature against userID
func authenticateRequestHandler(walletExt *services.Services, conn UserConn) {
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
		handleError(conn, walletExt.Logger(), fmt.Errorf("could not unmarshal request body - %w", err))
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

	// get an optional type of the message that was signed
	messageTypeValue := common.DefaultGatewayAuthMessageType
	if typeFromRequest, ok := reqJSONMap[common.JSONKeyType]; ok && typeFromRequest != "" {
		messageTypeValue = typeFromRequest
	}

	// check if a message type is valid
	messageType, ok := viewingkey.SignatureTypeMap[messageTypeValue]
	if !ok {
		handleError(conn, walletExt.Logger(), fmt.Errorf("invalid message type: %s", messageTypeValue))
	}

	// read userID from query params
	userID, err := getUserID(conn)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("malformed query: 'u' required - representing encryption token - %w", err))
		return
	}

	// check signature and add address and signature for that user
	err = walletExt.AddAddressToUser(userID, address, signature, messageType)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("internal error"))
		walletExt.Logger().Error(fmt.Sprintf("error adding address: %s to user: %s with signature: %s", address, userID, signature))
		return
	}
	err = conn.WriteResponse([]byte(common.SuccessMsg))
	if err != nil {
		walletExt.Logger().Error("error writing success response", log.ErrKey, err)
	}
}

// todo - is this needed?
// This function handles request to /query endpoint.
// In the query parameters address and userID are required. We check if provided address is registered for given userID
// and return true/false in json response
func queryRequestHandler(walletExt *services.Services, conn UserConn) {
	// read the request
	_, err := conn.ReadRequest()
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("error reading request: %w", err))
		return
	}

	userID, err := getUserID(conn)
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
	found, err := walletExt.UserHasAccount(userID, address)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("internal error"))
		walletExt.Logger().Error("error during checking if account exists for user", "userID", userID, log.ErrKey, err)
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
func revokeRequestHandler(walletExt *services.Services, conn UserConn) {
	// read the request
	_, err := conn.ReadRequest()
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("error reading request: %w", err))
		return
	}

	userID, err := getUserID(conn)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("user ('u') not found in query parameters"))
		walletExt.Logger().Info("user not found in the query params", log.ErrKey, err)
		return
	}

	// delete user and accounts associated with it from the database
	err = walletExt.Storage.DeleteUser(userID)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("internal error"))
		walletExt.Logger().Error("unable to delete user", "userID", userID, log.ErrKey, err)
		return
	}

	err = conn.WriteResponse([]byte(common.SuccessMsg))
	if err != nil {
		walletExt.Logger().Error("error writing success response", log.ErrKey, err)
	}
}

// Handles request to /health endpoint.
func healthRequestHandler(walletExt *services.Services, conn UserConn) {
	// read the request
	_, err := conn.ReadRequest()
	if err != nil {
		walletExt.Logger().Error("error reading request", log.ErrKey, err)
		return
	}

	// Use cache for health check response
	cacheKey := []byte("health_check")
	cacheCfg := &cache.Cfg{Type: cache.LatestBatch} // Short-living cache

	result, err := cache.WithCache(walletExt.RPCResponsesCache, cacheCfg, cacheKey, func() (*[]byte, error) {
		// TODO: connect to database and check if it is healthy
		response := []byte(common.SuccessMsg)
		return &response, nil
	})
	if err != nil {
		walletExt.Logger().Error("error getting health status", log.ErrKey, err)
		return
	}

	err = conn.WriteResponse(*result)
	if err != nil {
		walletExt.Logger().Error("error writing success response", log.ErrKey, err)
	}
}

// Handles request to /network-health endpoint.
func networkHealthRequestHandler(walletExt *services.Services, userConn UserConn) {
	// read the request
	_, err := userConn.ReadRequest()
	if err != nil {
		walletExt.Logger().Error("error reading request", log.ErrKey, err)
		return
	}

	// Use cache for network health check response
	cacheKey := []byte("network_health_check")
	cacheCfg := &cache.Cfg{Type: cache.LatestBatch} // Short-living cache

	result, err := cache.WithCache(walletExt.RPCResponsesCache, cacheCfg, cacheKey, func() (*[]byte, error) {
		// call `obscuro-health` rpc method to get the health status of the node
		healthStatus, err := walletExt.GetTenNodeHealthStatus()

		// create the response in the required format
		type HealthStatus struct {
			Errors        []string `json:"Errors"`
			OverallHealth bool     `json:"OverallHealth"`
		}

		errorStrings := make([]string, 0)
		if err != nil {
			errorStrings = append(errorStrings, err.Error())
		}
		healthStatusResponse := HealthStatus{
			Errors:        errorStrings,
			OverallHealth: healthStatus,
		}

		data, err := json.Marshal(map[string]interface{}{
			"id":      "1",
			"jsonrpc": "2.0",
			"result":  healthStatusResponse,
		})
		if err != nil {
			return nil, fmt.Errorf("error marshaling response: %w", err)
		}

		return &data, nil
	})
	if err != nil {
		walletExt.Logger().Error("error getting network health status", log.ErrKey, err)
		return
	}

	err = userConn.WriteResponse(*result)
	if err != nil {
		walletExt.Logger().Error("error writing success response", log.ErrKey, err)
	}
}

func networkConfigRequestHandler(walletExt *services.Services, userConn UserConn) {
	// read the request
	_, err := userConn.ReadRequest()
	if err != nil {
		walletExt.Logger().Error("error reading request", log.ErrKey, err)
		return
	}

	// Use cache for network config response
	cacheKey := []byte("network_config")
	cacheCfg := &cache.Cfg{Type: cache.LongLiving} // Long-living cache

	result, err := cache.WithCache(walletExt.RPCResponsesCache, cacheCfg, cacheKey, func() (*[]byte, error) {
		// Call the RPC method to get the network configuration
		networkConfig, err := walletExt.GetTenNetworkConfig()
		if err != nil {
			return nil, fmt.Errorf("error fetching network config: %w", err)
		}

		// Define a struct to represent the response
		type NetworkConfigResponse struct {
			NetworkConfigAddress            string            `json:"NetworkConfig"`
			EnclaveRegistryAddress          string            `json:"EnclaveRegistry"`
			DataAvailabilityRegistryAddress string            `json:"DataAvailabilityRegistry"`
			CrossChainAddress               string            `json:"CrossChain"`
			L1MessageBusAddress             string            `json:"L1MessageBus"`
			L2MessageBusAddress             string            `json:"L2MessageBus"`
			L1BridgeAddress                 string            `json:"L1Bridge"`
			L2BridgeAddress                 string            `json:"L2Bridge"`
			L1CrossChainMessengerAddress    string            `json:"L1CrossChainMessenger"`
			L2CrossChainMessengerAddress    string            `json:"L2CrossChainMessenger"`
			SystemContractsUpgrader         string            `json:"SystemContractsUpgrader"`
			L1StartHash                     string            `json:"L1StartHash"`
			AdditionalContracts             map[string]string `json:"AdditionalContracts"`
		}

		// Convert the TenNetworkInfo fields to strings
		additionalContracts := make(map[string]string)
		if len(networkConfig.AdditionalContracts) > 0 {
			for _, contract := range networkConfig.AdditionalContracts {
				additionalContracts[contract.Name] = contract.Addr.Hex()
			}
		}

		networkConfigResponse := NetworkConfigResponse{
			NetworkConfigAddress:            networkConfig.NetworkConfig.Hex(),
			EnclaveRegistryAddress:          networkConfig.EnclaveRegistry.Hex(),
			DataAvailabilityRegistryAddress: networkConfig.DataAvailabilityRegistry.Hex(),
			CrossChainAddress:               networkConfig.CrossChain.Hex(),
			L1MessageBusAddress:             networkConfig.L1MessageBus.Hex(),
			L2MessageBusAddress:             networkConfig.L2MessageBus.Hex(),
			L1BridgeAddress:                 networkConfig.L1Bridge.Hex(),
			L2BridgeAddress:                 networkConfig.L2Bridge.Hex(),
			L1CrossChainMessengerAddress:    networkConfig.L1CrossChainMessenger.Hex(),
			L2CrossChainMessengerAddress:    networkConfig.L2CrossChainMessenger.Hex(),
			SystemContractsUpgrader:         networkConfig.SystemContractsUpgrader.Hex(),
			L1StartHash:                     networkConfig.L1StartHash.Hex(),
			AdditionalContracts:             additionalContracts,
		}

		// Marshal the response into JSON format
		data, err := json.Marshal(networkConfigResponse)
		if err != nil {
			return nil, fmt.Errorf("error marshaling response: %w", err)
		}

		return &data, nil
	})
	if err != nil {
		walletExt.Logger().Error("error getting network config", log.ErrKey, err)
		return
	}

	// Write the response back to the user
	err = userConn.WriteResponse(*result)
	if err != nil {
		walletExt.Logger().Error("error writing success response", log.ErrKey, err)
	}
}

// Handles request to /version endpoint.
func versionRequestHandler(walletExt *services.Services, userConn UserConn) {
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

// getMessageRequestHandler handles request to /getmessage endpoint.
func getMessageRequestHandler(walletExt *services.Services, conn UserConn) {
	// read the request
	body, err := conn.ReadRequest()
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("error reading request: %w", err))
		return
	}

	// read body of the request
	var reqJSONMap map[string]interface{}
	err = json.Unmarshal(body, &reqJSONMap)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("could not unmarshal address request - %w", err))
		return
	}

	// get address from the request
	encryptionToken, ok := reqJSONMap[common.JSONKeyEncryptionToken]
	if !ok {
		handleError(conn, walletExt.Logger(), fmt.Errorf("encryptionToken field not found in the request"))
		return
	}
	if tokenStr, ok := encryptionToken.(string); !ok {
		handleError(conn, walletExt.Logger(), fmt.Errorf("encryptionToken field is not a string"))
		return
	} else if len(tokenStr) != common.MessageUserIDLen {
		handleError(conn, walletExt.Logger(), fmt.Errorf("encryptionToken field is not of correct length"))
		return
	}

	// get formats from the request, if present
	var formatsSlice []string
	if formatsInterface, ok := reqJSONMap[common.JSONKeyFormats]; ok {
		formats, ok := formatsInterface.([]interface{})
		if !ok {
			handleError(conn, walletExt.Logger(), fmt.Errorf("formats field is not an array"))
			return
		}

		for _, f := range formats {
			formatStr, ok := f.(string)
			if !ok {
				handleError(conn, walletExt.Logger(), fmt.Errorf("format value is not a string"))
				return
			}
			formatsSlice = append(formatsSlice, formatStr)
		}
	}

	userID := hexutils.HexToBytes(encryptionToken.(string))
	if len(userID) != viewingkey.UserIDLength {
		return
	}

	message, err := walletExt.GenerateUserMessageToSign(userID, formatsSlice)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("internal error"))
		walletExt.Logger().Error("error getting message", log.ErrKey, err)
		return
	}

	// create the response structure for EIP712 message where the message is a JSON object
	type JSONResponseEIP712 struct {
		Message json.RawMessage `json:"message"`
		Type    string          `json:"type"`
	}

	// create the response structure for personal sign message where the message is a string
	type JSONResponsePersonal struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	}

	// get string representation of the message format
	messageFormat := viewingkey.GetBestFormat(formatsSlice)
	messageFormatString := viewingkey.GetSignatureTypeString(messageFormat)
	responseBytes := []byte{}
	if messageFormat == viewingkey.PersonalSign {
		response := JSONResponsePersonal{
			Message: message,
			Type:    messageFormatString,
		}

		responseBytes, err = json.Marshal(response)
		if err != nil {
			handleError(conn, walletExt.Logger(), fmt.Errorf("error marshaling JSON response: %w", err))
			return
		}
	} else if messageFormat == viewingkey.EIP712Signature {
		var messageMap map[string]interface{}
		err = json.Unmarshal([]byte(message), &messageMap)
		if err != nil {
			handleError(conn, walletExt.Logger(), fmt.Errorf("error unmarshaling JSON: %w", err))
			return
		}

		if domainMap, ok := messageMap["domain"].(map[string]interface{}); ok {
			delete(domainMap, "salt")
		}

		if typesMap, ok := messageMap["types"].(map[string]interface{}); ok {
			delete(typesMap, "EIP712Domain")
		}

		// Marshal the modified map back to JSON
		modifiedMessage, err := json.Marshal(messageMap)
		if err != nil {
			handleError(conn, walletExt.Logger(), fmt.Errorf("error marshaling modified JSON: %w", err))
			return
		}

		response := JSONResponseEIP712{
			Message: modifiedMessage,
			Type:    messageFormatString,
		}

		responseBytes, err = json.Marshal(response)
		if err != nil {
			handleError(conn, walletExt.Logger(), fmt.Errorf("error marshaling JSON response: %w", err))
			return
		}
	}

	err = conn.WriteResponse(responseBytes)
	if err != nil {
		walletExt.Logger().Error("error writing success response", log.ErrKey, err)
	}
}

func listSKRequestHandler(walletExt *services.Services, conn UserConn) {
	withUser(walletExt, conn, func(user *common.GWUser) ([]byte, error) {
		if user.SessionKey == nil {
			return []byte{}, nil
		}
		return []byte(hexutils.BytesToHex(user.SessionKey.Account.Address.Bytes())), nil
	})
}

func createSKRequestHandler(walletExt *services.Services, conn UserConn) {
	withUser(walletExt, conn, func(user *common.GWUser) ([]byte, error) {
		sk, err := walletExt.SKManager.CreateSessionKey(user)
		if err != nil {
			handleError(conn, walletExt.Logger(), fmt.Errorf("could not create session key: %w", err))
			return nil, err
		}
		return []byte(hexutils.BytesToHex(sk.Account.Address.Bytes())), nil
	})
}

func deleteSKRequestHandler(walletExt *services.Services, conn UserConn) {
	withUser(walletExt, conn, func(user *common.GWUser) ([]byte, error) {
		res, err := walletExt.SKManager.DeleteSessionKey(user)
		return []byte{boolToByte(res)}, err
	})
}

func activateSKRequestHandler(walletExt *services.Services, conn UserConn) {
	withUser(walletExt, conn, func(user *common.GWUser) ([]byte, error) {
		res, err := walletExt.SKManager.ActivateSessionKey(user)
		return []byte{boolToByte(res)}, err
	})
}

func deactivateSKRequestHandler(walletExt *services.Services, conn UserConn) {
	withUser(walletExt, conn, func(user *common.GWUser) ([]byte, error) {
		res, err := walletExt.SKManager.DeactivateSessionKey(user)
		return []byte{boolToByte(res)}, err
	})
}

// extracts the user from the request, and writes the response to the connection
func withUser(walletExt *services.Services, conn UserConn, withUser func(user *common.GWUser) ([]byte, error)) {
	_, err := conn.ReadRequest()
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("error reading request: %w", err))
		return
	}

	userID, err := getUserID(conn)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("user ('u') not found in query parameters"))
		walletExt.Logger().Info("user not found in the query params", log.ErrKey, err)
		return
	}

	user, err := walletExt.Storage.GetUser(userID)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("could not get user: %w", err))
		return
	}

	resp, err := withUser(user)
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("could not process request: %w", err))
		return
	}

	err = conn.WriteResponse(resp)
	if err != nil {
		walletExt.Logger().Error("error writing success response", log.ErrKey, err)
	}
}

func boolToByte(res bool) byte {
	if res {
		return 1
	}
	return 0
}

func keyExchangeRequestHandler(walletExt *services.Services, conn UserConn) {
	// Read the request
	body, err := conn.ReadRequest()
	if err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("error reading request: %w", err))
		return
	}

	// Step 1: Deserialize the received message
	var receivedMessageOG keymanager.KeyExchangeRequest
	err = json.Unmarshal(body, &receivedMessageOG)
	if err != nil {
		walletExt.Logger().Error("OG: Failed to deserialize received message", log.ErrKey, err)
		handleError(conn, walletExt.Logger(), fmt.Errorf("failed to deserialize message: %w", err))
		return
	}

	// Step 2: Deserialize the public key
	receivedPubKey, err := keymanager.DeserializePublicKey(receivedMessageOG.PublicKey)
	if err != nil {
		walletExt.Logger().Error("OG: Failed to deserialize public key", log.ErrKey, err)
		handleError(conn, walletExt.Logger(), fmt.Errorf("failed to deserialize public key: %w", err))
		return
	}

	// Step 3: Deserialize the attestation report
	var receivedAttestation tencommon.AttestationReport
	if err := json.Unmarshal(receivedMessageOG.Attestation, &receivedAttestation); err != nil {
		handleError(conn, walletExt.Logger(), fmt.Errorf("error unmarshaling attestation report: %w", err))
		return
	}

	// Step 4: Verify the attestation report
	verifiedData, err := keymanager.VerifyReport(&receivedAttestation)
	if err != nil {
		walletExt.Logger().Error("OG: Failed to verify attestation report", log.ErrKey, err)
		handleError(conn, walletExt.Logger(), fmt.Errorf("failed to verify attestation report: %w", err))
		return
	}

	// Hash the received public key bytes
	pubKeyHash := sha256.Sum256(receivedMessageOG.PublicKey)

	// Only compare the first 32 bytes since verifiedData is padded to 64 bytes
	verifiedDataTruncated := verifiedData[:32]
	if bytes.Equal(verifiedDataTruncated, pubKeyHash[:]) {
		walletExt.Logger().Info("OG: Public keys match")
	} else {
		walletExt.Logger().Error("OG: Public keys do not match")
	}

	// Step 5 Encrypt the encryption key using the received public key
	encryptedKeyOG, err := keymanager.EncryptWithPublicKey(walletExt.Storage.GetEncryptionKey(), receivedPubKey)
	if err != nil {
		walletExt.Logger().Error("OG: Encryption failed", log.ErrKey, err)
		handleError(conn, walletExt.Logger(), fmt.Errorf("encryption failed: %w", err))
		return
	}

	// Step 6: Encode the encrypted encryption key to Base64
	encodedEncryptedKeyOG := keymanager.EncodeBase64(encryptedKeyOG)

	// Step 7: Create the response message containing the encrypted key
	messageOG := keymanager.KeyExchangeResponse{
		EncryptedKey: encodedEncryptedKeyOG,
	}

	// Step 8: Serialize the response message to JSON and send it back to the requester
	messageBytesOG, err := json.Marshal(messageOG)
	if err != nil {
		walletExt.Logger().Error("OG: Failed to serialize response message", log.ErrKey, err)
		handleError(conn, walletExt.Logger(), fmt.Errorf("failed to serialize response message: %w", err))
		return
	}
	walletExt.Logger().Info("Shared encrypted key with another gateway enclave")
	err = conn.WriteResponse(messageBytesOG)
	if err != nil {
		walletExt.Logger().Error("error writing response", log.ErrKey, err)
	}
}
