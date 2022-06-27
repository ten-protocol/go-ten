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

	// TODO - Use one flag for the shared HTTP + RPC host, and two other flags for the ports.
	nodeRPCHTTPAddressName    = "nodeRPCHTTPAddress"
	nodeRPCHTTPAddressDefault = "http://127.0.0.1:13000"
	nodeRPCHTTPAddressUsage   = "The address on which to connect to the node via RPC using HTTP"

	nodeRPCWebsocketAddressName    = "nodeRPCWebsocketAddress"
	nodeRPCWebsocketAddressDefault = "http://127.0.0.1:13001"
	nodeRPCWebsocketAddressUsage   = "The address on which to connect to the node via RPC using websockets"
)

func parseCLIArgs() walletextension.Config {
	walletExtensionPort := flag.Int(walletExtensionPortName, walletExtensionPortDefault, walletExtensionPortUsage)
	nodeRPCHTTPAddress := flag.String(nodeRPCHTTPAddressName, nodeRPCHTTPAddressDefault, nodeRPCHTTPAddressUsage)
	nodeRPCWebsocketAddress := flag.String(nodeRPCWebsocketAddressName, nodeRPCWebsocketAddressDefault, nodeRPCWebsocketAddressUsage)
	flag.Parse()

	return walletextension.Config{
		WalletExtensionPort:     *walletExtensionPort,
		NodeRPCHTTPAddress:      *nodeRPCHTTPAddress,
		NodeRPCWebsocketAddress: *nodeRPCWebsocketAddress,
	}
}
