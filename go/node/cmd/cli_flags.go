package main

// Flag names.
const (
	nodeNameFlag                = "node_name"
	nodeTypeFlag                = "node_type"
	isGenesisFlag               = "is_genesis"
	hostIDFlag                  = "host_id"
	isSGXEnabledFlag            = "is_sgx_enabled"
	enclaveDockerImageFlag      = "enclave_docker_image"
	hostDockerImageFlag         = "host_docker_image"
	l1WebsocketURLFlag          = "l1_ws_url"
	hostHTTPPortFlag            = "host_http_port"
	hostWSPortFlag              = "host_ws_port"
	hostP2PPortFlag             = "host_p2p_port"
	hostP2PHostFlag             = "host_p2p_host"
	hostP2PPublicAddrFlag       = "host_public_p2p_addr"
	enclaveHTTPPortFlag         = "enclave_http_port"
	enclaveWSPortFlag           = "enclave_WS_port"
	privateKeyFlag              = "private_key"
	sequencerIDFlag             = "sequencer_id"
	managementContractAddrFlag  = "management_contract_addr"
	messageBusContractAddrFlag  = "message_bus_contract_addr"
	l1StartBlockFlag            = "l1_start"
	pccsAddrFlag                = "pccs_addr"
	edgelessDBImageFlag         = "edgeless_db_image"
	isDebugNamespaceEnabledFlag = "is_debug_namespace_enabled"
	logLevelFlag                = "log_level"
	isInboundP2PDisabledFlag    = "is_inbound_p2p_disabled"
	batchIntervalFlag           = "batch_interval"
	rollupIntervalFlag          = "rollup_interval"
	l1ChainIDFlag               = "l1_chain_id"
)

// Returns a map of the flag usages.
// While we could just use constants instead of a map, this approach allows us to test that all the expected flags are defined.
func getFlagUsageMap() map[string]string {
	return map[string]string{
		nodeNameFlag:                "Specifies the node base name",
		nodeTypeFlag:                "The node's type (e.g. sequencer, validator)",
		isGenesisFlag:               "Wether the node is the genesis node of the network",
		hostIDFlag:                  "The 20 bytes of the address of the Obscuro host this enclave serves",
		isSGXEnabledFlag:            "Whether the it should run on an SGX is enabled CPU",
		enclaveDockerImageFlag:      "Docker image for the enclave",
		hostDockerImageFlag:         "Docker image for the host",
		l1WebsocketURLFlag:          "Layer 1 websocket RPC address",
		hostP2PPortFlag:             "Hosts p2p bound port",
		hostP2PPublicAddrFlag:       "Hosts public p2p host.",
		hostP2PHostFlag:             "Hosts p2p bound addr",
		enclaveHTTPPortFlag:         "Enclave's http bound port",
		enclaveWSPortFlag:           "Enclave's WS bound port",
		privateKeyFlag:              "L1 and L2 private key used in the node",
		sequencerIDFlag:             "The 20 bytes of the address of the sequencer for this network",
		managementContractAddrFlag:  "The management contract address on the L1",
		messageBusContractAddrFlag:  "The address of the L1 message bus contract owned by the management contract.",
		l1StartBlockFlag:            "The block hash on the L1 where the management contract was deployed",
		pccsAddrFlag:                "Sets the PCCS address",
		edgelessDBImageFlag:         "Sets the edgelessdb image",
		hostHTTPPortFlag:            "Host HTTPs bound port",
		hostWSPortFlag:              "Host WebSocket bound port",
		isDebugNamespaceEnabledFlag: "Enables the debug namespace for both enclave and host",
		logLevelFlag:                "Sets the log level 1-Error, 5-Trace",
		isInboundP2PDisabledFlag:    "Disables inbound p2p (for testing)",
		batchIntervalFlag:           "Duration between each batch. Can be formatted like 500ms or 1s",
		rollupIntervalFlag:          "Duration between each rollup. Can be formatted like 500ms or 1s",
		l1ChainIDFlag:               "Chain ID of the L1 network",
	}
}
