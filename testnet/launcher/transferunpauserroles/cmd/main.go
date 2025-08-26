package main

import (
	"fmt"
	"os"

	transferunpauserroles "github.com/ten-protocol/go-ten/testnet/launcher/transferunpauserroles"
)

func main() {
	cliConfig := ParseConfigCLI()

	// Validate required parameters
	if cliConfig.networkConfigAddr == "" {
		fmt.Println("Error: network_config_addr is required")
		os.Exit(1)
	}
	if cliConfig.multisigAddr == "" {
		fmt.Println("Error: multisig_addr is required")
		os.Exit(1)
	}
	//if cliConfig.merkleMessageBusAddr == "" {
	//	fmt.Println("Error: merkle_message_bus_addr is required")
	//	os.Exit(1)
	//}

	roleTransfer, err := transferunpauserroles.NewRoleTransfer(
		transferunpauserroles.NewRoleTransferConfig(
			transferunpauserroles.WithL1HTTPURL(cliConfig.l1HTTPURL),
			transferunpauserroles.WithPrivateKey(cliConfig.privateKey),
			transferunpauserroles.WithDockerImage(cliConfig.dockerImage),
			transferunpauserroles.WithNetworkConfigAddress(cliConfig.networkConfigAddr),
			transferunpauserroles.WithMultisigAddress(cliConfig.multisigAddr),
		),
	)
	if err != nil {
		fmt.Printf("unable to configure role transfer script - %v\n", err)
		os.Exit(1)
	}

	err = roleTransfer.Start()
	if err != nil {
		fmt.Printf("unable to start role transfer script - %v\n", err)
		os.Exit(1)
	}

	err = roleTransfer.WaitForFinish()
	if err != nil {
		fmt.Printf("unexpected error waiting for role transfer script to finish - %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Unpauser role transfer was successfully completed...")
	os.Exit(0)
}
