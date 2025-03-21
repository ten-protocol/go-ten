package main

// Flag names.
const (
	l1HTTPURLFlag           = "l1_http_url"
	privateKeyFlag          = "private_key"
	enclaveRegistryAddrFlag = "enclave_registry_addr"
	enclaveIDsFlag          = "enclave_ids"
	dockerImageFlag         = "docker_image"
	contractsEnvFileFlag    = "contracts_env_file"
	sequencerURLFlag        = "sequencer_url"
)

// Returns a map of the flag usages.
// While we could just use constants instead of a map, this approach allows us to test that all the expected flags are defined.
func getFlagUsageMap() map[string]string {
	return map[string]string{
		l1HTTPURLFlag:           "Layer 1 network http RPC addr",
		privateKeyFlag:          "L1 and L2 private key used in the node",
		enclaveRegistryAddrFlag: "L1 enclave registry contract address",
		enclaveIDsFlag:          "List of enclave public keys to grant sequencer role",
		dockerImageFlag:         "Docker image to run",
		contractsEnvFileFlag:    "If set, it will write the contract addresses to the file",
		sequencerURLFlag:        "Sequencer RPC URL to fetch enclave IDs (required if enclaveIDs are not provided)",
	}
}
