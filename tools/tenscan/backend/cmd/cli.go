package main

import (
	"flag"

	"github.com/ten-protocol/go-ten/tools/tenscan/backend/config"
)

func parseCLIArgs() *config.Config {
	defaultConfig := &config.Config{
		NodeHostAddress: "http://erpc.dev-testnet.ten.xyz:80",
		ServerAddress:   "0.0.0.0:80",
		LogPath:         "tenscan_logs.txt",
	}

	nodeHostAddress := flag.String(nodeHostAddressName, defaultConfig.NodeHostAddress, nodeHostAddressUsage)
	serverAddress := flag.String(serverAddressName, defaultConfig.ServerAddress, serverAddressUsage)
	logPath := flag.String(logPathName, defaultConfig.LogPath, logPathUsage)

	flag.Parse()

	return &config.Config{
		NodeHostAddress: *nodeHostAddress,
		ServerAddress:   *serverAddress,
		LogPath:         *logPath,
	}
}

const (
	nodeHostAddressName  = "nodeHostAddress"
	nodeHostAddressUsage = "The Obscuro Host Node address"

	serverAddressName  = "serverAddress"
	serverAddressUsage = "The address to serve tenscan on"

	logPathName  = "logPath"
	logPathUsage = "The path to use for tenscan's log file"
)
