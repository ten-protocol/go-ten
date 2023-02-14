package main

import (
	"fmt"
	"github.com/obscuronet/go-obscuro/testnet/launcher"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// TODO add config parsing
	testnet := launcher.NewTestnetLauncher()
	err := testnet.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Testnet start successfully")
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Press ctrl+c to stop...")
	<-done // Will block here until user hits ctrl+c
	// TODO add clean up / teardown
	os.Exit(0)

}
