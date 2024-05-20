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

func deployNode(rParams config.RunParams) {
	var err error
	// load flags with defaults from config / sub-configs
	rParams, nodeFlags, err := config.LoadFlagStrings(config.Node)
	if err != nil {
		if strings.Contains(err.Error(), "help requested") {
			return
		}
		panic(err)
	}

	if !validateNodeAction(rParams[node.Action]) {
		if rParams[node.Action] == "" {
			fmt.Printf("expected a node action string (%s) as the only argument after the flags but no argument provided\n",
				strings.Join(validNodeActions, ", "))
		} else {
			fmt.Printf("expected a node action string (%s) as the only argument after the flags but got %s\n",
				strings.Join(validNodeActions, ", "), rParams[node.Action])
		}
		os.Exit(1)
	}

	// retrieve node-config
	nodeConfig, err := ParseConfig(rParams)
	if err != nil {
		panic(err)
	}

	dockerNode := node.NewDockerNode(rParams, nodeConfig, nodeFlags)

	//// NETWORK CONFIGS INCLUDING NETWORK + NODE LEVEL details
	switch dockerNode.Action {
	case startAction:
		// write the network-level config to disk for future restarts
		err = node.WriteNetworkConfigToDisk(dockerNode.Cfg)
		if err != nil {
			panic(err)
		}
		err = dockerNode.Start()
	case upgradeAction:
		// load network-specific details from the initial node setup from disk
		var ntwCfg *config.NetworkInputConfig
		ntwCfg, err = node.ReadNetworkConfigFromDisk()
		if err != nil {
			panic(err)
		}

		err = dockerNode.Upgrade(ntwCfg)
	default:
		panic("unrecognized node action: " + dockerNode.Action)
	}
	if err != nil {
		panic(err)
	}
}

// ParseConfig returns a node.NodeConfig based on either the file identified by the `config` flag, or the flags with
// specific defaults (if the `config` flag isn't specified).
func ParseConfig(paths config.RunParams) (*config.NodeConfig, error) {
	inputCfg, err := config.LoadDefaultInputConfig(config.Node, paths)
	if err != nil {
		return nil, fmt.Errorf("issues loading default and override config from file: %w", err)
	}
	cfg := inputCfg.(*config.NodeConfig) // assert

	fs := flag.NewFlagSet(config.Node.String(), flag.ExitOnError)
	usageMap := config.FlagUsageMap()
	config.SetupFlagsFromStruct(cfg, fs, usageMap)

	// Remove command-line flags in the case both flags and env vars are set
	os.Args, err = config.EnvOrFlag(os.Args)
	if err != nil {
		return nil, fmt.Errorf("error resolving property collision between flags and env vars: %w", err)
	}

	// Parse command-line flags
	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, fmt.Errorf("error parsing flags: %w", err)
	}

	setCommonProperties(cfg)

	return cfg, nil
}

func setCommonProperties(cfg *config.NodeConfig) {
	overrideFields := []interface{}{
		&cfg.NetworkConfig,
		&cfg.NodeDetails,
		&cfg.NodeSettings,
	}

	for _, field := range overrideFields {
		config.ApplyOverrides(&cfg.HostConfig, field)
		config.ApplyOverrides(&cfg.EnclaveConfig, field)
	}
}

func validateNodeAction(action string) bool {
	for _, a := range validNodeActions {
		if a == action {
			return true
		}
	}
	return false
}
