package walletextension

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"io/ioutil"
	"net/http"
)

// todo - joel - encrypt with the address key on the return in the facade, and decrypt in the wallet extension
//  use a standard account key for now, since we don't have viewing keys

// WalletExtension is a server that handles the management of viewing keys and the forwarding of Ethereum JSON-RPC requests.
type WalletExtension struct {
	enclavePublicKey *ecdsa.PublicKey
}

func NewWalletExtension(enclavePublicKey *ecdsa.PublicKey) *WalletExtension {
	return &WalletExtension{enclavePublicKey: enclavePublicKey}
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

	// We unmarshall the JSON to inspect it.
	var jsonMap map[string]interface{}
	err = json.Unmarshal(body, &jsonMap)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fmt.Sprintf("Received request: %s", body))

	// We encrypt the JSON with the enclave's public key.
	eciesKey := ecies.ImportECDSAPublic(we.enclavePublicKey)
	encryptedBody, err := ecies.Encrypt(rand.Reader, eciesKey, body, nil, nil)
	if err != nil {
		fmt.Println(err)
	}

	// We forward the requests on to the Geth node.
	gethResp := forwardMsgOverWebsocket(obxFacadeWebsocketAddr, encryptedBody)
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
