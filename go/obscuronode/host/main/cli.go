package main

import (
	"flag"
	"strings"
)

const (
	// Flag names, defaults and usages.
	nodeIDName    = "nodeID"
	nodeIDDefault = ""
	nodeIDUsage   = "The 20 bytes of the node's address (default \"\")"

	genesisName    = "isGenesis"
	genesisDefault = true
	genesisUsage   = "Whether the node is the first node to join the network"

	gossipRoundNanosName    = "gossipRoundNanos"
	gossipRoundNanosDefault = 8333
	gossipRoundNanosUsage   = "The duration of the gossip round"

	rpcTimeoutSecsName    = "rpcTimeoutSecs"
	rpcTimeoutSecsDefault = 3
	rpcTimeoutSecsUsage   = "The timeout for host <-> enclave RPC communication"

	enclavePortName    = "enclavePort"
	enclavePortDefault = 11000
	enclavePortUsage   = "The port to use to connect to the Obscuro enclave service"

	ourP2PAddrName    = "ourP2PAddr"
	ourP2PAddrDefault = "localhost:10000"
	ourP2PAddrUsage   = "The P2P address for our node"

	peerP2PAddrsName    = "peerP2PAddresses"
	peerP2PAddrsDefault = ""
	peerP2PAddrsUsage   = "The P2P addresses of our peer nodes as a comma-separated list (default \"\")"
)

type hostConfig struct {
	nodeID           *string
	isGenesis        *bool
	gossipRoundNanos *uint64
	rpcTimeoutSecs   *uint64
	enclavePort      *uint64
	ourP2PAddr       *string
	peerP2PAddrs     []string
}

func parseCLIArgs() hostConfig {
	nodeID := flag.String(nodeIDName, nodeIDDefault, nodeIDUsage)
	isGenesis := flag.Bool(genesisName, genesisDefault, genesisUsage)
	gossipRoundNanos := flag.Uint64(gossipRoundNanosName, uint64(gossipRoundNanosDefault), gossipRoundNanosUsage)
	rpcTimeoutSecs := flag.Uint64(rpcTimeoutSecsName, rpcTimeoutSecsDefault, rpcTimeoutSecsUsage)
	enclavePort := flag.Uint64(enclavePortName, enclavePortDefault, enclavePortUsage)
	ourP2PAddr := flag.String(ourP2PAddrName, ourP2PAddrDefault, ourP2PAddrUsage)
	peerP2PAddrs := flag.String(peerP2PAddrsName, peerP2PAddrsDefault, peerP2PAddrsUsage)
	flag.Parse()

	return hostConfig{
		nodeID, isGenesis, gossipRoundNanos, rpcTimeoutSecs, enclavePort, ourP2PAddr, strings.Split(*peerP2PAddrs, ","),
	}
}
