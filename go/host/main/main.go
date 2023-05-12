package main

import (
	"fmt"
	"github.com/obscuronet/go-obscuro/go/common/container"
	hostcontainer "github.com/obscuronet/go-obscuro/go/host/container"
)

// Runs an Obscuro host as a standalone process.
func main() {
	parsedConfig, err := hostcontainer.ParseConfig()
	if err != nil {
		panic(fmt.Errorf("could not parse config. Cause: %w", err))
	}

	hostContainer := hostcontainer.NewHostContainerFromConfig(parsedConfig, nil)
	container.Serve(hostContainer)
}
