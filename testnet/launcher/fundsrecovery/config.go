package fundsrecovery

// Option is a function that applies configs to a Config Object
type Option = func(c *Config)

// Config holds the properties that configure the package
type Config struct {
	l1HTTPURL           string
	l1privateKey        string
	mgmtContractAddress string
	dockerImage         string
	accToPay            string
}

func NewFundsRecoveryConfig(opts ...Option) *Config {
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

func WithL1PrivateKey(s string) Option {
	return func(c *Config) {
		c.l1privateKey = s
	}
}

func WithMgmtContractAddress(s string) Option {
	return func(c *Config) {
		c.mgmtContractAddress = s
	}
}

func WithDockerImage(s string) Option {
	return func(c *Config) {
		c.dockerImage = s
	}
}

func WithAccToPay(s string) Option {
	return func(c *Config) {
		c.accToPay = s
	}
}
