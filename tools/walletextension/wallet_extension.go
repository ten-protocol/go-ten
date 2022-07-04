package walletextension

import (
	"context"
	"crypto/rand"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"

	"github.com/ethereum/go-ethereum/common"

	"github.com/obscuronet/obscuro-playground/go/rpcclientlib"

	"github.com/gorilla/websocket"

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

	reqJSONKeyMethod          = "method"
	reqJSONKeyParams          = "params"
	reqJSONKeyFrom            = "from"
	ReqJSONMethodGetBalance   = "eth_getBalance"
	ReqJSONMethodCall         = "eth_call"
	ReqJSONMethodGetTxReceipt = "eth_getTransactionReceipt"
	ReqJSONMethodSendRawTx    = "eth_sendRawTransaction"
	respJSONKeyErr            = "error"
	respJSONKeyMsg            = "message"
	RespJSONKeyResult         = "result"
	httpCodeErr               = 500

	// CORS-related constants.
	corsAllowOrigin  = "Access-Control-Allow-Origin"
	originAll        = "*"
	corsAllowMethods = "Access-Control-Allow-Methods"
	reqOptions       = "OPTIONS"
	corsAllowHeaders = "Access-Control-Allow-Headers"
	corsHeaders      = "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"

	Localhost         = "127.0.0.1"
	websocketProtocol = "ws://"

	// EnclavePublicKeyHex is the public key of the enclave.
	// TODO - Retrieve this key from the management contract instead.
	enclavePublicKeyHex = "034d3b7e63a8bcd532ee3d1d6ecad9d67fca7821981a044551f0f0cbec74d0bc5e"

	// ViewingKeySignedMsgPrefix is the prefix added when signing the viewing key in MetaMask using the personal_sign
	// API. Why is this needed? MetaMask has a security feature whereby if you ask it to sign something that looks like
	// a transaction using the personal_sign API, it modifies the data being signed. The goal is to prevent hackers
	// from asking a visitor to their website to personal_sign something that is actually a malicious transaction (e.g.
	// theft of funds). By adding a prefix, the viewing key bytes no longer looks like a transaction hash, and thus get
	// signed as-is.
	ViewingKeySignedMsgPrefix = "vk"
)

//go:embed static
var staticFiles embed.FS

// TODO - Display error in browser if Metamask is not enabled (i.e. `ethereum` object is not available in-browser).

// WalletExtension is a server that handles the management of viewing keys and the forwarding of Ethereum JSON-RPC requests.
type WalletExtension struct {
	enclavePublicKey *ecies.PublicKey // The public key used to encrypt requests for the enclave.
	hostAddr         string           // The address on which the Obscuro host can be reached.
	hostClient       rpcclientlib.Client
	// TODO - Support multiple viewing keys. This will require the enclave to attach metadata on encrypted results
	//  to indicate which viewing key they were encrypted with.
	viewingPublicKeyBytes  []byte
	viewingPrivateKeyEcies *ecies.PrivateKey
	// The address associated with the last viewing key submitted. Used to set missing `from` fields in `eth_call` requests.
	viewingPublicKeyAddress common.Address
	server                  *http.Server
}

func NewWalletExtension(config Config) *WalletExtension {
	enclavePublicKey, err := crypto.DecompressPubkey(common.Hex2Bytes(enclavePublicKeyHex))
	if err != nil {
		panic(err)
	}

	return &WalletExtension{
		enclavePublicKey: ecies.ImportECDSAPublic(enclavePublicKey),
		hostAddr:         config.NodeRPCWebsocketAddress,
		hostClient:       rpcclientlib.NewClient(config.NodeRPCHTTPAddress),
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

	we.server = &http.Server{Addr: hostAndPort, Handler: serveMux}

	err = we.server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (we *WalletExtension) Shutdown() {
	if we.server != nil {
		err := we.server.Shutdown(context.Background())
		if err != nil {
			fmt.Printf("could not shut down wallet extension: %s", err)
		}
	}
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

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not read JSON-RPC request body: %s", err))
		return
	}

	// We unmarshall the JSON request.
	var reqJSONMap map[string]interface{}
	err = json.Unmarshal(body, &reqJSONMap)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not unmarshall JSON-RPC request body to JSON: %s. "+
			"If you're trying to generate a viewing key, visit %s", err, pathViewingKeys))
		return
	}
	method := reqJSONMap[reqJSONKeyMethod]
	fmt.Printf("Received %s request from wallet: %s\n", method, body)

	reqJSONMap, err = we.ensureCallsHaveFromField(method, reqJSONMap)
	if err != nil {
		logAndSendErr(resp, err.Error())
		return
	}

	// We encrypt the request's params with the enclave's public key if it's a sensitive request.
	maybeEncryptedBody, err := we.encryptParamsIfNeeded(body, method, reqJSONMap)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not encrypt request parameters: %s", err))
		return
	}

	// We forward the request on to the Obscuro node.
	nodeResp, err := forwardMsgOverWebsocket(websocketProtocol+we.hostAddr, maybeEncryptedBody)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("received error response when forwarding request to node at %s: %s", we.hostAddr, err))
		return
	}

	// We unmarshall the JSON response.
	var respJSONMap map[string]interface{}
	err = json.Unmarshal(nodeResp, &respJSONMap)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not unmarshall enclave response to JSON: %s", err))
		return
	}

	// We report any errors from the request.
	if respJSONMap[respJSONKeyErr] != nil {
		logAndSendErr(resp, respJSONMap[respJSONKeyErr].(map[string]interface{})[respJSONKeyMsg].(string))
		return
	}

	// We decrypt the result field if it's encrypted.
	maybeDecryptedRespJSONMap, err := we.decryptResponseIfNeeded(method, respJSONMap)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not decrypt response: %s", err))
		return
	}

	// We marshal the response to present to the client.
	clientResponse, err := json.Marshal(maybeDecryptedRespJSONMap)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not marshal JSON response to present to the client: %s", err))
		return
	}
	fmt.Printf("Received %s response from Obscuro node: %s\n", method, strings.TrimSpace(string(clientResponse)))

	// We write the response to the client.
	_, err = resp.Write(clientResponse)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not write JSON-RPC response: %s", err))
		return
	}
}

// If an `eth_call` request doesn't have a `from` field, we won't be able to encrypt the response. In that case, we use
// the viewing key address as the `from` field to allow encryption and decryption.
func (we *WalletExtension) ensureCallsHaveFromField(method interface{}, reqJSONMap map[string]interface{}) (map[string]interface{}, error) {
	if method != ReqJSONMethodCall {
		// We only modify `eth_call` requests.
		return reqJSONMap, nil
	}

	params, ok := reqJSONMap[reqJSONKeyParams].([]interface{})
	if !ok {
		return nil, fmt.Errorf("params for %s request were malformed", method)
	}
	txCallParams, ok := params[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("params for %s request were malformed", method)
	}

	if txCallParams[reqJSONKeyFrom] != nil {
		// We only modify `eth_call` requests where the `from` field is not set.
		return reqJSONMap, nil
	}

	if we.viewingPublicKeyAddress == (common.Address{}) {
		return nil, fmt.Errorf("could not add `from` field to `eth_call` request as no viewing key has been generated")
	}

	txCallParams[reqJSONKeyFrom] = we.viewingPublicKeyAddress.Hex()
	params[0] = txCallParams
	reqJSONMap[reqJSONKeyParams] = params

	return reqJSONMap, nil
}

// Generates a new viewing key.
func (we *WalletExtension) handleGenerateViewingKey(resp http.ResponseWriter, _ *http.Request) {
	viewingKeyPrivate, err := crypto.GenerateKey()
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not generate new keypair: %s", err))
		return
	}
	we.viewingPublicKeyBytes = crypto.CompressPubkey(&viewingKeyPrivate.PublicKey)
	we.viewingPrivateKeyEcies = ecies.ImportECDSA(viewingKeyPrivate)

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
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not read viewing key and signature from client: %s", err))
		return
	}

	var reqJSONMap map[string]string
	err = json.Unmarshal(body, &reqJSONMap)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not unmarshall viewing key and signature from client to JSON: %s", err))
		return
	}

	// We drop the leading "0x", and transform the V from 27/28 to 0/1.
	signature, err := hex.DecodeString(reqJSONMap["signature"][2:])
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not decode signature from client to hex: %s", err))
		return
	}
	signature[64] -= 27

	// We recover the public key address.
	msgToSign := ViewingKeySignedMsgPrefix + hex.EncodeToString(we.viewingPublicKeyBytes)
	recoveredPublicKey, err := crypto.SigToPub(accounts.TextHash([]byte(msgToSign)), signature)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not recover public key from signature: %s", err))
		return
	}
	we.viewingPublicKeyAddress = crypto.PubkeyToAddress(*recoveredPublicKey)

	// We encrypt the viewing key bytes.
	encryptedViewingKeyBytes, err := ecies.Encrypt(rand.Reader, we.enclavePublicKey, we.viewingPublicKeyBytes, nil, nil)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not encrypt viewing key with enclave public key: %s", err))
		return
	}

	var rpcErr error
	err = we.hostClient.Call(&rpcErr, rpcclientlib.RPCAddViewingKey, encryptedViewingKeyBytes, signature)
	if err != nil {
		logAndSendErr(resp, fmt.Sprintf("could not add viewing key: %s", err))
		return
	}
}

// Logs the error message and sends it as an HTTP error.
func logAndSendErr(resp http.ResponseWriter, msg string) {
	fmt.Println(msg)
	http.Error(resp, msg, httpCodeErr)
}

// Config contains the configuration required by the WalletExtension.
type Config struct {
	WalletExtensionPort     int
	NodeRPCHTTPAddress      string
	NodeRPCWebsocketAddress string
}

func forwardMsgOverWebsocket(url string, msg []byte) ([]byte, error) {
	connection, resp, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	defer connection.Close()
	defer resp.Body.Close()

	err = connection.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		return nil, err
	}

	_, message, err := connection.ReadMessage()
	if err != nil {
		return nil, err
	}
	return message, nil
}

// Encrypts the request's params with the enclave public key if the request is sensitive.
func (we *WalletExtension) encryptParamsIfNeeded(body []byte, method interface{}, reqJSONMap map[string]interface{}) ([]byte, error) {
	if !isSensitive(method) {
		return body, nil
	}

	fmt.Println("🔒 Encrypting request from wallet with enclave public key.")
	params := reqJSONMap[reqJSONKeyParams]
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request params to JSON for encryption: %w", err)
	}
	encryptedParams, err := ecies.Encrypt(rand.Reader, we.enclavePublicKey, paramsJSON, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt request params with enclave public key: %w", err)
	}
	reqJSONMap[reqJSONKeyParams] = []interface{}{encryptedParams}
	body, err = json.Marshal(reqJSONMap)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request with encrypted params to JSON: %w", err)
	}

	return body, nil
}

func (we *WalletExtension) decryptResponseIfNeeded(method interface{}, respJSONMap map[string]interface{}) (map[string]interface{}, error) {
	if !isSensitive(method) {
		return respJSONMap, nil
	}

	if we.viewingPrivateKeyEcies == nil {
		return nil, fmt.Errorf("could not decrypt enclave response as no viewing key has been created")
	}

	fmt.Printf("🔐 Decrypting %s response from Obscuro node with viewing key.\n", method)
	encryptedResult := common.Hex2Bytes(respJSONMap[RespJSONKeyResult].(string))
	decryptedResult, err := we.viewingPrivateKeyEcies.Decrypt(encryptedResult, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("could not decrypt enclave response with viewing key: %w", err)
	}

	processedResult, err := processDecryptedResult(decryptedResult, method)
	if err != nil {
		return nil, fmt.Errorf("could not process decrypted enclave response: %w", err)
	}
	respJSONMap[RespJSONKeyResult] = processedResult

	return respJSONMap, nil
}

// Indicates whether the RPC method's requests and responses should be encrypted.
func isSensitive(method interface{}) bool {
	return method == ReqJSONMethodGetBalance || method == ReqJSONMethodCall || method == ReqJSONMethodGetTxReceipt || method == ReqJSONMethodSendRawTx
}

// Converts the decrypted result to its correct JSON representation.
func processDecryptedResult(decryptedResult []byte, method interface{}) (interface{}, error) {
	// This method returns a JSON map, rather than a string.
	if method == ReqJSONMethodGetTxReceipt {
		fields := map[string]interface{}{}
		err := json.Unmarshal(decryptedResult, &fields)
		if err != nil {
			return nil, err
		}
		return fields, nil
	}

	return string(decryptedResult), nil
}
