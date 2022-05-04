package main

import (
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/obscuronet/obscuro-playground/tools/walletextension"
)

func main() {
	enclavePrivateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	viewingKeyChannel := make(chan walletextension.ViewingKey)

	we := walletextension.NewWalletExtension(enclavePrivateKey, viewingKeyChannel)
	of := walletextension.NewObxFacade(enclavePrivateKey, viewingKeyChannel)

	go of.Serve("localhost:3001")
	we.Serve("localhost:3000")
}
