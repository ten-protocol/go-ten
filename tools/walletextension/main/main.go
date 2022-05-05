package main

import "github.com/obscuronet/obscuro-playground/tools/walletextension"

func main() {
	config := parseCLIArgs()
	stopNodesFunc := walletextension.StartWalletExtension(config)
	defer stopNodesFunc()
	select {}
}
