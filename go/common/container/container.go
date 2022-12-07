package container

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
)

type Status int

const (
	Running Status = iota
	Recovering
	Stopped
)

// Container is a Start-able server process that is expected to self-recover from any issues until Stop is called.
// Its exposes a Status method for monitoring. Both Enclave and Host implement Container.
//
// This abstraction can be started from a main() or controlled from another go process (allowing us to puppet the node components in simulations)
type Container interface {
	Start() error
	Stop() error
	Status() Status
}

// Serve is a convenience method to be called from the `main` runner for a container. It will attempt to cleanly shutdown
// the container on OS signal
// todo: expose the status to the operator from here (admin http service or a monitoring service)
func Serve(container Container) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-sigCh
		log.Printf("OS interrupt:%+v\n", oscall)
		cancel()
	}()

	err := container.Start()
	if err != nil {
		fmt.Printf("failed to start container - %s\n", err)
		os.Exit(1)
	}
	fmt.Println("Server started.")
	<-ctx.Done()
	fmt.Println("Stopping server...")
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("Failed to stop after 5 seconds. Exiting.")
		os.Exit(1)
	}()
	err = container.Stop()
	if err != nil {
		fmt.Printf("failed to stop gracefully - %s\n", err)
	}

	// Graceful shutdown complete
	os.Exit(0)
}
