package main

import (
	"fmt"
	"github.com/ten-protocol/go-ten/go/config"
)

var (
	action           = "action"
	startAction      = "start"
	upgradeAction    = "upgrade"
	validNodeActions = []string{startAction, upgradeAction}
)

// ParseConfig returns a node.DockerNode based on either the file identified by the `config` flag, or the flags with
// specific defaults (if the `config` flag isn't specified).
func ParseConfig(paths config.RunParams) (*config.NodeConfig, error) {
	inputCfg, err := config.LoadDefaultInputConfig(config.Node, paths)
	if err != nil {
		panic(fmt.Errorf("issues loading default and override config from file. Cause: %w", err))
	}
	return inputCfg.(*config.NodeConfig), nil // assert
}

func validateNodeAction(action string) bool {
	for _, a := range validNodeActions {
		if a == action {
			return true
		}
	}
	return false
}
