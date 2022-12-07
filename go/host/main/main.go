package main

import (
	"fmt"

	"github.com/obscuronet/go-obscuro/go/common/container"
	"github.com/obscuronet/go-obscuro/go/host/hostcontainer"
)

// Runs an Obscuro host as a standalone process.
func main() {
	parsedConfig, err := hostcontainer.ParseConfig()
	if err != nil {
		panic(fmt.Errorf("could not parse config. Cause: %w", err))
	}

	hostContainer := hostcontainer.NewHostContainer(parsedConfig)
	container.Serve(hostContainer)
}
