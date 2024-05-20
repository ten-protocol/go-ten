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
	EnclaveConfig EnclaveInputConfig `yaml:"enclave"`
	HostConfig    HostInputConfig    `yaml:"host"`
}

type NodeInputDetails struct {
	NodeName   string `yaml:"nodeName"`
	HostID     string `yaml:"hostID"`
	PrivateKey string `yaml:"privateKey"`
}

type NodeInputSettings struct {
	NodeType     string `yaml:"nodeType"`
	IsSGXEnabled bool   `yaml:"isSGXEnabled"`
	PccsAddr     string `yaml:"pccsAddr"`
	EnclaveDebug bool   `yaml:"enclaveDebug"`
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
	L1StartHash               string `yaml:"l1StartHash"`
	L1ChainID                 string `yaml:"l1ChainID"`
	ManagementContractAddress string `yaml:"managementContractAddress"`
	MessageBusAddress         string `yaml:"messageBusAddress"`
	SequencerP2PAddress       string `yaml:"sequencerP2PAddress"`
}

func (p *EnclaveInputConfig) GetNetwork() *NetworkInputConfig {
	return &NetworkInputConfig{
		ManagementContractAddress: p.ManagementContractAddress,
		MessageBusAddress:         p.MessageBusAddress,
		SequencerP2PAddress:       p.SequencerP2PAddress,
	}
}

func (p *EnclaveInputConfig) SetNetwork(config *NetworkInputConfig) {
	if config != nil {
		p.ManagementContractAddress = config.ManagementContractAddress
		p.MessageBusAddress = config.MessageBusAddress
		p.SequencerP2PAddress = config.SequencerP2PAddress
	}
}

func (p *HostInputConfig) GetNetwork() *NetworkInputConfig {
	return &NetworkInputConfig{
		ManagementContractAddress: p.ManagementContractAddress,
		MessageBusAddress:         p.MessageBusAddress,
		L1StartHash:               p.L1StartHash,
		SequencerP2PAddress:       p.SequencerP2PAddress,
	}
}

func (p *HostInputConfig) SetNetwork(config *NetworkInputConfig) {
	if config != nil {
		p.ManagementContractAddress = config.ManagementContractAddress
		p.MessageBusAddress = config.MessageBusAddress
		p.L1StartHash = config.L1StartHash
		p.SequencerP2PAddress = config.SequencerP2PAddress
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

func (n *NetworkInputConfig) GetNetwork() *NetworkInputConfig {
	return n
}

func (n *NetworkInputConfig) SetNetwork(config *NetworkInputConfig) {
	if config != nil {
		n.L1StartHash = config.L1StartHash
		n.L1ChainID = config.L1ChainID
		n.ManagementContractAddress = config.ManagementContractAddress
		n.MessageBusAddress = config.MessageBusAddress
		n.SequencerP2PAddress = config.SequencerP2PAddress
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
