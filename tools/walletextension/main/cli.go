package main

import (
	"flag"
	"fmt"

	"github.com/obscuronet/go-obscuro/tools/walletextension/config"
)

const (
	// Flag names, defaults and usages.
	walletExtensionHostName    = "host"
	walletExtensionHostDefault = "127.0.0.1"
	walletExtensionHostUsage   = "The host where the wallet extension should open the port."
	walletExtensionPortName    = "port"
	walletExtensionPortDefault = 3000
	walletExtensionPortUsage   = "The port on which to serve the wallet extension. Default: 3000."

	walletExtensionPortWSName    = "portWS"
	walletExtensionPortWSDefault = 3001
	walletExtensionPortWSUsage   = "The port on which to serve websocket JSON RPC requests. Default: 3001."

	nodeHostName    = "nodeHost"
	nodeHostDefault = "testnet.obscu.ro"
	nodeHostUsage   = "The host on which to connect to the Obscuro node. Default: `testnet.obscu.ro`."

	nodeHTTPPortName    = "nodePortHTTP"
	nodeHTTPPortDefault = 80
	nodeHTTPPortUsage   = "The port on which to connect to the Obscuro node via RPC over HTTP. Default: 80."

	nodeWebsocketPortName    = "nodePortWS"
	nodeWebsocketPortDefault = 81
	nodeWebsocketPortUsage   = "The port on which to connect to the Obscuro node via RPC over websockets. Default: 81."

	logPathName    = "logPath"
	logPathDefault = "wallet_extension_logs.txt"
	logPathUsage   = "The path to use for the wallet extension's log file"

	databasePathName    = "databasePath"
	databasePathDefault = ".obscuro/gateway_database.db"
	databasePathUsage   = "The path for the wallet extension's database file. Default: .obscuro/gateway_database.db"

	hostedName        = "hosted"
	hostedNameDefault = "local"
	hostedNameUsage   = "Select where Obscuro Gateway is running (needed for Metamask): local (default), dev-testnet, testnet" // Needed for Metamask to know where to connect to

	verboseFlagName    = "verbose"
	verboseFlagDefault = false
	verboseFlagUsage   = "Flag to enable verbose logging of wallet extension traffic"
)

func parseCLIArgs() config.Config {
	walletExtensionHost := flag.String(walletExtensionHostName, walletExtensionHostDefault, walletExtensionHostUsage)
	walletExtensionPort := flag.Int(walletExtensionPortName, walletExtensionPortDefault, walletExtensionPortUsage)
	walletExtensionPortWS := flag.Int(walletExtensionPortWSName, walletExtensionPortWSDefault, walletExtensionPortWSUsage)
	nodeHost := flag.String(nodeHostName, nodeHostDefault, nodeHostUsage)
	nodeHTTPPort := flag.Int(nodeHTTPPortName, nodeHTTPPortDefault, nodeHTTPPortUsage)
	nodeWebsocketPort := flag.Int(nodeWebsocketPortName, nodeWebsocketPortDefault, nodeWebsocketPortUsage)
	logPath := flag.String(logPathName, logPathDefault, logPathUsage)
	databasePath := flag.String(databasePathName, databasePathDefault, databasePathUsage)
	hosted := flag.String(hostedName, hostedNameDefault, hostedNameUsage)
	verboseFlag := flag.Bool(verboseFlagName, verboseFlagDefault, verboseFlagUsage)
	flag.Parse()

	return config.Config{
		WalletExtensionHost:     *walletExtensionHost,
		WalletExtensionPortHTTP: *walletExtensionPort,
		WalletExtensionPortWS:   *walletExtensionPortWS,
		NodeRPCHTTPAddress:      fmt.Sprintf("%s:%d", *nodeHost, *nodeHTTPPort),
		NodeRPCWebsocketAddress: fmt.Sprintf("%s:%d", *nodeHost, *nodeWebsocketPort),
		LogPath:                 *logPath,
		DBPathOverride:          *databasePath,
		Hosted:                  *hosted,
		VerboseFlag:             *verboseFlag,
	}
}
