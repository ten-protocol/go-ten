package hostrunner

import (
	"flag"
	"math/big"
	"strings"
	"time"

	"github.com/obscuronet/obscuro-playground/go/obscuronode/config"

	"github.com/ethereum/go-ethereum/common"
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

	defaultRPCTimeoutSecs = 3
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

func GetDefaultConfig() config.HostConfig {
	return config.HostConfig{
		ID:                    common.BytesToAddress([]byte("")),
		IsGenesis:             true,
		GossipRoundDuration:   8333,
		HasClientRPC:          true,
		ClientRPCAddress:      "127.0.0.1:13000",
		ClientRPCTimeout:      time.Duration(defaultRPCTimeoutSecs) * time.Second,
		EnclaveRPCAddress:     "127.0.0.1:11000",
		EnclaveRPCTimeout:     time.Duration(defaultRPCTimeoutSecs) * time.Second,
		P2PAddress:            "",
		AllP2PAddresses:       []string{},
		L1NodeHost:            "127.0.0.1",
		L1NodeWebsocketPort:   8546,
		RollupContractAddress: common.BytesToAddress([]byte("")),
		LogPath:               "host_logs.txt",
		PrivateKeyString:      "0000000000000000000000000000000000000000000000000000000000000001",
		ChainID:               *big.NewInt(1337),
	}
}

func ParseCLIArgs() config.HostConfig {
	defaultConfig := GetDefaultConfig()

	nodeID := flag.String(nodeIDName, "", nodeIDUsage)
	isGenesis := flag.Bool(isGenesisName, defaultConfig.IsGenesis, isGenesisUsage)
	gossipRoundNanos := flag.Uint64(gossipRoundNanosName, uint64(defaultConfig.GossipRoundDuration), gossipRoundNanosUsage)
	clientRPCAddress := flag.String(clientServerAddrName, defaultConfig.ClientRPCAddress, clientServerAddrUsage)
	clientRPCTimeoutSecs := flag.Uint64(clientRPCTimeoutSecsName, defaultRPCTimeoutSecs, clientRPCTimeoutSecsUsage)
	enclaveRPCAddress := flag.String(enclaveAddrName, defaultConfig.EnclaveRPCAddress, enclaveAddrUsage)
	enclaveRPCTimeoutSecs := flag.Uint64(enclaveRPCTimeoutSecsName, defaultRPCTimeoutSecs, enclaveRPCTimeoutSecsUsage)
	p2pAddress := flag.String(ourP2PAddrName, defaultConfig.P2PAddress, ourP2PAddrUsage)
	allP2PAddresses := flag.String(peerP2PAddrsName, "", peerP2PAddrsUsage)
	l1NodeHost := flag.String(ethClientHostName, defaultConfig.L1NodeHost, ethClientHostUsage)
	l1NodePort := flag.Uint64(ethClientPortName, uint64(defaultConfig.L1NodeWebsocketPort), ethClientPortUsage)
	rollupContractAddress := flag.String(contractAddrName, "", contractAddrUsage)
	logPath := flag.String(logPathName, defaultConfig.LogPath, logPathUsage)
	chainID := flag.Int64(chainIDName, defaultConfig.ChainID.Int64(), chainIDUsage)
	privateKeyStr := flag.String(privateKeyName, defaultConfig.PrivateKeyString, privateKeyUsage)

	flag.Parse()

	parsedP2PAddrs := strings.Split(*allP2PAddresses, ",")
	if *allP2PAddresses == "" {
		// We handle the special case of an empty list.
		parsedP2PAddrs = []string{}
	}

	return config.HostConfig{
		ID:                    common.BytesToAddress([]byte(*nodeID)),
		IsGenesis:             *isGenesis,
		GossipRoundDuration:   time.Duration(*gossipRoundNanos),
		HasClientRPC:          true,
		ClientRPCAddress:      *clientRPCAddress,
		ClientRPCTimeout:      time.Second * time.Duration(*enclaveRPCTimeoutSecs),
		EnclaveRPCAddress:     *enclaveRPCAddress,
		EnclaveRPCTimeout:     time.Second * time.Duration(*clientRPCTimeoutSecs),
		P2PAddress:            *p2pAddress,
		AllP2PAddresses:       parsedP2PAddrs,
		L1NodeHost:            *l1NodeHost,
		L1NodeWebsocketPort:   uint(*l1NodePort),
		RollupContractAddress: common.BytesToAddress([]byte(*rollupContractAddress)),
		PrivateKeyString:      *privateKeyStr,
		LogPath:               *logPath,
		ChainID:               *big.NewInt(*chainID),
	}
}
