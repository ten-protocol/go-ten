package main

import (
	"fmt"
	"github.com/ten-protocol/go-ten/go/config"

	"github.com/ten-protocol/go-ten/go/common/container"
	hostcontainer "github.com/ten-protocol/go-ten/go/host/container"
)

// Runs an Obscuro host as a standalone process.
func main() {
	parsedConfig, err := config.ParseHostConfig()
	if err != nil {
		panic(fmt.Errorf("could not parse config. Cause: %w", err))
	}

	hostContainer := hostcontainer.NewHostContainerFromConfig(parsedConfig, nil)
	container.Serve(hostContainer)
}
