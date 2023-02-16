package l1contractdeployer

type Option = func(c *Config)

type Config struct {
	l1Host      string
	privateKey  string
	l1Port      int
	dockerImage string
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

func WithPrivateKey(s string) Option {
	return func(c *Config) {
		c.privateKey = s
	}
}

func WithDockerImage(s string) Option {
	return func(c *Config) {
		c.dockerImage = s
	}
}
