package bridgetokenwhitelist

import "fmt"

type Config struct {
	tokenAddress      string
	tokenName         string
	tokenSymbol       string
	l1HTTPURL         string
	l2Host            string
	l2HTTPPort        int
	l2WSPort          int
	privateKey        string
	dockerImage       string
	networkConfigAddr string
}

func NewConfig(opts ...ConfigOption) *Config {
	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

type ConfigOption func(*Config)

func WithTokenAddress(address string) ConfigOption {
	return func(c *Config) {
		c.tokenAddress = address
	}
}

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

func WithL1HTTPURL(url string) ConfigOption {
	return func(c *Config) {
		c.l1HTTPURL = url
	}
}

func WithL2Host(host string) ConfigOption {
	return func(c *Config) {
		c.l2Host = host
	}
}

func WithL2HTTPPort(port int) ConfigOption {
	return func(c *Config) {
		c.l2HTTPPort = port
	}
}

func WithL2WSPort(port int) ConfigOption {
	return func(c *Config) {
		c.l2WSPort = port
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
	if c.tokenAddress == "" {
		return fmt.Errorf("token address is required")
	}
	if c.tokenName == "" {
		return fmt.Errorf("token name is required")
	}
	if c.tokenSymbol == "" {
		return fmt.Errorf("token symbol is required")
	}
	if c.l1HTTPURL == "" {
		return fmt.Errorf("L1 HTTP URL is required")
	}
	if c.l2Host == "" {
		return fmt.Errorf("L2 host is required")
	}
	if c.l2HTTPPort == 0 {
		return fmt.Errorf("L2 HTTP port is required")
	}
	if c.l2WSPort == 0 {
		return fmt.Errorf("L2 WS port is required")
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
	return fmt.Sprintf("Bridge Token Whitelist Config: tokenAddress=%s, tokenName=%s, tokenSymbol=%s, l1HTTPURL=%s, l2Host=%s, l2HTTPPort=%d, l2WSPort=%d, dockerImage=%s, networkConfigAddr=%s",
		c.tokenAddress, c.tokenName, c.tokenSymbol, c.l1HTTPURL, c.l2Host, c.l2HTTPPort, c.l2WSPort, c.dockerImage, c.networkConfigAddr)
}
