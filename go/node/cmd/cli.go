package main

import (
	"flag"
	"fmt"
	"github.com/ten-protocol/go-ten/go/config"
	"github.com/ten-protocol/go-ten/go/node"
	"os"
	"strings"
)

var (
	startAction      = "start"
	upgradeAction    = "upgrade"
	validNodeActions = []string{startAction, upgradeAction}
)

// ParseConfig returns a node.DockerNode based on either the file identified by the `config` flag, or the flags with
// specific defaults (if the `config` flag isn't specified).
func ParseConfig() (*node.DockerNode, error) {
	inputCfg, err := config.LoadDefaultInputConfig(config.Node)
	if err != nil {
		panic(fmt.Errorf("issues loading default and override config from file. Cause: %w", err))
	}
	cfg := inputCfg.(*config.NodeConfig) // assert

	action := flag.Arg(0)

	if !validateNodeAction(action) {
		if action == "" {
			fmt.Printf("expected a node action string (%s) as the only argument after the flags but no argument provided\n",
				strings.Join(validNodeActions, ", "))
		} else {
			fmt.Printf("expected a node action string (%s) as the only argument after the flags but got %s\n",
				strings.Join(validNodeActions, ", "), action)
		}
		os.Exit(1)
	}

	return node.NewDockerNode(action, cfg), nil
}

func validateNodeAction(action string) bool {
	for _, a := range validNodeActions {
		if a == action {
			return true
		}
	}
	return false
}
