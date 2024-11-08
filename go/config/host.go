package config

import "time"

// HostConfig is the configuration struct for the host service.
//
//	yaml: `host`
type HostConfig struct {
	DB      *HostDB      `mapstructure:"db"`
	Debug   *HostDebug   `mapstructure:"debug"`
	Enclave *HostEnclave `mapstructure:"enclave"`
	L1      *HostL1      `mapstructure:"l1"`
	Log     *HostLog     `mapstructure:"log"`
	P2P     *HostP2P     `mapstructure:"p2p"`
	RPC     *HostRPC     `mapstructure:"rpc"`
}

// HostDB contains the configuration for the host database.
//
//	yaml: `host.db`
type HostDB struct {
	// UseInMemory specifies whether the host should use an in-memory database.
	UseInMemory bool `mapstructure:"useInMemory"`
	// PostgresHost is the host address for Postgres DB instance (can be empty if using InMemory DB)
	PostgresHost string `mapstructure:"postgresHost"`
	// SqlitePath is the filepath for the sqlite DB persistence file (can be empty if a throwaway file in /tmp/ is acceptable or
	// if using InMemory DB)
	SqlitePath string `mapstructure:"sqlitePath"`
}

// HostLog contains the configuration for the host logger.
//
//	yaml: `host.log`
type HostLog struct {
	Level int    `mapstructure:"level"`
	Path  string `mapstructure:"path"`
}

// HostRPC contains the configuration for the host RPC server.
//
//	yaml: `host.rpc`
type HostRPC struct {
	Address    string `mapstructure:"address"`
	EnableHTTP bool   `mapstructure:"enableHTTP"`
	HTTPPort   uint64 `mapstructure:"httpPort"`
	EnableWS   bool   `mapstructure:"enableWS"`
	WSPort     uint64 `mapstructure:"wsPort"`
}

// HostP2P contains the configuration for the host P2P server.
//
//	yaml: `host.p2p`
type HostP2P struct {
	// IsDisabled specifies whether the host's P2P server should be disabled for incoming connections.
	IsDisabled bool `mapstructure:"disabled"`
	// BindAddress is the address to bind the P2P server to
	// (note: this is not the publicly advertised host address, which is currently on node config).
	BindAddress string        `mapstructure:"bindAddress"`
	Timeout     time.Duration `mapstructure:"timeout"`
}

// HostL1 contains the configuration for the host's L1 client and interactions.
//
//	yaml: `host.l1`
type HostL1 struct {
	WebsocketURL string `mapstructure:"wsURL"`
	// L1BeaconUrl of the beacon chain to fetch blob data
	L1BeaconUrl string `mapstructure:"beaconURL"`
	// L1BlobArchiveUrl of the blob archive to fetch expired blob data
	L1BlobArchiveUrl string `mapstructure:"blobArchiveURL"`
	// RPCTimeout is the timeout for L1 client operations.
	RPCTimeout time.Duration `mapstructure:"rpcTimeout"`
}

// HostEnclave contains the configuration for the host's enclave(s)
//
//	yaml: `host.enclave`
type HostEnclave struct {
	// RPCAddresses is a list of managed enclave RPC addresses.
	RPCAddresses []string      `mapstructure:"rpcAddresses"`
	RPCTimeout   time.Duration `mapstructure:"rpcTimeout"`
}

// HostDebug contains the configuration for the host's debug settings.
//
//	yaml: `host.debug`
type HostDebug struct {
	EnableMetrics        bool `mapstructure:"enableMetrics"`
	MetricsHTTPPort      uint `mapstructure:"metricsHTTPPort"`
	EnableProfiler       bool `mapstructure:"enableProfiler"`
	EnableDebugNamespace bool `mapstructure:"enableDebugNamespace"`
}
