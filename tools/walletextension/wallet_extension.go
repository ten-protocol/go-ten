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
	enclavePublicKey *ecies.PublicKey                              // The public key used to encrypt requests for the enclave.
	hostAddr         string                                        // The address on which the Obscuro host can be reached.
	accountClients   map[common.Address]*rpcclientlib.EncRPCClient // an encrypted RPC client per registered account
	unsignedVKs      map[common.Address]*rpcclientlib.ViewingKey   // map temporarily holding VKs that have been generated but not yet signed
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

	return &WalletExtension{
		enclavePublicKey: enclavePublicKey,
		hostAddr:         config.NodeRPCHTTPAddress,
		accountClients:   make(map[common.Address]*rpcclientlib.EncRPCClient),
		unsignedVKs:      make(map[common.Address]*rpcclientlib.ViewingKey),
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

	var rpcResp interface{}
	// proxyRequest will find the correct client to proxy the request (or try them all if appropriate)
	err = proxyRequest(rpcReq, rpcResp, we.accountClients)

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

func executeCall(client *rpcclientlib.EncRPCClient, req *rpcRequest, resp *interface{}) error {
	var err error
	if req.method == rpcclientlib.RPCCall {
		// RPCCall is a sensitive method that requires a viewing key lookup but the 'from' field is not mandatory in geth
		//	and is often not included from metamask etc. So we ensure it is populated here.
		account := client.Account()
		req.params, err = setCallFromFieldIfMissing(req.params, *account)
		if err != nil {
			return err
		}
	}

	return client.Call(resp, req.method, req.params...)
}

func parseRequest(body []byte) (*rpcRequest, error) {
	// We unmarshal the JSON request
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
	accAddress := common.HexToAddress(reqJSONMap[ReqJSONKeyAddress])
	we.unsignedVKs[accAddress] = &rpcclientlib.ViewingKey{
		Account:    &accAddress,
		PrivateKey: viewingPrivateKeyEcies,
		PublicKey:  viewingPublicKeyBytes,
		SignedKey:  nil, // we await a signature from the user before we can setup the EncRPCClient
	}

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
	accAddress := common.HexToAddress(reqJSONMap[ReqJSONKeyAddress])
	vk, found := we.unsignedVKs[accAddress]
	if !found {
		logAndSendErr(resp, fmt.Sprintf("no viewing key found to sign for acc=%s, please call %s to generate key before sending signature", accAddress, PathGenerateViewingKey))
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

	vk.SignedKey = signature
	// create an encrypted RPC client with the signed VK and register it with the enclave
	client, err := rpcclientlib.NewEncNetworkClient(we.hostAddr, vk)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("failed to create encrypted RPC client for acc=%s - %s", accAddress, err))
	}
	we.accountClients[accAddress] = client

	// finally we remove the VK from the pending 'unsigned VKs' map now the client has been created
	delete(we.unsignedVKs, accAddress)
}

// The enclave requires the `from` field to be set so that it can encrypt the response, but sources like MetaMask often
// don't set it. So we check whether it's present; if absent, we walk through the arguments in the request's `data`
// field, and if any of the arguments match our viewing key address, we set the `from` field to that address.
func setCallFromFieldIfMissing(args []interface{}, account common.Address) ([]interface{}, error) {
	callParams, err := parseParams(args)
	if err != nil {
		return nil, fmt.Errorf("could not parse eth_call params. Cause: %w", err)
	}

	// We only modify `eth_call` requests where the `from` field is not set.
	if callParams[reqJSONKeyFrom] != nil {
		return args, nil
	}

	callParams[reqJSONKeyFrom] = account
	args[0] = callParams
	return args, nil
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
