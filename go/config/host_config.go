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

// HostInputConfig contains the configuration that was parsed from a config file / command line to start the Obscuro host.
type HostInputConfig struct {
	// Whether the host is the genesis Obscuro node
	IsGenesis bool
	// The type of the node.
	NodeType common.NodeType
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
	// L1StartHash is the hash of the L1 block we can start streaming from for all Obscuro state (e.g. management contract deployment block)
	L1StartHash gethcommon.Hash

	// MetricsEnabled defines whether the metrics are enabled or not
	MetricsEnabled bool

	// MetricsHTTPPort sets the port where the http server is available
	MetricsHTTPPort uint
}

// ToHostConfig returns a HostConfig given a HostInputConfig
func (p HostInputConfig) ToHostConfig() *HostConfig {
	return &HostConfig{
		IsGenesis:              p.IsGenesis,
		NodeType:               p.NodeType,
		HasClientRPCHTTP:       p.HasClientRPCHTTP,
		ClientRPCPortHTTP:      p.ClientRPCPortHTTP,
		HasClientRPCWebsockets: p.HasClientRPCWebsockets,
		ClientRPCPortWS:        p.ClientRPCPortWS,
		ClientRPCHost:          p.ClientRPCHost,
		EnclaveRPCAddress:      p.EnclaveRPCAddress,
		P2PBindAddress:         p.P2PBindAddress,
		P2PPublicAddress:       p.P2PPublicAddress,
		L1NodeHost:             p.L1NodeHost,
		L1NodeWebsocketPort:    p.L1NodeWebsocketPort,
		EnclaveRPCTimeout:      p.EnclaveRPCTimeout,
		L1RPCTimeout:           p.L1RPCTimeout,
		P2PConnectionTimeout:   p.P2PConnectionTimeout,
		RollupContractAddress:  p.RollupContractAddress,
		LogLevel:               p.LogLevel,
		LogPath:                p.LogPath,
		PrivateKeyString:       p.PrivateKeyString,
		L1ChainID:              p.L1ChainID,
		ObscuroChainID:         p.ObscuroChainID,
		ProfilerEnabled:        p.ProfilerEnabled,
		L1StartHash:            p.L1StartHash,
		ID:                     gethcommon.Address{},
		MetricsEnabled:         p.MetricsEnabled,
		MetricsHTTPPort:        p.MetricsHTTPPort,
	}
}

// HostConfig contains the configuration used in the Obscuro host execution. Some fields are derived from the HostInputConfig.
type HostConfig struct {
	// Whether the host is the genesis Obscuro node
	IsGenesis bool
	// The type of the node.
	NodeType common.NodeType
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
	// L1StartHash is the hash of the L1 block we can start streaming from for all Obscuro state (e.g. management contract deployment block)
	L1StartHash gethcommon.Hash

	// The host's identity derived from the L1 Private Key
	ID gethcommon.Address

	// MetricsEnabled defines whether the metrics are enabled or not
	MetricsEnabled bool

	// MetricsHTTPPort sets the port where the http server is available
	MetricsHTTPPort uint
}

// DefaultHostParsedConfig returns a HostConfig with default values.
func DefaultHostParsedConfig() *HostInputConfig {
	return &HostInputConfig{
		IsGenesis:              true,
		NodeType:               common.Sequencer,
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
		EnclaveRPCTimeout:      time.Duration(defaultRPCTimeoutSecs) * time.Second,
		L1RPCTimeout:           time.Duration(defaultL1RPCTimeoutSecs) * time.Second,
		P2PConnectionTimeout:   time.Duration(defaultP2PTimeoutSecs) * time.Second,
		RollupContractAddress:  gethcommon.BytesToAddress([]byte("")),
		LogLevel:               int(log.LvlInfo),
		LogPath:                "sys_out",
		PrivateKeyString:       "0000000000000000000000000000000000000000000000000000000000000001",
		L1ChainID:              1337,
		ObscuroChainID:         777,
		ProfilerEnabled:        false,
		L1StartHash:            common.L1RootHash{}, // this hash will not be found, host will log a warning and then stream from L1 genesis
		MetricsEnabled:         true,
		MetricsHTTPPort:        14000,
	}
}
