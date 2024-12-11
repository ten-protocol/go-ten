package main

import (
	"flag"
)

// L1GrantSequencersConfigCLI represents the configurations needed to grant enclaves sequencer roles over CLI
type L1GrantSequencersConfigCLI struct {
	l1HTTPURL           string
	privateKey          string
	mgmtContractAddress string
	enclaveIDs          string
	dockerImage         string
	sequencerURL        string
}

// ParseConfigCLI returns a NodeConfigCLI based the cli params and defaults.
func ParseConfigCLI() *L1GrantSequencersConfigCLI {
	cfg := &L1GrantSequencersConfigCLI{}
	flagUsageMap := getFlagUsageMap()

	l1HTTPURL := flag.String(l1HTTPURLFlag, "http://eth2network:8025", flagUsageMap[l1HTTPURLFlag])
	privateKey := flag.String(privateKeyFlag, "", flagUsageMap[privateKeyFlag])
	mgmtContractAddress := flag.String(mgmtContractAddressFlag, "", flagUsageMap[mgmtContractAddressFlag])
	enclaveIDs := flag.String(enclaveIDsFlag, "", flagUsageMap[enclaveIDsFlag])
	dockerImage := flag.String(dockerImageFlag, "testnetobscuronet.azurecr.io/obscuronet/hardhatdeployer:latest", flagUsageMap[dockerImageFlag])
	sequencerURL := flag.String(sequencerURLFlag, "", flagUsageMap[sequencerURLFlag])
	flag.Parse()

	cfg.l1HTTPURL = *l1HTTPURL
	cfg.privateKey = *privateKey
	cfg.mgmtContractAddress = *mgmtContractAddress
	cfg.enclaveIDs = *enclaveIDs
	cfg.dockerImage = *dockerImage
	cfg.sequencerURL = *sequencerURL

	return cfg
}
