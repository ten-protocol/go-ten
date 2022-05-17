package main

import (
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host/hostrunner"
)

// Runs an Obscuro host as a standalone process.
func main() {
	config := hostrunner.ParseCLIArgs()
	hostrunner.RunHost(config)
}
