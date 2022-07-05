package config

import (
	"github.com/ethereum/go-ethereum/common"
)

// EnclaveConfig contains the full configuration for an Obscuro enclave service.
type EnclaveConfig struct {
	// The identity of the host the enclave service is tied to
	HostID common.Address
	// The peer-to-peer IP address of the host the enclave service is tied to
	HostAddress string
	// The address on which to serve requests
	Address string
	// The ID of the L1 chain
	L1ChainID int64
	// The ID of the Obscuro chain
	ObscuroChainID int64
	// Whether to produce a verified attestation report
	WillAttest bool
	// Whether to validate incoming L1 blocks
	ValidateL1Blocks bool
	// When validating incoming blocks, the genesis config for the L1 chain
	GenesisJSON []byte
	// Toggles the speculative execution background process
	SpeculativeExecution bool
	// The management contract address on the L1 network
	ManagementContractAddress common.Address
	// The addresses of ERC20 contracts to monitor on the L1 network
	ERC20ContractAddresses []*common.Address
	// LogLevel determines the verbosity of output logs
	LogLevel string
	// Whether to redirect the enclave's output to the log file.
	WriteToLogs bool
	// The path that the node's logs are written to
	LogPath string
	// Whether the enclave should use in-memory or persistent storage
	UseInMemoryDB bool
	// Whether the client and the enclave should encrypt sensitive requests and responses.
	// TODO - Consider removing this option and forcing the simulations to generate viewing keys.
	ViewingKeysEnabled bool
	// host address for the edgeless DB instance (can be empty if using InMemory DB or if attestation is disabled)
	EdgelessDBHost string
	// filepath for the sqlite DB persistence file (can be empty if a throwaway file in /tmp/ is acceptable or
	//	if using InMemory DB or if attestation is enabled)
	SqliteDBPath string
}

// DefaultEnclaveConfig returns an EnclaveConfig with default values.
func DefaultEnclaveConfig() EnclaveConfig {
	return EnclaveConfig{
		HostID:                    common.BytesToAddress([]byte("")),
		HostAddress:               "127.0.0.1:10000",
		Address:                   "127.0.0.1:11000",
		L1ChainID:                 1337,
		ObscuroChainID:            777,
		WillAttest:                false, // todo: attestation should be on by default before production release
		ValidateL1Blocks:          false,
		GenesisJSON:               nil,
		SpeculativeExecution:      false,
		ManagementContractAddress: common.BytesToAddress([]byte("")),
		ERC20ContractAddresses:    []*common.Address{},
		LogLevel:                  "debug",
		WriteToLogs:               false,
		LogPath:                   "enclave_logs.txt",
		UseInMemoryDB:             true, // todo: persistence should be on by default before production release
		ViewingKeysEnabled:        true,
		EdgelessDBHost:            "",
		SqliteDBPath:              "",
	}
}
