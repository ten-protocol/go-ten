package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ten-protocol/go-ten/tools/faucet/container"
)

// local execution: PORT=80 go run . --nodeHost 127.0.0.1 --pk 0x8dfb8083da6275ae3e4f41e3e8a8c19d028d32c9247e24530933782f2a05035b --jwtSecret This_is_the_secret
func main() {
	cfg := parseCLIArgs()

	if cfg.PK == "" {
		panic("no key loaded")
	}
	if cfg.JWTSecret == "" {
		panic("no jwt secret loaded")
	}

	faucetContainer, err := container.NewFaucetContainerFromConfig(cfg)
	if err != nil {
		panic(err)
	}

	err = faucetContainer.Start()
	if err != nil {
		panic(err)
	}
	fmt.Printf("💡 Faucet started on http://%s:%d \n", cfg.Host, cfg.HTTPPort)
	
	// Create a channel to receive signals
	signalCh := make(chan os.Signal, 1)

	// Notify the channel for interrupt signals
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	// Wait for an interrupt signal
	<-signalCh

	fmt.Println("Shutting down")

	err = faucetContainer.Stop()
	if err != nil {
		panic(err)
	}
}
