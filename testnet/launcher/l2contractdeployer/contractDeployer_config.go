package l2contractdeployer

type Option = func(c *Config)

type Config struct {
	l1Host            string
	l1privateKey      string
	l1Port            int
	l2Port            int
	l2Host            string
	l2PrivateKey      string
	hocPKString       string
	pocPKString       string
	messageBusAddress string
}

func NewContractDeployerConfig(opts ...Option) *Config {
	defaultConfig := &Config{}

	for _, opt := range opts {
		opt(defaultConfig)
	}

	return defaultConfig
}

func WithL1Host(s string) Option {
	return func(c *Config) {
		c.l1Host = s
	}
}

func WithL1Port(i int) Option {
	return func(c *Config) {
		c.l1Port = i
	}
}

func WithL1PrivateKey(s string) Option {
	return func(c *Config) {
		c.l1privateKey = s
	}
}

func WithL2Port(i int) Option {
	return func(c *Config) {
		c.l2Port = i
	}
}

func WithL2Host(s string) Option {
	return func(c *Config) {
		c.l2Host = s
	}
}

func WithMessageBusContractAddress(s string) Option {
	return func(c *Config) {
		c.messageBusAddress = s
	}
}

func WithL2PrivateKey(s string) Option {
	return func(c *Config) {
		c.l2PrivateKey = s
	}
}

func WithHocPKString(s string) Option {
	return func(c *Config) {
		c.hocPKString = s
	}
}
func WithPocPKString(s string) Option {
	return func(c *Config) {
		c.pocPKString = s
	}
}
