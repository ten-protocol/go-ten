package main

import (
	"fmt"
	"os"

	directupgrade "github.com/ten-protocol/go-ten/testnet/launcher/directupgrade"
)

func main() {
	cliConfig := ParseConfigCLI()

	// Validate required parameters
	if cliConfig.networkConfigAddr == "" {
		fmt.Println("Error: network_config_addr is required")
		os.Exit(1)
	}
	if cliConfig.multisigAddress == "" {
		fmt.Println("Error: multisig_address is required")
		os.Exit(1)
	}

	directUpgrade, err := directupgrade.NewDirectUpgrade(
		directupgrade.NewDirectUpgradeConfig(
			directupgrade.WithL1HTTPURL(cliConfig.l1HTTPURL),
			directupgrade.WithPrivateKey(cliConfig.privateKey),
			directupgrade.WithDockerImage(cliConfig.dockerImage),
			directupgrade.WithNetworkConfigAddress(cliConfig.networkConfigAddr),
			directupgrade.WithMultisigAddress(cliConfig.multisigAddress),
		),
	)
	if err != nil {
		fmt.Printf("unable to configure direct upgrade script - %v\n", err)
		os.Exit(1)
	}

	err = directUpgrade.Start()
	if err != nil {
		fmt.Printf("unable to start direct upgrade script - %v\n", err)
		os.Exit(1)
	}

	err = directUpgrade.WaitForFinish()
	if err != nil {
		fmt.Printf("unexpected error waiting for direct upgrade script to finish - %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Direct upgrades were successfully completed...")
	os.Exit(0)
}