package launcher

type Eth2NetworkOption = func(c *Eth2NetworkConfig)

type Eth2NetworkConfig struct {
	gethHTTPPort   int
	gethWSPort     int
	prefundedAddrs []string
}

func NewEth2NetworkConfig(opts ...Eth2NetworkOption) *Eth2NetworkConfig {
	defaultConfig := &Eth2NetworkConfig{}

	for _, opt := range opts {
		opt(defaultConfig)
	}

	return defaultConfig
}

func WithGethHTTPStartPort(i int) Eth2NetworkOption {
	return func(c *Eth2NetworkConfig) {
		c.gethHTTPPort = i
	}
}

func WithGethWSStartPort(i int) Eth2NetworkOption {
	return func(c *Eth2NetworkConfig) {
		c.gethWSPort = i
	}
}

func WithGethPrefundedAddrs(addrs []string) Eth2NetworkOption {
	return func(c *Eth2NetworkConfig) {
		c.prefundedAddrs = addrs
	}
}
