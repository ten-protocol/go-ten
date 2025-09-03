package main

import (
	"flag"
)

// PauseAllContractsConfigCLI represents the configurations needed to pause all contracts over CLI
type PauseAllContractsConfigCLI struct {
	l1HTTPURL            string
	privateKey           string
	networkConfigAddr    string
	merkleMessageBusAddr string
	dockerImage          string
	action               string
}

// ParseConfigCLI returns a PauseAllContractsConfigCLI based the cli params and defaults.
func ParseConfigCLI() *PauseAllContractsConfigCLI {
	cfg := &PauseAllContractsConfigCLI{}
	flagUsageMap := getFlagUsageMap()

	l1HTTPURL := flag.String(l1HTTPURLFlag, "http://eth2network:8025", flagUsageMap[l1HTTPURLFlag])
	privateKey := flag.String(privateKeyFlag, "f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb", flagUsageMap[privateKeyFlag])
	networkConfigAddr := flag.String(networkConfigAddrFlag, "", flagUsageMap[networkConfigAddrFlag])
	merkleMessageBusAddr := flag.String(merkleMessageBusAddrFlag, "", flagUsageMap[merkleMessageBusAddrFlag])
	dockerImage := flag.String(dockerImageFlag, "testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest", flagUsageMap[dockerImageFlag])
	action := flag.String(actionFlag, "PAUSE", flagUsageMap[actionFlag])
	flag.Parse()

	cfg.l1HTTPURL = *l1HTTPURL
	cfg.privateKey = *privateKey
	cfg.networkConfigAddr = *networkConfigAddr
	cfg.merkleMessageBusAddr = *merkleMessageBusAddr
	cfg.dockerImage = *dockerImage
	cfg.action = *action

	return cfg
}
