package main

import (
	"github.com/obscuronet/obscuro-playground/tools/walletextension"
)

func main() {
	we := walletextension.NewWalletExtension()
	of := walletextension.NewObxFacade()

	go of.Serve("localhost:3001")
	we.Serve("localhost:3000")
}
