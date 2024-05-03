package main

import (
	"fmt"
	"github.com/ten-protocol/go-ten/go/config"
)

// ParseConfig returns an L1 Deployer connfig based the cli params and defaults.
func ParseConfig(paths config.RunParams) (*config.L1ContractDeployerConfig, error) {
	inputCfg, err := config.LoadDefaultInputConfig(config.L1Deployer, paths)
	if err != nil {
		panic(fmt.Errorf("issues loading default and override config from file. Cause: %w", err))
	}
	return inputCfg.(*config.L1ContractDeployerConfig), nil // assert
}
