package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/obscuro-playground/integration/gethnetwork"
	"github.com/obscuronet/obscuro-playground/tools/walletextension"
	"strconv"
)

const (
	localhost = "localhost:"
	// TODO - Parameterise these ports.
	walletExtensionPort = 3000
	obscuroFacadePort   = 3001
	gethHTTPPort        = 3002
	gethWebsocketPort   = 8546
)

func main() {
	config := parseCLIArgs()

	gethWebsocketAddr := "ws://localhost:" + strconv.Itoa(gethWebsocketPort)
	if *config.localNetwork {
		gethBinaryPath, err := gethnetwork.EnsureBinariesExist(gethnetwork.LatestVersion)
		if err != nil {
			panic(err)
		}

		network := gethnetwork.NewGethNetwork(gethHTTPPort, gethBinaryPath, 1, 1, config.prefundedAccounts)
		defer network.StopNodes()
		fmt.Println("Local Geth network started.")

		gethWebsocketAddr = "ws://localhost:" + strconv.Itoa(int(network.WebSocketPorts[0]))
	}

	enclavePrivateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	viewingKeyChannel := make(chan walletextension.ViewingKey)

	obscuroFacadeAddr := localhost + strconv.Itoa(obscuroFacadePort)
	walletExtensionAddr := localhost + strconv.Itoa(walletExtensionPort)
	walletExtension := walletextension.NewWalletExtension(enclavePrivateKey, obscuroFacadeAddr, viewingKeyChannel)
	obscuroFacade := walletextension.NewObscuroFacade(enclavePrivateKey, gethWebsocketAddr, viewingKeyChannel)

	go obscuroFacade.Serve(obscuroFacadeAddr)
	fmt.Println("Obscuro facade started.")
	go walletExtension.Serve(walletExtensionAddr)
	fmt.Printf("Wallet extension started.\n💡 Visit %s/viewingkeys/ to generate an ephemeral viewing key. "+
		"Without a viewing key, you will not be able to decrypt the enclave's secure responses to your "+
		"eth_getBalance and eth_getStorageAt requests.\n", walletExtensionAddr)
	select {}
}
