package walletextension

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/obscuro-playground/integration/gethnetwork"
	"strconv"

	"github.com/gorilla/websocket"
)

const (
	reqJSONKeyMethod          = "method"
	reqJSONMethodGetBalance   = "eth_getBalance"
	reqJSONMethodGetStorageAt = "eth_getStorageAt"
	respJSONKeyErr            = "error"
	pathRoot                  = "/"
	httpCodeErr               = 500

	localhost = "localhost:"
	// TODO - Parameterise these ports.
	walletExtensionPort = 3000
	obscuroFacadePort   = 3001
	gethHTTPPort        = 3002
	gethWebsocketPort   = 8546
)

// ViewingKey is the packet of data sent to the enclave when storing a new viewing key.
type ViewingKey struct {
	viewingKeyPublic *ecdsa.PublicKey
	signedBytes      []byte
}

// RunConfig contains the configuration required by StartWalletExtension.
type RunConfig struct {
	LocalNetwork      bool
	PrefundedAccounts []string
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

// StartWalletExtension starts the wallet extension and Obscuro facade, and optionally a local Ethereum network.
func StartWalletExtension(config RunConfig) {
	gethWebsocketAddr := "ws://localhost:" + strconv.Itoa(gethWebsocketPort)
	if config.LocalNetwork {
		gethBinaryPath, err := gethnetwork.EnsureBinariesExist(gethnetwork.LatestVersion)
		if err != nil {
			panic(err)
		}

		network := gethnetwork.NewGethNetwork(gethHTTPPort, gethBinaryPath, 1, 1, config.PrefundedAccounts)
		defer network.StopNodes()
		fmt.Println("Local Geth network started.")

		gethWebsocketAddr = "ws://localhost:" + strconv.Itoa(int(network.WebSocketPorts[0]))
	}

	enclavePrivateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	viewingKeyChannel := make(chan ViewingKey)

	obscuroFacadeAddr := localhost + strconv.Itoa(obscuroFacadePort)
	walletExtensionAddr := localhost + strconv.Itoa(walletExtensionPort)
	walletExtension := NewWalletExtension(enclavePrivateKey, obscuroFacadeAddr, viewingKeyChannel)
	obscuroFacade := NewObscuroFacade(enclavePrivateKey, gethWebsocketAddr, viewingKeyChannel)

	go obscuroFacade.Serve(obscuroFacadeAddr)
	fmt.Println("Obscuro facade started.")
	go walletExtension.Serve(walletExtensionAddr)
	fmt.Printf("Wallet extension started.\n💡 Visit %s/viewingkeys/ to generate an ephemeral viewing key. "+
		"Without a viewing key, you will not be able to decrypt the enclave's secure responses to your "+
		"eth_getBalance and eth_getStorageAt requests.\n", walletExtensionAddr)
}
