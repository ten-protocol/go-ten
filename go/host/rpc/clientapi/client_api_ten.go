package clientapi

import (
	"context"
	"fmt"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/responses"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/host"
)

// TenAPI implements Ten-specific JSON RPC operations.
type TenAPI struct {
	host   host.Host
	rpcKey []byte
	logger gethlog.Logger
}

func NewTenAPI(host host.Host, logger gethlog.Logger) *TenAPI {
	return &TenAPI{
		host:   host,
		logger: logger,
	}
}

// Version returns the protocol version of the Obscuro network.
func (api *TenAPI) Version() string {
	return fmt.Sprintf("%d", api.host.Config().TenChainID)
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

func (api *TenAPI) RpcKey() ([]byte, error) {
	if api.rpcKey != nil {
		return api.rpcKey, nil
	}
	var err error
	api.rpcKey, err = api.host.EnclaveClient().RPCEncryptionKey(context.Background())
	if err != nil {
		return nil, err
	}
	return api.rpcKey, nil
}

func (api *TenAPI) EncryptedRPC(ctx context.Context, encryptedParams common.EncryptedRPCRequest) (responses.EnclaveResponse, error) {
	var enclaveResponse *responses.EnclaveResponse
	var sysError error
	if encryptedParams.IsTx {
		enclaveResponse, sysError = api.host.SubmitAndBroadcastTx(ctx, encryptedParams.Req)
	} else {
		enclaveResponse, sysError = api.host.EnclaveClient().EncryptedRPC(ctx, encryptedParams.Req)
	}
	if sysError != nil {
		api.logger.Error("Enclave System Error. Function EncryptedRPC", log.ErrKey, sysError)
		return responses.EnclaveResponse{
			Err: &responses.InternalErrMsg,
		}, nil
	}
	return *enclaveResponse, nil
}

// ChecksumFormattedTenNetworkConfig serialises the addresses as EIP55 checksum addresses.
type ChecksumFormattedTenNetworkConfig struct {
	ManagementContractAddress       gethcommon.AddressEIP55
	L1StartHash                     gethcommon.Hash
	MessageBusAddress               gethcommon.AddressEIP55
	L2MessageBusAddress             gethcommon.AddressEIP55
	ImportantContracts              map[string]gethcommon.AddressEIP55 // map of contract name to address
	TransactionPostProcessorAddress gethcommon.AddressEIP55
	PublicSystemContracts           map[string]gethcommon.AddressEIP55
}

func checksumFormatted(info *common.TenNetworkInfo) *ChecksumFormattedTenNetworkConfig {
	importantContracts := make(map[string]gethcommon.AddressEIP55)
	for name, addr := range info.ImportantContracts {
		importantContracts[name] = gethcommon.AddressEIP55(addr)
	}

	publicSystemContracts := make(map[string]gethcommon.AddressEIP55)
	for name, addr := range info.PublicSystemContracts {
		publicSystemContracts[name] = gethcommon.AddressEIP55(addr)
	}

	return &ChecksumFormattedTenNetworkConfig{
		ManagementContractAddress:       gethcommon.AddressEIP55(info.ManagementContractAddress),
		L1StartHash:                     info.L1StartHash,
		MessageBusAddress:               gethcommon.AddressEIP55(info.MessageBusAddress),
		L2MessageBusAddress:             gethcommon.AddressEIP55(info.L2MessageBusAddress),
		ImportantContracts:              importantContracts,
		TransactionPostProcessorAddress: gethcommon.AddressEIP55(info.TransactionPostProcessorAddress),
		PublicSystemContracts:           publicSystemContracts,
	}
}
