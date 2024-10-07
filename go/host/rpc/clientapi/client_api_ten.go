package clientapi

import (
	"context"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/host"
)

// TenAPI implements Ten-specific JSON RPC operations.
type TenAPI struct {
	host host.Host
}

func NewTenAPI(host host.Host) *TenAPI {
	return &TenAPI{
		host: host,
	}
}

// Health returns the health status of TEN host + enclave + db
func (api *TenAPI) Health(ctx context.Context) (*host.HealthCheck, error) {
	return api.host.HealthCheck(ctx)
}

// Config returns the config status of TEN host + enclave + db
func (api *TenAPI) Config() (*ChecksumFormattedTenNetworkConfig, error) {
	config, err := api.host.TenConfig()
	if err != nil {
		return nil, err
	}
	return checksumFormatted(config), nil
}

// ChecksumFormattedTenNetworkConfig serialises the addresses as EIP55 checksum addresses.
type ChecksumFormattedTenNetworkConfig struct {
	ManagementContractAddress       gethcommon.AddressEIP55
	L1StartHash                     gethcommon.Hash
	MessageBusAddress               gethcommon.AddressEIP55
	L2MessageBusAddress             gethcommon.AddressEIP55
	ImportantContracts              map[string]gethcommon.AddressEIP55 // map of contract name to address
	TransactionPostProcessorAddress gethcommon.AddressEIP55
}

func checksumFormatted(info *common.TenNetworkInfo) *ChecksumFormattedTenNetworkConfig {
	importantContracts := make(map[string]gethcommon.AddressEIP55)
	for name, addr := range info.ImportantContracts {
		importantContracts[name] = gethcommon.AddressEIP55(addr)
	}
	return &ChecksumFormattedTenNetworkConfig{
		ManagementContractAddress:       gethcommon.AddressEIP55(info.ManagementContractAddress),
		L1StartHash:                     info.L1StartHash,
		MessageBusAddress:               gethcommon.AddressEIP55(info.MessageBusAddress),
		L2MessageBusAddress:             gethcommon.AddressEIP55(info.L2MessageBusAddress),
		ImportantContracts:              importantContracts,
		TransactionPostProcessorAddress: gethcommon.AddressEIP55(info.TransactionPostProcessorAddress),
	}
}
