package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ten-protocol/go-ten/go/config"
	l1cd "github.com/ten-protocol/go-ten/testnet/launcher/l1contractdeployer"
)

func main() {
	tenCfg, err := config.LoadTenConfig()
	if err != nil {
		fmt.Println("Error loading ten config:", err)
		os.Exit(1)
	}

	deployerConfig := l1cd.NewContractDeployerConfig(tenCfg)
	l1ContractDeployer, err := l1cd.NewDockerContractDeployer(deployerConfig)
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
	fmt.Printf("BRIDGECONTRACTADDR=%s\n", networkConfig.BridgeAddress)
	fmt.Printf("L1START=%s\n", networkConfig.L1StartHash)

	// the responsibility of writing to disk is outside the deployers domain
	if tenCfg.Deployment.OutputEnvFile != "" {
		envFile := fmt.Sprintf("ENCLAVEREGISTRYADDR=%s\nCROSSCHAINADDR=%s\nDAREGISTRYADDR=%s\nNETWORKCONFIGADDR=%s\nMSGBUSCONTRACTADDR=%s\nBRIDGECONTRACTADDR=%s\nL1START=%s\n",
			networkConfig.EnclaveRegistryAddress,
			networkConfig.CrossChainAddress,
			networkConfig.DataAvailabilityRegistryAddress,
			networkConfig.NetworkConfigAddress,
			networkConfig.MessageBusAddress,
			networkConfig.BridgeAddress,
			networkConfig.L1StartHash)

		// Write the content to a new file or override the existing file
		err = os.WriteFile(tenCfg.Deployment.OutputEnvFile, []byte(envFile), 0o644) //nolint:gosec
		if err != nil {
			fmt.Println("Error writing file:", err)
			os.Exit(1)
		}
	}

	// Store in Azure Key Vault if configured
	if cliConfig.azureKeyVaultURL != "" {
		if err := l1cd.StoreNetworkCfgInKeyVault(context.Background(), cliConfig.azureKeyVaultURL, cliConfig.azureKeyVaultEnv, networkConfig); err != nil {
			fmt.Printf("Failed to store contracts in Azure Key Vault: %v\n", err)
			os.Exit(1)
		}
	}

	os.Exit(0)
}
