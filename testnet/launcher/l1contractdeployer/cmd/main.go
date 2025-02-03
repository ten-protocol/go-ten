package main

import (
	"context"
	"fmt"
	"os"

	l1cd "github.com/ten-protocol/go-ten/testnet/launcher/l1contractdeployer"
)

func main() {
	cliConfig := ParseConfigCLI()

	l1ContractDeployer, err := l1cd.NewDockerContractDeployer(
		l1cd.NewContractDeployerConfig(
			l1cd.WithL1HTTPURL(cliConfig.l1HTTPURL),     // "http://eth2network:8025"
			l1cd.WithPrivateKey(cliConfig.privateKey),   //"f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb"),
			l1cd.WithDockerImage(cliConfig.dockerImage), //"testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest"
			l1cd.WithAzureKeyVaultURL(cliConfig.azureKeyVaultURL),
		),
	)
	if err != nil {
		fmt.Printf("unable to configure l1 contract deployer - %s\n", err)
		os.Exit(1)
	}

	err = l1ContractDeployer.Start()
	if err != nil {
		fmt.Printf("unable to start l1 contract deployer - %s\n", err)
		os.Exit(1)
	}

	networkConfig, err := l1ContractDeployer.RetrieveL1ContractAddresses()
	if err != nil {
		fmt.Printf("unable to fetch l1 contract addresses - %s", err)
		os.Exit(0)
	}
	fmt.Println("L1 Contracts were successfully deployed...")

	fmt.Printf("ENCLAVEREGISTRYADDR=%s\n", networkConfig.EnclaveRegistryAddress)
	fmt.Printf("CROSSCHAINADDR=%s\n", networkConfig.CrossChainAddress)
	fmt.Printf("DAREGISTRYADDR=%s\n", networkConfig.DataAvailabilityRegistryAddress)
	fmt.Printf("NETWORKCONFIGADDR=%s\n", networkConfig.NetworkConfigAddress)
	fmt.Printf("MSGBUSCONTRACTADDR=%s\n", networkConfig.MessageBusAddress)
	fmt.Printf("L1START=%s\n", networkConfig.L1StartHash)

	// the responsibility of writing to disk is outside the deployers domain
	if cliConfig.contractsEnvFile != "" {
		envFile := fmt.Sprintf("ENCLAVEREGISTRYADDR=%s\nCROSSCHAINADDR=%s\nDAREGISTRYADDR=%s\nNETWORKCONFIGADDR=%s\nMSGBUSCONTRACTADDR=%s\nL1START=%s\n",
			networkConfig.EnclaveRegistryAddress,
			networkConfig.CrossChainAddress,
			networkConfig.DataAvailabilityRegistryAddress,
			networkConfig.NetworkConfigAddress,
			networkConfig.MessageBusAddress,
			networkConfig.L1StartHash)

		// Write the content to a new file or override the existing file
		err = os.WriteFile(cliConfig.contractsEnvFile, []byte(envFile), 0o644) //nolint:gosec
		if err != nil {
			fmt.Println("Error writing file:", err)
			os.Exit(1)
		}
	}

	// Store in Azure Key Vault if configured
	if cliConfig.azureKeyVaultURL != "" {
		if err := l1cd.StoreNetworkCfgInKeyVault(context.Background(), cliConfig.azureKeyVaultURL, networkConfig); err != nil {
			fmt.Printf("Failed to store contracts in Azure Key Vault: %v\n", err)
			os.Exit(1)
		}
	}

	os.Exit(0)
}
