package walletextension

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// WalletExtension is a server that handles the management of viewing keys and the forwarding of Ethereum JSON-RPC requests.
type WalletExtension struct{}

func NewWalletExtension() *WalletExtension {
	return &WalletExtension{}
}

func (we WalletExtension) Serve(hostAndPort string) {
	serveMux := http.NewServeMux()

	// Handles Ethereum JSON-RPC requests received over HTTP.
	serveMux.HandleFunc("/", handleHttpEthJson)

	// Handles the management of viewing keys.
	serveMux.Handle("/viewingkeys/", http.StripPrefix("/viewingkeys/", http.FileServer(http.Dir("./tools/walletextension/static"))))
	serveMux.HandleFunc("/getViewingKey", handleGetViewingKey)

	err := http.ListenAndServe(hostAndPort, serveMux)
	if err != nil {
		panic(err)
	}
}

// Reads the Ethereum JSON-RPC request, and forwards it to the Geth node over a websocket.
func handleHttpEthJson(resp http.ResponseWriter, req *http.Request) {
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

	// We forward the requests on to the Geth node.
	gethResp := forwardMsgOverWebsocket(obxFacadeWebsocketAddr, body)
	_, err = resp.Write(gethResp)
	if err != nil {
		fmt.Println(err)
	}
}

// Returns a new viewing key.
func handleGetViewingKey(resp http.ResponseWriter, _ *http.Request) {
	// todo - generate viewing key properly
	// todo - store private key
	_, err := resp.Write([]byte("dummyViewingKey"))
	if err != nil {
		fmt.Println(err)
	}
}
