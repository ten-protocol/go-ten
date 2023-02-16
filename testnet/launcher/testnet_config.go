package launcher

type Option = func(c *Config)

type Config struct {
	numberNodes int
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
