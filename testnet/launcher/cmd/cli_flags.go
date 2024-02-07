package main

// Flag names.
const (
	validatorEnclaveDockerImageFlag = "validator-enclave-docker-image"
	validatorEnclaveDebugFlag       = "validator-enclave-debug"
	sequencerEnclaveDockerImageFlag = "sequencer-enclave-docker-image"
	sequencerEnclaveDebugFlag       = "sequencer-enclave-debug"
	contractDeployerDockerImageFlag = "contract-deployer-docker-image"
	contractDeployerDebugFlag       = "contract-deployer-debug"
	isSGXEnabledFlag                = "is-sgx-enabled"
	logLevelFlag                    = "log-level"
)

// Returns a map of the flag usages.
// While we could just use constants instead of a map, this approach allows us to test that all the expected flags are defined.
func getFlagUsageMap() map[string]string {
	return map[string]string{
		validatorEnclaveDockerImageFlag: "The docker image that runs the validator enclave",
		validatorEnclaveDebugFlag:       "Enables the use of DLV to debug the validator enclave",
		sequencerEnclaveDockerImageFlag: "The docker image that runs the sequencer enclave",
		sequencerEnclaveDebugFlag:       "Enables the use of DLV to debug the sequencer enclave",
		contractDeployerDockerImageFlag: "The docker image that runs the contract deployer",
		contractDeployerDebugFlag:       "Enables the use of node inspector to debug the contract deployer",
		isSGXEnabledFlag:                "Enables the SGX usage",
		logLevelFlag:                    "Log level for all network",
	}
}
