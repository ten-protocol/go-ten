package gateway

import "time"

// Option is a function that applies configs to a Config Object
type Option = func(c *Config)

// Config holds the properties that configure the package
type Config struct {
	tenNodeHost              string
	tenNodeHTTPPort          int
	tenNodeWSPort            int
	gatewayHTTPPort          int
	gatewayWSPort            int
	rateLimitUserComputeTime time.Duration
	dockerImage              string
	chainID                  int64
}

func NewGatewayConfig(opts ...Option) *Config {
	defaultConfig := &Config{}

	for _, opt := range opts {
		opt(defaultConfig)
	}

	return defaultConfig
}

func WithTenNodeHost(s string) Option {
	return func(c *Config) {
		c.tenNodeHost = s
	}
}

func WithDockerImage(s string) Option {
	return func(c *Config) {
		c.dockerImage = s
	}
}

func WithTenNodeHTTPPort(i int) Option {
	return func(c *Config) {
		c.tenNodeHTTPPort = i
	}
}

func WithTenNodeWSPort(i int) Option {
	return func(c *Config) {
		c.tenNodeWSPort = i
	}
}

func WithGatewayHTTPPort(i int) Option {
	return func(c *Config) {
		c.gatewayHTTPPort = i
	}
}

func WithGatewayWSPort(i int) Option {
	return func(c *Config) {
		c.gatewayWSPort = i
	}
}

func WithRateLimitUserComputeTime(d time.Duration) Option {
	return func(c *Config) {
		c.rateLimitUserComputeTime = d
	}
}

func WithChainID(chainID int64) Option {
	return func(c *Config) {
		c.chainID = chainID
	}
}
