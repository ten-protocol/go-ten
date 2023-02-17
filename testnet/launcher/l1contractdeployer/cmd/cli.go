package main

import (
	"flag"
)

// L1ContractDeployerConfigCLI represents the configurations passed into the deployer over CLI
type L1ContractDeployerConfigCLI struct {
	l1Host           string
	l1HTTPPort       int
	privateKey       string
	dockerImage      string
	contractsEnvFile string
}

// ParseConfigCLI returns a NodeConfigCLI based the cli params and defaults.
func ParseConfigCLI() *L1ContractDeployerConfigCLI {
	cfg := &L1ContractDeployerConfigCLI{}
	flagUsageMap := getFlagUsageMap()

	l1Host := flag.String(l1HostFlag, "eth2network", flagUsageMap[l1HostFlag])
	l1HTTPPort := flag.Int(l1HTTPPortFlag, 8025, flagUsageMap[l1HTTPPortFlag])
	privateKey := flag.String(privateKeyFlag, "", flagUsageMap[privateKeyFlag])
	dockerImage := flag.String(dockerImageFlag, "testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest", flagUsageMap[dockerImageFlag])
	contractsEnvFile := flag.String(contractsEnvFileFlag, "", flagUsageMap[contractsEnvFileFlag])
	flag.Parse()

	cfg.l1Host = *l1Host
	cfg.l1HTTPPort = *l1HTTPPort
	cfg.privateKey = *privateKey
	cfg.dockerImage = *dockerImage
	cfg.contractsEnvFile = *contractsEnvFile

	return cfg
}
