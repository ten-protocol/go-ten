package main

import (
	"flag"
	"fmt"

	"github.com/ten-protocol/go-ten/tools/walletextension/config"
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
	nodeHostDefault = "erpc.sepolia-testnet.obscu.ro"
	nodeHostUsage   = "The host on which to connect to the Obscuro node. Default: `erpc.sepolia-testnet.obscu.ro`."

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

	verboseFlagName    = "verbose"
	verboseFlagDefault = false
	verboseFlagUsage   = "Flag to enable verbose logging of wallet extension traffic"

	dbTypeFlagName    = "dbType"
	dbTypeFlagDefault = "sqlite"
	dbTypeFlagUsage   = "Defined the db type (sqlite or mariaDB)"

	dbConnectionURLFlagName    = "dbConnectionURL"
	dbConnectionURLFlagDefault = ""
	dbConnectionURLFlagUsage   = "If dbType is set to mariaDB, this must be set. ex: obscurouser:password@tcp(127.0.0.1:3306)/ogdb"

	tenChainIDName      = "tenChainID"
	tenChainIDDefault   = 443
	tenChainIDFlagUsage = "ChainID of Ten network that the gateway is communicating with"
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
	verboseFlag := flag.Bool(verboseFlagName, verboseFlagDefault, verboseFlagUsage)
	dbType := flag.String(dbTypeFlagName, dbTypeFlagDefault, dbTypeFlagUsage)
	dbConnectionURL := flag.String(dbConnectionURLFlagName, dbConnectionURLFlagDefault, dbConnectionURLFlagUsage)
	tenChainID := flag.Int(tenChainIDName, tenChainIDDefault, tenChainIDFlagUsage)
	flag.Parse()

	return config.Config{
		WalletExtensionHost:     *walletExtensionHost,
		WalletExtensionPortHTTP: *walletExtensionPort,
		WalletExtensionPortWS:   *walletExtensionPortWS,
		NodeRPCHTTPAddress:      fmt.Sprintf("%s:%d", *nodeHost, *nodeHTTPPort),
		NodeRPCWebsocketAddress: fmt.Sprintf("%s:%d", *nodeHost, *nodeWebsocketPort),
		LogPath:                 *logPath,
		DBPathOverride:          *databasePath,
		VerboseFlag:             *verboseFlag,
		DBType:                  *dbType,
		DBConnectionURL:         *dbConnectionURL,
		TenChainID:              *tenChainID,
	}
}
