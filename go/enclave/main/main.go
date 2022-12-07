package main

import (
	"fmt"
)

// Runs an Obscuro enclave as a standalone process.
func main() {
	config, err := enclavecontainer.ParseConfig()
	if err != nil {
		panic(fmt.Errorf("could not parse config. Cause: %w", err))
	}
	enclavecontainer.RunEnclave(config)
}
