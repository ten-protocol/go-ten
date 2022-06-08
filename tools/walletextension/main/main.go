package main

import (
	"fmt"

	"github.com/obscuronet/obscuro-playground/tools/walletextension"
)

func main() {
	config := parseCLIArgs()
	walletExtension := walletextension.NewWalletExtension(config)
	defer walletExtension.Shutdown()
	walletExtension.Serve(fmt.Sprintf("%s:%d", walletextension.Localhost, config.WalletExtensionPort))
}
