package walletextension

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/gorilla/websocket"
	"net/http"
)

// ObxFacade is a server that inverts the encryption and decryption performed by WalletExtension, so that the forwarded
// Ethereum JSON-RPC requests can be understood by a regular Geth node.
type ObxFacade struct {
	enclavePrivateKey *ecdsa.PrivateKey
}

func NewObxFacade(enclavePrivateKey *ecdsa.PrivateKey) *ObxFacade {
	return &ObxFacade{enclavePrivateKey: enclavePrivateKey}
}

func (of ObxFacade) Serve(hostAndPort string) {
	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/", of.handleWSEthJson)

	err := http.ListenAndServe(hostAndPort, serveMux)
	if err != nil {
		panic(err)
	}
}

func (of ObxFacade) handleWSEthJson(resp http.ResponseWriter, req *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	connection, err := upgrader.Upgrade(resp, req, nil)
	if err != nil {
		fmt.Println(err)
	}

	// We read the message from the wallet extension.
	_, encryptedMessage, err := connection.ReadMessage()
	if err != nil {
		fmt.Println(err)
	}

	// We decrypt the JSON with the enclave's private key.
	eciesKey := ecies.ImportECDSA(of.enclavePrivateKey)
	message, err := eciesKey.Decrypt(encryptedMessage, nil, nil)
	if err != nil {
		fmt.Println(err)
	}

	// We forward the message to the Geth node.
	gethResp := forwardMsgOverWebsocket(gethWebsocketAddr, message)

	// We write the message back to the wallet extension.
	err = connection.WriteMessage(websocket.TextMessage, gethResp)
	if err != nil {
		fmt.Println(err)
	}
}
