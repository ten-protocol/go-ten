package main

import (
	"github.com/obscuronet/obscuro-playground/tools/walletextension"
)

func main() {
	cp := walletextension.NewWalletExtension()
	cp.Serve()
}
