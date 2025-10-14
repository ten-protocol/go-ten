package main

import (
	"flag"
)

// FundsRecoveryConfigCLI represents the configurations passed into the deployer over CLI
type FundsRecoveryConfigCLI struct {
	l1HTTPURL   string
	privateKey  string
	dockerImage string
	networkConfigAddr string
	receiverAddress string
}

// ParseConfigCLI returns a NodeConfigCLI based the cli params and defaults.
func ParseConfigCLI() *FundsRecoveryConfigCLI {
	cfg := &FundsRecoveryConfigCLI{}
	flagUsageMap := getFlagUsageMap()

	l1HTTPURL := flag.String(l1HTTPURLFlag, "http://eth2network:8025 ", flagUsageMap[l1HTTPURLFlag])
	privateKey := flag.String(privateKeyFlag, "f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb", flagUsageMap[privateKeyFlag])
	dockerImage := flag.String(dockerImageFlag, "testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest", flagUsageMap[dockerImageFlag])
	networkConfigAddr := flag.String(networkConfigAddrFlag, "", flagUsageMap[networkConfigAddrFlag])
	receiverAddress := flag.String(receiverAddressFlag, "0xeA052c9635F1647A8a199c2315B9A66ce7d1e2a7", flagUsageMap[receiverAddressFlag]) //gnosis sepolia eth address is default

	flag.Parse()

	cfg.l1HTTPURL = *l1HTTPURL
	cfg.privateKey = *privateKey
	cfg.dockerImage = *dockerImage
	cfg.networkConfigAddr = *networkConfigAddr
	cfg.receiverAddress = *receiverAddress

	return cfg
}
