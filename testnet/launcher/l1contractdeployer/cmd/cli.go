package main

import (
	"flag"
)

// L1ContractDeployerConfigCLI represents the configurations passed into the deployer over CLI
type L1ContractDeployerConfigCLI struct {
	l1HTTPURL        string
	privateKey       string
	dockerImage      string
	contractsEnvFile string
	azureKeyVaultURL string
}

// ParseConfigCLI returns a NodeConfigCLI based the cli params and defaults.
func ParseConfigCLI() *L1ContractDeployerConfigCLI {
	cfg := &L1ContractDeployerConfigCLI{}
	flagUsageMap := getFlagUsageMap()

	l1HTTPURL := flag.String(l1HTTPURLFlag, "http://eth2network:8025", flagUsageMap[l1HTTPURLFlag])
	privateKey := flag.String(privateKeyFlag, "", flagUsageMap[privateKeyFlag])
	dockerImage := flag.String(dockerImageFlag, "testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest", flagUsageMap[dockerImageFlag])
	contractsEnvFile := flag.String(contractsEnvFileFlag, "", flagUsageMap[contractsEnvFileFlag])
	azureKeyVaultURL := flag.String(azureKeyVaultURLFlag, "", flagUsageMap[azureKeyVaultURLFlag])
	flag.Parse()

	cfg.l1HTTPURL = *l1HTTPURL
	cfg.privateKey = *privateKey
	cfg.dockerImage = *dockerImage
	cfg.contractsEnvFile = *contractsEnvFile
	cfg.azureKeyVaultURL = *azureKeyVaultURL

	return cfg
}
