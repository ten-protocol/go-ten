package main

import (
	"flag"

	"github.com/obscuronet/obscuro-playground/tools/walletextension"
)

const (
	// Flag names, defaults and usages.
	startPortName    = "startPort"
	startPortDefault = 3000
	startPortUsage   = "The first port to allocate. Ports will be allocated incrementally from this port as needed"
)

func parseCLIArgs() walletextension.RunConfig {
	startPort := flag.Int(startPortName, startPortDefault, startPortUsage)
	flag.Parse()

	return walletextension.RunConfig{
		StartPort: *startPort,
	}
}
