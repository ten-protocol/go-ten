package main

import (
	"fmt"
	"github.com/ten-protocol/go-ten/go/config"
)

// ParseConfig returns a NodeConfig based the cli params and defaults.
func ParseConfig(confFiles map[string]string) (map[string]*config.NodeConfig, error) {
	var configs = make(map[string]*config.NodeConfig)
	for k, v := range confFiles {
		inputCfg, err := config.LoadConfigFromFile(config.Node, v)
		configs[k] = inputCfg.(*config.NodeConfig) // assert
		if err != nil {
			panic(fmt.Errorf("issues loading default and override config from file. Cause: %w", err))
		}
	}

	return configs, nil
}
