package main

import "github.com/obscuronet/go-obscuro/tools/walletextension"

func main() {
	config := parseCLIArgs()
	stopHandle := walletextension.StartWalletExtension(config)
	defer stopHandle()
	select {}
}
