package main

import (
	"fmt"

	"github.com/obscuronet/obscuro-playground/tools/walletextension"
)

func main() {
	config := parseCLIArgs()
	walletExtension := walletextension.NewWalletExtension(config)
	defer walletExtension.Shutdown()

	walletExtensionAddr := fmt.Sprintf("%s:%d", walletextension.Localhost, config.WalletExtensionPort)
	go walletExtension.Serve(walletExtensionAddr)
	fmt.Printf("Wallet extension started.\n💡 Visit %s/viewingkeys/ to generate an ephemeral viewing key. "+
		"Without a viewing key, you will not be able to decrypt the enclave's secure responses to sensitive requests.\n", walletExtensionAddr)

	select {}
}
