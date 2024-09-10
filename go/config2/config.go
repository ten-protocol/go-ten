package config2

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// TenConfig is the top-level configuration struct for all TEN services.
// Fields are organised hierarchically, not all fields and structures will be used by all services.
//
// The default config is defined in yaml files formatted hierarchically like:
// ```
// network:
//
//	chainId: 1234
//	l1Contracts:
//	  managementContract: "0x1234567890abcdef1234567890abcdef12345678"
//	  ...
//
// node:
//
//	nodeType: "validator"
//	...
//
// ```
// The config can be overridden by other yaml files formatted the same.
// But fields can also be overridden directly per yaml spec, e.g.:
// ```
// network.chainId: 5678
// node.nodeType: "sequencer"
// ```
//
// Fields can also be overridden by environment variables, with the key being the flattened path of the field, like:
// ```
// export NETWORK_CHAINID=5678
// export NODE_NODETYPE=sequencer
// ```
// For ease of reading only the top-level struct is defined in this file, the nested structs are defined in their own files.
type TenConfig struct {
	Network *NetworkConfig `mapstructure:"network"`
	Node    *NodeConfig    `mapstructure:"node"`
	Host    *HostConfig    `mapstructure:"host"`
	Other   *HostConfig    `mapstructure:"other"`
}

func (t *TenConfig) PrettyPrint() {
	output, err := yaml.Marshal(t)
	if err != nil {
		fmt.Printf("Error printing config as YAML: %v\n", err)
		return
	}
	fmt.Println(string(output))
}
