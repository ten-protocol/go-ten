package walletextension

import (
	"context"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-kit/kit/transport/http/jsonrpc"
	"github.com/obscuronet/go-obscuro/go/rpcclientlib"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
)

const (
	pathRoot               = "/"
	PathReady              = "/ready/"
	pathViewingKeys        = "/viewingkeys/"
	PathGenerateViewingKey = "/generateviewingkey/"
	PathSubmitViewingKey   = "/submitviewingkey/"
	staticDir              = "static"

	reqJSONKeyID        = "id"
	reqJSONKeyMethod    = "method"
	reqJSONKeyParams    = "params"
	ReqJSONKeyAddress   = "address"
	ReqJSONKeySignature = "signature"
	resJSONKeyID        = "id"
	resJSONKeyRPCVer    = "jsonrpc"
	RespJSONKeyResult   = "result"
	httpCodeErr         = 500

	// eth_call related constants
	reqCallJSONKeyFrom  = "from"
	reqCallJSONKeyData  = "data"
	ethCallPaddedArgLen = 64
	ethCallAddrPadding  = "000000000000000000000000"

	// CORS-related constants.
	corsAllowOrigin  = "Access-Control-Allow-Origin"
	originAll        = "*"
	corsAllowMethods = "Access-Control-Allow-Methods"
	reqOptions       = "OPTIONS"
	corsAllowHeaders = "Access-Control-Allow-Headers"
	corsHeaders      = "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"

	// EnclavePublicKeyHex is the public key of the enclave.
	// TODO - Retrieve this key from the management contract instead.
	enclavePublicKeyHex = "034d3b7e63a8bcd532ee3d1d6ecad9d67fca7821981a044551f0f0cbec74d0bc5e"
)

//go:embed static
var staticFiles embed.FS

// WalletExtension is a server that handles the management of viewing keys and the forwarding of Ethereum JSON-RPC requests.
type WalletExtension struct {
	enclavePublicKey *ecies.PublicKey // The public key used to encrypt requests for the enclave.
	hostAddr         string           // The address on which the Obscuro host can be reached.
	hostClient       *rpcclientlib.ViewingKeyClient
	viewingKeyRepo   *ViewingKeyRepository // stores accounts and viewing keys that are using the wallet extension
	server           *http.Server
}

type rpcRequest struct {
	id     interface{} // can be string or int
	method string
	params []interface{}
}

func NewWalletExtension(config Config) *WalletExtension {
	enclPubECDSA, err := crypto.DecompressPubkey(common.Hex2Bytes(enclavePublicKeyHex))
	if err != nil {
		panic(err)
	}
	enclavePublicKey := ecies.ImportECDSAPublic(enclPubECDSA)

	setLogs(config.LogPath)

	rpcClient, err := rpcclientlib.NewNetworkClient(config.NodeRPCHTTPAddress)
	if err != nil {
		panic(err)
	}
	vkRepo := NewViewingKeyRepository()
	client, err := rpcclientlib.NewViewingKeyClient(rpcClient, vkRepo)
	if err != nil {
		panic(err)
	}
	return &WalletExtension{
		enclavePublicKey: enclavePublicKey,
		hostAddr:         config.NodeRPCWebsocketAddress,
		hostClient:       client,
		viewingKeyRepo:   vkRepo,
	}
}

// Serve listens for and serves Ethereum JSON-RPC requests and viewing-key generation requests.
func (we *WalletExtension) Serve(hostAndPort string) {
	serveMux := http.NewServeMux()

	// Handles Ethereum JSON-RPC requests received over HTTP.
	serveMux.HandleFunc(pathRoot, we.handleHTTPEthJSON)
	serveMux.HandleFunc(PathReady, we.handleReady)
	serveMux.HandleFunc(PathGenerateViewingKey, we.handleGenerateViewingKey)
	serveMux.HandleFunc(PathSubmitViewingKey, we.handleSubmitViewingKey)

	// Serves the web assets for the management of viewing keys.
	noPrefixStaticFiles, err := fs.Sub(staticFiles, staticDir)
	if err != nil {
		panic(fmt.Sprintf("could not serve static files. Cause: %s", err))
	}
	serveMux.Handle(pathViewingKeys, http.StripPrefix(pathViewingKeys, http.FileServer(http.FS(noPrefixStaticFiles))))

	we.server = &http.Server{Addr: hostAndPort, Handler: serveMux, ReadHeaderTimeout: 10 * time.Second}

	err = we.server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (we *WalletExtension) Shutdown() {
	if we.server != nil {
		err := we.server.Shutdown(context.Background())
		if err != nil {
			log.Warn("could not shut down wallet extension: %s\n", err)
		}
	}
}

// Sets the log file.
func setLogs(logPath string) {
	if logPath == "" {
		return
	}
	logFile, err := os.Create(logPath)
	if err != nil {
		panic(fmt.Sprintf("could not create log file. Cause: %s", err))
	}
	log.OutputToFile(logFile)
}

// Used to check whether the server is ready.
func (we *WalletExtension) handleReady(http.ResponseWriter, *http.Request) {}

// Encrypts Ethereum JSON-RPC request, forwards it to the Obscuro node over a websocket, and decrypts the response if needed.
func (we *WalletExtension) handleHTTPEthJSON(resp http.ResponseWriter, req *http.Request) {
	// We enable CORS, as required by some browsers (e.g. Firefox).
	resp.Header().Set(corsAllowOrigin, originAll)
	if (*req).Method == reqOptions {
		resp.Header().Set(corsAllowMethods, reqOptions)
		resp.Header().Set(corsAllowHeaders, corsHeaders)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not read JSON-RPC request body: %s", err))
		return
	}

	rpcReq, err := parseRequest(body)
	if err != nil {
		logAndSendErr(resp, err.Error())
		return
	}

	if rpcReq.method == rpcclientlib.RPCCall {
		// RPCCall is a sensitive method that requires a viewing key lookup but the 'from' field is not mandatory in geth
		//	and is often not included from metamask etc. So we ensure it is populated here.
		rpcReq.params, err = we.setFromFieldIfMissing(rpcReq.params)
		if err != nil {
			logAndSendErr(resp, fmt.Sprintf("failed to set eth_call `from` field if it was missing - %s", err))
			return
		}
	}

	var rpcResp interface{}
	err = we.hostClient.Call(&rpcResp, rpcReq.method, rpcReq.params...)
	if err != nil {
		// if err was for a nil response then we will return an RPC result of null to the caller (this is a valid "not-found" response for some methods)
		if !errors.Is(err, rpcclientlib.ErrNilResponse) {
			logAndSendErr(resp, fmt.Sprintf("rpc request failed: %s", err))
			return
		}
	}

	respMap := make(map[string]interface{})
	respMap[resJSONKeyID] = rpcReq.id
	respMap[resJSONKeyRPCVer] = jsonrpc.Version
	respMap[RespJSONKeyResult] = rpcResp

	// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-658.md
	// TODO fix this upstream on the decode
	if result, found := respMap["result"]; found { //nolint
		if resultMap, ok := result.(map[string]interface{}); ok {
			if val, foundRoot := resultMap["root"]; foundRoot {
				if val == "0x" {
					respMap["result"].(map[string]interface{})["root"] = nil
				}
			}
		}
	}

	rpcRespToSend, err := json.Marshal(respMap)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("failed to remarshal RPC response to return to caller: %s", err))
	}
	log.Info("Forwarding %s response from Obscuro node: %s", rpcReq.method, rpcRespToSend)

	// We write the response to the client.
	_, err = resp.Write(rpcRespToSend)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not write JSON-RPC response: %s", err))
		return
	}
}

func parseRequest(body []byte) (*rpcRequest, error) {
	// We unmarshall the JSON request
	var reqJSONMap map[string]json.RawMessage
	err := json.Unmarshal(body, &reqJSONMap)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal JSON-RPC request body to JSON: %s. "+
			"If you're trying to generate a viewing key, visit %s", err, pathViewingKeys)
	}

	var reqID interface{}
	err = json.Unmarshal(reqJSONMap[reqJSONKeyID], &reqID)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal id from JSON-RPC request body: %w", err)
	}
	var method string
	err = json.Unmarshal(reqJSONMap[reqJSONKeyMethod], &method)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal method string from JSON-RPC request body: %w", err)
	}
	log.Info("Received %s request from wallet: %s", method, body)

	// we extract the params into a JSON list
	var params []interface{}
	err = json.Unmarshal(reqJSONMap[reqJSONKeyParams], &params)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal params list from JSON-RPC request body: %w", err)
	}

	return &rpcRequest{
		id:     reqID,
		method: method,
		params: params,
	}, nil
}

// Generates a new viewing key.
func (we *WalletExtension) handleGenerateViewingKey(resp http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not read viewing key and signature from client: %s", err))
		return
	}

	var reqJSONMap map[string]string
	err = json.Unmarshal(body, &reqJSONMap)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not unmarshal viewing key and signature from client to JSON: %s", err))
		return
	}

	viewingKeyPrivate, err := crypto.GenerateKey()
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not generate new keypair: %s", err))
		return
	}
	viewingPublicKeyBytes := crypto.CompressPubkey(&viewingKeyPrivate.PublicKey)
	viewingPrivateKeyEcies := ecies.ImportECDSA(viewingKeyPrivate)
	we.viewingKeyRepo.SetViewingKey(common.HexToAddress(reqJSONMap[ReqJSONKeyAddress]), viewingPrivateKeyEcies, viewingPublicKeyBytes)

	// We return the hex of the viewing key's public key for MetaMask to sign over.
	viewingKeyBytes := crypto.CompressPubkey(&viewingKeyPrivate.PublicKey)
	viewingKeyHex := hex.EncodeToString(viewingKeyBytes)
	_, err = resp.Write([]byte(viewingKeyHex))
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not return viewing key public key hex to client: %s", err))
		return
	}
}

// Submits the viewing key and signed bytes to the enclave.
func (we *WalletExtension) handleSubmitViewingKey(resp http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not read viewing key and signature from client: %s", err))
		return
	}

	var reqJSONMap map[string]string
	err = json.Unmarshal(body, &reqJSONMap)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not unmarshal viewing key and signature from client to JSON: %s", err))
		return
	}

	//  We drop the leading "0x".
	signature, err := hex.DecodeString(reqJSONMap[ReqJSONKeySignature][2:])
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not decode signature from client to hex: %s", err))
		return
	}

	// We transform the V from 27/28 to 0/1. This same change is made in Geth internals, for legacy reasons to be able
	// to recover the address: https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L452-L459
	signature[64] -= 27

	publicViewingKey := we.viewingKeyRepo.viewingKeysPublic[common.HexToAddress(reqJSONMap[ReqJSONKeyAddress])]
	// TODO: Store signatures to be able to resubmit keys if they are evicted by the node?
	// We return the hex of the viewing key's public key for MetaMask to sign over.
	err = we.hostClient.RegisterViewingKeyWithEnclave(publicViewingKey, signature)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("RPC request to register viewing key failed: %s", err))
		return
	}
}

// The enclave requires the `from` field to be set so that it can encrypt the response, but sources like MetaMask often
// don't set it. So we check whether it's present; if absent, we walk through the arguments in the request's `data`
// field, and if any of the arguments match our viewing key address, we set the `from` field to that address.
func (we *WalletExtension) setFromFieldIfMissing(args []interface{}) ([]interface{}, error) {
	callParams, err := parseCallParams(args)
	if err != nil {
		return nil, fmt.Errorf("could not parse eth_call params. Cause: %w", err)
	}

	// We only modify `eth_call` requests where the `from` field is not set.
	if callParams[reqCallJSONKeyFrom] != nil {
		return args, nil
	}

	fromAddress, err := we.viewingKeyRepo.suggestFromAddressForEthCall(callParams)
	if err != nil {
		return nil, err
	}

	callParams[reqCallJSONKeyFrom] = fromAddress
	args[0] = callParams
	return args, nil
}

// Parses the eth_call params into a map.
func parseCallParams(args []interface{}) (map[string]interface{}, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("expected eth_call params to have a 'from' field but no params found")
	}

	callParams, ok := args[0].(map[string]interface{})
	if !ok {
		callParamsJSON, ok := args[0].([]byte)
		if !ok {
			return nil, fmt.Errorf("expected eth_call first param to be a map or json encoded bytes but "+
				"was %t", args[0])
		}

		err := json.Unmarshal(callParamsJSON, &callParams)
		if err != nil {
			return nil, fmt.Errorf("expected eth_call first param to be a map or json encoded bytes, "+
				"failed to decode: %w", err)
		}
	}

	return callParams, nil
}

// Extracts the arguments from the request's `data` field. If any of them, after removing padding, match the viewing
// key address, we return that address. Otherwise, we return nil.
func searchDataFieldForFrom(callParams map[string]interface{}, viewingKeysPrivate map[common.Address]*ecies.PrivateKey) (*common.Address, error) {
	// We ensure that the `data` field is present.
	data := callParams[reqCallJSONKeyData]
	if data == nil {
		return nil, fmt.Errorf("eth_call request did not have its `data` field set")
	}
	dataString, ok := data.(string)
	if !ok {
		return nil, fmt.Errorf("eth_call request's `data` field was not of the expected type `string`")
	}

	// We check that the data field is long enough before removing the leading "0x" (1 bytes/2 chars) and the method ID
	// (4 bytes/8 chars).
	if len(dataString) < 10 {
		return nil, nil //nolint:nilnil
	}
	dataString = dataString[10:]

	// We split up the arguments in the `data` field.
	var dataArgs []string
	for i := 0; i < len(dataString); i += ethCallPaddedArgLen {
		if i+ethCallPaddedArgLen > len(dataString) {
			break
		}
		dataArgs = append(dataArgs, dataString[i:i+ethCallPaddedArgLen])
	}

	// We iterate over the arguments, looking for an argument that matches the viewing key address. If we find one, we
	// set the `from` field to that address.
	for _, dataArg := range dataArgs {
		// If the argument doesn't have the correct padding, it's not an address.
		if !strings.HasPrefix(dataArg, ethCallAddrPadding) {
			continue
		}

		maybeAddress := common.HexToAddress(dataArg[len(ethCallAddrPadding):])
		if _, ok := viewingKeysPrivate[maybeAddress]; ok {
			return &maybeAddress, nil
		}
	}

	return nil, nil //nolint:nilnil
}

// Logs the error message and sends it as an HTTP error.
func logAndSendErr(resp http.ResponseWriter, msg string) {
	log.Error(msg)
	fmt.Println(msg)
	http.Error(resp, msg, httpCodeErr)
}

// Config contains the configuration required by the WalletExtension.
type Config struct {
	WalletExtensionPort     int
	NodeRPCHTTPAddress      string
	NodeRPCWebsocketAddress string
	LogPath                 string
}
