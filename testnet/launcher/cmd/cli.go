package main

import (
	"fmt"
	"github.com/ten-protocol/go-ten/go/config"
)

// ParseConfig returns a NodeConfig based the cli params and defaults.
func ParseConfig(paths config.RunParams) (*config.TestnetConfig, error) {
	inputCfg, err := config.LoadDefaultInputConfig(config.Node, paths)
	if err != nil {
		panic(fmt.Errorf("issues loading default and override config from file. Cause: %w", err))
	}
	return inputCfg.(*config.TestnetConfig), nil // assert
}
