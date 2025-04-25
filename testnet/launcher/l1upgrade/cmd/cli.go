package main

import (
	"flag"
)

// L1UpgradeContractsConfigCLI represents the configurations needed to upgrade L1 contracts over CLI
type L1UpgradeContractsConfigCLI struct {
	l1HTTPURL         string
	privateKey        string
	networkConfigAddr string
	dockerImage       string
}

// ParseConfigCLI returns a NodeConfigCLI based the cli params and defaults.
func ParseConfigCLI() *L1UpgradeContractsConfigCLI {
	cfg := &L1UpgradeContractsConfigCLI{}
	flagUsageMap := getFlagUsageMap()

	l1HTTPURL := flag.String(l1HTTPURLFlag, "http://eth2network:8025", flagUsageMap[l1HTTPURLFlag])
	privateKey := flag.String(privateKeyFlag, "", flagUsageMap[privateKeyFlag])
	neworkConfigAddr := flag.String(networkConfigAddrFlag, "", flagUsageMap[networkConfigAddrFlag])
	dockerImage := flag.String(dockerImageFlag, "testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest", flagUsageMap[dockerImageFlag])
	flag.Parse()

	cfg.l1HTTPURL = *l1HTTPURL
	cfg.privateKey = *privateKey
	cfg.networkConfigAddr = *neworkConfigAddr
	cfg.dockerImage = *dockerImage

	return cfg
}
