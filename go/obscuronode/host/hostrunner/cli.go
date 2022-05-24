package hostrunner

import (
	"flag"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/host"
)

const (
	// Flag names, defaults and usages.
	nodeIDName  = "nodeID"
	nodeIDUsage = "The 20 bytes of the node's address"

	isGenesisName  = "isGenesis"
	isGenesisUsage = "Whether the node is the first node to join the network"

	gossipRoundNanosName  = "gossipRoundNanos"
	gossipRoundNanosUsage = "The duration of the gossip round"

	enclaveRPCTimeoutSecsName  = "enclaveRPCTimeoutSecs"
	enclaveRPCTimeoutSecsUsage = "The timeout for host <-> enclave RPC communication"

	enclaveAddrName  = "enclaveAddress"
	enclaveAddrUsage = "The address to use to connect to the Obscuro enclave service"

	ourP2PAddrName  = "ourP2PAddr"
	ourP2PAddrUsage = "The P2P address for our node"

	peerP2PAddrsName  = "peerP2PAddresses"
	peerP2PAddrsUsage = "The P2P addresses of our peer nodes as a comma-separated list"

	clientServerAddrName  = "clientServerAddress"
	clientServerAddrUsage = "The address on which to listen for client application RPC requests"

	clientRPCTimeoutSecsName  = "clientRPCTimeoutSecs"
	clientRPCTimeoutSecsUsage = "The timeout for client <-> host RPC communication"

	privateKeyName  = "privateKey"
	privateKeyUsage = "The private key for the L1 node account"

	contractAddrName  = "contractAddress"
	contractAddrUsage = "The management contract address on the L1"

	ethClientHostName  = "ethClientHost"
	ethClientHostUsage = "The host on which to connect to the Ethereum client"

	ethClientPortName  = "ethClientPort"
	ethClientPortUsage = "The port on which to connect to the Ethereum client"

	logPathName  = "logPath"
	logPathUsage = "The path to use for the host's log file"

	chainIDName  = "chainID"
	chainIDUsage = "The ID of the L1 chain"
)

type DefaultHostConfig struct {
	NodeID                string
	IsGenesis             bool
	GossipRoundNanos      uint64
	EnclaveRPCTimeoutSecs uint64
	ClientRPCTimeoutSecs  uint64
	EnclaveAddr           string
	OurP2PAddr            string
	PeerP2PAddrs          []string
	ClientServerAddr      string
	PrivateKeyString      string
	ContractAddress       string
	EthClientHost         string
	EthClientPort         uint64
	LogPath               string
	ChainID               int64
}

func GetDefaults() DefaultHostConfig {
	return DefaultHostConfig{
		NodeID:                "",
		IsGenesis:             true,
		GossipRoundNanos:      8333,
		EnclaveRPCTimeoutSecs: 3,
		ClientRPCTimeoutSecs:  3,
		EnclaveAddr:           "127.0.0.1:11000",
		OurP2PAddr:            "",
		PeerP2PAddrs:          []string{},
		ClientServerAddr:      "127.0.0.1:13000",
		PrivateKeyString:      "0000000000000000000000000000000000000000000000000000000000000001",
		ContractAddress:       "",
		EthClientHost:         "127.0.0.1",
		EthClientPort:         8546,
		LogPath:               "host_logs.txt",
		ChainID:               1337,
	}
}

func ParseCLIArgs() host.Config {
	defaultConfig := GetDefaults()

	nodeID := flag.String(nodeIDName, defaultConfig.NodeID, nodeIDUsage)
	isGenesis := flag.Bool(isGenesisName, defaultConfig.IsGenesis, isGenesisUsage)
	gossipRoundNanos := flag.Uint64(gossipRoundNanosName, defaultConfig.GossipRoundNanos, gossipRoundNanosUsage)
	enclaveRPCTimeoutSecs := flag.Uint64(enclaveRPCTimeoutSecsName, defaultConfig.EnclaveRPCTimeoutSecs, enclaveRPCTimeoutSecsUsage)
	clientRPCTimeoutSecs := flag.Uint64(clientRPCTimeoutSecsName, defaultConfig.ClientRPCTimeoutSecs, clientRPCTimeoutSecsUsage)
	enclaveAddr := flag.String(enclaveAddrName, defaultConfig.EnclaveAddr, enclaveAddrUsage)
	ourP2PAddr := flag.String(ourP2PAddrName, defaultConfig.OurP2PAddr, ourP2PAddrUsage)
	peerP2PAddrs := flag.String(peerP2PAddrsName, "", peerP2PAddrsUsage)
	clientServerAddr := flag.String(clientServerAddrName, defaultConfig.ClientServerAddr, clientServerAddrUsage)
	privateKeyStr := flag.String(privateKeyName, defaultConfig.PrivateKeyString, privateKeyUsage)
	contractAddress := flag.String(contractAddrName, defaultConfig.ContractAddress, contractAddrUsage)
	ethClientHost := flag.String(ethClientHostName, defaultConfig.EthClientHost, ethClientHostUsage)
	ethClientPort := flag.Uint64(ethClientPortName, defaultConfig.EthClientPort, ethClientPortUsage)
	logPath := flag.String(logPathName, defaultConfig.LogPath, logPathUsage)
	chainID := flag.Int64(chainIDName, defaultConfig.ChainID, chainIDUsage)

	flag.Parse()

	parsedP2PAddrs := strings.Split(*peerP2PAddrs, ",")
	if *peerP2PAddrs == "" {
		// We handle the special case of an empty list.
		parsedP2PAddrs = []string{}
	}

	return host.Config{
		ID:                    common.BytesToAddress([]byte(*nodeID)),
		IsGenesis:             *isGenesis,
		GossipRoundDuration:   time.Duration(*gossipRoundNanos),
		HasClientRPC:          true,
		ClientRPCAddress:      *clientServerAddr,
		ClientRPCTimeout:      time.Second * time.Duration(*enclaveRPCTimeoutSecs),
		EnclaveRPCAddress:     *enclaveAddr,
		EnclaveRPCTimeout:     time.Second * time.Duration(*clientRPCTimeoutSecs),
		P2PAddress:            *ourP2PAddr,
		AllP2PAddresses:       parsedP2PAddrs,
		L1NodeHost:            *ethClientHost,
		L1NodeWebsocketPort:   uint(*ethClientPort),
		RollupContractAddress: common.BytesToAddress([]byte(*contractAddress)),
		PrivateKeyString:      *privateKeyStr,
		LogPath:               *logPath,
		ChainID:               *chainID,
	}
}
