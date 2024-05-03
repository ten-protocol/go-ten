package config

import (
	"fmt"
	"time"

	"github.com/ten-protocol/go-ten/go/common"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// HostInputConfig used for parsing default config or partial override from yaml file.
type HostInputConfig struct {
	IsGenesis                 bool     `yaml:"isGenesis"`
	NodeType                  string   `yaml:"nodeType"`
	ClientRPCPortHTTP         uint64   `yaml:"clientRPCPortHTTP"`
	ClientRPCPortWS           uint64   `yaml:"clientRPCPortWS"`
	ClientRPCHost             string   `yaml:"clientRPCHost"`
	EnclaveRPCAddresses       []string `yaml:"enclaveRPCAddresses"`
	P2PBindAddress            string   `yaml:"p2pBindAddress"`
	P2PPublicAddress          string   `yaml:"p2pPublicAddress"`
	L1WebsocketURL            string   `yaml:"l1WebsocketURL"`
	EnclaveRPCTimeout         int      `yaml:"enclaveRPCTimeout"`
	L1RPCTimeout              int      `yaml:"l1RPCTimeout"`
	P2PConnectionTimeout      int      `yaml:"p2pConnectionTimeout"`
	ManagementContractAddress string   `yaml:"managementContractAddress"`
	MessageBusAddress         string   `yaml:"messageBusAddress"`
	LogLevel                  int      `yaml:"logLevel"`
	LogPath                   string   `yaml:"logPath"`
	PrivateKey                string   `yaml:"privateKey"`
	L1ChainID                 int64    `yaml:"l1ChainID"`
	TenChainID                int64    `yaml:"tenChainID"`
	ProfilerEnabled           bool     `yaml:"profilerEnabled"`
	L1StartHash               string   `yaml:"l1StartHash"`
	SequencerP2PAddress       string   `yaml:"sequencerP2PAddress"`
	MetricsEnabled            bool     `yaml:"metricsEnabled"`
	MetricsHTTPPort           uint     `yaml:"metricsHTTPPort"`
	UseInMemoryDB             bool     `yaml:"useInMemoryDB"`
	PostgresDBHost            string   `yaml:"postgresDBHost"`
	SqliteDBPath              string   `yaml:"sqliteDBPath"`
	LevelDBPath               string   `yaml:"levelDBPath"`
	DebugNamespaceEnabled     bool     `yaml:"debugNamespaceEnabled"`
	BatchInterval             int      `yaml:"batchInterval"`
	MaxBatchInterval          int      `yaml:"maxBatchInterval"`
	RollupInterval            int      `yaml:"rollupInterval"`
	L1BlockTime               int      `yaml:"l1BlockTime"`
	IsInboundP2PDisabled      bool     `yaml:"isInboundP2PDisabled"`
	MaxRollupSize             uint64   `yaml:"maxRollupSize"`
}

// ToHostConfig Generates an HostConfig from flags or yaml to one with proper typing
func (p *HostInputConfig) ToHostConfig() (*HostConfig, error) {
	// calculated
	nodeType, err := common.ToNodeType(p.NodeType)
	if err != nil {
		return nil, fmt.Errorf("unrecognized node type %s: %w", p.NodeType, err)
	}

	hostConfig := &HostConfig{
		IsGenesis:             p.IsGenesis,
		NodeType:              nodeType,
		ClientRPCPortHTTP:     p.ClientRPCPortHTTP,
		ClientRPCPortWS:       p.ClientRPCPortWS,
		ClientRPCHost:         p.ClientRPCHost,
		EnclaveRPCAddresses:   p.EnclaveRPCAddresses,
		P2PBindAddress:        p.P2PBindAddress,
		P2PPublicAddress:      p.P2PPublicAddress,
		L1WebsocketURL:        p.L1WebsocketURL,
		LogLevel:              p.LogLevel,
		LogPath:               p.LogPath,
		PrivateKey:            p.PrivateKey,
		L1ChainID:             p.L1ChainID,
		TenChainID:            p.TenChainID,
		SequencerP2PAddress:   p.SequencerP2PAddress,
		ProfilerEnabled:       p.ProfilerEnabled,
		MetricsEnabled:        p.MetricsEnabled,
		MetricsHTTPPort:       p.MetricsHTTPPort,
		UseInMemoryDB:         p.UseInMemoryDB,
		PostgresDBHost:        p.PostgresDBHost,
		SqliteDBPath:          p.SqliteDBPath,
		DebugNamespaceEnabled: p.DebugNamespaceEnabled,
		IsInboundP2PDisabled:  p.IsInboundP2PDisabled,
		MaxRollupSize:         p.MaxRollupSize,
	}

	// boolean
	hostConfig.HasClientRPCHTTP = p.ClientRPCPortHTTP != 0
	hostConfig.HasClientRPCWebsockets = p.ClientRPCPortWS != 0

	// durations
	hostConfig.EnclaveRPCTimeout = time.Duration(p.EnclaveRPCTimeout) * time.Second
	hostConfig.L1RPCTimeout = time.Duration(p.L1RPCTimeout) * time.Second
	hostConfig.P2PConnectionTimeout = time.Duration(p.P2PConnectionTimeout) * time.Second
	hostConfig.BatchInterval = time.Duration(p.BatchInterval) * time.Second
	hostConfig.MaxBatchInterval = time.Duration(p.MaxBatchInterval) * time.Second
	hostConfig.RollupInterval = time.Duration(p.RollupInterval) * time.Second
	hostConfig.L1BlockTime = time.Duration(p.L1BlockTime) * time.Second

	// address
	hostConfig.ManagementContractAddress = gethcommon.HexToAddress(p.ManagementContractAddress)
	hostConfig.MessageBusAddress = gethcommon.HexToAddress(p.MessageBusAddress)
	hostConfig.L1StartHash = gethcommon.HexToHash(p.L1StartHash)

	return hostConfig, nil
}

// HostConfig contains the full configuration for a Ten Host service.
type HostConfig struct {
	/////
	// TEN NETWORK CONFIG (these properties are the same for all obscuro nodes on the network)
	/////

	// The ID of the L1 chain
	L1ChainID int64
	// The ID of the Obscuro chain
	TenChainID int64
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

	/////
	// NODE CONFIG
	/////

	// ID, the host's identity derived from the L1 Private Key
	ID gethcommon.Address
	// PrivateKey, the stringified private key for the host's L1 wallet
	PrivateKey string
	// IsGenesis, whether the host is the genesis Obscuro node
	IsGenesis bool
	// NodeType, the type of the node.
	NodeType common.NodeType
	// LogLevel determines the verbosity of output logs
	LogLevel int
	// LogPath, the path that the node's logs are written to
	LogPath string
	// UseInMemoryDB, whether the host should use in-memory or persistent storage
	UseInMemoryDB bool
	// Host address for Postgres DB instance (can be empty if using InMemory DB or if attestation is disabled)
	PostgresDBHost string
	// SqliteDBPath filepath for the sqlite DB persistence file (can be empty if a throwaway file in /tmp/ is acceptable or if using InMemory DB)
	SqliteDBPath string

	//////
	// NODE NETWORKING
	//////

	// HasClientRPCHTTP whether to serve client RPC requests over HTTP
	HasClientRPCHTTP bool
	// ClientRPCPortHTTP, port on which to handle HTTP client RPC requests
	ClientRPCPortHTTP uint64
	// HasClientRPCWebsockets, whether to serve client RPC requests over websockets
	HasClientRPCWebsockets bool
	// ClientRPCPortWS, port on which to handle websocket client RPC requests
	ClientRPCPortWS uint64
	// ClientRPCHost on which to handle client RPC requests
	ClientRPCHost string
	// EnclaveRPCAddresses, addresses on which to connect to the node's enclaves (HA setups may have multiple)
	EnclaveRPCAddresses []string
	// P2PBindAddress is the address where the P2P server is bound to
	P2PBindAddress string
	// P2PPublicAddress is the advertised P2P server address
	P2PPublicAddress string
	// L1WebsocketURL is the RPC address for interactions with the L1
	L1WebsocketURL string
	// EnclaveRPCTimeout duration for RPC requests to the enclave service
	EnclaveRPCTimeout time.Duration
	// L1RPCTimeout duration for connecting to, and communicating with, the L1 node
	L1RPCTimeout time.Duration
	// P2PConnectionTimeout duration for messaging between hosts.
	P2PConnectionTimeout time.Duration
	// ProfilerEnabled starts a profiler instance
	ProfilerEnabled bool
	// MetricsEnabled defines whether the metrics are enabled or not
	MetricsEnabled bool
	// MetricsHTTPPort sets the port where the http server is available
	MetricsHTTPPort uint
	// DebugNamespaceEnabled enables the debug namespace handler in the host rpc server
	DebugNamespaceEnabled bool
	// IsInboundP2PDisabled, whether p2p is enabled or not
	IsInboundP2PDisabled bool
}
