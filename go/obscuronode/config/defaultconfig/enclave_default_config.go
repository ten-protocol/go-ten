package defaultconfig

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/config"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/evm"
)

// DefaultEnclaveConfig returns an EnclaveConfig with default values.
func DefaultEnclaveConfig() config.EnclaveConfig {
	return config.EnclaveConfig{
		HostID:                    common.BytesToAddress([]byte("")),
		HostAddress:               "127.0.0.1:10000",
		Address:                   "127.0.0.1:11000",
		L1ChainID:                 1337,
		ObscuroChainID:            777,
		WillAttest:                false,
		ValidateL1Blocks:          false,
		GenesisJSON:               nil,
		SpeculativeExecution:      false,
		ManagementContractAddress: common.BytesToAddress([]byte("")),
		ERC20ContractAddresses:    []*common.Address{&evm.WBtcContract, &evm.WEthContract},
		WriteToLogs:               false,
		LogPath:                   "enclave_logs.txt",
		UseInMemoryDB:             true,
		ViewingKeysEnabled:        true,
		EdgelessDBHost:            "",
	}
}
