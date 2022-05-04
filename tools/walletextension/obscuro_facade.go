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

// ObscuroFacade is a server that inverts the encryption and decryption performed by WalletExtension, so that the forwarded
// Ethereum JSON-RPC requests can be understood by a regular Geth node.
type ObscuroFacade struct {
	enclavePrivateKey *ecdsa.PrivateKey
	gethWebsocketAddr string
	viewingKeyChannel <-chan ViewingKey
	viewingKey        *ecdsa.PublicKey
}

func NewObscuroFacade(
	enclavePrivateKey *ecdsa.PrivateKey,
	gethWebsocketAddr string,
	viewingKeyChannel <-chan ViewingKey,
) *ObscuroFacade {
	return &ObscuroFacade{
		enclavePrivateKey: enclavePrivateKey,
		gethWebsocketAddr: gethWebsocketAddr,
		viewingKeyChannel: viewingKeyChannel}
}

// Serve listens for and serves Ethereum JSON-RPC requests from the wallet extension.
func (of *ObscuroFacade) Serve(hostAndPort string) {
	// We listen for the account viewing key.
	go func() {
		viewingKey := <-of.viewingKeyChannel
		// TODO - Verify signed bytes.
		of.viewingKey = viewingKey.viewingKeyPublic
	}()

	serveMux := http.NewServeMux()

	serveMux.HandleFunc(pathRoot, of.handleWSEthJson)

	err := http.ListenAndServe(hostAndPort, serveMux)
	if err != nil {
		panic(err)
	}
}

func (of *ObscuroFacade) handleWSEthJson(resp http.ResponseWriter, req *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	connection, err := upgrader.Upgrade(resp, req, nil)
	if err != nil {
		http.Error(resp, fmt.Sprintf("could not upgrade connection to a websocket connection: %v", err), httpCodeErr)
		return
	}

	// We read the message from the wallet extension.
	_, encryptedMessage, err := connection.ReadMessage()
	if err != nil {
		msg := fmt.Sprintf("could not read Ethereum JSON-RPC request: %v", err)
		_ = connection.WriteMessage(websocket.TextMessage, []byte(msg))
		return
	}

	// We decrypt the JSON with the enclave's private key.
	eciesPrivateKey := ecies.ImportECDSA(of.enclavePrivateKey)
	message, err := eciesPrivateKey.Decrypt(encryptedMessage, nil, nil)
	if err != nil {
		msg := fmt.Sprintf("could not decrypt Ethereum JSON-RPC request with enclave public key: %v", err)
		_ = connection.WriteMessage(websocket.TextMessage, []byte(msg))
		return
	}

	// We forward the message to the Geth node.
	gethResp := forwardMsgOverWebsocket(of.gethWebsocketAddr, message)

	// We unmarshall the JSON request to inspect it.
	var reqJsonMap map[string]interface{}
	err = json.Unmarshal(message, &reqJsonMap)
	if err != nil {
		msg := fmt.Sprintf("could not unmarshall Ethereum JSON-RPC request to JSON: %v", err)
		_ = connection.WriteMessage(websocket.TextMessage, []byte(msg))
		return
	}

	// We encrypt the response if needed.
	method := reqJsonMap[reqJsonKeyMethod]
	if method == reqJsonMethodGetBalance || method == reqJsonMethodGetStorageAt {
		if of.viewingKey == nil {
			msg := fmt.Sprintf("enclave could not respond securely to %s request because there is no viewing key for the account", method)
			_ = connection.WriteMessage(websocket.TextMessage, []byte(msg))
			return
		}

		eciesPublicKey := ecies.ImportECDSAPublic(of.viewingKey)
		gethResp, err = ecies.Encrypt(rand.Reader, eciesPublicKey, gethResp, nil, nil)
		if err != nil {
			msg := fmt.Sprintf("could not encrypt Ethereum JSON-RPC response with viewing key: %v", err)
			_ = connection.WriteMessage(websocket.TextMessage, []byte(msg))
			return
		}
	}

	// We write the message back to the wallet extension.
	err = connection.WriteMessage(websocket.TextMessage, gethResp)
	if err != nil {
		msg := fmt.Sprintf("could not write JSON-RPC response: %v", err)
		_ = connection.WriteMessage(websocket.TextMessage, []byte(msg))
		return
	}
}
