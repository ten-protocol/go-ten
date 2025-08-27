package main

import (
	"flag"
)

// RoleTransferConfigCLI represents the configurations needed to transfer unpauser roles over CLI
type RoleTransferConfigCLI struct {
	l1HTTPURL            string
	privateKey           string
	networkConfigAddr    string
	multisigAddr         string
	merkleMessageBusAddr string
	dockerImage          string
}

// ParseConfigCLI returns a RoleTransferConfigCLI based the cli params and defaults.
func ParseConfigCLI() *RoleTransferConfigCLI {
	cfg := &RoleTransferConfigCLI{}
	flagUsageMap := getFlagUsageMap()

	l1HTTPURL := flag.String(l1HTTPURLFlag, "http://eth2network:8025", flagUsageMap[l1HTTPURLFlag])
	privateKey := flag.String(privateKeyFlag, "f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb", flagUsageMap[privateKeyFlag])
	networkConfigAddr := flag.String(networkConfigAddrFlag, "", flagUsageMap[networkConfigAddrFlag])
	multisigAddr := flag.String(multisigAddrFlag, "", flagUsageMap[multisigAddrFlag])
	merkleMessageBusAddr := flag.String(merkleMessageBusAddrFlag, "", flagUsageMap[merkleMessageBusAddrFlag])
	dockerImage := flag.String(dockerImageFlag, "testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest", flagUsageMap[dockerImageFlag])
	flag.Parse()

	cfg.l1HTTPURL = *l1HTTPURL
	cfg.privateKey = *privateKey
	cfg.networkConfigAddr = *networkConfigAddr
	cfg.multisigAddr = *multisigAddr
	cfg.merkleMessageBusAddr = *merkleMessageBusAddr
	cfg.dockerImage = *dockerImage

	return cfg
}
