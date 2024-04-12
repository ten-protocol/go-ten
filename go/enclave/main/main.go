package main

import (
	"fmt"
	"github.com/ten-protocol/go-ten/go/common/container"
	enclavecontainer "github.com/ten-protocol/go-ten/go/enclave/container"
)

// Runs an TEN enclave as a standalone process.
func main() {
	parsedConfig, err := enclavecontainer.ParseConfig()
	if err != nil {
		panic(fmt.Errorf("could not parse config. Cause: %w", err))
	}

	enclaveContainer := enclavecontainer.NewEnclaveContainerFromConfig(parsedConfig)
	container.Serve(enclaveContainer)
}
