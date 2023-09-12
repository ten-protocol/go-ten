package main

import (
	"flag"
)

// L1ContractDeployerConfigCLI represents the configurations passed into the deployer over CLI
type L1ContractDeployerConfigCLI struct {
	l1HTTPRPCAddr    string
	privateKey       string
	dockerImage      string
	contractsEnvFile string
}

// ParseConfigCLI returns a NodeConfigCLI based the cli params and defaults.
func ParseConfigCLI() *L1ContractDeployerConfigCLI {
	cfg := &L1ContractDeployerConfigCLI{}
	flagUsageMap := getFlagUsageMap()

	l1HTTPRPCAddr := flag.String(l1HTTPRPCAddressFlag, "http://eth2network:8025", flagUsageMap[l1HTTPRPCAddressFlag])
	privateKey := flag.String(privateKeyFlag, "", flagUsageMap[privateKeyFlag])
	dockerImage := flag.String(dockerImageFlag, "testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest", flagUsageMap[dockerImageFlag])
	contractsEnvFile := flag.String(contractsEnvFileFlag, "", flagUsageMap[contractsEnvFileFlag])
	flag.Parse()

	cfg.l1HTTPRPCAddr = *l1HTTPRPCAddr
	cfg.privateKey = *privateKey
	cfg.dockerImage = *dockerImage
	cfg.contractsEnvFile = *contractsEnvFile

	return cfg
}
