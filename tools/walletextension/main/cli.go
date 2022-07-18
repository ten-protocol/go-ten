package main

import (
	"flag"

	"github.com/obscuronet/go-obscuro/tools/walletextension"
)

const (
	// Flag names, defaults and usages.
	walletExtensionPortName    = "port"
	walletExtensionPortDefault = 3000
	walletExtensionPortUsage   = "The port on which to serve the wallet extension"

	// TODO - Use one flag for the shared HTTP + RPC host, and two other flags for the ports.
	nodeRPCHTTPAddressName    = "nodeRPCHTTPAddress"
	nodeRPCHTTPAddressDefault = "testnet.obscu.ro:13000"
	nodeRPCHTTPAddressUsage   = "The address on which to connect to the node via RPC using HTTP"

	nodeRPCWebsocketAddressName    = "nodeRPCWebsocketAddress"
	nodeRPCWebsocketAddressDefault = "testnet.obscu.ro:13001"
	nodeRPCWebsocketAddressUsage   = "The address on which to connect to the node via RPC using websockets"

	logPathName    = "logPath"
	logPathDefault = "wallet_extension_logs.txt"
	logPathUsage   = "The path to use for the wallet extension's log file"
)

func parseCLIArgs() walletextension.Config {
	walletExtensionPort := flag.Int(walletExtensionPortName, walletExtensionPortDefault, walletExtensionPortUsage)
	nodeRPCHTTPAddress := flag.String(nodeRPCHTTPAddressName, nodeRPCHTTPAddressDefault, nodeRPCHTTPAddressUsage)
	nodeRPCWebsocketAddress := flag.String(nodeRPCWebsocketAddressName, nodeRPCWebsocketAddressDefault, nodeRPCWebsocketAddressUsage)
	logPath := flag.String(logPathName, logPathDefault, logPathUsage)
	flag.Parse()

	return walletextension.Config{
		WalletExtensionPort:     *walletExtensionPort,
		NodeRPCHTTPAddress:      *nodeRPCHTTPAddress,
		NodeRPCWebsocketAddress: *nodeRPCWebsocketAddress,
		LogPath:                 *logPath,
	}
}
