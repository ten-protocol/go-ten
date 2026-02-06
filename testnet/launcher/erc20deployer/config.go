package erc20deployer

import "fmt"

type Config struct {
	tokenName          string
	tokenSymbol        string
	tokenDecimals      string
	tokenSupply        string
	l1HTTPURL          string
	privateKey         string
	dockerImage        string
	networkConfigAddr  string
}

func NewConfig(opts ...ConfigOption) *Config {
	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

type ConfigOption func(*Config)

func WithTokenName(name string) ConfigOption {
	return func(c *Config) {
		c.tokenName = name
	}
}

func WithTokenSymbol(symbol string) ConfigOption {
	return func(c *Config) {
		c.tokenSymbol = symbol
	}
}

func WithTokenDecimals(decimals string) ConfigOption {
	return func(c *Config) {
		c.tokenDecimals = decimals
	}
}

func WithTokenSupply(supply string) ConfigOption {
	return func(c *Config) {
		c.tokenSupply = supply
	}
}

func WithL1HTTPURL(url string) ConfigOption {
	return func(c *Config) {
		c.l1HTTPURL = url
	}
}

func WithPrivateKey(key string) ConfigOption {
	return func(c *Config) {
		c.privateKey = key
	}
}

func WithDockerImage(image string) ConfigOption {
	return func(c *Config) {
		c.dockerImage = image
	}
}

func WithNetworkConfigAddress(addr string) ConfigOption {
	return func(c *Config) {
		c.networkConfigAddr = addr
	}
}

func (c *Config) Validate() error {
	if c.tokenName == "" {
		return fmt.Errorf("token name is required")
	}
	if c.tokenSymbol == "" {
		return fmt.Errorf("token symbol is required")
	}
	if c.tokenDecimals == "" {
		return fmt.Errorf("token decimals is required")
	}
	if c.l1HTTPURL == "" {
		return fmt.Errorf("L1 HTTP URL is required")
	}
	if c.privateKey == "" {
		return fmt.Errorf("private key is required")
	}
	if c.dockerImage == "" {
		return fmt.Errorf("docker image is required")
	}
	if c.networkConfigAddr == "" {
		return fmt.Errorf("network config address is required")
	}
	return nil
}

func (c *Config) String() string {
	return fmt.Sprintf("ERC20 Deployer Config: tokenName=%s, tokenSymbol=%s, tokenDecimals=%s, tokenSupply=%s, l1HTTPURL=%s, dockerImage=%s, networkConfigAddr=%s",
		c.tokenName, c.tokenSymbol, c.tokenDecimals, c.tokenSupply, c.l1HTTPURL, c.dockerImage, c.networkConfigAddr)
}
