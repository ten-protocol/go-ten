package walletextension

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
)

const gethWebsocketAddr = "ws://localhost:8546"

type WalletExtension struct{}

func NewWalletExtension() *WalletExtension {
	return &WalletExtension{}
}

func (cp WalletExtension) Serve() {
	// RPC request handler.
	// This is a helpful resource: https://docs.alchemy.com/alchemy/apis/ethereum/eth_chainid.
	http.HandleFunc("/", handleEthJsonReq)

	// Web app handler.
	http.Handle("/register/", http.StripPrefix("/register/", http.FileServer(http.Dir("./tools/walletextension/static"))))
	http.HandleFunc("/manageViewingKeys", handleManageViewingKeys)

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}

func handleEthJsonReq(resp http.ResponseWriter, req *http.Request) {
	var jsonMap map[string]interface{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(body, &jsonMap)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(fmt.Sprintf("Received %s request.", jsonMap["method"]))

	// We forward requests on to the Geth node.
	gethResp := forwardMsgOverWebsocket(gethWebsocketAddr, body)
	_, err = resp.Write(gethResp)
	if err != nil {
		fmt.Println(err)
	}
}

func handleManageViewingKeys(resp http.ResponseWriter, _ *http.Request) {
	// todo - generate viewing key properly
	// todo - store private key
	_, err := resp.Write([]byte("dummyViewingKey"))
	if err != nil {
		fmt.Println(err)
	}
}

func forwardMsgOverWebsocket(url string, req []byte) []byte {
	connection, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		fmt.Println(err)
	}

	err = connection.WriteMessage(websocket.TextMessage, req)
	if err != nil {
		fmt.Println(err)
	}

	_, message, err := connection.ReadMessage()
	if err != nil {
		fmt.Println(err)
	}
	return message
}
