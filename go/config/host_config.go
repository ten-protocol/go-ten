package config

import (
	"time"

	"github.com/obscuronet/go-obscuro/go/common"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

const (
	defaultRPCTimeoutSecs   = 10
	defaultL1RPCTimeoutSecs = 15
	defaultP2PTimeoutSecs   = 10
)

// HostConfig contains the full configuration for an Obscuro host.
type HostConfig struct {
	// The host's identity
	ID gethcommon.Address
	// Whether the host is the genesis Obscuro node
	IsGenesis bool
	// The type of the node.
	NodeType common.NodeType
	// Duration of the gossip round
	GossipRoundDuration time.Duration
	// Whether to serve client RPC requests over HTTP
	HasClientRPCHTTP bool
	// Port on which to handle HTTP client RPC requests
	ClientRPCPortHTTP uint64
	// Whether to serve client RPC requests over websockets
	HasClientRPCWebsockets bool
	// Port on which to handle websocket client RPC requests
	ClientRPCPortWS uint64
	// Host on which to handle client RPC requests
	ClientRPCHost string
	// Address on which to connect to the enclave
	EnclaveRPCAddress string
	// P2PBindAddress is the address where the P2P server is bound to
	P2PBindAddress string
	// P2PPublicAddress is the advertised P2P server address
	P2PPublicAddress string
	// The host of the connected L1 node
	L1NodeHost string
	// The websocket port of the connected L1 node
	L1NodeWebsocketPort uint
	// Timeout duration for RPC requests from client applications
	ClientRPCTimeout time.Duration
	// Timeout duration for RPC requests to the enclave service
	EnclaveRPCTimeout time.Duration
	// Timeout duration for connecting to, and communicating with, the L1 node
	L1RPCTimeout time.Duration
	// Timeout duration for messaging between hosts.
	P2PConnectionTimeout time.Duration
	// The rollup contract address on the L1 network
	RollupContractAddress gethcommon.Address
	// LogLevel determines the verbosity of output logs
	LogLevel int
	// The path that the node's logs are written to
	LogPath string
	// The stringified private key for the host's L1 wallet
	PrivateKeyString string
	// The ID of the L1 chain
	L1ChainID int64
	// The ID of the Obscuro chain
	ObscuroChainID int64
	// ProfilerEnabled starts a profiler instance
	ProfilerEnabled bool
}

// DefaultHostConfig returns a HostConfig with default values.
func DefaultHostConfig() HostConfig {
	return HostConfig{
		ID:                     gethcommon.BytesToAddress([]byte("")),
		IsGenesis:              true,
		NodeType:               common.Aggregator,
		GossipRoundDuration:    8333,
		HasClientRPCHTTP:       true,
		ClientRPCPortHTTP:      13000,
		HasClientRPCWebsockets: true,
		ClientRPCPortWS:        13001,
		ClientRPCHost:          "127.0.0.1",
		EnclaveRPCAddress:      "127.0.0.1:11000",
		P2PBindAddress:         "0.0.0.0:10000",
		P2PPublicAddress:       "127.0.0.1:10000",
		L1NodeHost:             "127.0.0.1",
		L1NodeWebsocketPort:    8546,
		ClientRPCTimeout:       time.Duration(defaultRPCTimeoutSecs) * time.Second,
		EnclaveRPCTimeout:      time.Duration(defaultRPCTimeoutSecs) * time.Second,
		L1RPCTimeout:           time.Duration(defaultL1RPCTimeoutSecs) * time.Second,
		P2PConnectionTimeout:   time.Duration(defaultP2PTimeoutSecs) * time.Second,
		RollupContractAddress:  gethcommon.BytesToAddress([]byte("")),
		LogLevel:               int(log.LvlInfo),
		LogPath:                "",
		PrivateKeyString:       "0000000000000000000000000000000000000000000000000000000000000001",
		L1ChainID:              1337,
		ObscuroChainID:         777,
		ProfilerEnabled:        false,
	}
}
