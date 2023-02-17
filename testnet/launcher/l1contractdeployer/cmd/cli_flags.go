package main

// Flag names.
const (
	l1HostFlag           = "l1_host"
	l1HTTPPortFlag       = "l1_http_port"
	privateKeyFlag       = "private_key"
	dockerImageFlag      = "docker_image"
	contractsEnvFileFlag = "contracts_env_file"
)

// Returns a map of the flag usages.
// While we could just use constants instead of a map, this approach allows us to test that all the expected flags are defined.
func getFlagUsageMap() map[string]string {
	return map[string]string{
		l1HostFlag:           "Layer 1 network host addr",
		l1HTTPPortFlag:       "Layer 1 network HTTP port",
		privateKeyFlag:       "L1 and L2 private key used in the node",
		dockerImageFlag:      "Docker image to run",
		contractsEnvFileFlag: "If set, it will write the contract addresses to the file",
	}
}
