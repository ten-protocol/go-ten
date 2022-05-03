package walletextension

import (
	"fmt"
	"github.com/gorilla/websocket"
)

const obxFacadeWebsocketAddr = "ws://localhost:3001"
const gethWebsocketAddr = "ws://localhost:8546"

const reqJsonKeyMethod = "method"
const reqJsonMethodGetBalance = "eth_getBalance"
const reqJsonMethodGetStorageAt = "eth_getStorageAt"

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
