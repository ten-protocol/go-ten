package main

import (
	"flag"

	"github.com/obscuronet/go-obscuro/tools/obscuroscan_v2/backend/config"
)

func parseCLIArgs() *config.Config {
	defaultConfig := &config.Config{
		NodeHostAddress: "http://127.0.0.1:37801",
		ServerAddress:   "0.0.0.0:43910",
		LogPath:         "obscuroscan_logs.txt",
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
	serverAddressUsage = "The address to serve Obscuroscan on"

	logPathName  = "logPath"
	logPathUsage = "The path to use for Obscuroscan's log file"
)
