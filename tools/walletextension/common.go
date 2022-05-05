package walletextension

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/gorilla/websocket"
)

const (
	reqJSONKeyMethod          = "method"
	reqJSONMethodGetBalance   = "eth_getBalance"
	reqJSONMethodGetStorageAt = "eth_getStorageAt"
	pathRoot                  = "/"
	httpCodeErr               = 500
)

// ViewingKey is the packet of data sent to the enclave when storing a new viewing key.
type ViewingKey struct {
	viewingKeyPublic *ecdsa.PublicKey
	signedBytes      []byte
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
		fmt.Println(err)
		return nil, err
	}
	return message, nil
}
