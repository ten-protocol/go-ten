package hostrunner

import (
	"flag"
	"strings"
)

const (
	// Flag names, defaults and usages.
	nodeIDName  = "nodeID"
	nodeIDUsage = "The 20 bytes of the node's address (default \"\")"

	genesisName  = "isGenesis"
	genesisUsage = "Whether the node is the first node to join the network"

	gossipRoundNanosName  = "gossipRoundNanos"
	gossipRoundNanosUsage = "The duration of the gossip round"

	rpcTimeoutSecsName  = "rpcTimeoutSecs"
	rpcTimeoutSecsUsage = "The timeout for host <-> enclave RPC communication"

	enclaveAddrName  = "enclaveAddress"
	enclaveAddrUsage = "The address to use to connect to the Obscuro enclave service"

	ourP2PAddrName    = "ourP2PAddr"
	ourP2PAddrDefault = "localhost:10000"
	ourP2PAddrUsage   = "The P2P address for our node"

	peerP2PAddrsName  = "peerP2PAddresses"
	peerP2PAddrsUsage = "The P2P addresses of our peer nodes as a comma-separated list (default \"\")"

	clientServerAddrName  = "clientServerAddress"
	clientServerAddrUsage = "The address on which to listen for client application RPC requests"

	privateKeyName  = "privateKey"
	privateKeyUsage = "The private key for the L1 node account"

	contractAddrName  = "contractAddress"
	contractAddrUsage = "The management contract address on the L1"
)

type HostConfig struct {
	NodeID           string
	IsGenesis        bool
	GossipRoundNanos uint64
	RPCTimeoutSecs   uint64
	EnclaveAddr      string
	OurP2PAddr       string
	PeerP2PAddrs     []string
	ClientServerAddr string
	PrivateKeyString string
	ContractAddress  string
}

func DefaultHostConfig() HostConfig {
	return HostConfig{
		NodeID:           "",
		IsGenesis:        true,
		GossipRoundNanos: 8333,
		RPCTimeoutSecs:   3,
		EnclaveAddr:      "localhost:11000",
		OurP2PAddr:       "",
		PeerP2PAddrs:     []string{},
		ClientServerAddr: "localhost:12000",
		PrivateKeyString: "0000000000000000000000000000000000000000000000000000000000000001",
		ContractAddress:  "",
	}
}

func ParseCLIArgs() HostConfig {
	defaultConfig := DefaultHostConfig()

	// TODO - Only provide defaults for certain flags. Some flags cannot be defaulted meaningfully (e.g. privateKeyString).
	nodeID := flag.String(nodeIDName, defaultConfig.NodeID, nodeIDUsage)
	isGenesis := flag.Bool(genesisName, defaultConfig.IsGenesis, genesisUsage)
	gossipRoundNanos := flag.Uint64(gossipRoundNanosName, defaultConfig.GossipRoundNanos, gossipRoundNanosUsage)
	rpcTimeoutSecs := flag.Uint64(rpcTimeoutSecsName, defaultConfig.RPCTimeoutSecs, rpcTimeoutSecsUsage)
	enclaveAddr := flag.String(enclaveAddrName, defaultConfig.EnclaveAddr, enclaveAddrUsage)
	ourP2PAddr := flag.String(ourP2PAddrName, defaultConfig.OurP2PAddr, ourP2PAddrUsage)
	peerP2PAddrs := flag.String(peerP2PAddrsName, "", peerP2PAddrsUsage)
	clientServerAddr := flag.String(clientServerAddrName, defaultConfig.ClientServerAddr, clientServerAddrUsage)
	privateKeyStr := flag.String(privateKeyName, defaultConfig.PrivateKeyString, privateKeyUsage)
	contractAddress := flag.String(contractAddrName, defaultConfig.ContractAddress, contractAddrUsage)
	flag.Parse()

	return HostConfig{
		NodeID:           *nodeID,
		IsGenesis:        *isGenesis,
		GossipRoundNanos: *gossipRoundNanos,
		RPCTimeoutSecs:   *rpcTimeoutSecs,
		EnclaveAddr:      *enclaveAddr,
		OurP2PAddr:       *ourP2PAddr,
		PeerP2PAddrs:     strings.Split(*peerP2PAddrs, ","),
		ClientServerAddr: *clientServerAddr,
		PrivateKeyString: *privateKeyStr,
		ContractAddress:  *contractAddress,
	}
}
