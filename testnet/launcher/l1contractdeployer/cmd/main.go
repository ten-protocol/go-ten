package main

import (
	"fmt"
	"github.com/ten-protocol/go-ten/go/common/docker"
	"github.com/ten-protocol/go-ten/go/config"
	"os"

	l1cd "github.com/ten-protocol/go-ten/testnet/launcher/l1contractdeployer"
)

func main() {
	var err error
	// load flags with defaults from config / sub-configs
	rParams, _, err := config.LoadFlagStrings(config.L1Deployer)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting a L1 contract deployer...")
	conf, err := ParseConfig(rParams)

	l1ContractDeployer, err := l1cd.NewDockerContractDeployer(conf)
	if err != nil {
		fmt.Println("unable to configure l1 contract deployer - %w", err)
		os.Exit(1)
	}

	// remove if exists
	_ = docker.RemoveContainer(conf.ContainerName)

	err = l1ContractDeployer.Start()
	if err != nil {
		fmt.Println("unable to start l1 contract deployer - %w", err)
		os.Exit(1)
	}

	networkConfig, err := l1ContractDeployer.RetrieveL1ContractAddresses()
	if err != nil {
		fmt.Println("unable to fetch l1 contract addresses - %w", err)
		os.Exit(1)
	}
	networkDetails := networkConfig.GetNetwork()
	fmt.Println("L1 Contracts were successfully deployed...")

	fmt.Printf("MGMTCONTRACTADDR=%s\n", networkDetails.ManagementContractAddress)
	fmt.Printf("MSGBUSCONTRACTADDR=%s\n", networkDetails.MessageBusAddress)
	fmt.Printf("L1START=%s\n", networkDetails.L1StartHash)

	// the responsibility of writing to disk is outside the deployers domain
	if conf.ContractsEnvFile != "" {
		envFile := fmt.Sprintf("MGMTCONTRACTADDR=%s\nMSGBUSCONTRACTADDR=%s\nL1START=%s\n",
			networkDetails.ManagementContractAddress, networkDetails.MessageBusAddress, networkDetails.L1StartHash)

		// Write the content to a new file or override the existing file
		err = os.WriteFile(conf.ContractsEnvFile, []byte(envFile), 0o644) //nolint:gosec
		if err != nil {
			fmt.Println("Error writing file:", err)
			os.Exit(1)
		}
	}
	os.Exit(0)
}
