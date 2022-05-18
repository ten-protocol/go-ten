package main

import "github.com/obscuronet/obscuro-playground/tools/walletextension"

func main() {
	config := parseCLIArgs()
	stopHandle := walletextension.StartWalletExtension(config)
	defer stopHandle()
	select {}
}
