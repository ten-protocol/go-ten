package main

import (
	"flag"
)

// TestnetConfigCLI represents the configurations passed into the testnet over CLI
type TestnetConfigCLI struct {
	numberNodes int
}

// ParseConfigCLI returns a NodeConfigCLI based the cli params and defaults.
func ParseConfigCLI() *TestnetConfigCLI {
	cfg := &TestnetConfigCLI{}
	flagUsageMap := getFlagUsageMap()

	numberNodes := flag.Int(numberNodesFlag, 2, flagUsageMap[numberNodesFlag])
	flag.Parse()

	cfg.numberNodes = *numberNodes

	return cfg
}
