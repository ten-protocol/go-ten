package bridgetokenwhitelist

import "fmt"

type Config struct {
	tokenAddress      string
	tokenName         string
	tokenSymbol       string
	l1HTTPURL         string
	l2RPCURL          string
	l2Nonce           string
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

func WithL2RPCURL(url string) ConfigOption {
	return func(c *Config) {
		c.l2RPCURL = url
	}
}

func WithL2Nonce(nonce string) ConfigOption {
	return func(c *Config) {
		c.l2Nonce = nonce
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
	if c.l2RPCURL == "" {
		return fmt.Errorf("L2 RPC URL is required")
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
	return fmt.Sprintf("Bridge Token Whitelist Config: tokenAddress=%s, tokenName=%s, tokenSymbol=%s, l1HTTPURL=%s, l2RPCURL=%s, dockerImage=%s, networkConfigAddr=%s",
		c.tokenAddress, c.tokenName, c.tokenSymbol, c.l1HTTPURL, c.l2RPCURL, c.dockerImage, c.networkConfigAddr)
}
