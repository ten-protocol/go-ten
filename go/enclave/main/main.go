package main

import (
	"fmt"

	"github.com/obscuronet/go-obscuro/go/common/container"
	"github.com/obscuronet/go-obscuro/go/enclave/enclavecontainer"
)

// Runs an Obscuro enclave as a standalone process.
func main() {
	config, err := enclavecontainer.ParseConfig()
	if err != nil {
		panic(fmt.Errorf("could not parse config. Cause: %w", err))
	}

	enclaveContainer := enclavecontainer.NewEnclaveContainer(config)
	container.Serve(enclaveContainer)
}
