package main

// Flag names.
const (
	upgradeScriptFlag     = "upgrade_script"
	l1HTTPURLFlag         = "l1_http_url"
	privateKeyFlag        = "private_key"
	networkConfigAddrFlag = "network_config_addr"
	dockerImageFlag       = "docker_image"
)

// Returns a map of the flag usages.
// While we could just use constants instead of a map, this approach allows us to test that all the expected flags are defined.
func getFlagUsageMap() map[string]string {
	return map[string]string{
		upgradeScriptFlag:     "Hardhat script to run for upgrading L1 contracts (within the `/upgrade` folder, e.g. all_contracts_upgrade.ts)",
		l1HTTPURLFlag:         "Layer 1 network http RPC addr",
		privateKeyFlag:        "L1 and L2 private key used in the node",
		networkConfigAddrFlag: "L1 enclave registry contract address",
		dockerImageFlag:       "Docker image to run",
	}
}
