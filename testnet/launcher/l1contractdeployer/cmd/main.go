package main

import (
	"fmt"
	"os"

	l1cd "github.com/obscuronet/go-obscuro/testnet/launcher/l1contractdeployer"
)

func main() {
	cliConfig := ParseConfigCLI()

	l1ContractDeployer, err := l1cd.NewDockerContractDeployer(
		l1cd.NewContractDeployerConfig(
			l1cd.WithL1Host(cliConfig.l1Addr),
			l1cd.WithL1Port(cliConfig.l1HTTPPort),       // 8025
			l1cd.WithPrivateKey(cliConfig.privateKey),   //"f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb"),
			l1cd.WithDockerImage(cliConfig.dockerImage), //"testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest"
		),
	)
	if err != nil {
		fmt.Println("unable to configure l1 contract deployer - %w", err)
		os.Exit(1)
	}

	err = l1ContractDeployer.Start()
	if err != nil {
		fmt.Println("unable to start l1 contract deployer - %w", err)
		os.Exit(1)
	}

	managementContractAddr, messageBusContractAddr, err := l1ContractDeployer.RetrieveL1ContractAddresses()
	if err != nil {
		fmt.Println("unable to fetch l1 contract addresses - %w", err)
		os.Exit(1)
	}
	fmt.Println("L1 Contracts were successfully deployed...")

	fmt.Printf("MGMTCONTRACTADDR=%s\n", managementContractAddr)
	fmt.Printf("MSGBUSCONTRACTADDR=%s\n", messageBusContractAddr)

	// the responsibility of writing to disk is outside the deployers domain
	if cliConfig.contractsEnvFile != "" {
		envFile := fmt.Sprintf("MGMTCONTRACTADDR=%s\nMSGBUSCONTRACTADDR=%s\n", managementContractAddr, messageBusContractAddr)

		// Write the content to a new file or override the existing file
		err = os.WriteFile(cliConfig.contractsEnvFile, []byte(envFile), 0o644) //nolint:gosec
		if err != nil {
			fmt.Println("Error writing file:", err)
			os.Exit(1)
		}
	}
	os.Exit(0)
}
