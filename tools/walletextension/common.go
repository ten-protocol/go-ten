package walletextension

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/gorilla/websocket"
)

const (
	reqJsonKeyMethod          = "method"
	reqJsonMethodGetBalance   = "eth_getBalance"
	reqJsonMethodGetStorageAt = "eth_getStorageAt"
	pathRoot                  = "/"
	httpCodeErr               = 500
)

// ViewingKey is the packet of data sent to the enclave when storing a new viewing key.
type ViewingKey struct {
	viewingKeyPublic *ecdsa.PublicKey
	signedBytes      []byte
}

func forwardMsgOverWebsocket(url string, msg []byte) []byte {
	connection, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		fmt.Println(err)
	}

	err = connection.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		fmt.Println(err)
	}

	_, message, err := connection.ReadMessage()
	if err != nil {
		fmt.Println(err)
	}
	return message
}
