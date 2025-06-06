package main

// Flag names.
const (
	l1HTTPURLFlag   = "l1_http_url"
	privateKeyFlag  = "private_key"
	dockerImageFlag = "docker_image"
)

// Returns a map of the flag usages.
// While we could just use constants instead of a map, this approach allows us to test that all the expected flags are defined.
func getFlagUsageMap() map[string]string {
	return map[string]string{
		l1HTTPURLFlag:   "Layer 1 network http RPC addr",
		privateKeyFlag:  "L1 mgmt contract owning key",
		dockerImageFlag: "Docker image to run",
	}
}
