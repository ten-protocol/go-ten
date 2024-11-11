package config

import (
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/config"
)

// For now, this is the bridge between TenConfig and the config used internally by the host service.

// HostConfig contains the configuration used in the Obscuro host execution.
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
	// L1BeaconUrl of the beacon chain to fetch blob data
	L1BeaconUrl string
	// L1BlobArchiveUrl of the blob archive to fetch expired blob data
	L1BlobArchiveUrl string
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

func HostConfigFromTenConfig(tenCfg *config.TenConfig) *HostConfig {
	return &HostConfig{
		PrivateKeyString: tenCfg.Node.PrivateKeyString,
		IsGenesis:        tenCfg.Node.IsGenesis,
		NodeType:         tenCfg.Node.NodeType,

		L1ChainID:      tenCfg.Network.L1.ChainID,
		ObscuroChainID: tenCfg.Network.ChainID,

		L1StartHash:               tenCfg.Network.L1.StartHash,
		L1BlockTime:               tenCfg.Network.L1.BlockTime,
		SequencerP2PAddress:       tenCfg.Network.Sequencer.P2PAddress,
		ManagementContractAddress: tenCfg.Network.L1.L1Contracts.ManagementContract,
		MessageBusAddress:         tenCfg.Network.L1.L1Contracts.MessageBusContract,

		BatchInterval:      tenCfg.Network.Batch.Interval,
		MaxBatchInterval:   tenCfg.Network.Batch.MaxInterval,
		RollupInterval:     tenCfg.Network.Rollup.Interval,
		MaxRollupSize:      tenCfg.Network.Rollup.MaxSize,
		CrossChainInterval: tenCfg.Network.CrossChain.Interval,

		LogLevel: tenCfg.Host.Log.Level,
		LogPath:  tenCfg.Host.Log.Path,

		UseInMemoryDB:  tenCfg.Host.DB.UseInMemory,
		PostgresDBHost: tenCfg.Host.DB.PostgresHost,
		SqliteDBPath:   tenCfg.Host.DB.SqlitePath,

		HasClientRPCHTTP:       tenCfg.Host.RPC.EnableHTTP,
		ClientRPCPortHTTP:      tenCfg.Host.RPC.HTTPPort,
		HasClientRPCWebsockets: tenCfg.Host.RPC.EnableWS,
		ClientRPCPortWS:        tenCfg.Host.RPC.WSPort,
		ClientRPCHost:          tenCfg.Host.RPC.Address,

		EnclaveRPCAddresses: tenCfg.Host.Enclave.RPCAddresses,
		EnclaveRPCTimeout:   tenCfg.Host.Enclave.RPCTimeout,

		IsInboundP2PDisabled: tenCfg.Host.P2P.IsDisabled,
		P2PBindAddress:       tenCfg.Host.P2P.BindAddress,
		P2PConnectionTimeout: tenCfg.Host.P2P.Timeout,
		P2PPublicAddress:     tenCfg.Node.HostAddress,

		L1WebsocketURL:   tenCfg.Host.L1.WebsocketURL,
		L1BeaconUrl:      tenCfg.Host.L1.L1BeaconUrl,
		L1BlobArchiveUrl: tenCfg.Host.L1.L1BlobArchiveUrl,
		L1RPCTimeout:     tenCfg.Host.L1.RPCTimeout,

		ProfilerEnabled:       tenCfg.Host.Debug.EnableProfiler,
		MetricsEnabled:        tenCfg.Host.Debug.EnableMetrics,
		MetricsHTTPPort:       tenCfg.Host.Debug.MetricsHTTPPort,
		DebugNamespaceEnabled: tenCfg.Host.Debug.EnableDebugNamespace,
	}
}
