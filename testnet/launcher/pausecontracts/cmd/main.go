package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ten-protocol/go-ten/testnet/launcher/pausecontracts"
)

func main() {
	cliConfig := ParseConfigCLI()

	// Validate required parameters
	if cliConfig.networkConfigAddr == "" {
		fmt.Println("Error: network_config_addr is required")
		os.Exit(1)
	}
	if cliConfig.merkleMessageBusAddr == "" {
		fmt.Println("Error: merkle_message_bus_addr is required")
		os.Exit(1)
	}

	// Validate action parameter
	action := strings.ToUpper(cliConfig.action)
	if action != "PAUSE" && action != "UNPAUSE" {
		fmt.Println("Error: action must be either 'PAUSE' or 'UNPAUSE'")
		os.Exit(1)
	}

	pauseAllContracts, err := pausecontracts.NewPauseAllContracts(
		pausecontracts.NewPauseAllContractsConfig(
			pausecontracts.WithL1HTTPURL(cliConfig.l1HTTPURL),
			pausecontracts.WithPrivateKey(cliConfig.privateKey),
			pausecontracts.WithDockerImage(cliConfig.dockerImage),
			pausecontracts.WithNetworkConfigAddress(cliConfig.networkConfigAddr),
			pausecontracts.WithMerkleTreeMessageBus(cliConfig.merkleMessageBusAddr),
			pausecontracts.WithAction(action),
		),
	)
	if err != nil {
		fmt.Printf("unable to configure pause all contracts script - %v\n", err)
		os.Exit(1)
	}

	err = pauseAllContracts.Start()
	if err != nil {
		fmt.Printf("unable to start pause all contracts script - %v\n", err)
		os.Exit(1)
	}

	err = pauseAllContracts.WaitForFinish()
	if err != nil {
		fmt.Printf("unexpected error waiting for pause all contracts script to finish - %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s all contracts was successfully completed...\n", action)
	os.Exit(0)
}
