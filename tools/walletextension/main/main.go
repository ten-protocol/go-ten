package main

import "github.com/obscuronet/obscuro-playground/tools/walletextension"

func main() {
	config := parseCLIArgs()
	walletextension.StartWalletExtension(config)
	select {}
}
