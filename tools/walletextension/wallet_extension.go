package walletextension

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	pathViewingKeys        = "/viewingkeys/"
	pathGenerateViewingKey = "/generateviewingkey/"
	pathSubmitViewingKey   = "/submitviewingkey/"
	staticDir              = "./tools/walletextension/static"
)

// WalletExtension is a server that handles the management of viewing keys and the forwarding of Ethereum JSON-RPC requests.
type WalletExtension struct {
	enclavePrivateKey *ecdsa.PrivateKey
	// todo - support multiple viewing keys. this will require the enclave to attach metadata on encrypted results
	//  to indicate which key they were encrypted with
	viewingKeyPrivate *ecdsa.PrivateKey
	// todo - replace this channel with port-based communication with the enclave
	viewingKeyChannel chan<- ViewingKey
}

func NewWalletExtension(enclavePrivateKey *ecdsa.PrivateKey, viewingKeyChannel chan<- ViewingKey) *WalletExtension {
	return &WalletExtension{enclavePrivateKey: enclavePrivateKey, viewingKeyChannel: viewingKeyChannel}
}

func (we *WalletExtension) Serve(hostAndPort string) {
	serveMux := http.NewServeMux()

	// Handles Ethereum JSON-RPC requests received over HTTP.
	serveMux.HandleFunc(pathRoot, we.handleHttpEthJson)

	// Handles the management of viewing keys.
	serveMux.Handle(pathViewingKeys, http.StripPrefix(pathViewingKeys, http.FileServer(http.Dir(staticDir))))
	serveMux.HandleFunc(pathGenerateViewingKey, we.handleGenerateViewingKey)
	serveMux.HandleFunc(pathSubmitViewingKey, we.handleSubmitViewingKey)

	err := http.ListenAndServe(hostAndPort, serveMux)
	if err != nil {
		panic(err)
	}
}

// Encrypts Ethereum JSON-RPC request, forwards it to the Geth node over a websocket, and decrypts the response if needed.
func (we *WalletExtension) handleHttpEthJson(resp http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("could not read JSON-RPC request body: %v\n", err)
		return // todo - return error response
	}

	// We unmarshall the JSON request.
	var reqJsonMap map[string]interface{}
	err = json.Unmarshal(body, &reqJsonMap)
	if err != nil {
		fmt.Printf("could not unmarshall JSON-RPC request body to JSON: %v\n", err)
		return // todo - return error response
	}
	fmt.Println(fmt.Sprintf("Received request from wallet: %s", body))

	// We encrypt the JSON with the enclave's public key.
	fmt.Println("🔒 Encrypting request from wallet with enclave public key.")
	eciesPublicKey := ecies.ImportECDSAPublic(&we.enclavePrivateKey.PublicKey)
	encryptedBody, err := ecies.Encrypt(rand.Reader, eciesPublicKey, body, nil, nil)
	if err != nil {
		fmt.Printf("could not encrypt request with enclave public key: %v\n", err)
		return // todo - return error response
	}

	// We forward the request on to the Geth node.
	gethResp := forwardMsgOverWebsocket(obxFacadeWebsocketAddr, encryptedBody)

	// We decrypt the response if it's encrypted.
	method := reqJsonMap[reqJsonKeyMethod]
	if method == reqJsonMethodGetBalance || method == reqJsonMethodGetStorageAt {
		fmt.Println(fmt.Sprintf("🔐 Decrypting %s response from Geth node with viewing key.", method))
		eciesPrivateKey := ecies.ImportECDSA(we.viewingKeyPrivate)
		gethResp, err = eciesPrivateKey.Decrypt(gethResp, nil, nil)
		if err != nil {
			fmt.Printf("could not decrypt enclave response with viewing key: %v\n", err)
			return // todo - return error response
		}
	}

	// We unmarshall the JSON response.
	var respJsonMap map[string]interface{}
	err = json.Unmarshal(gethResp, &respJsonMap)
	if err != nil {
		fmt.Printf("could not unmarshall enclave response to JSON: %v\n", err)
		return // todo - return error response
	}
	fmt.Println(fmt.Sprintf("Received response from Geth node: %s", strings.TrimSpace(string(gethResp))))

	// We write the response to the client.
	_, err = resp.Write(gethResp)
	if err != nil {
		fmt.Printf("could not write JSON-RPC response: %v\n", err)
		return // todo - return error response
	}
}

// Generates a new viewing key.
func (we *WalletExtension) handleGenerateViewingKey(resp http.ResponseWriter, _ *http.Request) {
	viewingKeyPrivate, err := crypto.GenerateKey()
	if err != nil {
		fmt.Printf("could not generate new keypair: %v\n", err)
		return // todo - return error response
	}
	we.viewingKeyPrivate = viewingKeyPrivate

	// We return the hex of the viewing key's public key for MetaMask to sign over.
	viewingKeyBytes := crypto.CompressPubkey(&viewingKeyPrivate.PublicKey)
	viewingKeyHex := hex.EncodeToString(viewingKeyBytes)
	_, err = resp.Write([]byte(viewingKeyHex))
	if err != nil {
		fmt.Printf("could not return viewing key public key hex to client: %v\n", err)
		// todo - return error response
	}
}

// Submits the viewing key and signed bytes to the enclave.
func (we *WalletExtension) handleSubmitViewingKey(_ http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("could not read viewing key and signed bytes from client: %v\n", err)
		return // todo - return error response
	}

	var reqJsonMap map[string]interface{}
	err = json.Unmarshal(body, &reqJsonMap)
	if err != nil {
		fmt.Printf("could not unmarshall viewing key and signed bytes from client to JSON: %v\n", err)
		return // todo - return error response
	}
	signedBytes := []byte(reqJsonMap["signedBytes"].(string))

	viewingKey := ViewingKey{viewingKeyPublic: &we.viewingKeyPrivate.PublicKey, signedBytes: signedBytes}
	we.viewingKeyChannel <- viewingKey
}
