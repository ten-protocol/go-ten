package l2contractdeployer

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

// Option is a function that applies configs to a Config Object
type Option = func(c *Config)

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

func NewContractDeployerConfig(opts ...Option) *Config {
	defaultConfig := &Config{
		FaucetPrefundAmount: "10000",
	}

	for _, opt := range opts {
		opt(defaultConfig)
	}

	return defaultConfig
}

func WithL1HTTPURL(s string) Option {
	return func(c *Config) {
		c.L1HTTPURL = s
	}
}

func WithL1PrivateKey(s string) Option {
	return func(c *Config) {
		c.L1PrivateKey = s
	}
}

func WithL2WSPort(i int) Option {
	return func(c *Config) {
		c.L2Port = i
	}
}

func WithL2Host(s string) Option {
	return func(c *Config) {
		c.L2Host = s
	}
}

func WithEnclaveRegistryAddress(s string) Option {
	return func(c *Config) {
		c.EnclaveRegistryAddress = s
	}
}

func WithCrossChainAddress(s string) Option {
	return func(c *Config) {
		c.CrossChainAddress = s
	}
}

func WithDataAvailabilityRegistryAddress(s string) Option {
	return func(c *Config) {
		c.DaRegistryAddress = s
	}
}

func WithNetworkConfigAddress(s string) Option {
	return func(c *Config) {
		c.NetworkConfigAddress = s
	}
}

func WithMessageBusContractAddress(s string) Option {
	return func(c *Config) {
		c.MessageBusAddress = s
	}
}

func WithL2PrivateKey(s string) Option {
	return func(c *Config) {
		c.L2PrivateKey = s
	}
}

func WithDockerImage(s string) Option {
	return func(c *Config) {
		c.DockerImage = s
	}
}

func WithFaucetFunds(f string) Option {
	return func(c *Config) {
		c.FaucetPrefundAmount = f
	}
}

func WithDebugEnabled(b bool) Option {
	return func(c *Config) {
		c.DebugEnabled = b
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
