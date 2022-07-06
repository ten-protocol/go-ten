package main

import (
	"github.com/obscuronet/obscuro-playground/go/host/hostrunner"
)

// Runs an Obscuro host as a standalone process.
func main() {
	config := hostrunner.ParseConfig()
	hostrunner.RunHost(config)
}
