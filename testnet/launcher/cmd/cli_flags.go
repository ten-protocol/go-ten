package main

// Flag names.
const (
	validatorEnclaveDockerImageFlag = "validator-enclave-docker-image"
	validatorEnclaveDebugFlag       = "validator-enclave-debug"
	sequencerEnclaveDockerImageFlag = "sequencer-enclave-docker-image"
	sequencerEnclaveDebugFlag       = "sequencer-enclave-debug"
	isSGXEnabledFlag                = "is-sgx-enabled"
)

// Returns a map of the flag usages.
// While we could just use constants instead of a map, this approach allows us to test that all the expected flags are defined.
func getFlagUsageMap() map[string]string {
	return map[string]string{
		validatorEnclaveDockerImageFlag: "The docker image that runs the validator enclave",
		validatorEnclaveDebugFlag:       "Enables the use of DLV to debug the validator enclave",
		sequencerEnclaveDockerImageFlag: "The docker image that runs the sequencer enclave",
		sequencerEnclaveDebugFlag:       "Enables the use of DLV to debug the sequencer enclave",
		isSGXEnabledFlag:                "Enables the SGX usage",
	}
}
