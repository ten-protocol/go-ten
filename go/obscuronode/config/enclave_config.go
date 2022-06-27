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
}
