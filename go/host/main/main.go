package main

import (
	"fmt"
	"os"

	"github.com/ten-protocol/go-ten/go/common/container"
	"github.com/ten-protocol/go-ten/go/config"
	hostconfig "github.com/ten-protocol/go-ten/go/host/config"
	hostcontainer "github.com/ten-protocol/go-ten/go/host/container"
)

// Runs an Obscuro host as a standalone process.
func main() {
	tenCfg, err := config.LoadTenConfig()
	if err != nil {
		fmt.Println("Error loading ten config:", err)
		os.Exit(1)
	}
	fmt.Println("Starting host with the following TenConfig:")
	tenCfg.PrettyPrint() // dump config to stdout

	hostCfg := hostconfig.HostConfigFromTenConfig(tenCfg)
	hostContainer := hostcontainer.NewHostContainerFromConfig(hostCfg, nil)
	container.Serve(hostContainer)
}
