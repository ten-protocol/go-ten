package main

import (
	"github.com/obscuronet/obscuro-playground/go/enclave/enclaverunner"
)

// Runs an Obscuro enclave as a standalone process.
func main() {
	config := enclaverunner.ParseConfig()
	// We set the logs outside of `RunEnclave` so we can override the logging in tests.
	enclaverunner.SetLogs(config.WriteToLogs, config.LogPath)
	enclaverunner.RunEnclave(config)
}
