package config

import "time"

// EnclaveConfig is the configuration struct for the enclave service.
//
//	yaml: `enclave`
type EnclaveConfig struct {
	// EnableAttestation specifies whether the enclave will produce verified attestation report.
	EnableAttestation         bool `mapstructure:"enableAttestation"`
	StoreExecutedTransactions bool `mapstructure:"storeExecutedTransactions"`

	DB    *EnclaveDB    `mapstructure:"db"`
	Debug *EnclaveDebug `mapstructure:"debug"`
	Log   *EnclaveLog   `mapstructure:"log"`
	RPC   *EnclaveRPC   `mapstructure:"rpc"`
}

// EnclaveDB contains the configuration for the enclave database.
//
//	yaml: `enclave.db`
type EnclaveDB struct {
	// UseInMemory specifies whether the enclave should use an in-memory database.
	UseInMemory bool `mapstructure:"useInMemory"`
	// EdgelessDBHost is the host address for the edgeless DB instance (can be empty if using InMemory DB or if attestation is disabled).
	EdgelessDBHost string `mapstructure:"edgelessDBHost"`
	// SqliteDBPath is the filepath for the sqlite DB persistence file (can be empty if a throwaway file in /tmp/ is acceptable or
	// if using InMemory DB or if attestation is enabled).
	SqlitePath string `mapstructure:"sqlitePath"`
}

// EnclaveDebug contains the configuration for the enclave debug.
//
//	yaml: `enclave.debug`
type EnclaveDebug struct {
	EnableDebugNamespace bool `mapstructure:"enableDebugNamespace"`
	EnableProfiler       bool `mapstructure:"enableProfiler"`
}

// EnclaveLog contains the configuration for the enclave logger.
//
//	yaml: `enclave.log`
type EnclaveLog struct {
	Level int    `mapstructure:"level"`
	Path  string `mapstructure:"path"`
}

// EnclaveRPC contains the configuration for the enclave RPC server.
//
//	yaml: `enclave.rpc`
type EnclaveRPC struct {
	BindAddress string `mapstructure:"bindAddress"`
	// Timeout - calls that are longer than this will be cancelled, to prevent resource starvation
	// (normally, the context is propagated from the host, but in some cases like the evm, we have to create a context)
	Timeout time.Duration `mapstructure:"timeout"`
}
