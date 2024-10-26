package main

import (
	"fmt"
	"os"

	"github.com/ten-protocol/go-ten/go/common/container"
	"github.com/ten-protocol/go-ten/go/config"
	enclaveconfig "github.com/ten-protocol/go-ten/go/enclave/config"
	enclavecontainer "github.com/ten-protocol/go-ten/go/enclave/container"
)

// Runs an Obscuro enclave as a standalone process.
func main() {
	tenCfg, err := config.LoadTenConfigForEnv("local")
	if err != nil {
		fmt.Println("Error loading ten config:", err)
		os.Exit(1)
	}

	enclaveConfig := enclaveconfig.EnclaveConfigFromTenConfig(tenCfg)

	enclaveContainer := enclavecontainer.NewEnclaveContainerFromConfig(enclaveConfig)
	container.Serve(enclaveContainer)
}
