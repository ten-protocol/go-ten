package clientapi

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
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

type CrossChainProof struct {
	Proof hexutil.Bytes
	Root  gethcommon.Hash
}

func (api *TenAPI) GetCrossChainProof(_ context.Context, messageType string, crossChainMessage gethcommon.Hash) (CrossChainProof, error) {
	proof, root, err := api.host.Storage().FetchCrossChainProof(messageType, crossChainMessage)
	if err != nil {
		return CrossChainProof{}, err
	}
	encodedProof, err := rlp.EncodeToBytes(proof)
	if err != nil {
		return CrossChainProof{}, err
	}
	return CrossChainProof{
		Proof: encodedProof,
		Root:  root,
	}, nil
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
	NetworkConfigAddress      gethcommon.AddressEIP55
	EnclaveRegistry           gethcommon.AddressEIP55
	CrossChain                gethcommon.AddressEIP55
	RollupContract            gethcommon.AddressEIP55
	L1MessageBus              gethcommon.AddressEIP55
	L2MessageBus              gethcommon.AddressEIP55
	L1Bridge                  gethcommon.AddressEIP55
	L2Bridge                  gethcommon.AddressEIP55
	L1CrossChainMessenger     gethcommon.AddressEIP55
	L2CrossChainMessenger     gethcommon.AddressEIP55
	TransactionsPostProcessor gethcommon.AddressEIP55
	L1StartHash               gethcommon.Hash
	PublicSystemContracts     map[string]gethcommon.AddressEIP55
	AdditionalContracts       []*common.NamedAddress
}

func checksumFormatted(info *common.TenNetworkInfo) *ChecksumFormattedTenNetworkConfig {
	additionalContracts := info.AdditionalContracts

	publicSystemContracts := make(map[string]gethcommon.AddressEIP55)
	for name, addr := range info.PublicSystemContracts {
		publicSystemContracts[name] = gethcommon.AddressEIP55(addr)
	}

	return &ChecksumFormattedTenNetworkConfig{
		NetworkConfigAddress:      gethcommon.AddressEIP55(info.NetworkConfigAddress),
		EnclaveRegistry:           gethcommon.AddressEIP55(info.EnclaveRegistry),
		CrossChain:                gethcommon.AddressEIP55(info.CrossChain),
		RollupContract:            gethcommon.AddressEIP55(info.RollupContract),
		L1MessageBus:              gethcommon.AddressEIP55(info.L1MessageBus),
		L2MessageBus:              gethcommon.AddressEIP55(info.L2MessageBus),
		L1Bridge:                  gethcommon.AddressEIP55(info.L1Bridge),
		L2Bridge:                  gethcommon.AddressEIP55(info.L2Bridge),
		L1CrossChainMessenger:     gethcommon.AddressEIP55(info.L1CrossChainMessenger),
		L2CrossChainMessenger:     gethcommon.AddressEIP55(info.L2CrossChainMessenger),
		TransactionsPostProcessor: gethcommon.AddressEIP55(info.TransactionsPostProcessor),
		L1StartHash:               info.L1StartHash,
		PublicSystemContracts:     publicSystemContracts,
		AdditionalContracts:       additionalContracts,
	}
}
