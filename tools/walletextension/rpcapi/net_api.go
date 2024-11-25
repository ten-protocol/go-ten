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
	return UnauthenticatedTenRPCCall[string](ctx, api.we, &cache.Cfg{Type: cache.LongLiving}, "net_version")
}

type ConfigResponseJson struct {
	ManagementContractAddress       string
	MessageBusAddress               string
	L2MessageBusAddress             string
	TransactionPostProcessorAddress string
}

func (api *NetAPI) Config(ctx context.Context) (*ConfigResponseJson, error) {
	return UnauthenticatedTenRPCCall[ConfigResponseJson](ctx, api.we, &cache.Cfg{Type: cache.LongLiving}, "ten_config")
}
