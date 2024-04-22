package clientapi

import (
	"context"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/host"
)

// ObscuroAPI implements Obscuro-specific JSON RPC operations.
type ObscuroAPI struct {
	host host.Host
}

func NewObscuroAPI(host host.Host) *ObscuroAPI {
	return &ObscuroAPI{
		host: host,
	}
}

// Health returns the health status of obscuro host + enclave + db
func (api *ObscuroAPI) Health(ctx context.Context) (*host.HealthCheck, error) {
	return api.host.HealthCheck(ctx)
}

// Config returns the config status of obscuro host + enclave + db
func (api *ObscuroAPI) Config() (*ChecksumFormattedObscuroNetworkConfig, error) {
	config, err := api.host.ObscuroConfig()
	if err != nil {
		return nil, err
	}
	return checksumFormatted(config), nil
}

// ChecksumFormattedObscuroNetworkConfig serialises the addresses as EIP55 checksum addresses.
type ChecksumFormattedObscuroNetworkConfig struct {
	ManagementContractAddress gethcommon.AddressEIP55
	L1StartHash               gethcommon.Hash
	SequencerID               gethcommon.AddressEIP55
	MessageBusAddress         gethcommon.AddressEIP55
	L2MessageBusAddress       gethcommon.AddressEIP55
	ImportantContracts        map[string]gethcommon.AddressEIP55 // map of contract name to address
}

func checksumFormatted(info *common.ObscuroNetworkInfo) *ChecksumFormattedObscuroNetworkConfig {
	importantContracts := make(map[string]gethcommon.AddressEIP55)
	for name, addr := range info.ImportantContracts {
		importantContracts[name] = gethcommon.AddressEIP55(addr)
	}
	return &ChecksumFormattedObscuroNetworkConfig{
		ManagementContractAddress: gethcommon.AddressEIP55(info.ManagementContractAddress),
		L1StartHash:               info.L1StartHash,
		SequencerID:               gethcommon.AddressEIP55(info.SequencerID),
		MessageBusAddress:         gethcommon.AddressEIP55(info.MessageBusAddress),
		L2MessageBusAddress:       gethcommon.AddressEIP55(info.L2MessageBusAddress),
		ImportantContracts:        importantContracts,
	}
}
