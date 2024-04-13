package config

// NetworkConfig interface for configs that contain network details
type NetworkConfig interface {
	GetNetwork() *NetworkInputConfig
	SetNetwork(config *NetworkInputConfig)
}

type NodeConfig struct {
	NetworkConfig NetworkInputConfig `yaml:"networkConfig"`
	NodeDetails   NodeInputDetails   `yaml:"nodeDetails"`
	NodeSettings  NodeInputSettings  `yaml:"nodeSettings"`
}

type NodeInputDetails struct {
	NodeName          string `yaml:"nodeName"`
	HostID            string `yaml:"hostID"`
	PrivateKey        string `yaml:"privateKey"`
	L1WebsocketURL    string `yaml:"l1WebsocketURL"`
	P2pPublicAddress  string `yaml:"p2pPublicAddress"`
	ClientRPCPortHTTP int    `yaml:"clientRPCPortHTTP"`
	ClientRPCPortWS   int    `yaml:"clientRPCPortWS"`
}

type NodeInputSettings struct {
	NodeType              string `yaml:"nodeType"`
	IsSGXEnabled          bool   `yaml:"isSGXEnabled"`
	PccsAddr              string `yaml:"pccsAddr"`
	DebugNamespaceEnabled bool   `yaml:"debugNamespaceEnabled"`
	LogLevel              int    `yaml:"logLevel"`
	ProfilerEnabled       bool   `yaml:"profilerEnabled"`
	HostUseInMemoryDB     bool   `yaml:"useInMemoryDB"`
	HostPostgresDBHost    string `yaml:"hostPostgresDBHost"`
	HostImage             string `yaml:"hostImage"`
	EnclaveImage          string `yaml:"enclaveImage"`
	EdgelessDBImage       string `yaml:"edgelessDBImage"`
}

// NetworkInputConfig handles higher level configuration, note there is no need
// for an underlying `NetworkConfig` struct because typing only applies to the
// derived types for HostConfig and EnclaveConfig
type NetworkInputConfig struct {
	ManagementContractAddress string `yaml:"managementContractAddress"`
	MessageBusAddress         string `yaml:"messageBusAddress"`
	L1StartHash               string `yaml:"l1StartHash"`
	SequencerID               string `yaml:"sequencerID"`
}

func (p *HostInputConfig) GetNetwork() *NetworkInputConfig {
	return &NetworkInputConfig{
		ManagementContractAddress: p.ManagementContractAddress,
		MessageBusAddress:         p.MessageBusAddress,
		L1StartHash:               p.L1StartHash,
	}
}

func (p *HostInputConfig) SetNetwork(config *NetworkInputConfig) {
	if config != nil {
		p.ManagementContractAddress = config.ManagementContractAddress
		p.MessageBusAddress = config.MessageBusAddress
		p.L1StartHash = config.L1StartHash
	}
}

func (n *NodeConfig) GetNetwork() *NetworkInputConfig {
	return &n.NetworkConfig
}

func (n *NodeConfig) SetNetwork(config *NetworkInputConfig) {
	if config != nil {
		n.NetworkConfig = *config
	}
}
