package main

import (
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/enclaverunner"
)

// Runs an Obscuro enclave as a standalone process.
func main() {
	config := enclaverunner.ParseCLIArgs()
	enclaverunner.RunEnclave(config)
}
