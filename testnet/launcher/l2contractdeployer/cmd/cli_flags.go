package main

// Flag names.
const (
	l1HostFlag                 = "l1_host"
	l1HTTPPortFlag             = "l1_http_port"
	privateKeyFlag             = "private_key"
	dockerImageFlag            = "docker_image"
	l2HostFlag                 = "l2_host"
	l2WSPortFlag               = "l2_ws_port"
	messageBusContractAddrFlag = "message_bus_contract_addr"
	l2privateKeyFlag           = "l2_private_key"
	l2HOCPrivateKeyFlag        = "l2_hoc_private_key"
	l2POCPrivateKeyFlag        = "l2_poc_private_key"
)

// Returns a map of the flag usages.
// While we could just use constants instead of a map, this approach allows us to test that all the expected flags are defined.
func getFlagUsageMap() map[string]string {
	return map[string]string{
		l1HostFlag:                 "Layer 1 network host addr",
		l1HTTPPortFlag:             "Layer 1 network HTTP port",
		privateKeyFlag:             "L1 and L2 private key used in the node",
		dockerImageFlag:            "Docker image to run",
		l2HostFlag:                 "Layer 2 network host addr",
		l2WSPortFlag:               "Layer 2 network WebSocket port",
		messageBusContractAddrFlag: "Message bus contract address",
		l2privateKeyFlag:           "Layer 2 private key",
		l2HOCPrivateKeyFlag:        "Layer 2 HOC contract private key",
		l2POCPrivateKeyFlag:        "Layer 2 POC contract private key",
	}
}
