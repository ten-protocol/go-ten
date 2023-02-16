package launcher

type Option = func(c *Config)

type Config struct {
	numberNodes        int
	enclaveDockerImage string
	enclaveDebug       bool
}

func NewTestnetConfig(opts ...Option) *Config {
	defaultConfig := &Config{}

	for _, opt := range opts {
		opt(defaultConfig)
	}

	return defaultConfig
}

func WithNumberNodes(i int) Option {
	return func(c *Config) {
		c.numberNodes = i
	}
}

func WithEnclaveDockerImage(s string) Option {
	return func(c *Config) {
		c.enclaveDockerImage = s
	}
}

func WithEnclaveDebug(b bool) Option {
	return func(c *Config) {
		c.enclaveDebug = b
	}
}
