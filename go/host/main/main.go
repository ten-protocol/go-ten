package main

import (
	"github.com/ten-protocol/go-ten/go/common/container"
	"github.com/ten-protocol/go-ten/go/config"
	hostcontainer "github.com/ten-protocol/go-ten/go/host/container"
)

// Runs an TEN host as a standalone process.
func main() {
	var err error
	// load flags with defaults from config / sub-configs
	rParams, _, err := config.LoadFlagStrings(config.Host)
	if err != nil {
		panic(err)
	}

	parsedConfig, err := hostcontainer.ParseConfig(rParams)
	if err != nil {
		panic("error loading default configurations: %s" + err.Error())
	}

	hostContainer := hostcontainer.NewHostContainerFromConfig(parsedConfig, nil)
	container.Serve(hostContainer)
}
