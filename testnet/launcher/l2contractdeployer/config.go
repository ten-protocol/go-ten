package l2contractdeployer

import (
	"fmt"

	"github.com/ten-protocol/go-ten/go/config"
	"gopkg.in/yaml.v2"
)

// Config holds the properties that configure the package
type Config struct {
	L1HTTPURL              string `yaml:"l1_http_url"`
	L1PrivateKey           string `yaml:"l1_private_key"`
	L2Port                 int    `yaml:"l2_port"`
	L2Host                 string `yaml:"l2_host"`
	L2PrivateKey           string `yaml:"l2_private_key"`
	EnclaveRegistryAddress string `yaml:"enclave_registry_address"`
	CrossChainAddress      string `yaml:"cross_chain_address"`
	DaRegistryAddress      string `yaml:"da_registry_address"`
	NetworkConfigAddress   string `yaml:"network_config_address"`
	MessageBusAddress      string `yaml:"message_bus_address"`
	DockerImage            string `yaml:"docker_image"`
	FaucetPrefundAmount    string `yaml:"faucet_prefund_amount"`
	DebugEnabled           bool   `yaml:"debug_enabled"`
}

func NewContractDeployerConfig(tenCfg *config.TenConfig) *Config {
	return &Config{
		L1HTTPURL:              tenCfg.Deployment.L1Deploy.RPCAddress,
		L1PrivateKey:           tenCfg.Deployment.L1Deploy.DeployerPK,
		L2Port:                 tenCfg.Deployment.L2Deploy.WSPort,
		L2Host:                 tenCfg.Deployment.L2Deploy.RPCAddress,
		L2PrivateKey:           tenCfg.Deployment.L2Deploy.DeployerPK,
		EnclaveRegistryAddress: tenCfg.Network.L1.L1Contracts.EnclaveRegistryContract.Hex(),
		CrossChainAddress:      tenCfg.Network.L1.L1Contracts.CrossChainContract.Hex(),
		DaRegistryAddress:      tenCfg.Deployment.L1Deploy.DARegistry.Hex(),
		NetworkConfigAddress:   tenCfg.Network.L1.L1Contracts.NetworkConfigContract.Hex(),
		MessageBusAddress:      tenCfg.Network.L1.L1Contracts.MessageBusContract.Hex(),
		DockerImage:            tenCfg.Deployment.DockerImage,
		FaucetPrefundAmount:    tenCfg.Deployment.L2Deploy.FaucetPrefund,
		DebugEnabled:           tenCfg.Deployment.DebugEnabled,
	}
}

func (c *Config) Obfuscate() string {
	configCopy := *c

	// Mask both private keys
	if configCopy.L1PrivateKey != "" {
		configCopy.L1PrivateKey = "****"
	}
	if configCopy.L2PrivateKey != "" {
		configCopy.L2PrivateKey = "****"
	}

	output, err := yaml.Marshal(&configCopy)
	if err != nil {
		return fmt.Sprintf("Error marshaling config to YAML: %v", err)
	}
	return string(output)
}
