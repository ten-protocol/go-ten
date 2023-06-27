package main

import (
	"flag"
)

// TestnetConfigCLI represents the configurations passed into the testnet over CLI
type TestnetConfigCLI struct {
	validatorEnclaveDockerImage string
	validatorEnclaveDebug       bool
	sequencerEnclaveDockerImage string
	sequencerEnclaveDebug       bool

	isSGXEnabled bool

	edgelessDBDockerImage string
}

// ParseConfigCLI returns a NodeConfigCLI based the cli params and defaults.
func ParseConfigCLI() *TestnetConfigCLI {
	cfg := &TestnetConfigCLI{}
	flagUsageMap := getFlagUsageMap()

	validatorEnclaveDockerImage := flag.String(validatorEnclaveDockerImageFlag, "testnetobscuronet.azurecr.io/obscuronet/enclave:latest", flagUsageMap[validatorEnclaveDockerImageFlag])
	validatorEnclaveDebug := flag.Bool(validatorEnclaveDebugFlag, false, flagUsageMap[validatorEnclaveDebugFlag])
	sequencerEnclaveDockerImage := flag.String(sequencerEnclaveDockerImageFlag, "testnetobscuronet.azurecr.io/obscuronet/enclave:latest", flagUsageMap[sequencerEnclaveDockerImageFlag])
	sequencerEnclaveDebug := flag.Bool(sequencerEnclaveDebugFlag, false, flagUsageMap[sequencerEnclaveDebugFlag])

	edgelessDBDockerImage := flag.String(edgelessDBDockerImageFlag, "ghcr.io/edgelesssys/edgelessdb-sgx-4gb:v0.3.2", flagUsageMap[edgelessDBDockerImageFlag])

	isSGXEnabled := flag.Bool(isSGXEnabledFlag, false, flagUsageMap[isSGXEnabledFlag])
	flag.Parse()

	cfg.validatorEnclaveDockerImage = *validatorEnclaveDockerImage
	cfg.sequencerEnclaveDockerImage = *sequencerEnclaveDockerImage
	cfg.validatorEnclaveDebug = *validatorEnclaveDebug
	cfg.sequencerEnclaveDebug = *sequencerEnclaveDebug

	cfg.isSGXEnabled = *isSGXEnabled

	cfg.edgelessDBDockerImage = *edgelessDBDockerImage

	return cfg
}
