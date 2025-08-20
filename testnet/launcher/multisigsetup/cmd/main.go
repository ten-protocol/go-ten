package main

import (
	"fmt"
	"os"

	multisigsetup "github.com/ten-protocol/go-ten/testnet/launcher/multisigsetup"
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

	multisigSetup, err := multisigsetup.NewMultisigSetup(
		multisigsetup.NewMultisigSetupConfig(
			multisigsetup.WithL1HTTPURL(cliConfig.l1HTTPURL),
			multisigsetup.WithPrivateKey(cliConfig.privateKey),
			multisigsetup.WithDockerImage(cliConfig.dockerImage),
			multisigsetup.WithNetworkConfigAddress(cliConfig.networkConfigAddr),
			multisigsetup.WithMultisigAddress(cliConfig.multisigAddress),
		),
	)
	if err != nil {
		fmt.Printf("unable to configure multisig setup script - %v\n", err)
		os.Exit(1)
	}

	err = multisigSetup.Start()
	if err != nil {
		fmt.Printf("unable to start multisig setup script - %v\n", err)
		os.Exit(1)
	}

	err = multisigSetup.WaitForFinish()
	if err != nil {
		fmt.Printf("unexpected error waiting for multisig setup script to finish - %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Multisig setup was successfully completed...")
	os.Exit(0)
}
