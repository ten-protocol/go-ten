package main

import (
	"fmt"

	"github.com/ten-protocol/go-ten/go/common/container"
	tenflag "github.com/ten-protocol/go-ten/go/common/flag"
	"github.com/ten-protocol/go-ten/go/config"
	enclavecontainer "github.com/ten-protocol/go-ten/go/enclave/container"
)

// Runs an TEN enclave as a standalone process.
func main() {
	// fetch and parse flags
	flags := config.EnclaveFlags                       // fetch the flags that enclave requires
	err := tenflag.CreateCLIFlags(config.EnclaveFlags) // using tenflag convert those flags into the golang flags package ( go flags is a singlen )
	if err != nil {
		panic(fmt.Errorf("could not create CLI flags. Cause: %w", err))
	}

	tenflag.Parse() // parse the golang flags package defined flags from CLI

	enclaveConfig, err := config.NewConfigFromFlags(flags)
	if err != nil {
		panic(fmt.Errorf("unable to create config from flags - %w", err))
	}

	enclaveContainer := enclavecontainer.NewEnclaveContainerFromConfig(enclaveConfig)
	container.Serve(enclaveContainer)
}
