package eth2network

// Option is a function that applies configs to a Config Object
type Option = func(c *Config)

// Config holds the properties that configure the package
type Config struct {
	gethHTTPPort    int
	gethWSPort      int
	gethNetworkPort int
	gethRPCPort     int
	prysmP2PPort    int
	prysmRPCPort    int
	prefundedAddrs  []string
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

func WithGethNetworkStartPort(i int) Option {
	return func(c *Config) {
		c.gethNetworkPort = i
	}
}

func WithGethRPCStartPort(i int) Option {
	return func(c *Config) {
		c.gethRPCPort = i
	}
}

func WithPrysmP2PStartPort(i int) Option {
	return func(c *Config) {
		c.prysmP2PPort = i
	}
}

func WithPrysmRPCStartPort(i int) Option {
	return func(c *Config) {
		c.prysmRPCPort = i
	}
}

func WithGethPrefundedAddrs(addrs []string) Option {
	return func(c *Config) {
		c.prefundedAddrs = addrs
	}
}
