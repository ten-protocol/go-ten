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

	nodeRPCWebsocketAddressName    = "nodeRPCWebsocketAddress"
	nodeRPCWebsocketAddressDefault = "127.0.0.1:13000"
	nodeRPCWebsocketAddressUsage   = "The address on which to connect to the node via RPC using websockets"
)

func parseCLIArgs() walletextension.Config {
	walletExtensionPort := flag.Int(walletExtensionPortName, walletExtensionPortDefault, walletExtensionPortUsage)
	nodeRPCAddress := flag.String(nodeRPCWebsocketAddressName, nodeRPCWebsocketAddressDefault, nodeRPCWebsocketAddressUsage)
	flag.Parse()

	return walletextension.Config{
		WalletExtensionPort:     *walletExtensionPort,
		NodeRPCWebsocketAddress: *nodeRPCAddress,
	}
}
