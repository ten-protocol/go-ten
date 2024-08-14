package main

import (
	"flag"
	"strings"

	gethlog "github.com/ethereum/go-ethereum/log"
)

const (
	// Flag names and usages.
	chainIDName  = "chainId"
	chainIDUsage = "The chain Id to use by the eth2 network"

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

	prysmBeaconGatewayStartPortName  = "prysmBeaconGatewayStartPort"
	prysmBeaconGatewayStartPortUsage = "The gateway port to connect to prysm"

	prysmBeaconP2PStartPortName  = "prysmBeaconP2PtartPort"
	prysmBeaconP2PStartPortUsage = "The p2p udp prysm port"

	prefundedAddrsName  = "prefundedAddrs"
	prefundedAddrsUsage = "The addresses to prefund as a comma-separated list"

	logLevelName  = "logLevel"
	logLevelUsage = "logLevel"
)

type ethConfig struct {
	chainID                     int
	gethHTTPStartPort           int
	gethWSStartPort             int
	gethAuthRPCStartPort        int
	gethNetworkStartPort        int
	prysmBeaconRPCStartPort     int
	prysmBeaconP2PStartPort     int
	prysmBeaconGatewayStartPort int
	logLevel                    int
	prefundedAddrs              []string
}

func defaultConfig() *ethConfig {
	return &ethConfig{
		chainID:                     1337,
		gethHTTPStartPort:           12000,
		gethWSStartPort:             12100,
		gethAuthRPCStartPort:        12200,
		gethNetworkStartPort:        12300,
		prysmBeaconRPCStartPort:     12400,
		prysmBeaconP2PStartPort:     12500,
		prysmBeaconGatewayStartPort: 12600,
		prefundedAddrs:              []string{},
		logLevel:                    int(gethlog.LvlDebug),
	}
}

func parseCLIArgs() *ethConfig {
	defaultConfig := defaultConfig()

	chainID := flag.Int(chainIDName, defaultConfig.chainID, chainIDUsage)
	gethHTTPPort := flag.Int(gethHTTPStartPortName, defaultConfig.gethHTTPStartPort, gethHTTPStartPortUsage)
	gethWSPort := flag.Int(websocketStartPortName, defaultConfig.gethWSStartPort, websocketStartPortUsage)
	gethAuthRPCStartPort := flag.Int(gethAuthRPCStartPortName, defaultConfig.gethAuthRPCStartPort, gethAuthRPCStartPortUsage)
	gethNetworkStartPort := flag.Int(gethNetworkStartPortName, defaultConfig.gethNetworkStartPort, gethNetworkStartPortUsage)
	prysmBeaconP2PStartPort := flag.Int(prysmBeaconP2PStartPortName, defaultConfig.prysmBeaconP2PStartPort, prysmBeaconP2PStartPortUsage)
	prysmBeaconRPCStartPort := flag.Int(prysmBeaconRPCStartPortName, defaultConfig.prysmBeaconRPCStartPort, prysmBeaconRPCStartPortUsage)
	prysmBeaconGatewayStartPort := flag.Int(prysmBeaconGatewayStartPortName, defaultConfig.prysmBeaconGatewayStartPort, prysmBeaconGatewayStartPortUsage)
	logLevel := flag.Int(logLevelName, defaultConfig.logLevel, logLevelUsage)
	prefundedAddrs := flag.String(prefundedAddrsName, "", prefundedAddrsUsage)

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
		chainID:                     *chainID,
		gethHTTPStartPort:           *gethHTTPPort,
		gethWSStartPort:             *gethWSPort,
		gethAuthRPCStartPort:        *gethAuthRPCStartPort,
		gethNetworkStartPort:        *gethNetworkStartPort,
		prysmBeaconRPCStartPort:     *prysmBeaconRPCStartPort,
		prysmBeaconP2PStartPort:     *prysmBeaconP2PStartPort,
		prysmBeaconGatewayStartPort: *prysmBeaconGatewayStartPort,
		logLevel:                    *logLevel,
		prefundedAddrs:              parsedPrefundedAddrs,
	}
}
