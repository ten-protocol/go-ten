package main

import (
	"fmt"

	"github.com/ten-protocol/go-ten/go/common/container"
	"github.com/ten-protocol/go-ten/go/config2"
	enclavecontainer "github.com/ten-protocol/go-ten/go/enclave/container"
)

// Runs an Obscuro enclave as a standalone process.
func main() {
	tenCfg, err := config2.Load([]string{})
	if err != nil {
		panic(fmt.Errorf("unable to load Ten config - %w", err))
	}

	enclaveConfig := enclaveconfig.EnclaveConfigFromTenConfig(tenCfg)

	enclaveContainer := enclavecontainer.NewEnclaveContainerFromConfig(enclaveConfig)
	container.Serve(enclaveContainer)
}
