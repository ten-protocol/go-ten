package container

import (
	"flag"
	"fmt"
	"github.com/ten-protocol/go-ten/go/config"
	"os"
)

// ParseConfig returns a config.HostConfig based on either the file identified by the `config` flag, or the flags with
// specific defaults (if the `config` flag isn't specified).
func ParseConfig(paths config.ConfPaths) (*config.HostConfig, error) {
	inputCfg, err := config.LoadDefaultInputConfig(config.Host, paths)
	if err != nil {
		return nil, fmt.Errorf("issues loading default and override config from file: %w", err)
	}
	cfg := inputCfg.(*config.HostInputConfig) // assert

	fs := flag.NewFlagSet(config.Host.String(), flag.ExitOnError)
	usageMap := config.FlagUsageMap()
	config.SetupFlagsFromStruct(cfg, fs, usageMap)

	// Parse command-line flags
	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, fmt.Errorf("error parsing flags: %w", err)
	}

	hostConfig, err := cfg.ToHostConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to convert HostInputConfig to HostConfig: %w", err)
	}
	return hostConfig, nil
}
