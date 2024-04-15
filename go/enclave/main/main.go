package main

import (
	"fmt"
	"github.com/ten-protocol/go-ten/go/common/container"
	"github.com/ten-protocol/go-ten/go/config"
	enclavecontainer "github.com/ten-protocol/go-ten/go/enclave/container"
)

// Runs an TEN enclave as a standalone process.
func main() {
	rParams, _, err := config.LoadFlagStrings(config.Enclave)
	if err != nil {
		panic(fmt.Errorf("could not parse config. Cause: %w", err))
	}

	parsedConfig, err := enclavecontainer.ParseConfig(rParams)
	if err != nil {
		panic("error loading default configurations: %s" + err.Error())
	}

	enclaveContainer := enclavecontainer.NewEnclaveContainerFromConfig(parsedConfig)
	container.Serve(enclaveContainer)
}
