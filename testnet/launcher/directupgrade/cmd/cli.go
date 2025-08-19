package main

import (
	"flag"
)

// DirectUpgradeConfigCLI represents the configurations needed to run direct upgrades over CLI
type DirectUpgradeConfigCLI struct {
	l1HTTPURL         string
	privateKey        string
	networkConfigAddr string
	multisigAddress   string
	dockerImage       string
}

// ParseConfigCLI returns a DirectUpgradeConfigCLI based the cli params and defaults.
func ParseConfigCLI() *DirectUpgradeConfigCLI {
	cfg := &DirectUpgradeConfigCLI{}
	flagUsageMap := getFlagUsageMap()

	l1HTTPURL := flag.String(l1HTTPURLFlag, "http://eth2network:8025", flagUsageMap[l1HTTPURLFlag])
	privateKey := flag.String(privateKeyFlag, "f52e5418e349dccdda29b6ac8b0abe6576bb7713886aa85abea6181ba731f9bb", flagUsageMap[privateKeyFlag])
	networkConfigAddr := flag.String(networkConfigAddrFlag, "", flagUsageMap[networkConfigAddrFlag])
	multisigAddress := flag.String(multisigAddressFlag, "", flagUsageMap[multisigAddressFlag])
	dockerImage := flag.String(dockerImageFlag, "testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest", flagUsageMap[dockerImageFlag])
	flag.Parse()

	cfg.l1HTTPURL = *l1HTTPURL
	cfg.privateKey = *privateKey
	cfg.networkConfigAddr = *networkConfigAddr
	cfg.multisigAddress = *multisigAddress
	cfg.dockerImage = *dockerImage

	return cfg
}