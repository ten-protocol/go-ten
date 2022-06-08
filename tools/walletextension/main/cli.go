package main

import (
	"flag"

	"github.com/obscuronet/obscuro-playground/tools/walletextension"
)

const (
	// Flag names, defaults and usages.
	walletExtensionPortName    = "port"
	walletExtensionPortDefault = 3000
	walletExtensionPortUsage   = "The port on which to serve the wallet extension"

	nodeRPCAddressName    = "nodeRPCAddress"
	nodeRPCAddressDefault = "127.0.0.1:13000"
	nodeRPCAddressUsage   = "The address on which to connect to the node via RPC"
)

func parseCLIArgs() walletextension.Config {
	walletExtensionPort := flag.Int(walletExtensionPortName, walletExtensionPortDefault, walletExtensionPortUsage)
	nodeRPCAddress := flag.String(nodeRPCAddressName, nodeRPCAddressDefault, nodeRPCAddressUsage)
	flag.Parse()

	return walletextension.Config{
		WalletExtensionPort: *walletExtensionPort,
		NodeRPCAddress:      *nodeRPCAddress,
	}
}
