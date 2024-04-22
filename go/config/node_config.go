package config

import (
	"fmt"
	"reflect"
	"strings"
)

// NetworkConfig interface for configs that contain network details
type NetworkConfig interface {
	// GetNetwork returns a subset struct representing network
	GetNetwork() *NetworkInputConfig
	// SetNetwork sets a network struct
	SetNetwork(config *NetworkInputConfig)
}

type NodeConfig struct {
	NetworkConfig NetworkInputConfig `yaml:"networkConfig"`
	NodeDetails   NodeInputDetails   `yaml:"nodeDetails"`
	NodeSettings  NodeInputSettings  `yaml:"nodeSettings"`
	NodeImages    NodeInputImages    `yaml:"nodeImages"`
}

type NodeInputDetails struct {
	NodeName            string   `yaml:"nodeName"`
	HostID              string   `yaml:"hostID"`
	PrivateKey          string   `yaml:"privateKey"`
	L1WebsocketURL      string   `yaml:"l1WebsocketURL"`
	P2pPublicAddress    string   `yaml:"p2pPublicAddress"`
	P2pBindAddress      string   `yaml:"p2pBindAddress"`
	ClientRPCPortHTTP   int      `yaml:"clientRPCPortHTTP"`
	ClientRPCPortWS     int      `yaml:"clientRPCPortWS"`
	EnclaveRPCAddresses []string `yaml:"enclaveRPCAddresses"`
}

type NodeInputSettings struct {
	NodeType              string `yaml:"nodeType"`
	IsSGXEnabled          bool   `yaml:"isSGXEnabled"`
	IsGenesis             bool   `yaml:"isGenesis"`
	PccsAddr              string `yaml:"pccsAddr"`
	DebugNamespaceEnabled bool   `yaml:"debugNamespaceEnabled"`
	LogLevel              int    `yaml:"logLevel"`
	ProfilerEnabled       bool   `yaml:"profilerEnabled"`
	UseInMemoryDB         bool   `yaml:"useInMemoryDB"`
	PostgresDBHost        string `yaml:"postgresDBHost"`
}

type NodeInputImages struct {
	HostImage       string `yaml:"hostImage"`
	EnclaveImage    string `yaml:"enclaveImage"`
	EdgelessDBImage string `yaml:"edgelessDBImage"`
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

func (n *TestnetConfig) GetNetwork() *NetworkInputConfig {
	return &n.Network
}

func (n *TestnetConfig) SetNetwork(config *NetworkInputConfig) {
	if config != nil {
		n.Network = *config
	}
}

// GetConfigAsEnvVars returns a set of environment variables from overrides/flags in a NodeConfig that is valid for
// a TypConfig (Host | Enclave)
// This makes it easier to pass custom configurations to docker containers.
func (n *NodeConfig) GetConfigAsEnvVars(t TypeConfig) map[string]string {
	envVars := make(map[string]string)
	// Reflect on the NodeConfig to access its fields.
	val := reflect.ValueOf(n).Elem()

	// Iterate over each field in NodeConfig.
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldValue := val.Field(i)

		// Skip NodeImages section based on field name.
		if field.Name == "NodeImages" {
			continue
		}

		// Now, reflect on each sub-struct like NetworkConfig, NodeDetails, NodeSettings
		if fieldValue.Kind() == reflect.Struct {
			for j := 0; j < fieldValue.NumField(); j++ {
				subField := fieldValue.Type().Field(j)
				subFieldValue := fieldValue.Field(j)
				yamlTag := subField.Tag.Get("yaml")

				// Check if this field's yaml tag is part of the relevant configuration
				if FlagsByService[t][yamlTag] {
					// Create environment variable key by transforming yaml tag to uppercase
					envKey := strings.ToUpper(strings.ReplaceAll(yamlTag, "-", "_"))
					envVars[envKey] = fmt.Sprintf("%v", subFieldValue.Interface())
				}
			}
		}
	}
	return envVars
}
