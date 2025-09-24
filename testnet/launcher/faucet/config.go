package faucet

// Option is a function that applies configs to a Config Object
type Option = func(c *Config)

// Config holds the properties that configure the package
type Config struct {
	tenNodeHost   string
	tenNodePort   int
	faucetPort    int
	faucetPrivKey string
	dockerImage   string
	chainID       int64
}

func NewFaucetConfig(opts ...Option) *Config {
	defaultConfig := &Config{}

	for _, opt := range opts {
		opt(defaultConfig)
	}

	return defaultConfig
}

func WithTenNodeHost(s string) Option {
	return func(c *Config) {
		c.tenNodeHost = s
	}
}

func WithFaucetPrivKey(s string) Option {
	return func(c *Config) {
		c.faucetPrivKey = s
	}
}

func WithDockerImage(s string) Option {
	return func(c *Config) {
		c.dockerImage = s
	}
}

func WithTenNodePort(i int) Option {
	return func(c *Config) {
		c.tenNodePort = i
	}
}

func WithFaucetPort(i int) Option {
	return func(c *Config) {
		c.faucetPort = i
	}
}

func WithChainID(chainID int64) Option {
	return func(c *Config) {
		c.chainID = chainID
	}
}
