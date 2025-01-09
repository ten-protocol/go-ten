package main

import (
	"flag"
)

// L2ContractDeployerConfigCLI represents the configurations passed into the deployer over CLI
type L2ContractDeployerConfigCLI struct {
	l1HTTPURL              string
	privateKey             string
	dockerImage            string
	l2Host                 string
	l2WSPort               int
	managementContractAddr string
	messageBusContractAddr string
	l2PrivateKey           string
	l2HOCPrivateKey        string
	l2POCPrivateKey        string
	faucetFunding          string
	challengePeriod        int
}

// ParseConfigCLI returns a NodeConfigCLI based the cli params and defaults.
func ParseConfigCLI() *L2ContractDeployerConfigCLI {
	cfg := &L2ContractDeployerConfigCLI{}
	flagUsageMap := getFlagUsageMap()

	l1HTTPURL := flag.String(l1HTTPURLFlag, "http://eth2network:8025", flagUsageMap[l1HTTPURLFlag])
	privateKey := flag.String(privateKeyFlag, "", flagUsageMap[privateKeyFlag])
	dockerImage := flag.String(dockerImageFlag, "testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest", flagUsageMap[dockerImageFlag])
	l2Host := flag.String(l2HostFlag, "", flagUsageMap[l2HostFlag])
	l2WSPort := flag.Int(l2WSPortFlag, 9000, flagUsageMap[l2WSPortFlag])
	managementContractAddr := flag.String(managementContractAddrFlag, "", flagUsageMap[managementContractAddrFlag])
	messageBusContractAddr := flag.String(messageBusContractAddrFlag, "", flagUsageMap[messageBusContractAddrFlag])
	l2PrivateKey := flag.String(l2privateKeyFlag, "", flagUsageMap[l2privateKeyFlag])
	l2HOCPrivateKey := flag.String(l2HOCPrivateKeyFlag, "", flagUsageMap[l2HOCPrivateKeyFlag])
	l2POCPrivateKey := flag.String(l2POCPrivateKeyFlag, "", flagUsageMap[l2POCPrivateKeyFlag])
	faucetFunds := flag.String(faucetFundingFlag, "0", flagUsageMap[faucetFundingFlag])
	challengePeriod := flag.Int(challengePeriodFlag, 0, flagUsageMap[challengePeriodFlag])

	flag.Parse()

	cfg.l1HTTPURL = *l1HTTPURL
	cfg.privateKey = *privateKey
	cfg.dockerImage = *dockerImage
	cfg.l2Host = *l2Host
	cfg.l2WSPort = *l2WSPort
	cfg.managementContractAddr = *managementContractAddr
	cfg.messageBusContractAddr = *messageBusContractAddr
	cfg.l2PrivateKey = *l2PrivateKey
	cfg.l2HOCPrivateKey = *l2POCPrivateKey
	cfg.l2POCPrivateKey = *l2HOCPrivateKey
	cfg.faucetFunding = *faucetFunds
	cfg.challengePeriod = *challengePeriod

	return cfg
}
