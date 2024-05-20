package node

import (
	"github.com/ten-protocol/go-ten/go/config"
)

type Node interface {
	Start() error
	Stop() error
	Upgrade(networkCfg *config.NetworkConfig) error
}
