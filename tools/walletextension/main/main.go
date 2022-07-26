package main

import (
	"encoding/json"
	"fmt"
	"github.com/obscuronet/go-obscuro/tools/walletextension"
)

const (
	localhost = "127.0.0.1"
)

func main() {
	config := parseCLIArgs()
	walletExtension := walletextension.NewWalletExtension(config)
	defer walletExtension.Shutdown()

	walletExtensionAddr := fmt.Sprintf("%s:%d", localhost, config.WalletExtensionPort)
	go walletExtension.Serve(walletExtensionAddr)
	fmt.Printf("Wallet extension started.\n💡 Visit http://%s/viewingkeys/ to generate an ephemeral viewing key.\n", walletExtensionAddr)
	fmt.Println()
	s, _ := json.MarshalIndent(config, "", "  ")
	fmt.Printf("Wallet extension config: \n%s", string(s))

	select {}
}
