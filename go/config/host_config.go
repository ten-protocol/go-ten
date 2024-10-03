package config

import (
	"time"

	"github.com/ten-protocol/go-ten/go/common"

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
	EnclaveRPCAddresses []string
	// P2PBindAddress is the address where the P2P server is bound to
	P2PBindAddress string
	// P2PPublicAddress is the advertised P2P server address
	P2PPublicAddress string
	// L1WebsocketURL is the RPC address for interactions with the L1
	L1WebsocketURL string
	// Timeout duration for RPC requests to the enclave service
	EnclaveRPCTimeout time.Duration
	// Timeout duration for connecting to, and communicating with, the L1 node
	L1RPCTimeout time.Duration
	// Timeout duration for messaging between hosts.
	P2PConnectionTimeout time.Duration
	// P2P address of network sequencer node
	SequencerP2PAddress string
	// The rollup contract address on the L1 network
	ManagementContractAddress gethcommon.Address
	// The message bus contract address on the L1 network
	MessageBusAddress gethcommon.Address
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
	// L1BeaconUrl of the beacon chain to fetch blob data
	L1BeaconUrl string
	// L1BlobArchiveUrl of the blob archive to fetch expired blob data
	L1BlobArchiveUrl string

	// MetricsEnabled defines whether the metrics are enabled or not
	MetricsEnabled bool

	// MetricsHTTPPort sets the port where the http server is available
	MetricsHTTPPort uint

	// UseInMemoryDB sets whether the host should use in-memory or persistent storage
	UseInMemoryDB bool

	// PostgresDBHost db url for connecting to Postgres host database
	PostgresDBHost string

	// DebugNamespaceEnabled enables the debug namespace handler in the host rpc server
	DebugNamespaceEnabled bool

	// Min interval before creating the next batch (only used by Sequencer nodes)
	BatchInterval time.Duration

	// MaxBatchInterval is the max interval between batches, if this is set higher than BatchInterval, the host will
	// not create empty batches until the MaxBatchInterval is reached or a transaction is received.
	MaxBatchInterval time.Duration

	// Min interval before creating the next rollup (only used by Sequencer nodes)
	RollupInterval time.Duration

	// The expected time between blocks on the L1 network
	L1BlockTime time.Duration

	// CrossChainInterval - The interval at which the host will check for new cross chain data to submit
	CrossChainInterval time.Duration

	// Whether inbound p2p is enabled or not
	IsInboundP2PDisabled bool

	// MaxRollupSize specifies the threshold size which the sequencer-host publishes a rollup
	MaxRollupSize uint64
}

// ToHostConfig returns a HostConfig given a HostInputConfig
func (p HostInputConfig) ToHostConfig() *HostConfig {
	return &HostConfig{
		IsGenesis:                 p.IsGenesis,
		NodeType:                  p.NodeType,
		HasClientRPCHTTP:          p.HasClientRPCHTTP,
		ClientRPCPortHTTP:         p.ClientRPCPortHTTP,
		HasClientRPCWebsockets:    p.HasClientRPCWebsockets,
		ClientRPCPortWS:           p.ClientRPCPortWS,
		ClientRPCHost:             p.ClientRPCHost,
		EnclaveRPCAddresses:       p.EnclaveRPCAddresses,
		P2PBindAddress:            p.P2PBindAddress,
		P2PPublicAddress:          p.P2PPublicAddress,
		L1WebsocketURL:            p.L1WebsocketURL,
		EnclaveRPCTimeout:         p.EnclaveRPCTimeout,
		L1RPCTimeout:              p.L1RPCTimeout,
		P2PConnectionTimeout:      p.P2PConnectionTimeout,
		ManagementContractAddress: p.ManagementContractAddress,
		MessageBusAddress:         p.MessageBusAddress,
		LogLevel:                  p.LogLevel,
		LogPath:                   p.LogPath,
		PrivateKeyString:          p.PrivateKeyString,
		L1ChainID:                 p.L1ChainID,
		ObscuroChainID:            p.ObscuroChainID,
		ProfilerEnabled:           p.ProfilerEnabled,
		L1StartHash:               p.L1StartHash,
		SequencerP2PAddress:       p.SequencerP2PAddress,
		ID:                        gethcommon.Address{},
		MetricsEnabled:            p.MetricsEnabled,
		MetricsHTTPPort:           p.MetricsHTTPPort,
		UseInMemoryDB:             p.UseInMemoryDB,
		PostgresDBHost:            p.PostgresDBHost,
		DebugNamespaceEnabled:     p.DebugNamespaceEnabled,
		BatchInterval:             p.BatchInterval,
		MaxBatchInterval:          p.MaxBatchInterval,
		RollupInterval:            p.RollupInterval,
		L1BlockTime:               p.L1BlockTime,
		CrossChainInterval:        p.CrossChainInterval,
		IsInboundP2PDisabled:      p.IsInboundP2PDisabled,
		MaxRollupSize:             p.MaxRollupSize,
		L1BeaconUrl:               p.L1BeaconUrl,
	}
}

// HostConfig contains the configuration used in the Obscuro host execution. Some fields are derived from the HostInputConfig.
type HostConfig struct {
	/////
	// OBSCURO NETWORK CONFIG (these properties are the same for all obscuro nodes on the network)
	/////

	// The ID of the L1 chain
	L1ChainID int64
	// The ID of the Obscuro chain
	ObscuroChainID int64
	// L1StartHash is the hash of the L1 block we can start streaming from for all Obscuro state (e.g. management contract deployment block)
	L1StartHash gethcommon.Hash
	// The address of the sequencer node's P2P server
	SequencerP2PAddress string
	// The rollup contract address on the L1 network
	ManagementContractAddress gethcommon.Address
	// The message bus contract address on the L1 network
	MessageBusAddress gethcommon.Address
	// Min interval before creating the next batch (only used by Sequencer nodes)
	BatchInterval time.Duration
	// MaxBatchInterval is the max interval between batches, if this is set higher than BatchInterval, the host will
	// not create empty batches until the MaxBatchInterval is reached or a transaction is received.
	MaxBatchInterval time.Duration
	// Min interval before creating the next rollup (only used by Sequencer nodes)
	RollupInterval time.Duration
	// MaxRollupSize is the max size of the rollup
	MaxRollupSize uint64
	// The expected time between blocks on the L1 network
	L1BlockTime time.Duration
	// CrossChainInterval - The interval at which the host will check for new cross chain data to submit
	CrossChainInterval time.Duration
	// L1BeaconUrl of the beacon chain to fetch blob data
	L1BeaconUrl string
	// L1BlobArchiveUrl of the blob archive to fetch expired blob data
	L1BlobArchiveUrl string

	/////
	// NODE CONFIG
	/////

	// The host's identity derived from the L1 Private Key
	ID gethcommon.Address
	// The stringified private key for the host's L1 wallet
	PrivateKeyString string
	// Whether the host is the genesis Obscuro node
	IsGenesis bool
	// The type of the node.
	NodeType common.NodeType
	// LogLevel determines the verbosity of output logs
	LogLevel int
	// The path that the node's logs are written to
	LogPath string
	// Whether the host should use in-memory or persistent storage
	UseInMemoryDB bool
	// Host address for Postgres DB instance (can be empty if using InMemory DB or if attestation is disabled)
	PostgresDBHost string
	// filepath for the sqlite DB persistence file (can be empty if a throwaway file in /tmp/ is acceptable or
	// if using InMemory DB)
	SqliteDBPath string

	//////
	// NODE NETWORKING
	//////

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
	// Addresses on which to connect to the node's enclaves (HA setups may have multiple)
	EnclaveRPCAddresses []string
	// P2PBindAddress is the address where the P2P server is bound to
	P2PBindAddress string
	// P2PPublicAddress is the advertised P2P server address
	P2PPublicAddress string
	// L1WebsocketURL is the RPC address for interactions with the L1
	L1WebsocketURL string
	// Timeout duration for RPC requests to the enclave service
	EnclaveRPCTimeout time.Duration
	// Timeout duration for connecting to, and communicating with, the L1 node
	L1RPCTimeout time.Duration
	// Timeout duration for messaging between hosts.
	P2PConnectionTimeout time.Duration
	// ProfilerEnabled starts a profiler instance
	ProfilerEnabled bool
	// MetricsEnabled defines whether the metrics are enabled or not
	MetricsEnabled bool
	// MetricsHTTPPort sets the port where the http server is available
	MetricsHTTPPort uint
	// DebugNamespaceEnabled enables the debug namespace handler in the host rpc server
	DebugNamespaceEnabled bool
	// Whether p2p is enabled or not
	IsInboundP2PDisabled bool
}

// DefaultHostParsedConfig returns a HostConfig with default values.
func DefaultHostParsedConfig() *HostInputConfig {
	return &HostInputConfig{
		IsGenesis:                 true,
		NodeType:                  common.Sequencer,
		HasClientRPCHTTP:          true,
		ClientRPCPortHTTP:         80,
		HasClientRPCWebsockets:    true,
		ClientRPCPortWS:           81,
		ClientRPCHost:             "127.0.0.1",
		EnclaveRPCAddresses:       []string{"127.0.0.1:11000"},
		P2PBindAddress:            "0.0.0.0:10000",
		P2PPublicAddress:          "127.0.0.1:10000",
		L1WebsocketURL:            "ws://127.0.0.1:8546",
		EnclaveRPCTimeout:         time.Duration(defaultRPCTimeoutSecs) * time.Second,
		L1RPCTimeout:              time.Duration(defaultL1RPCTimeoutSecs) * time.Second,
		P2PConnectionTimeout:      time.Duration(defaultP2PTimeoutSecs) * time.Second,
		ManagementContractAddress: gethcommon.BytesToAddress([]byte("")),
		MessageBusAddress:         gethcommon.BytesToAddress([]byte("")),
		LogLevel:                  int(log.LvlInfo),
		LogPath:                   "",
		PrivateKeyString:          "0000000000000000000000000000000000000000000000000000000000000001",
		L1ChainID:                 1337,
		ObscuroChainID:            443,
		ProfilerEnabled:           false,
		L1StartHash:               common.L1BlockHash{}, // this hash will not be found, host will log a warning and then stream from L1 genesis
		SequencerP2PAddress:       "127.0.0.1:10000",
		MetricsEnabled:            true,
		MetricsHTTPPort:           14000,
		UseInMemoryDB:             true,
		DebugNamespaceEnabled:     false, BatchInterval: 1 * time.Second,
		MaxBatchInterval:     1 * time.Second,
		RollupInterval:       5 * time.Second,
		L1BlockTime:          15 * time.Second,
		IsInboundP2PDisabled: false,
		MaxRollupSize:        1024 * 128,
		CrossChainInterval:   6 * time.Second,
		// L1BeaconUrl:          "127.0.0.1:12600", // default port for the beacon chain if not specified
		L1BeaconUrl: "eth2network:12600", // local testnet defaults here
	}
}
