package config

import (
	"math/big"

	"github.com/obscuronet/go-obscuro/go/common"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/go-obscuro/go/common/log"

	gethlog "github.com/ethereum/go-ethereum/log"
)

// EnclaveConfig contains the full configuration for an Obscuro enclave service.
type EnclaveConfig struct {
	// The identity of the host the enclave service is tied to
	HostID gethcommon.Address
	// The peer-to-peer IP address of the host the enclave service is tied to
	HostAddress string
	// The address on which to serve requests
	Address string
	// The type of the node.
	NodeType common.NodeType
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
	// The management contract address on the L1 network
	ManagementContractAddress gethcommon.Address
	// LogLevel determines the verbosity of output logs
	LogLevel int
	// The path that the enclave's logs are written to
	LogPath string
	// Whether the enclave should use in-memory or persistent storage
	UseInMemoryDB bool
	// host address for the edgeless DB instance (can be empty if using InMemory DB or if attestation is disabled)
	EdgelessDBHost string
	// filepath for the sqlite DB persistence file (can be empty if a throwaway file in /tmp/ is acceptable or
	//	if using InMemory DB or if attestation is enabled)
	SqliteDBPath string
	// ProfilerEnabled starts a profiler instance
	ProfilerEnabled bool
	// MinGasPrice is the minimum gas price for mining a transaction
	MinGasPrice *big.Int
	// MessageBus L1 Address
	MessageBusAddress gethcommon.Address
	// The identity of the sequencer for the network
	SequencerID gethcommon.Address
	// A json string that specifies the prefunded addresses at the genesis of the Obscuro network
	ObscuroGenesis string
}

// DefaultEnclaveConfig returns an EnclaveConfig with default values.
func DefaultEnclaveConfig() EnclaveConfig {
	return EnclaveConfig{
		HostID:                    gethcommon.BytesToAddress([]byte("")),
		HostAddress:               "127.0.0.1:10000",
		Address:                   "127.0.0.1:11000",
		NodeType:                  common.Sequencer,
		L1ChainID:                 1337,
		ObscuroChainID:            777,
		WillAttest:                false, // todo: attestation should be on by default before production release
		ValidateL1Blocks:          false,
		GenesisJSON:               nil,
		ManagementContractAddress: gethcommon.BytesToAddress([]byte("")),
		LogLevel:                  int(gethlog.LvlInfo),
		LogPath:                   log.SysOut,
		UseInMemoryDB:             true, // todo: persistence should be on by default before production release
		EdgelessDBHost:            "",
		SqliteDBPath:              "",
		ProfilerEnabled:           false,
		MinGasPrice:               big.NewInt(1),
		SequencerID:               gethcommon.BytesToAddress([]byte("")),
		ObscuroGenesis:            "",
	}
}
