package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ten-protocol/go-ten/tools/gateway"

	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/tools/gateway/common"
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
	logger := log.New(log.GatewayCmp, config.LogLevel, config.LogPath)

	logger.Info("Welcome to the TEN gateway")
	logger.Info("Starting with following config", "config", string(jsonConfig))

	// Start the gateway right away
	gatewayContainer := gateway.NewContainerFromConfig(config, logger)
	err := gatewayContainer.Start()
	if err != nil {
		logger.Error("Failed to start gateway", "error", err)
		os.Exit(1)
	}

	gatewayAddr := fmt.Sprintf("%s:%d", common.Localhost, config.GatewayPortHTTP)
	fmt.Println("TEN gateway started") // We expect stdout message in some tests
	logger.Info("TEN gateway started: ", "url", fmt.Sprintf("http://%s/v1/network-config", gatewayAddr))

	select {}
}
