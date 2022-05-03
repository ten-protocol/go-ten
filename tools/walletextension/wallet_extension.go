package walletextension

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"io/ioutil"
	"net/http"
	"strings"
)

// todo - joel - encrypt with the address key on the return in the facade, and decrypt in the wallet extension
//  use a standard account key for now, since we don't have viewing keys

// WalletExtension is a server that handles the management of viewing keys and the forwarding of Ethereum JSON-RPC requests.
type WalletExtension struct {
	enclavePrivateKey *ecdsa.PrivateKey
}

func NewWalletExtension(enclavePrivateKey *ecdsa.PrivateKey) *WalletExtension {
	return &WalletExtension{enclavePrivateKey: enclavePrivateKey}
}

func (we WalletExtension) Serve(hostAndPort string) {
	serveMux := http.NewServeMux()

	// Handles Ethereum JSON-RPC requests received over HTTP.
	serveMux.HandleFunc("/", we.handleHttpEthJson)

	// Handles the management of viewing keys.
	serveMux.Handle("/viewingkeys/", http.StripPrefix("/viewingkeys/", http.FileServer(http.Dir("./tools/walletextension/static"))))
	serveMux.HandleFunc("/getViewingKey", we.handleGetViewingKey)

	err := http.ListenAndServe(hostAndPort, serveMux)
	if err != nil {
		panic(err)
	}
}

// Reads the Ethereum JSON-RPC request, and forwards it to the Geth node over a websocket.
func (we WalletExtension) handleHttpEthJson(resp http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
	}

	// We unmarshall the JSON request to inspect it.
	var reqJsonMap map[string]interface{}
	err = json.Unmarshal(body, &reqJsonMap)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fmt.Sprintf("Received request from wallet: %s", body))

	// We encrypt the JSON with the enclave's public key.
	fmt.Println("🔒 Encrypting request from wallet with enclave public key.")
	eciesPublicKey := ecies.ImportECDSAPublic(&we.enclavePrivateKey.PublicKey)
	encryptedBody, err := ecies.Encrypt(rand.Reader, eciesPublicKey, body, nil, nil)
	if err != nil {
		fmt.Println(err)
	}

	// We forward the request on to the Geth node.
	gethResp := forwardMsgOverWebsocket(obxFacadeWebsocketAddr, encryptedBody)

	// We decrypt the response if it's encrypted.
	method := reqJsonMap["method"]
	if method == "eth_getBalance" || method == "eth_getStorageAt" {
		fmt.Println(fmt.Sprintf("🔐 Decrypting %s response from Geth node with viewing key.", method))
		eciesPrivateKey := ecies.ImportECDSA(we.enclavePrivateKey)
		gethResp, err = eciesPrivateKey.Decrypt(gethResp, nil, nil)
		if err != nil {
			fmt.Println(err)
		}
	}

	// We unmarshall the JSON response to inspect it.
	var respJsonMap map[string]interface{}
	err = json.Unmarshal(gethResp, &respJsonMap)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fmt.Sprintf("Received response from Geth node: %s", strings.TrimSpace(string(gethResp))))

	// We write the response to the client.
	_, err = resp.Write(gethResp)
	if err != nil {
		fmt.Println(err)
	}
}

// Returns a new viewing key.
func (we WalletExtension) handleGetViewingKey(resp http.ResponseWriter, _ *http.Request) {
	// todo - generate viewing key properly
	// todo - store private key
	_, err := resp.Write([]byte("dummyViewingKey"))
	if err != nil {
		fmt.Println(err)
	}
}
