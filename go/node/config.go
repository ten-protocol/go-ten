package node

// Option is a function that applies configs to a Config Object
type Option = func(c *Config)

// Config holds the properties that configure the package
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
	hostWSPort                int
	enclaveWSPort             int
	messageBusContractAddress string
	managementContractAddr    string
	l1WSPort                  int
	hostP2PHost               string
	hostPublicP2PAddr         string
	pccsAddr                  string
	edgelessDBImage           string
	enclaveDebug              bool
	nodeName                  string
}

func NewNodeConfig(opts ...Option) *Config {
	defaultConfig := &Config{}

	for _, opt := range opts {
		opt(defaultConfig)
	}

	return defaultConfig
}

func WithNodeName(s string) Option {
	return func(c *Config) {
		c.nodeName = s
	}
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

func WithEnclaveDebug(b bool) Option {
	return func(c *Config) {
		c.enclaveDebug = b
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

func WithHostP2PHost(s string) Option {
	return func(c *Config) {
		c.hostP2PHost = s
	}
}

func WithHostPublicP2PAddr(s string) Option {
	return func(c *Config) {
		c.hostPublicP2PAddr = s
	}
}

func WithHostHTTPPort(i int) Option {
	return func(c *Config) {
		c.hostHTTPPort = i
	}
}

func WithHostWSPort(i int) Option {
	return func(c *Config) {
		c.hostWSPort = i
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
