package config

import (
	"github.com/ten-protocol/go-ten/go/common"
)

// NodeConfig contains the configuration for this Ten node (maybe relevant to multiple processes, host and enclave)
//
//	yaml: `node`
type NodeConfig struct {
	NodeType common.NodeType `mapstructure:"nodeType"`
	// Name of the node, used by orchestrator to name the containers etc., mostly useful for local testnets
	Name string `mapstructure:"name"`
	// Arbitrary identification of the node. Useful for debugging.
	ID string `mapstructure:"id"`
	// The public peer-to-peer IP address of the host
	// todo: does host address still need to exist for the enclave to sign over or does the enclave ID cover the usages?
	HostAddress string `mapstructure:"hostAddress"`
	// The stringified private key for the host's L1 wallet
	PrivateKeyString string `mapstructure:"privateKey"`
	// Whether the host is the genesis Obscuro node
	IsGenesis bool `mapstructure:"isGenesis"`
}
