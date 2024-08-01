package main

import (
	"fmt"
	"github.com/ten-protocol/go-ten/go/common/container"
	tenflag "github.com/ten-protocol/go-ten/go/common/flag"
	"github.com/ten-protocol/go-ten/go/config"
	enclavecontainer "github.com/ten-protocol/go-ten/go/enclave/container"
	"os"
	"runtime/pprof"
	"time"
)

// Runs an Obscuro enclave as a standalone process.
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

	// temporary code to help identify OOM
	go func() {
		for {
			heap, err := os.Open(fmt.Sprintf("heap_%s.pprof", time.Now().Format(time.RFC3339)))
			if err != nil {
				panic(fmt.Errorf("could not open CPU profile: %w", err))
			}
			err = pprof.WriteHeapProfile(heap)
			if err != nil {
				panic(fmt.Errorf("could not write CPU profile: %w", err))
			}
			stack, err := os.Open(fmt.Sprintf("stack_%s.pprof", time.Now().Format(time.RFC3339)))
			if err != nil {
				panic(fmt.Errorf("could not open CPU profile: %w", err))
			}
			err = pprof.Lookup("goroutine").WriteTo(stack, 1)
			if err != nil {
				panic(fmt.Errorf("could not write CPU profile: %w", err))
			}
			time.Sleep(1 * time.Hour)
		}
	}()

	enclaveContainer := enclavecontainer.NewEnclaveContainerFromConfig(enclaveConfig)
	container.Serve(enclaveContainer)
}
