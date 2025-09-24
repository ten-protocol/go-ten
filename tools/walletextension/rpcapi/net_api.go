package rpcapi

import (
	"context"

	"github.com/ten-protocol/go-ten/tools/walletextension/cache"

	"github.com/ten-protocol/go-ten/tools/walletextension/services"
)

type NetAPI struct {
	we *services.Services
}

func NewNetAPI(we *services.Services) *NetAPI {
	return &NetAPI{we}
}

func (api *NetAPI) Version(ctx context.Context) (*string, error) {
	return UnauthenticatedTenRPCCall[string](ctx, api.we, &cache.Cfg{Type: cache.LongLiving}, "ten_version")
}

type ConfigResponseJson struct {
	NetworkConfigAddress            string            `json:"NetworkConfig"`
	EnclaveRegistryAddress          string            `json:"EnclaveRegistry"`
	DataAvailabilityRegistryAddress string            `json:"DataAvailabilityRegistry"`
	CrossChainAddress               string            `json:"CrossChain"`
	L1MessageBusAddress             string            `json:"L1MessageBus"`
	L2MessageBusAddress             string            `json:"L2MessageBus"`
	TransactionPostProcessorAddress string            `json:"TransactionsPostProcessor"`
	L1Bridge                        string            `json:"L1Bridge"`
	L2Bridge                        string            `json:"L2Bridge"`
	L1CrossChainMessenger           string            `json:"L1CrossChainMessenger"`
	L2CrossChainMessenger           string            `json:"L2CrossChainMessenger"`
	SystemContractsUpgrader         string            `json:"SystemContractsUpgrader"`
	L1StartHash                     string            `json:"L1StartHash"`
	PublicSystemContracts           map[string]string `json:"PublicSystemContracts"`
	AdditionalContracts             interface{}       `json:"AdditionalContracts"`
}

func (api *NetAPI) Config(ctx context.Context) (*ConfigResponseJson, error) {
	return UnauthenticatedTenRPCCall[ConfigResponseJson](ctx, api.we, &cache.Cfg{Type: cache.LongLiving}, "ten_config")
}
