package launcher

// Option is a function that applies configs to a Config Object
type Option = func(c *Config)

// Config holds the properties that configure the package
type Config struct {
	sequencerEnclaveDockerImage string
	sequencerEnclaveDebug       bool

	validatorEnclaveDockerImage string
	validatorEnclaveDebug       bool

	contractDeployerDockerImage string
	contractDeployerDebug       bool

	isSGXEnabled bool

	logLevel int
}

func NewTestnetConfig(opts ...Option) *Config {
	defaultConfig := &Config{}

	for _, opt := range opts {
		opt(defaultConfig)
	}

	return defaultConfig
}

func WithSequencerEnclaveDockerImage(s string) Option {
	return func(c *Config) {
		c.sequencerEnclaveDockerImage = s
	}
}

func WithSequencerEnclaveDebug(b bool) Option {
	return func(c *Config) {
		c.sequencerEnclaveDebug = b
	}
}

func WithValidatorEnclaveDockerImage(s string) Option {
	return func(c *Config) {
		c.validatorEnclaveDockerImage = s
	}
}

func WithValidatorEnclaveDebug(b bool) Option {
	return func(c *Config) {
		c.validatorEnclaveDebug = b
	}
}

func WithSGXEnabled(b bool) Option {
	return func(c *Config) {
		c.isSGXEnabled = b
	}
}

func WithContractDeployerDockerImage(s string) Option {
	return func(c *Config) {
		c.contractDeployerDockerImage = s
	}
}

func WithContractDeployerDebug(b bool) Option {
	return func(c *Config) {
		c.contractDeployerDebug = b
	}
}

func WithLogLevel(i int) Option {
	return func(c *Config) {
		c.logLevel = i
	}
}
