package main

import (
	"flag"
)

// FundsRecoveryConfigCLI represents the configurations passed into the deployer over CLI
type FundsRecoveryConfigCLI struct {
	l1HTTPURL   string
	privateKey  string
	dockerImage string
}

// ParseConfigCLI returns a NodeConfigCLI based the cli params and defaults.
func ParseConfigCLI() *FundsRecoveryConfigCLI {
	cfg := &FundsRecoveryConfigCLI{}
	flagUsageMap := getFlagUsageMap()

	l1HTTPURL := flag.String(l1HTTPURLFlag, "", flagUsageMap[l1HTTPURLFlag])
	privateKey := flag.String(privateKeyFlag, "", flagUsageMap[privateKeyFlag])
	dockerImage := flag.String(dockerImageFlag, "", flagUsageMap[dockerImageFlag])

	flag.Parse()

	cfg.l1HTTPURL = *l1HTTPURL
	cfg.privateKey = *privateKey
	cfg.dockerImage = *dockerImage

	return cfg
}
