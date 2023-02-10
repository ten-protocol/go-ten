package main

import (
	"flag"
	"strings"

	"github.com/obscuronet/go-obscuro/go/common/log"

	gethlog "github.com/ethereum/go-ethereum/log"
)

const (
	// Flag names and usages.
	numNodesName  = "numNodes"
	numNodesUsage = "The number of nodes on the network"

	gethHTTPStartPortName  = "gethHTTPStartPort"
	gethHTTPStartPortUsage = "The initial port to start allocating ports from"

	websocketStartPortName  = "gethWSStartPort"
	websocketStartPortUsage = "The initial port to start allocating websocket ports from"

	gethAuthRPCStartPortName  = "gethAuthRPCStartPort"
	gethAuthRPCStartPortUsage = "The initial port to start allocating geth auth rpc ports"

	gethNetworkStartPortName  = "gethNetworkStartPort"
	gethNetworkStartPortUsage = "The initial port to start allocating geths network/gossip ports"

	prysmBeaconRPCStartPortName  = "prysmBeaconRPCStartPort"
	prysmBeaconRPCStartPortUsage = "The initial port to start allocating prysm rpc port"

	prefundedAddrsName  = "prefundedAddrs"
	prefundedAddrsUsage = "The addresses to prefund as a comma-separated list"

	blockTimeSecsName  = "blockTimeSecs"
	blockTimeSecsUsage = "The block time in seconds"

	chainIDName  = "chainId"
	chainIDUsage = "The chain Id to use by the eth2 network"

	onlyDownloadName  = "onlyDownload"
	onlyDownloadUsage = "Only downloads the necessary files doesn't start the network"

	logLevelName  = "logLevel"
	logLevelUsage = "logLevel"

	logPathName  = "logPath"
	logPathUsage = "logPath"
)

type ethConfig struct {
	numNodes                int
	gethHTTPStartPort       int
	gethWSStartPort         int
	gethAuthRPCStartPort    int
	gethNetworkStartPort    int
	prysmBeaconRPCStartPort int
	prysmBeaconP2PStartPort int
	blockTimeSecs           int
	logLevel                int
	chainID                 int
	onlyDownload            bool
	logPath                 string
	prefundedAddrs          []string
}

func defaultConfig() *ethConfig {
	return &ethConfig{
		chainID:                 1337,
		numNodes:                1,
		gethHTTPStartPort:       12000,
		gethWSStartPort:         12100,
		gethAuthRPCStartPort:    12200,
		gethNetworkStartPort:    12300,
		prysmBeaconRPCStartPort: 12400,
		prysmBeaconP2PStartPort: 12500,
		onlyDownload:            false,
		prefundedAddrs:          []string{},
		blockTimeSecs:           1,
		logPath:                 log.SysOut,
		logLevel:                int(gethlog.LvlDebug),
	}
}

func parseCLIArgs() *ethConfig {
	defaultConfig := defaultConfig()

	onlyDownload := flag.Bool(onlyDownloadName, defaultConfig.onlyDownload, onlyDownloadUsage)
	numNodes := flag.Int(numNodesName, defaultConfig.numNodes, numNodesUsage)
	startPort := flag.Int(gethHTTPStartPortName, defaultConfig.gethHTTPStartPort, gethHTTPStartPortUsage)
	websocketStartPort := flag.Int(websocketStartPortName, defaultConfig.gethWSStartPort, websocketStartPortUsage)
	prefundedAddrs := flag.String(prefundedAddrsName, "", prefundedAddrsUsage)
	blockTimeSecs := flag.Int(blockTimeSecsName, defaultConfig.blockTimeSecs, blockTimeSecsUsage)
	logLevel := flag.Int(logLevelName, defaultConfig.logLevel, logLevelUsage)
	logPath := flag.String(logPathName, defaultConfig.logPath, logPathUsage)
	chainID := flag.Int(chainIDName, defaultConfig.chainID, chainIDUsage)

	gethAuthRPCStartPort := flag.Int(gethAuthRPCStartPortName, defaultConfig.gethAuthRPCStartPort, gethAuthRPCStartPortUsage)
	gethNetworkStartPort := flag.Int(gethNetworkStartPortName, defaultConfig.gethNetworkStartPort, gethNetworkStartPortUsage)
	prysmBeaconRPCStartPort := flag.Int(prysmBeaconRPCStartPortName, defaultConfig.prysmBeaconRPCStartPort, prysmBeaconRPCStartPortUsage)

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

	return &ethConfig{
		numNodes:                *numNodes,
		chainID:                 *chainID,
		gethHTTPStartPort:       *startPort,
		gethWSStartPort:         *websocketStartPort,
		prefundedAddrs:          parsedPrefundedAddrs,
		blockTimeSecs:           *blockTimeSecs,
		logLevel:                *logLevel,
		logPath:                 *logPath,
		gethAuthRPCStartPort:    *gethAuthRPCStartPort,
		gethNetworkStartPort:    *gethNetworkStartPort,
		prysmBeaconRPCStartPort: *prysmBeaconRPCStartPort,
		onlyDownload:            *onlyDownload,
	}
}
