package main

import (
	"flag"
	"strings"
)

const (
	// Flag names, defaults and usages.
	numNodesName  = "numNodes"
	numNodesUsage = "The number of nodes on the network"

	startPortName  = "startPort"
	startPortUsage = "The initial port to start allocating ports from"

	prefundedAddrsName  = "prefundedAddrs"
	prefundedAddrsUsage = "The addresses to prefund as a comma-separated list"
)

type GethConfig struct {
	NumNodes       int
	StartPort      int
	PrefundedAddrs []string
}

func DefaultHostConfig() GethConfig {
	return GethConfig{
		NumNodes:       1,
		StartPort:      12000,
		PrefundedAddrs: []string{},
	}
}

func parseCLIArgs() GethConfig {
	defaultConfig := DefaultHostConfig()

	numNodes := flag.Int(numNodesName, defaultConfig.NumNodes, numNodesUsage)
	startPort := flag.Int(startPortName, defaultConfig.StartPort, startPortUsage)
	prefundedAddrs := flag.String(prefundedAddrsName, "", prefundedAddrsUsage)

	flag.Parse()

	parsedPrefundedAddrs := strings.Split(*prefundedAddrs, ",")
	if *prefundedAddrs == "" {
		// We handle the special case of an empty list.
		parsedPrefundedAddrs = []string{}
	}

	return GethConfig{
		NumNodes:       *numNodes,
		StartPort:      *startPort,
		PrefundedAddrs: parsedPrefundedAddrs,
	}
}
