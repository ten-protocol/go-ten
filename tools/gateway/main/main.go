package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ten-protocol/go-ten/tools/walletextension"

	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/tools/walletextension/common"
)

func main() {
	config := parseCLIArgs()
	jsonConfig, _ := json.MarshalIndent(config, "", "  ")

	// Setup logging first
	if config.LogPath != log.SysOut {
		_, err := os.Create(config.LogPath)
		if err != nil {
			panic(fmt.Sprintf("could not create log file. Cause: %s", err))
		}
	}
	logger := log.New(log.WalletExtCmp, config.LogLevel, config.LogPath)

	logger.Info("Welcome to the TEN gateway")
	logger.Info("Starting with following config", "config", string(jsonConfig))

	// Start the wallet extension right away
	walletExtContainer := walletextension.NewContainerFromConfig(config, logger)
	err := walletExtContainer.Start()
	if err != nil {
		logger.Error("Failed to start wallet extension", "error", err)
		os.Exit(1)
	}

	walletExtensionAddr := fmt.Sprintf("%s:%d", common.Localhost, config.WalletExtensionPortHTTP)
	fmt.Println("TEN gateway started") // We expect stdout message in some tests
	logger.Info("TEN gateway started: ", "url", fmt.Sprintf("http://%s/v1/network-config", walletExtensionAddr))

	select {}
}
