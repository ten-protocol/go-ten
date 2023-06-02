package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/ethereum/go-ethereum/crypto/ecies"

	"github.com/ethereum/go-ethereum/crypto"

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
			Name: common.PathAuthenticate,
			Func: httpHandler(walletExt, authenticateRequestHandler),
		},

		{
			Name: common.PathJoin,
			Func: httpHandler(walletExt, joinRequestHandler),
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

		{
			Name: common.PathAuthenticate,
			Func: wsHandler(walletExt, authenticateRequestHandler),
		},

		{
			Name: common.PathJoin,
			Func: wsHandler(walletExt, joinRequestHandler),
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
		return
	}
}

func authenticateRequestHandler(walletExt *walletextension.WalletExtension, userConn userconn.UserConn) {
	// check if the text is well-formed and extract userID and address
	body, err := userConn.ReadRequest()
	if err != nil {
		return
	}

	var reqJSONMap map[string]string
	err = json.Unmarshal(body, &reqJSONMap)
	if err != nil {
		userConn.HandleError(fmt.Sprintf("could not unmarshal viewing key and signature from client to JSON: %s", err))
		return
	}

	accAddress := gethcommon.HexToAddress(reqJSONMap[common.JSONKeyAddress])

	signature, err := hex.DecodeString(reqJSONMap[common.JSONKeySignature][2:])
	if err != nil {
		userConn.HandleError(fmt.Sprintf("could not decode signature from client to hex: %s", err))
		return
	}

	message, ok := reqJSONMap[common.JSONKeyMessage]
	if !ok || message == "" {
		userConn.HandleError("message not found in the request")
		return
	}

	userID, err := getUser(userConn.ReadRequestParams())
	if err != nil {
		userConn.HandleError("userID not found in the request")
		return
	}

	// check the userID corresponds to the one in text
	messageUserID := ""
	// todo @ziga ( check if messageAddress needs to be included in the request)
	// messageAddress := ""
	regex := regexp.MustCompile(`^Register\s(\w+)\sfor\s(\w+)$`)
	if regex.MatchString(message) {
		params := regex.FindStringSubmatch(message)
		messageUserID = params[1]
		// messageAddress = params[2]
	} else {
		userConn.HandleError(fmt.Sprintf("Submitted message is not in the correct format: %s", message))
	}

	if userID != messageUserID || messageUserID == "" {
		userConn.HandleError(fmt.Sprintf("User in submitted message (%s) does not match user provided in the request (%s)", messageUserID, userID))
	}

	// check the signature if it corresponds to the address and is valid
	vk, found := walletExt.UnsignedVKs[accAddress]
	if !found {
		userConn.HandleError(fmt.Sprintf("no viewing key found to sign for acc=%s, please visit /join/ before sending signature", accAddress))
		return
	}

	// We transform the V from 27/28 to 0/1. This same change is made in Geth internals, for legacy reasons to be able
	// to recover the address: https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L452-L459
	signature[64] -= 27

	messageHash := crypto.Keccak256Hash([]byte(message))
	signatureNoRecoverID := signature[:len(signature)-1]
	verified := crypto.VerifySignature(vk.PublicKey, messageHash.Bytes(), signatureNoRecoverID)

	if !verified {
		userConn.HandleError("signature verification was not successful for your account")
		return
	}

	// save the text+signature against the userID
	// todo: where should I save the text? In another column in the database?
	vk.SignedKey = signature
	err = walletExt.Storage.SaveUserVK(userID, vk, message)
	if err != nil {
		userConn.HandleError("error saving viewing key")
		return
	}
}

func joinRequestHandler(walletExt *walletextension.WalletExtension, userConn userconn.UserConn) {
	body, err := userConn.ReadRequest()
	if err != nil {
		return
	}
	var reqJSONMap map[string]string
	err = json.Unmarshal(body, &reqJSONMap)
	if err != nil {
		userConn.HandleError("could not unmarshal /join request")
		return
	}

	reqAddress, ok := reqJSONMap[common.JSONKeyAddress]
	if !ok || reqAddress == "" {
		userConn.HandleError("message not found in the request")
		return
	}

	viewingKeyPrivate, err := crypto.GenerateKey()
	if err != nil {
		userConn.HandleError(fmt.Sprintf("could not generate new keypair: %s", err))
		return
	}

	viewingPublicKeyBytes := crypto.CompressPubkey(&viewingKeyPrivate.PublicKey)
	viewingPrivateKeyEcies := ecies.ImportECDSA(viewingKeyPrivate)
	accAddress := gethcommon.HexToAddress(reqAddress)
	// todo (@ziga) remove unsigedVKs and do everything with the database
	walletExt.UnsignedVKs[accAddress] = &rpc.ViewingKey{
		Account:    &accAddress,
		PrivateKey: viewingPrivateKeyEcies,
		PublicKey:  viewingPublicKeyBytes,
		SignedKey:  nil, // we await a signature from the user before we can set up the EncRPCClient
	}

	userID := crypto.Keccak256Hash(viewingPublicKeyBytes)
	err = userConn.WriteResponse(userID.Bytes())
	if err != nil {
		return
	}
}
