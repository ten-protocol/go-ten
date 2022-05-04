package main

import (
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/obscuro-playground/tools/walletextension"
)

const (
	walletExtensionAddr = "localhost:3000"
	obscuroFacadeAddr   = "localhost:3001"
	gethWebsocketAddr   = "ws://localhost:8546"
)

func main() {
	enclavePrivateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	viewingKeyChannel := make(chan walletextension.ViewingKey)

	walletExtension := walletextension.NewWalletExtension(enclavePrivateKey, obscuroFacadeAddr, viewingKeyChannel)
	obscuroFacade := walletextension.NewObscuroFacade(enclavePrivateKey, gethWebsocketAddr, viewingKeyChannel)

	go obscuroFacade.Serve(obscuroFacadeAddr)
	walletExtension.Serve(walletExtensionAddr)
}
