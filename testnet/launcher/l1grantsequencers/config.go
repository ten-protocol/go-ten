package l1grantsequencers

// Option is a function that applies configs to a Config Object
type Option = func(c *Config)

// Config holds the properties that configure the package
type Config struct {
	l1HTTPURL           string
	privateKey          string
	mgmtContractAddress string
	enclaveIDs          string
	dockerImage         string
	sequencerURL        string

	// debugEnabled        bool
}

func NewGrantSequencerConfig(opts ...Option) *Config {
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

func WithMgmtContractAddress(s string) Option {
	return func(c *Config) {
		c.mgmtContractAddress = s
	}
}

func WithEnclaveIDs(s string) Option {
	return func(c *Config) {
		c.enclaveIDs = s
	}
}

func WithDockerImage(s string) Option {
	return func(c *Config) {
		c.dockerImage = s
	}
}

func WithSequencerURL(s string) Option {
	return func(c *Config) {
		c.sequencerURL = s
	}
}

//func WithDebugEnabled(b bool) Option {
//	return func(c *Config) {
//		c.debugEnabled = b
//	}
//}
