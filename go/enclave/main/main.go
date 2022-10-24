package main

import (
	"fmt"

	"github.com/obscuronet/go-obscuro/go/enclave/enclaverunner"
)

// Runs an Obscuro enclave as a standalone process.
func main() {
	config, err := enclaverunner.ParseConfig()
	if err != nil {
		panic(fmt.Errorf("could not parse config. Cause: %w", err))
	}
	enclaverunner.RunEnclave(config)
}
