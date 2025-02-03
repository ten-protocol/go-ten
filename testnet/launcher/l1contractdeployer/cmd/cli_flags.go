package main

// Flag names.
const (
	l1HTTPURLFlag        = "l1_http_url"
	privateKeyFlag       = "private_key"
	dockerImageFlag      = "docker_image"
	contractsEnvFileFlag = "contracts_env_file"
	azureKeyVaultURLFlag = "azure_key_vault_url"
)

// Returns a map of the flag usages.
// While we could just use constants instead of a map, this approach allows us to test that all the expected flags are defined.
func getFlagUsageMap() map[string]string {
	return map[string]string{
		l1HTTPURLFlag:        "Layer 1 network http RPC addr",
		privateKeyFlag:       "L1 and L2 private key used in the node",
		dockerImageFlag:      "Docker image to run",
		contractsEnvFileFlag: "If set, it will write the contract addresses to the file",
		azureKeyVaultURLFlag: "If set, contract addresses will be stored in Azure Key Vault",
	}
}
