package main

import (
	"fmt"

	"github.com/ten-protocol/go-ten/go/common/container"
	tenflag "github.com/ten-protocol/go-ten/go/common/flag"
	"github.com/ten-protocol/go-ten/go/config"
	enclavecontainer "github.com/ten-protocol/go-ten/go/enclave/container"
)

// Runs an Obscuro enclave as a standalone process.
func main() {
	// fetch and parse flags
	flags := config.EnclaveFlags()
	err := tenflag.CreateCLIFlags(flags)
	if err != nil {
		panic(fmt.Errorf("could not parse config. Cause: %w", err))
	}

	tenflag.Parse()

	enclaveConfig, err := config.NewConfigFromFlags(flags)
	if err != nil {
		panic(fmt.Errorf("unable to create config from flags - %w", err))
	}

	enclaveContainer := enclavecontainer.NewEnclaveContainerFromConfig(enclaveConfig)
	container.Serve(enclaveContainer)
}
