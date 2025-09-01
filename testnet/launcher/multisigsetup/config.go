package multisigsetup

// Option is a function that applies configs to a Config Object
type Option = func(c *Config)

// Config holds the properties that configure the package
type Config struct {
	l1HTTPURL            string
	privateKey           string
	networkConfigAddress string
	multisigAddress      string
	proxyAdminAddress    string
	dockerImage          string
}

func NewMultisigSetupConfig(opts ...Option) *Config {
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

func WithNetworkConfigAddress(s string) Option {
	return func(c *Config) {
		c.networkConfigAddress = s
	}
}

func WithMultisigAddress(s string) Option {
	return func(c *Config) {
		c.multisigAddress = s
	}
}
func WithProxyAdminAddress(s string) Option {
	return func(c *Config) {
		c.proxyAdminAddress = s
	}
}

func WithDockerImage(s string) Option {
	return func(c *Config) {
		c.dockerImage = s
	}
}
