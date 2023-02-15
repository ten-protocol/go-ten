package node

type Option = func(c *Config)

type Config struct {
	isGenesis                 bool
	sgxEnabled                bool
	enclaveImage              string
	hostImage                 string
	nodeType                  string
	l1Host                    string
	sequencerID               string
	privateKey                string
	hostP2PPort               int
	hostID                    string
	hostHTTPPort              int
	enclaveHTTPPort           int
	messageBusContractAddress string
	managementContractAddr    string
	enclaveWSPort             int
	l1WSPort                  int
	hostP2PAddr               string
	pccsAddr                  string
	edgelessDBImage           string
}

func NewNodeConfig(opts ...Option) *Config {
	defaultConfig := &Config{}

	for _, opt := range opts {
		opt(defaultConfig)
	}

	return defaultConfig
}

func WithNodeType(nodeType string) Option {
	return func(c *Config) {
		c.nodeType = nodeType
	}
}

func WithGenesis(b bool) Option {
	return func(c *Config) {
		c.isGenesis = b
	}
}

func WithSGXEnabled(b bool) Option {
	return func(c *Config) {
		c.sgxEnabled = b
	}
}

func WithEnclaveImage(s string) Option {
	return func(c *Config) {
		c.enclaveImage = s
	}
}

func WithHostImage(s string) Option {
	return func(c *Config) {
		c.hostImage = s
	}
}

func WithMessageBusContractAddress(s string) Option {
	return func(c *Config) {
		c.messageBusContractAddress = s
	}
}

func WithManagementContractAddress(s string) Option {
	return func(c *Config) {
		c.managementContractAddr = s
	}
}

func WithSequencerID(s string) Option {
	return func(c *Config) {
		c.sequencerID = s
	}
}

func WithHostID(s string) Option {
	return func(c *Config) {
		c.hostID = s
	}
}

func WithPrivateKey(s string) Option {
	return func(c *Config) {
		c.privateKey = s
	}
}

func WithEnclaveWSPort(i int) Option {
	return func(c *Config) {
		c.enclaveWSPort = i
	}
}

func WithEnclaveHTTPPort(i int) Option {
	return func(c *Config) {
		c.enclaveHTTPPort = i
	}
}

func WithL1WSPort(i int) Option {
	return func(c *Config) {
		c.l1WSPort = i
	}
}

func WithL1Host(s string) Option {
	return func(c *Config) {
		c.l1Host = s
	}
}

func WithHostP2PPort(i int) Option {
	return func(c *Config) {
		c.hostP2PPort = i
	}
}

func WithHostP2PAddr(s string) Option {
	return func(c *Config) {
		c.hostP2PAddr = s
	}
}

func WithEdgelessDBImage(s string) Option {
	return func(c *Config) {
		c.edgelessDBImage = s
	}
}

func WithPCCSAddr(s string) Option {
	return func(c *Config) {
		c.pccsAddr = s
	}
}
