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

	enclaveAddrName    = "enclaveAddress"
	enclaveAddrDefault = "localhost:11000"
	enclaveAddrUsage   = "The address to use to connect to the Obscuro enclave service"

	ourP2PAddrName    = "ourP2PAddr"
	ourP2PAddrDefault = "localhost:10000"
	ourP2PAddrUsage   = "The P2P address for our node"

	peerP2PAddrsName    = "peerP2PAddresses"
	peerP2PAddrsDefault = ""
	peerP2PAddrsUsage   = "The P2P addresses of our peer nodes as a comma-separated list (default \"\")"

	privateKeyName    = "privateKey"
	privateKeyDefault = "5dbbff1b5ff19f1ad6ea656433be35f6846e890b3f3ec6ef2b2e2137a8cab4ae" // ganache default key
	privateKeyUsage   = "The private key for the L1 node account"

	contractAddrName    = "contractAddress"
	contractAddrDefault = "0x63eA4F258791234606F9d11C0C66D83c99bfC7CA" // ganache contract deploy address
	contractAddrUsage   = "The rollup contract address"
)

type hostConfig struct {
	nodeID           *string
	isGenesis        *bool
	gossipRoundNanos *uint64
	rpcTimeoutSecs   *uint64
	enclaveAddr      *string
	ourP2PAddr       *string
	peerP2PAddrs     []string
	privateKeyString *string
	contractAddress  *string
}

func parseCLIArgs() hostConfig {
	nodeID := flag.String(nodeIDName, nodeIDDefault, nodeIDUsage)
	isGenesis := flag.Bool(genesisName, genesisDefault, genesisUsage)
	gossipRoundNanos := flag.Uint64(gossipRoundNanosName, uint64(gossipRoundNanosDefault), gossipRoundNanosUsage)
	rpcTimeoutSecs := flag.Uint64(rpcTimeoutSecsName, rpcTimeoutSecsDefault, rpcTimeoutSecsUsage)
	enclaveAddr := flag.String(enclaveAddrName, enclaveAddrDefault, enclaveAddrUsage)
	ourP2PAddr := flag.String(ourP2PAddrName, ourP2PAddrDefault, ourP2PAddrUsage)
	peerP2PAddrs := flag.String(peerP2PAddrsName, peerP2PAddrsDefault, peerP2PAddrsUsage)
	privateKeyStr := flag.String(privateKeyName, privateKeyDefault, privateKeyUsage)
	contractAddress := flag.String(contractAddrName, contractAddrDefault, contractAddrUsage)
	flag.Parse()

	return hostConfig{
		nodeID:           nodeID,
		isGenesis:        isGenesis,
		gossipRoundNanos: gossipRoundNanos,
		rpcTimeoutSecs:   rpcTimeoutSecs,
		enclaveAddr:      enclaveAddr,
		ourP2PAddr:       ourP2PAddr,
		peerP2PAddrs:     strings.Split(*peerP2PAddrs, ","),
		privateKeyString: privateKeyStr,
		contractAddress:  contractAddress,
	}
}
