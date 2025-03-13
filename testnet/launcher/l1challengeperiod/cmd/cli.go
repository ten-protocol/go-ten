package main

import (
	"flag"
)

// L1ChallengePeriodConfigCLI represents the configurations needed to grant enclaves sequencer roles over CLI
type L1ChallengePeriodConfigCLI struct {
	l1HTTPURL         string
	privateKey        string
	crossChainAddress string
	dockerImage       string
	challengePeriod   int
}

// ParseConfigCLI returns a NodeConfigCLI based the cli params and defaults.
func ParseConfigCLI() *L1ChallengePeriodConfigCLI {
	cfg := &L1ChallengePeriodConfigCLI{}
	flagUsageMap := getFlagUsageMap()

	l1HTTPURL := flag.String(l1HTTPURLFlag, "http://eth2network:8025", flagUsageMap[l1HTTPURLFlag])
	privateKey := flag.String(privateKeyFlag, "", flagUsageMap[privateKeyFlag])
	crossChainAddress := flag.String(crossChainContractAddressFlag, "", flagUsageMap[crossChainContractAddressFlag])
	dockerImage := flag.String(dockerImageFlag, "testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest", flagUsageMap[dockerImageFlag])
	challengePeriod := flag.Int(challengePeriodFlag, 0, flagUsageMap[challengePeriodFlag])
	flag.Parse()

	cfg.l1HTTPURL = *l1HTTPURL
	cfg.privateKey = *privateKey
	cfg.crossChainAddress = *crossChainAddress
	cfg.dockerImage = *dockerImage
	cfg.challengePeriod = *challengePeriod

	return cfg
}
