package eth2network

type Option = func(c *Config)

type Config struct {
	gethHTTPPort   int
	gethWSPort     int
	prefundedAddrs []string
}

func NewEth2NetworkConfig(opts ...Option) *Config {
	defaultConfig := &Config{}

	for _, opt := range opts {
		opt(defaultConfig)
	}

	return defaultConfig
}

func WithGethHTTPStartPort(i int) Option {
	return func(c *Config) {
		c.gethHTTPPort = i
	}
}

func WithGethWSStartPort(i int) Option {
	return func(c *Config) {
		c.gethWSPort = i
	}
}

func WithGethPrefundedAddrs(addrs []string) Option {
	return func(c *Config) {
		c.prefundedAddrs = addrs
	}
}
