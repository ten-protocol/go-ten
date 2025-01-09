package main

// Flag names.
const (
	l1HTTPURLFlag              = "l1_http_url"
	privateKeyFlag             = "private_key"
	dockerImageFlag            = "docker_image"
	l2HostFlag                 = "l2_host"
	l2WSPortFlag               = "l2_ws_port"
	managementContractAddrFlag = "management_contract_addr"
	messageBusContractAddrFlag = "message_bus_contract_addr"
	l2privateKeyFlag           = "l2_private_key"
	l2HOCPrivateKeyFlag        = "l2_hoc_private_key"
	l2POCPrivateKeyFlag        = "l2_poc_private_key"
	faucetFundingFlag          = "faucet_funds"
	challengePeriodFlag        = "challenge_period"
)

// Returns a map of the flag usages.
// While we could just use constants instead of a map, this approach allows us to test that all the expected flags are defined.
func getFlagUsageMap() map[string]string {
	return map[string]string{
		l1HTTPURLFlag:              "Layer 1 network http RPC addr",
		privateKeyFlag:             "L1 and L2 private key used in the node",
		dockerImageFlag:            "Docker image to run",
		l2HostFlag:                 "Layer 2 network host addr",
		l2WSPortFlag:               "Layer 2 network WebSocket port",
		managementContractAddrFlag: "Management contract address",
		messageBusContractAddrFlag: "Message bus contract address",
		l2privateKeyFlag:           "Layer 2 private key",
		l2HOCPrivateKeyFlag:        "Layer 2 HOC contract private key",
		l2POCPrivateKeyFlag:        "Layer 2 POC contract private key",
		faucetFundingFlag:          "How much funds should the faucet account receive",
		challengePeriodFlag:        "Delay when adding message bus root",
	}
}
