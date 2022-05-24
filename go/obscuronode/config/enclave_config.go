package config

import (
	"github.com/ethereum/go-ethereum/common"
)

// EnclaveConfig contains the full configuration for an Obscuro enclave service.
type EnclaveConfig struct {
	// The identity of the host the enclave service is tied to
	HostID common.Address
	// The address on which to serve requests
	Address string
	// The ID of the L1 chain
	ChainID int64
	// Whether to validate incoming L1 blocks
	ValidateL1Blocks bool
	// When validating incoming blocks, the genesis config for the L1 chain
	GenesisJSON []byte
}
