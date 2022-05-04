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

	serveMux.HandleFunc("/", of.handleWSEthJson)

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
		fmt.Println(err)
		return
	}

	// We read the message from the wallet extension.
	_, encryptedMessage, err := connection.ReadMessage()
	if err != nil {
		fmt.Println(err)
		return
	}

	// We decrypt the JSON with the enclave's private key.
	eciesPrivateKey := ecies.ImportECDSA(of.enclavePrivateKey)
	message, err := eciesPrivateKey.Decrypt(encryptedMessage, nil, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// We forward the message to the Geth node.
	gethResp := forwardMsgOverWebsocket(gethWebsocketAddr, message)

	// We unmarshall the JSON request to inspect it.
	var reqJsonMap map[string]interface{}
	err = json.Unmarshal(message, &reqJsonMap)
	if err != nil {
		fmt.Println(err)
		return
	}

	// We encrypt the response if needed.
	method := reqJsonMap[reqJsonKeyMethod]
	if method == reqJsonMethodGetBalance || method == reqJsonMethodGetStorageAt {
		if of.viewingKey == nil {
			fmt.Printf("Could not respond securely to %s request because there is no viewing key for the account.\n", method)
			return
		}

		eciesPublicKey := ecies.ImportECDSAPublic(of.viewingKey)
		gethResp, err = ecies.Encrypt(rand.Reader, eciesPublicKey, gethResp, nil, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// We write the message back to the wallet extension.
	err = connection.WriteMessage(websocket.TextMessage, gethResp)
	if err != nil {
		fmt.Println(err)
		return
	}
}
