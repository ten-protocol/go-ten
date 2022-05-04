package walletextension

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/gorilla/websocket"
	"net/http"
)

// ObxFacade is a server that inverts the encryption and decryption performed by WalletExtension, so that the forwarded
// Ethereum JSON-RPC requests can be understood by a regular Geth node.
type ObxFacade struct {
	enclavePrivateKey *ecdsa.PrivateKey
	viewingKeyChannel <-chan ViewingKey
	viewingKey        *ecdsa.PublicKey
}

func NewObxFacade(enclavePrivateKey *ecdsa.PrivateKey, viewingKeyChannel <-chan ViewingKey) *ObxFacade {
	return &ObxFacade{enclavePrivateKey: enclavePrivateKey, viewingKeyChannel: viewingKeyChannel}
}

func (of *ObxFacade) Serve(hostAndPort string) {
	// We listen for the account viewing key.
	go func() {
		viewingKey := <-of.viewingKeyChannel
		// todo - verify signed bytes
		of.viewingKey = viewingKey.viewingKeyPublic
	}()

	serveMux := http.NewServeMux()

	serveMux.HandleFunc(pathRoot, of.handleWSEthJson)

	err := http.ListenAndServe(hostAndPort, serveMux)
	if err != nil {
		panic(err)
	}
}

func (of *ObxFacade) handleWSEthJson(resp http.ResponseWriter, req *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	connection, err := upgrader.Upgrade(resp, req, nil)
	if err != nil {
		fmt.Printf("could not upgrade connection to a websocket connection: %v\n", err)
		return // todo - return error response
	}

	// We read the message from the wallet extension.
	_, encryptedMessage, err := connection.ReadMessage()
	if err != nil {
		fmt.Printf("could not read Ethereum JSON-RPC request: %v\n", err)
		return // todo - return error response
	}

	// We decrypt the JSON with the enclave's private key.
	eciesPrivateKey := ecies.ImportECDSA(of.enclavePrivateKey)
	message, err := eciesPrivateKey.Decrypt(encryptedMessage, nil, nil)
	if err != nil {
		fmt.Printf("could not decrypt Ethereum JSON-RPC request with enclave public key: %v\n", err)
		return // todo - return error response
	}

	// We forward the message to the Geth node.
	gethResp := forwardMsgOverWebsocket(gethWebsocketAddr, message)

	// We unmarshall the JSON request to inspect it.
	var reqJsonMap map[string]interface{}
	err = json.Unmarshal(message, &reqJsonMap)
	if err != nil {
		fmt.Printf("could not unmarshall Ethereum JSON-RPC request to JSON: %v\n", err)
		return // todo - return error response
	}

	// We encrypt the response if needed.
	method := reqJsonMap[reqJsonKeyMethod]
	if method == reqJsonMethodGetBalance || method == reqJsonMethodGetStorageAt {
		if of.viewingKey == nil {
			fmt.Printf("could not respond securely to %s request because there is no viewing key for the account.\n", method)
			return // todo - return error response
		}

		eciesPublicKey := ecies.ImportECDSAPublic(of.viewingKey)
		gethResp, err = ecies.Encrypt(rand.Reader, eciesPublicKey, gethResp, nil, nil)
		if err != nil {
			fmt.Printf("could not encrypt Ethereum JSON-RPC response with viewing key: %v\n", err)
			return // todo - return error response
		}
	}

	// We write the message back to the wallet extension.
	err = connection.WriteMessage(websocket.TextMessage, gethResp)
	if err != nil {
		fmt.Printf("could not write JSON-RPC response: %v\n", err)
		return // todo - return error response
	}
}
