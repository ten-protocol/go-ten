package main

import (
	"github.com/ten-protocol/go-ten/go/common/container"
	"github.com/ten-protocol/go-ten/go/config2"
	hostcontainer "github.com/ten-protocol/go-ten/go/host/container"
)

// Runs an Obscuro host as a standalone process.
func main() {
	tenCfg := config2.Load()

	hostCfg := hostconfig.HostConfigFromTenConfig(tenCfg)
	hostContainer := hostcontainer.NewHostContainerFromConfig(hostCfg, nil)
	container.Serve(hostContainer)
}
