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
	nodeIDName  = "id"
	nodeIDUsage = "The 20 bytes of the node's address"

	isGenesisName  = "isGenesis"
	isGenesisUsage = "Whether the node is the first node to join the network"

	gossipRoundNanosName  = "gossipRoundNanos"
	gossipRoundNanosUsage = "The duration of the gossip round"

	clientRPCAddressName  = "clientRPCAddress"
	clientRPCAddressUsage = "The address on which to listen for client application RPC requests"

	clientRPCTimeoutSecsName  = "clientRPCTimeoutSecs"
	clientRPCTimeoutSecsUsage = "The timeout for client <-> host RPC communication"

	enclaveRPCAddressName  = "enclaveRPCAddress"
	enclaveRPCAddressUsage = "The address to use to connect to the Obscuro enclave service"

	enclaveRPCTimeoutSecsName  = "enclaveRPCTimeoutSecs"
	enclaveRPCTimeoutSecsUsage = "The timeout for host <-> enclave RPC communication"

	p2pAddressName  = "p2pAddress"
	p2pAddressUsage = "The P2P address for our node"

	peerP2PAddressesName = "peerP2PAddresses"
	peerP2PAddrsUsage    = "The P2P addresses of our peer nodes as a comma-separated list"

	l1NodeHostName  = "l1NodeHost"
	l1NodeHostUsage = "The host on which to connect to the Ethereum client"

	l1NodePortName  = "l1NodePort"
	l1NodePortUsage = "The port on which to connect to the Ethereum client"

	rollupContractAddrName  = "rollupContractAddress"
	rollupContractAddrUsage = "The management contract address on the L1"

	logPathName  = "logPath"
	logPathUsage = "The path to use for the host's log file"

	privateKeyName  = "privateKey"
	privateKeyUsage = "The private key for the L1 node account"

	chainIDName  = "chainID"
	chainIDUsage = "The ID of the L1 chain"

	defaultRPCTimeoutSecs = 3
)

func DefaultHostConfig() config.HostConfig {
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
	defaultConfig := DefaultHostConfig()

	nodeID := flag.String(nodeIDName, "", nodeIDUsage)
	isGenesis := flag.Bool(isGenesisName, defaultConfig.IsGenesis, isGenesisUsage)
	gossipRoundNanos := flag.Uint64(gossipRoundNanosName, uint64(defaultConfig.GossipRoundDuration), gossipRoundNanosUsage)
	clientRPCAddress := flag.String(clientRPCAddressName, defaultConfig.ClientRPCAddress, clientRPCAddressUsage)
	clientRPCTimeoutSecs := flag.Uint64(clientRPCTimeoutSecsName, defaultRPCTimeoutSecs, clientRPCTimeoutSecsUsage)
	enclaveRPCAddress := flag.String(enclaveRPCAddressName, defaultConfig.EnclaveRPCAddress, enclaveRPCAddressUsage)
	enclaveRPCTimeoutSecs := flag.Uint64(enclaveRPCTimeoutSecsName, defaultRPCTimeoutSecs, enclaveRPCTimeoutSecsUsage)
	p2pAddress := flag.String(p2pAddressName, defaultConfig.P2PAddress, p2pAddressUsage)
	allP2PAddresses := flag.String(peerP2PAddressesName, "", peerP2PAddrsUsage)
	l1NodeHost := flag.String(l1NodeHostName, defaultConfig.L1NodeHost, l1NodeHostUsage)
	l1NodePort := flag.Uint64(l1NodePortName, uint64(defaultConfig.L1NodeWebsocketPort), l1NodePortUsage)
	rollupContractAddress := flag.String(rollupContractAddrName, "", rollupContractAddrUsage)
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
