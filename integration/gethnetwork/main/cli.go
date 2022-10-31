package main

import (
	"flag"
	"strings"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/obscuronet/go-obscuro/go/common/log"
)

const (
	// Flag names and usages.
	numNodesName  = "numNodes"
	numNodesUsage = "The number of nodes on the network"

	startPortName  = "startPort"
	startPortUsage = "The initial port to start allocating ports from"

	websocketStartPortName  = "websocketStartPort"
	websocketStartPortUsage = "The initial port to start allocating websocket ports from"

	prefundedAddrsName  = "prefundedAddrs"
	prefundedAddrsUsage = "The addresses to prefund as a comma-separated list"

	blockTimeSecsName  = "blockTimeSecs"
	blockTimeSecsUsage = "The block time in seconds"

	logLevelName  = "logLevel"
	logLevelUsage = "logLevel"

	logPathName  = "logPath"
	logPathUsage = "logPath"
)

type gethConfig struct {
	numNodes           int
	startPort          int
	websocketStartPort int
	prefundedAddrs     []string
	blockTimeSecs      int
	logLevel           int
	logPath            string
}

func defaultHostConfig() gethConfig {
	return gethConfig{
		numNodes:           1,
		startPort:          12000,
		websocketStartPort: 12100,
		prefundedAddrs:     []string{},
		blockTimeSecs:      1,
		logPath:            log.SysOut,
		logLevel:           int(gethlog.LvlDebug),
	}
}

func parseCLIArgs() gethConfig {
	defaultConfig := defaultHostConfig()

	numNodes := flag.Int(numNodesName, defaultConfig.numNodes, numNodesUsage)
	startPort := flag.Int(startPortName, defaultConfig.startPort, startPortUsage)
	websocketStartPort := flag.Int(websocketStartPortName, defaultConfig.websocketStartPort, websocketStartPortUsage)
	prefundedAddrs := flag.String(prefundedAddrsName, "", prefundedAddrsUsage)
	blockTimeSecs := flag.Int(blockTimeSecsName, defaultConfig.blockTimeSecs, blockTimeSecsUsage)
	logLevel := flag.Int(logLevelName, defaultConfig.logLevel, logLevelUsage)
	logPath := flag.String(logPathName, defaultConfig.logPath, logPathUsage)

	flag.Parse()

	addrs := *prefundedAddrs
	// When running locally, we don't have to add the quotes around the prefunded addresses.
	// This is stripping them away in case they were added.
	if strings.HasPrefix(addrs, "'") {
		addrs = addrs[1 : len(addrs)-1]
	}
	parsedPrefundedAddrs := strings.Split(addrs, ",")
	if addrs == "" {
		// We handle the special case of an empty list.
		parsedPrefundedAddrs = []string{}
	}

	return gethConfig{
		numNodes:           *numNodes,
		startPort:          *startPort,
		websocketStartPort: *websocketStartPort,
		prefundedAddrs:     parsedPrefundedAddrs,
		blockTimeSecs:      *blockTimeSecs,
		logLevel:           *logLevel,
		logPath:            *logPath,
	}
}
