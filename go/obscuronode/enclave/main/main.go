package main

import (
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/enclaverunner"
	"os"
)

// Runs an Obscuro enclave as a standalone process.
func main() {
	for _, arg := range os.Args {
		println(arg)
	}
	config := enclaverunner.ParseConfig()
	println("jjj ", config.ViewingKeysEnabled)
	// We set the logs outside of `RunEnclave` so we can override the logging in tests.
	enclaverunner.SetLogs(config.WriteToLogs, config.LogPath)
	enclaverunner.RunEnclave(config)
}
