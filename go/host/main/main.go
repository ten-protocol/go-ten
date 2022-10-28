package main

import (
	"fmt"

	"github.com/obscuronet/go-obscuro/go/host/hostrunner"
)

// Runs an Obscuro host as a standalone process.
func main() {
	config, err := hostrunner.ParseConfig()
	if err != nil {
		panic(fmt.Errorf("could not parse config. Cause: %w", err))
	}
	hostrunner.RunHost(config)
}
