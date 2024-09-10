package main

import (
	"fmt"

	"github.com/ten-protocol/go-ten/go/common/container"
	"github.com/ten-protocol/go-ten/go/config2"
	hostcontainer "github.com/ten-protocol/go-ten/go/host/container"
)

// Runs an Obscuro host as a standalone process.
func main() {
	cfg := config2.Load()
	hostCfg, err := hostcontainer.ReadHostConfig(cfg)
	if err != nil {
		panic(fmt.Errorf("could not parse config. Cause: %w", err))
	}

	hostContainer := hostcontainer.NewHostContainerFromConfig(hostCfg, nil)
	container.Serve(hostContainer)
}
