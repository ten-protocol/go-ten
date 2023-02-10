package walletextension

import (
	"context"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"sync/atomic"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	gethlog "github.com/ethereum/go-ethereum/log"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/go-kit/kit/transport/http/jsonrpc"
	"github.com/obscuronet/go-obscuro/go/common/httputil"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/rpc"
	"github.com/obscuronet/go-obscuro/tools/walletextension/accountmanager"
	"github.com/obscuronet/go-obscuro/tools/walletextension/common"
	"github.com/obscuronet/go-obscuro/tools/walletextension/persistence"
	"github.com/obscuronet/go-obscuro/tools/walletextension/userconn"
)

const (
	pathRoot               = "/"
	PathReady              = "/ready/"
	pathViewingKeys        = "/viewingkeys/"
	PathGenerateViewingKey = "/generateviewingkey/"
	PathSubmitViewingKey   = "/submitviewingkey/"
	staticDir              = "static"
	wsProtocol             = "ws://"

	successMsg = "success"
)

var ErrSubscribeFailHTTP = fmt.Sprintf("received an %s request but the connection does not support subscriptions", rpc.Subscribe)

//go:embed static
var staticFiles embed.FS

// WalletExtension is a server that handles the management of viewing keys and the forwarding of Ethereum JSON-RPC requests.
type WalletExtension struct {
	hostAddr           string // The address on which the Obscuro host can be reached.
	accountManager     accountmanager.AccountManager
	unsignedVKs        map[gethcommon.Address]*rpc.ViewingKey // Map temporarily holding VKs that have been generated but not yet signed
	serverHTTPShutdown func(ctx context.Context) error
	serverWSShutdown   func(ctx context.Context) error
	persistence        *persistence.Persistence
	logger             gethlog.Logger
	isShutDown         atomicBool
}

type atomicBool int32

func (b *atomicBool) isSet() bool { return atomic.LoadInt32((*int32)(b)) != 0 }
func (b *atomicBool) setTrue()    { atomic.StoreInt32((*int32)(b), 1) }

func NewWalletExtension(config Config, logger gethlog.Logger) *WalletExtension {
	unauthedClient, err := rpc.NewNetworkClient(wsProtocol + config.NodeRPCWebsocketAddress)
	if err != nil {
		logger.Crit("unable to create temporary client for request ", log.ErrKey, err)
	}

	walletExtension := &WalletExtension{
		hostAddr:       wsProtocol + config.NodeRPCWebsocketAddress,
		unsignedVKs:    make(map[gethcommon.Address]*rpc.ViewingKey),
		accountManager: accountmanager.NewAccountManager(unauthedClient, logger),
		persistence:    persistence.NewPersistence(config.NodeRPCWebsocketAddress, config.PersistencePathOverride, logger),
		logger:         logger,
	}

	// We reload the existing viewing keys from persistence.
	for accountAddr, viewingKey := range walletExtension.persistence.LoadViewingKeys() {
		// create an encrypted RPC client with the signed VK and register it with the enclave
		// TODO - Create the clients lazily, to reduce connections to the host.
		client, err := rpc.NewEncNetworkClient(walletExtension.hostAddr, viewingKey, logger)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to create encrypted RPC client for persisted account %s", accountAddr), log.ErrKey, err)
			continue
		}
		walletExtension.accountManager.AddClient(accountAddr, client)
	}

	return walletExtension
}

// Serve listens for and serves Ethereum JSON-RPC requests and viewing-key generation requests.
func (we *WalletExtension) Serve(host string, httpPort int, wsPort int) {
	httpServer := we.createHTTPServer(host, httpPort)
	wsServer := we.createWSServer(host, wsPort)

	go func() {
		err := wsServer.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	err := httpServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (we *WalletExtension) Shutdown() {
	we.isShutDown.setTrue()
	if we.serverHTTPShutdown != nil {
		err := we.serverHTTPShutdown(context.Background())
		if err != nil {
			we.logger.Warn("could not shut down wallet extension", log.ErrKey, err)
		}
	}

	if we.serverWSShutdown != nil {
		err := we.serverWSShutdown(context.Background())
		if err != nil {
			we.logger.Warn("could not shut down wallet extension", log.ErrKey, err)
		}
	}
}

func (we *WalletExtension) createHTTPServer(host string, httpPort int) *http.Server {
	serveMuxHTTP := http.NewServeMux()

	// Handles Ethereum JSON-RPC requests received over HTTP.
	serveMuxHTTP.HandleFunc(pathRoot, we.handleEthJSONHTTP)
	serveMuxHTTP.HandleFunc(PathReady, we.handleReady)
	serveMuxHTTP.HandleFunc(PathGenerateViewingKey, we.handleGenerateViewingKeyHTTP)
	serveMuxHTTP.HandleFunc(PathSubmitViewingKey, we.handleSubmitViewingKeyHTTP)

	// Serves the web assets for the management of viewing keys.
	noPrefixStaticFiles, err := fs.Sub(staticFiles, staticDir)
	if err != nil {
		panic(fmt.Sprintf("could not serve static files. Cause: %s", err))
	}
	serveMuxHTTP.Handle(pathViewingKeys, http.StripPrefix(pathViewingKeys, http.FileServer(http.FS(noPrefixStaticFiles))))

	server := &http.Server{Addr: fmt.Sprintf("%s:%d", host, httpPort), Handler: serveMuxHTTP, ReadHeaderTimeout: 10 * time.Second}
	we.serverHTTPShutdown = server.Shutdown
	return server
}

func (we *WalletExtension) createWSServer(host string, wsPort int) *http.Server {
	serveMuxWS := http.NewServeMux()

	// Handles Ethereum JSON-RPC requests received over websockets.
	serveMuxWS.HandleFunc(pathRoot, we.handleEthJSONWS)
	serveMuxWS.HandleFunc(PathReady, we.handleReady)
	serveMuxWS.HandleFunc(PathGenerateViewingKey, we.handleGenerateViewingKeyWS)
	serveMuxWS.HandleFunc(PathSubmitViewingKey, we.handleSubmitViewingKeyWS)

	server := &http.Server{Addr: fmt.Sprintf("%s:%d", host, wsPort), Handler: serveMuxWS, ReadHeaderTimeout: 10 * time.Second}
	we.serverWSShutdown = server.Shutdown
	return server
}

// Used to check whether the server is ready.
func (we *WalletExtension) handleReady(resp http.ResponseWriter, req *http.Request) {
	if httputil.EnableCORS(resp, req) {
		return
	}
}

func (we *WalletExtension) handleEthJSONHTTP(resp http.ResponseWriter, req *http.Request) {
	we.handleRequestHTTP(resp, req, we.handleEthJSON)
}

func (we *WalletExtension) handleEthJSONWS(resp http.ResponseWriter, req *http.Request) {
	we.handleRequestWS(resp, req, we.handleEthJSON)
}

func (we *WalletExtension) handleGenerateViewingKeyHTTP(resp http.ResponseWriter, req *http.Request) {
	we.handleRequestHTTP(resp, req, we.handleGenerateViewingKey)
}

func (we *WalletExtension) handleGenerateViewingKeyWS(resp http.ResponseWriter, req *http.Request) {
	we.handleRequestWS(resp, req, we.handleGenerateViewingKey)
}

func (we *WalletExtension) handleSubmitViewingKeyHTTP(resp http.ResponseWriter, req *http.Request) {
	we.handleRequestHTTP(resp, req, we.handleSubmitViewingKey)
}

func (we *WalletExtension) handleSubmitViewingKeyWS(resp http.ResponseWriter, req *http.Request) {
	we.handleRequestWS(resp, req, we.handleSubmitViewingKey)
}

// Creates an HTTP connection to handle the request.
func (we *WalletExtension) handleRequestHTTP(resp http.ResponseWriter, req *http.Request, fun func(conn userconn.UserConn)) {
	if we.isShutDown.isSet() {
		return
	}
	if httputil.EnableCORS(resp, req) {
		return
	}
	userConn := userconn.NewUserConnHTTP(resp, req, we.logger)
	fun(userConn)
}

// Creates a websocket connection to handle the request.
func (we *WalletExtension) handleRequestWS(resp http.ResponseWriter, req *http.Request, fun func(conn userconn.UserConn)) {
	if we.isShutDown.isSet() {
		return
	}
	userConn, err := userconn.NewUserConnWS(resp, req, we.logger)
	if err != nil {
		return
	}
	// We handle requests in a loop until the connection is closed on the client side.
	for !userConn.IsClosed() {
		fun(userConn)
	}
}

// Encrypts the Ethereum JSON-RPC request, forwards it to the Obscuro node over a websocket, and decrypts the response if needed.
func (we *WalletExtension) handleEthJSON(userConn userconn.UserConn) {
	body, err := userConn.ReadRequest()
	if err != nil {
		return
	}

	rpcReq, err := we.parseRequest(body)
	if err != nil {
		userConn.HandleError(err.Error())
		return
	}

	if rpcReq.Method == rpc.Subscribe && !userConn.SupportsSubscriptions() {
		userConn.HandleError(ErrSubscribeFailHTTP)
		return
	}

	respMap := make(map[string]interface{})
	// all responses must contain the request id. Both successful and unsuccessful.
	respMap[common.JSONKeyRPCVersion] = jsonrpc.Version
	respMap[common.JSONKeyID] = rpcReq.ID

	// proxyRequest will find the correct client to proxy the request (or try them all if appropriate)
	var rpcResp interface{}
	err = we.accountManager.ProxyRequest(rpcReq, &rpcResp, userConn)

	if err != nil && !errors.Is(err, rpc.ErrNilResponse) {
		createErrorResponse(respMap, err)
	} else if errors.Is(err, rpc.ErrNilResponse) {
		// if err was for a nil response then we will return an RPC result of null to the caller (this is a valid "not-found" response for some methods)
		respMap[common.JSONKeyResult] = nil
	} else {
		respMap[common.JSONKeyResult] = rpcResp

		// TODO fix this upstream on the decode
		// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-658.md
		adjustStateRoot(rpcResp, respMap)
	}

	rpcRespToSend, err := json.Marshal(respMap)
	if err != nil {
		userConn.HandleError(fmt.Sprintf("failed to remarshal RPC response to return to caller: %s", err))
		return
	}
	we.logger.Info(fmt.Sprintf("Forwarding %s response from Obscuro node: %s", rpcReq.Method, rpcRespToSend))

	err = userConn.WriteResponse(rpcRespToSend)
	if err != nil {
		return
	}
}

func createErrorResponse(respMap map[string]interface{}, err error) {
	errMap := make(map[string]interface{})
	respMap[common.JSONKeyErr] = errMap

	errMap[common.JSONKeyMessage] = err.Error()

	var e gethrpc.Error
	ok := errors.As(err, &e)
	if ok {
		errMap[common.JSONKeyCode] = e.ErrorCode()
	}

	var de gethrpc.DataError
	ok = errors.As(err, &de)
	if ok {
		errMap[common.JSONKeyData] = de.ErrorData()
	}
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

func (we *WalletExtension) parseRequest(body []byte) (*accountmanager.RPCRequest, error) {
	// We unmarshal the JSON request
	var reqJSONMap map[string]json.RawMessage
	err := json.Unmarshal(body, &reqJSONMap)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal JSON-RPC request body to JSON: %s. "+
			"If you're trying to generate a viewing key, visit %s", err, pathViewingKeys)
	}

	reqID := reqJSONMap[common.JSONKeyID]
	var method string
	err = json.Unmarshal(reqJSONMap[common.JSONKeyMethod], &method)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal method string from JSON-RPC request body: %w", err)
	}
	we.logger.Info(fmt.Sprintf("Received %s request from wallet: %s", method, body))

	// we extract the params into a JSON list
	var params []interface{}
	err = json.Unmarshal(reqJSONMap[common.JSONKeyParams], &params)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal params list from JSON-RPC request body: %w", err)
	}

	return &accountmanager.RPCRequest{
		ID:     reqID,
		Method: method,
		Params: params,
	}, nil
}

// Generates a new viewing key.
func (we *WalletExtension) handleGenerateViewingKey(userConn userconn.UserConn) {
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

	viewingKeyPrivate, err := crypto.GenerateKey()
	if err != nil {
		userConn.HandleError(fmt.Sprintf("could not generate new keypair: %s", err))
		return
	}
	viewingPublicKeyBytes := crypto.CompressPubkey(&viewingKeyPrivate.PublicKey)
	viewingPrivateKeyEcies := ecies.ImportECDSA(viewingKeyPrivate)
	accAddress := gethcommon.HexToAddress(reqJSONMap[common.JSONKeyAddress])
	we.unsignedVKs[accAddress] = &rpc.ViewingKey{
		Account:    &accAddress,
		PrivateKey: viewingPrivateKeyEcies,
		PublicKey:  viewingPublicKeyBytes,
		SignedKey:  nil, // we await a signature from the user before we can set up the EncRPCClient
	}

	// We return the hex of the viewing key's public key for MetaMask to sign over.
	viewingKeyBytes := crypto.CompressPubkey(&viewingKeyPrivate.PublicKey)
	viewingKeyHex := hex.EncodeToString(viewingKeyBytes)
	err = userConn.WriteResponse([]byte(viewingKeyHex))
	if err != nil {
		return
	}
}

// Submits the viewing key and signed bytes to the enclave.
func (we *WalletExtension) handleSubmitViewingKey(userConn userconn.UserConn) {
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
	vk, found := we.unsignedVKs[accAddress]
	if !found {
		userConn.HandleError(fmt.Sprintf("no viewing key found to sign for acc=%s, please call %s to generate key before sending signature", accAddress, PathGenerateViewingKey))
		return
	}

	//  We drop the leading "0x".
	signature, err := hex.DecodeString(reqJSONMap[common.JSONKeySignature][2:])
	if err != nil {
		userConn.HandleError(fmt.Sprintf("could not decode signature from client to hex: %s", err))
		return
	}

	// We transform the V from 27/28 to 0/1. This same change is made in Geth internals, for legacy reasons to be able
	// to recover the address: https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L452-L459
	signature[64] -= 27

	vk.SignedKey = signature
	// create an encrypted RPC client with the signed VK and register it with the enclave
	// TODO - Create the clients lazily, to reduce connections to the host.
	client, err := rpc.NewEncNetworkClient(we.hostAddr, vk, we.logger)
	if err != nil {
		userConn.HandleError(fmt.Sprintf("failed to create encrypted RPC client for account %s. Cause: %s", accAddress, err))
		return
	}
	we.accountManager.AddClient(accAddress, client)

	we.persistence.PersistViewingKey(vk)
	// finally we remove the VK from the pending 'unsigned VKs' map now the client has been created
	delete(we.unsignedVKs, accAddress)

	err = userConn.WriteResponse([]byte(successMsg))
	if err != nil {
		return
	}
}

// Config contains the configuration required by the WalletExtension.
type Config struct {
	WalletExtensionHost     string
	WalletExtensionPort     int
	WalletExtensionPortWS   int
	NodeRPCHTTPAddress      string
	NodeRPCWebsocketAddress string
	LogPath                 string
	PersistencePathOverride string // Overrides the persistence file location. Used in tests.
}
