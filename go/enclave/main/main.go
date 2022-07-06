package main

import (
	"github.com/obscuronet/obscuro-playground/go/enclave/enclaverunner"
)

// Runs an Obscuro enclave as a standalone process.
func main() {
	config := enclaverunner.ParseConfig()
	enclaverunner.RunEnclave(config)
}
