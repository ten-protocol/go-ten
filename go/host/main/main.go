package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"github.com/ten-protocol/go-ten/go/common/container"
	hostcontainer "github.com/ten-protocol/go-ten/go/host/container"
)

// Runs an Obscuro host as a standalone process.
func main() {
	parsedConfig, err := hostcontainer.ParseConfig()
	if err != nil {
		panic(fmt.Errorf("could not parse config. Cause: %w", err))
	}

	err = os.Mkdir("/data", os.FileMode(0777))
	if err != nil {
		panic(fmt.Errorf("could not create /data directory: %w", err))
	}

	// temporary code to help identify OOM
	go func() {
		for {
			heap, err := os.Create(fmt.Sprintf("/data/heap_%d.pprof", time.Now().UnixMilli()))
			if err != nil {
				panic(fmt.Errorf("could not open heap profile: %w", err))
			}
			err = pprof.WriteHeapProfile(heap)
			if err != nil {
				panic(fmt.Errorf("could not write CPU profile: %w", err))
			}
			stack, err := os.Create(fmt.Sprintf("/data/stack_%d.pprof", time.Now().UnixMilli()))
			if err != nil {
				panic(fmt.Errorf("could not open stack profile: %w", err))
			}
			err = pprof.Lookup("goroutine").WriteTo(stack, 1)
			if err != nil {
				panic(fmt.Errorf("could not write CPU profile: %w", err))
			}
			time.Sleep(1 * time.Hour)
		}
	}()

	hostContainer := hostcontainer.NewHostContainerFromConfig(parsedConfig, nil)
	container.Serve(hostContainer)
}
