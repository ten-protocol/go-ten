package main

import (
	"flag"
	"fmt"

	"github.com/obscuronet/go-obscuro/tools/walletextension"
)

const (
	// Flag names, defaults and usages.
	walletExtensionPortName    = "port"
	walletExtensionPortDefault = 3000
	walletExtensionPortUsage   = "The port on which to serve the wallet extension. Default: 3000."

	nodeHostName    = "nodeHost"
	nodeHostDefault = "testnet.obscu.ro"
	nodeHostUsage   = "The host on which to connect to the Obscuro node. Default: `testnet.obscu.ro`."

	nodeHTTPPortName    = "nodePortHTTP"
	nodeHTTPPortDefault = 13000
	nodeHTTPPortUsage   = "The port on which to connect to the Obscuro node via RPC over HTTP. Default: 13000."

	nodeWebsocketPortName    = "nodePortWS"
	nodeWebsocketPortDefault = 13001
	nodeWebsocketPortUsage   = "The port on which to connect to the Obscuro node via RPC over websockets. Default: 13001."

	logPathName    = "logPath"
	logPathDefault = "wallet_extension_logs.txt"
	logPathUsage   = "The path to use for the wallet extension's log file"
)

func parseCLIArgs() walletextension.Config {
	walletExtensionPort := flag.Int(walletExtensionPortName, walletExtensionPortDefault, walletExtensionPortUsage)
	nodeHost := flag.String(nodeHostName, nodeHostDefault, nodeHostUsage)
	nodeHTTPPort := flag.Int(nodeHTTPPortName, nodeHTTPPortDefault, nodeHTTPPortUsage)
	nodeWebsocketPort := flag.Int(nodeWebsocketPortName, nodeWebsocketPortDefault, nodeWebsocketPortUsage)
	logPath := flag.String(logPathName, logPathDefault, logPathUsage)
	flag.Parse()

	return walletextension.Config{
		WalletExtensionPort:     *walletExtensionPort,
		NodeRPCHTTPAddress:      fmt.Sprintf("http://%s:%d", *nodeHost, *nodeHTTPPort),
		NodeRPCWebsocketAddress: fmt.Sprintf("ws://%s:%d", *nodeHost, *nodeWebsocketPort),
		LogPath:                 *logPath,
	}
}
