package launcher

import (
	"github.com/ten-protocol/go-ten/go/config"
)

// Option is a function that applies configs to a Config Object
type Option = func(c *Config)

// Config holds the properties that configure the package
type Config struct {
	Nodes map[string]*config.NodeConfig
}
