package pausecontracts

// Option is a function that applies configs to a Config Object
type Option = func(c *Config)

// Config holds the properties that configure the package
type Config struct {
	l1HTTPURL            string
	privateKey           string
	networkConfigAddr    string
	merkleMessageBusAddr string
	dockerImage          string
	action               string
}

func NewPauseAllContractsConfig(opts ...Option) *Config {
	defaultConfig := &Config{
		action: "PAUSE", // Default to PAUSE
	}

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

func WithNetworkConfigAddress(s string) Option {
	return func(c *Config) {
		c.networkConfigAddr = s
	}
}

func WithMerkleTreeMessageBus(s string) Option {
	return func(c *Config) {
		c.merkleMessageBusAddr = s
	}
}

func WithDockerImage(s string) Option {
	return func(c *Config) {
		c.dockerImage = s
	}
}

func WithAction(s string) Option {
	return func(c *Config) {
		c.action = s
	}
}
