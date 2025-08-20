package main

import (
	"flag"
)

// MultisigSetupConfigCLI represents the configurations needed to setup multisig control over CLI
type MultisigSetupConfigCLI struct {
	l1HTTPURL         string
	privateKey        string
	networkConfigAddr string
	multisigAddress   string
	dockerImage       string
}

// ParseConfigCLI returns a MultisigSetupConfigCLI based the cli params and defaults.
func ParseConfigCLI() *MultisigSetupConfigCLI {
	cfg := &MultisigSetupConfigCLI{}
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
