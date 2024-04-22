package l1contractdeployer

// Config holds the properties that configure the package
type L struct {
	l1HTTPURL    string
	privateKey   string
	dockerImage  string
	debugEnabled bool
}

func NewContractDeployerConfig(opts ...Option) *Config {
	defaultConfig := &Config{}

	for _, opt := range opts {
		opt(defaultConfig)
	}

	return defaultConfig
}

func WithL1HTTPURL(s string) Option {
	return func(c *Config) {
		c.l1HTTPURL = s
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

func WithDebugEnabled(b bool) Option {
	return func(c *Config) {
		c.debugEnabled = b
	}
}
