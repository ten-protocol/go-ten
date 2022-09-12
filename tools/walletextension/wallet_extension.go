package walletextension

import (
	"context"
	"embed"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/obscuronet/go-obscuro/tools/walletextension/readwriter"

	"github.com/obscuronet/go-obscuro/go/common/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-kit/kit/transport/http/jsonrpc"
	"github.com/obscuronet/go-obscuro/go/rpc"

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
	obscuroDirName         = ".obscuro"

	reqJSONKeyID        = "id"
	reqJSONKeyMethod    = "method"
	reqJSONKeyParams    = "params"
	ReqJSONKeyAddress   = "address"
	ReqJSONKeySignature = "signature"
	resJSONKeyID        = "id"
	resJSONKeyRPCVer    = "jsonrpc"
	RespJSONKeyResult   = "result"

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

	persistenceFileName      = "wallet_extension_persistence"
	persistenceNumComponents = 4
	persistenceIdxHost       = 0
	persistenceIdxAccount    = 1
	persistenceIdxViewingKey = 2
	persistenceIdxSignedKey  = 3
)

//go:embed static
var staticFiles embed.FS

// WalletExtension is a server that handles the management of viewing keys and the forwarding of Ethereum JSON-RPC requests.
type WalletExtension struct {
	enclavePublicKey *ecies.PublicKey // The public key used to encrypt requests for the enclave.
	hostAddr         string           // The address on which the Obscuro host can be reached.
	// TODO - Create two types of clients - WS clients, and HTTP clients - to not create WS clients unnecessarily.
	accountClients  map[common.Address]*rpc.EncRPCClient // An encrypted RPC client per registered account
	unauthedClient  rpc.Client                           // Unauthenticated client used for non-sensitive requests if no encrypted clients exist.
	unsignedVKs     map[common.Address]*rpc.ViewingKey   // Map temporarily holding VKs that have been generated but not yet signed
	serverHTTP      *http.Server
	serverWS        *http.Server
	persistencePath string // The path of the file used to store the submitted viewing keys
}

type rpcRequest struct {
	id     interface{} // can be string or int
	method string
	params []interface{}
}

// TODO - We default to websocket clients to the host, since these can handle any kind of connection to the wallet
//  extension (HTTP or WS). Consider creating the clients on-the-fly, based on the request type.

func NewWalletExtension(config Config) *WalletExtension {
	setUpLogs(config.LogPath)

	enclPubECDSA, err := crypto.DecompressPubkey(common.Hex2Bytes(enclavePublicKeyHex))
	if err != nil {
		log.Panic("%s", err)
	}
	enclavePublicKey := ecies.ImportECDSAPublic(enclPubECDSA)

	unauthedClient, err := rpc.NewNetworkClient(config.NodeRPCWebsocketAddress)
	if err != nil {
		log.Panic("unable to create temporary client for request - %s", err)
	}

	walletExtension := &WalletExtension{
		enclavePublicKey: enclavePublicKey,
		hostAddr:         config.NodeRPCWebsocketAddress,
		accountClients:   make(map[common.Address]*rpc.EncRPCClient),
		unsignedVKs:      make(map[common.Address]*rpc.ViewingKey),
		unauthedClient:   unauthedClient,
		persistencePath:  setUpPersistence(config.PersistencePathOverride),
	}

	// We reload the existing viewing keys from persistence.
	for accountAddr, viewingKey := range walletExtension.loadViewingKeys() {
		// create an encrypted RPC client with the signed VK and register it with the enclave
		// TODO - Create the clients lazily, to reduce connections to the host.
		client, err := rpc.NewEncNetworkClient(walletExtension.hostAddr, viewingKey)
		if err != nil {
			log.Error("failed to create encrypted RPC client for persisted account %s. Cause: %s", accountAddr, err)
			continue
		}
		walletExtension.accountClients[accountAddr] = client
	}

	return walletExtension
}

// Serve listens for and serves Ethereum JSON-RPC requests and viewing-key generation requests.
func (we *WalletExtension) Serve(host string, httpPort int, wsPort int) {
	we.createHTTPServer(host, httpPort)
	we.createWSServer(host, wsPort)

	go func() {
		err := we.serverWS.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	err := we.serverHTTP.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (we *WalletExtension) Shutdown() {
	if we.serverHTTP != nil {
		err := we.serverHTTP.Shutdown(context.Background())
		if err != nil {
			log.Warn("could not shut down wallet extension: %s\n", err)
		}
	}

	if we.serverWS != nil {
		err := we.serverWS.Shutdown(context.Background())
		if err != nil {
			log.Warn("could not shut down wallet extension: %s\n", err)
		}
	}
}

func (we *WalletExtension) createHTTPServer(host string, httpPort int) {
	serveMuxHTTP := http.NewServeMux()

	// Handles Ethereum JSON-RPC requests received over HTTP.
	serveMuxHTTP.HandleFunc(pathRoot, we.handleEthJSONHTTP)
	serveMuxHTTP.HandleFunc(PathReady, we.handleReady)
	serveMuxHTTP.HandleFunc(PathGenerateViewingKey, we.handleGenerateViewingKey)
	serveMuxHTTP.HandleFunc(PathSubmitViewingKey, we.handleSubmitViewingKey)

	// Serves the web assets for the management of viewing keys.
	noPrefixStaticFiles, err := fs.Sub(staticFiles, staticDir)
	if err != nil {
		panic(fmt.Sprintf("could not serve static files. Cause: %s", err))
	}
	serveMuxHTTP.Handle(pathViewingKeys, http.StripPrefix(pathViewingKeys, http.FileServer(http.FS(noPrefixStaticFiles))))

	we.serverHTTP = &http.Server{Addr: fmt.Sprintf("%s:%d", host, httpPort), Handler: serveMuxHTTP, ReadHeaderTimeout: 10 * time.Second}
}

func (we *WalletExtension) createWSServer(host string, wsPort int) {
	serveMuxWS := http.NewServeMux()
	serveMuxWS.HandleFunc(pathRoot, we.handleEthJSONWS)
	we.serverWS = &http.Server{Addr: fmt.Sprintf("%s:%d", host, wsPort), Handler: serveMuxWS, ReadHeaderTimeout: 10 * time.Second}
}

// Sets up the log file.
func setUpLogs(logPath string) {
	if logPath == "" {
		return
	}
	logFile, err := os.Create(logPath)
	if err != nil {
		panic(fmt.Sprintf("could not create log file. Cause: %s", err))
	}
	log.OutputToFile(logFile)
}

// Sets up the persistence file and returns its path. Defaults to the user's home directory if the path is empty.
func setUpPersistence(persistenceFilePath string) string {
	// We set the default if the persistence file is not overridden.
	if persistenceFilePath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic("cannot create persistence file as user's home directory is not defined")
		}
		obscuroDir := filepath.Join(homeDir, obscuroDirName)
		err = os.MkdirAll(obscuroDir, 0o777)
		if err != nil {
			panic(fmt.Sprintf("could not create %s directory in user's home directory", obscuroDirName))
		}

		persistenceFilePath = filepath.Join(obscuroDir, persistenceFileName)
	}

	_, err := os.OpenFile(persistenceFilePath, os.O_CREATE|os.O_RDONLY, 0o644)
	if err != nil {
		panic(fmt.Sprintf("could not create persistence file. Cause: %s", err))
	}

	return persistenceFilePath
}

// Used to check whether the server is ready.
func (we *WalletExtension) handleReady(http.ResponseWriter, *http.Request) {}

// Handles the Ethereum JSON-RPC request over HTTP.
func (we *WalletExtension) handleEthJSONHTTP(resp http.ResponseWriter, req *http.Request) {
	if we.enableCORS(resp, req) {
		return
	}
	readWriter := readwriter.NewHTTPReadWriter(resp, req)
	we.handleEthJSON(readWriter)
}

// Handles the Ethereum JSON-RPC request over websockets.
func (we *WalletExtension) handleEthJSONWS(resp http.ResponseWriter, req *http.Request) {
	readWriter, err := readwriter.NewWSReadWriter(resp, req)
	if err != nil {
		return
	}
	we.handleEthJSON(readWriter)
}

// Encrypts the Ethereum JSON-RPC request, forwards it to the Obscuro node over a websocket, and decrypts the response if needed.
func (we *WalletExtension) handleEthJSON(readWriter readwriter.ReadWriter) {
	body, err := readWriter.ReadRequest()
	if err != nil {
		readWriter.HandleError(err.Error())
		return
	}

	rpcReq, err := parseRequest(body)
	if err != nil {
		readWriter.HandleError(err.Error())
		return
	}

	if rpcReq.method == rpc.RPCSubscribe && !readWriter.SupportsSubscriptions() {
		readWriter.HandleError(fmt.Sprintf("received an %s request but the connection does not support subscriptions", rpc.RPCSubscribe))
	}

	var rpcResp interface{}
	// proxyRequest will find the correct client to proxy the request (or try them all if appropriate)
	err = proxyRequest(rpcReq, &rpcResp, we)
	if err != nil {
		// if err was for a nil response then we will return an RPC result of null to the caller (this is a valid "not-found" response for some methods)
		if !errors.Is(err, rpc.ErrNilResponse) {
			readWriter.HandleError(fmt.Sprintf("rpc request failed: %s", err))
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
		readWriter.HandleError(fmt.Sprintf("failed to remarshal RPC response to return to caller: %s", err))
		return
	}
	log.Info("Forwarding %s response from Obscuro node: %s", rpcReq.method, rpcRespToSend)

	err = readWriter.WriteResponse(rpcRespToSend)
	if err != nil {
		readWriter.HandleError(err.Error())
		return
	}
}

// Enables CORS, as required by some browsers (e.g. Firefox). Returns true if CORS was enabled.
func (we *WalletExtension) enableCORS(resp http.ResponseWriter, req *http.Request) bool {
	resp.Header().Set(corsAllowOrigin, originAll)
	if (*req).Method == reqOptions {
		resp.Header().Set(corsAllowMethods, reqOptions)
		resp.Header().Set(corsAllowHeaders, corsHeaders)
		return true
	}
	return false
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
	readWriter := readwriter.NewHTTPReadWriter(resp, req)

	body, err := readWriter.ReadRequest()
	if err != nil {
		readWriter.HandleError(err.Error())
		return
	}

	var reqJSONMap map[string]string
	err = json.Unmarshal(body, &reqJSONMap)
	if err != nil {
		readWriter.HandleError(fmt.Sprintf("could not unmarshal viewing key and signature from client to JSON: %s", err))
		return
	}

	viewingKeyPrivate, err := crypto.GenerateKey()
	if err != nil {
		readWriter.HandleError(fmt.Sprintf("could not generate new keypair: %s", err))
		return
	}
	viewingPublicKeyBytes := crypto.CompressPubkey(&viewingKeyPrivate.PublicKey)
	viewingPrivateKeyEcies := ecies.ImportECDSA(viewingKeyPrivate)
	accAddress := common.HexToAddress(reqJSONMap[ReqJSONKeyAddress])
	we.unsignedVKs[accAddress] = &rpc.ViewingKey{
		Account:    &accAddress,
		PrivateKey: viewingPrivateKeyEcies,
		PublicKey:  viewingPublicKeyBytes,
		SignedKey:  nil, // we await a signature from the user before we can set up the EncRPCClient
	}

	// We return the hex of the viewing key's public key for MetaMask to sign over.
	viewingKeyBytes := crypto.CompressPubkey(&viewingKeyPrivate.PublicKey)
	viewingKeyHex := hex.EncodeToString(viewingKeyBytes)
	err = readWriter.WriteResponse([]byte(viewingKeyHex))
	if err != nil {
		readWriter.HandleError(fmt.Sprintf("could not return viewing key public key hex to client: %s", err))
		return
	}
}

// Submits the viewing key and signed bytes to the enclave.
func (we *WalletExtension) handleSubmitViewingKey(resp http.ResponseWriter, req *http.Request) {
	readWriter := readwriter.NewHTTPReadWriter(resp, req)

	body, err := readWriter.ReadRequest()
	if err != nil {
		readWriter.HandleError(err.Error())
		return
	}

	var reqJSONMap map[string]string
	err = json.Unmarshal(body, &reqJSONMap)
	if err != nil {
		readWriter.HandleError(fmt.Sprintf("could not unmarshal viewing key and signature from client to JSON: %s", err))
		return
	}
	accAddress := common.HexToAddress(reqJSONMap[ReqJSONKeyAddress])
	vk, found := we.unsignedVKs[accAddress]
	if !found {
		readWriter.HandleError(fmt.Sprintf("no viewing key found to sign for acc=%s, please call %s to generate key before sending signature", accAddress, PathGenerateViewingKey))
		return
	}

	//  We drop the leading "0x".
	signature, err := hex.DecodeString(reqJSONMap[ReqJSONKeySignature][2:])
	if err != nil {
		readWriter.HandleError(fmt.Sprintf("could not decode signature from client to hex: %s", err))
		return
	}

	// We transform the V from 27/28 to 0/1. This same change is made in Geth internals, for legacy reasons to be able
	// to recover the address: https://github.com/ethereum/go-ethereum/blob/55599ee95d4151a2502465e0afc7c47bd1acba77/internal/ethapi/api.go#L452-L459
	signature[64] -= 27

	vk.SignedKey = signature
	// create an encrypted RPC client with the signed VK and register it with the enclave
	// TODO - Create the clients lazily, to reduce connections to the host.
	client, err := rpc.NewEncNetworkClient(we.hostAddr, vk)
	if err != nil {
		readWriter.HandleError(fmt.Sprintf("failed to create encrypted RPC client for account %s. Cause: %s", accAddress, err))
	}
	we.accountClients[accAddress] = client

	we.persistViewingKey(vk)
	// finally we remove the VK from the pending 'unsigned VKs' map now the client has been created
	delete(we.unsignedVKs, accAddress)
}

// The enclave requires the `from` field to be set so that it can encrypt the response, but sources like MetaMask often
// don't set it. So we check whether it's present; if absent, we walk through the arguments in the request's `data`
// field, and if any of the arguments match our viewing key address, we set the `from` field to that address.
// TODO - Move this method into multi_acc_helper.go.
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

// Stores a viewing key to disk.
// TODO - Move the persistence-related methods onto a separate struct.
func (we *WalletExtension) persistViewingKey(viewingKey *rpc.ViewingKey) {
	viewingPrivateKeyBytes := crypto.FromECDSA(viewingKey.PrivateKey.ExportECDSA())

	record := []string{
		we.hostAddr,
		viewingKey.Account.Hex(),
		// We encode the bytes as hex to ensure there are no unintentional line breaks to make parsing the file harder.
		hex.EncodeToString(viewingPrivateKeyBytes),
		hex.EncodeToString(viewingKey.SignedKey),
	}

	persistenceFile, err := os.OpenFile(we.persistencePath, os.O_APPEND|os.O_WRONLY, 0o644)
	defer persistenceFile.Close() //nolint:staticcheck
	if err != nil {
		log.Error("could not open persistence file. Cause: %s", err)
	}

	writer := csv.NewWriter(persistenceFile)
	defer writer.Flush()
	err = writer.Write(record)
	if err != nil {
		log.Error("failed to write viewing key to persistence file. Cause: %s", err)
	}
}

// Loads any viewing keys from disk. Viewing keys for other hosts are ignored.
func (we *WalletExtension) loadViewingKeys() map[common.Address]*rpc.ViewingKey {
	viewingKeys := make(map[common.Address]*rpc.ViewingKey)

	persistenceFile, err := os.OpenFile(we.persistencePath, os.O_RDONLY, 0o644)
	defer persistenceFile.Close() //nolint:staticcheck
	if err != nil {
		log.Error("could not open persistence file. Cause: %s", err)
	}

	reader := csv.NewReader(persistenceFile)
	records, err := reader.ReadAll()
	if err != nil {
		log.Error("could not read records from persistence file. Cause: %s", err)
	}

	for _, record := range records {
		// TODO - Determine strategy for invalid persistence entries - delete? Warn? Shutdown? For now, we log a warning.
		if len(record) != persistenceNumComponents {
			log.Warn("persistence file entry did not have expected number of components: %s", record)
			continue
		}

		hostAddr := record[persistenceIdxHost]
		if hostAddr != we.hostAddr {
			log.Info("skipping persistence file entry for another host. Current host is %s, entry was for %s", we.hostAddr, hostAddr)
			continue
		}

		account := common.HexToAddress(record[persistenceIdxAccount])
		viewingKeyPrivateHex := record[persistenceIdxViewingKey]
		viewingKeyPrivateBytes, err := hex.DecodeString(viewingKeyPrivateHex)
		if err != nil {
			log.Warn("could not decode the following viewing private key from hex in the persistence file: %s", viewingKeyPrivateHex)
			continue
		}
		viewingKeyPrivate, err := crypto.ToECDSA(viewingKeyPrivateBytes)
		if err != nil {
			log.Warn("could not convert the following viewing private key bytes to ECDSA in the persistence file: %s", viewingKeyPrivateHex)
			continue
		}
		signedKeyHex := record[persistenceIdxSignedKey]
		signedKey, err := hex.DecodeString(signedKeyHex)
		if err != nil {
			log.Warn("could not decode the following signed key from hex in the persistence file: %s", signedKeyHex)
			continue
		}

		viewingKey := rpc.ViewingKey{
			Account:    &account,
			PrivateKey: ecies.ImportECDSA(viewingKeyPrivate),
			PublicKey:  crypto.CompressPubkey(&viewingKeyPrivate.PublicKey),
			SignedKey:  signedKey,
		}
		viewingKeys[account] = &viewingKey
	}

	logReRegisteredViewingKeys(viewingKeys)

	return viewingKeys
}

// Logs and prints the accounts for which we are re-registering viewing keys.
func logReRegisteredViewingKeys(viewingKeys map[common.Address]*rpc.ViewingKey) {
	if len(viewingKeys) == 0 {
		return
	}

	var accounts []string //nolint:prealloc
	for account := range viewingKeys {
		accounts = append(accounts, account.Hex())
	}

	msg := fmt.Sprintf("Re-registering persisted viewing keys for the following addresses: %s",
		strings.Join(accounts, ", "))
	log.Info(msg)
	fmt.Println(msg)
}

// Config contains the configuration required by the WalletExtension.
type Config struct {
	WalletExtensionPort     int
	WalletExtensionPortWS   int
	NodeRPCHTTPAddress      string
	NodeRPCWebsocketAddress string
	LogPath                 string
	PersistencePathOverride string // Overrides the persistence file location. Used in tests.
}
