package rpcapi

import (
	"context"
)

type NetAPI struct {
	we *Services
}

func NewNetAPI(we *Services) *NetAPI {
	return &NetAPI{we}
}

func (api *NetAPI) Version(ctx context.Context) (*string, error) {
	return UnauthenticatedTenRPCCall[string](ctx, api.we, &CacheCfg{CacheType: LongLiving}, "net_version")
}

type ConfigResponseJson struct {
	ManagementContractAddress       string
	MessageBusAddress               string
	L2MessageBusAddress             string
	TransactionPostProcessorAddress string
}

func (api *NetAPI) Config(ctx context.Context) (*ConfigResponseJson, error) {
	return UnauthenticatedTenRPCCall[ConfigResponseJson](ctx, api.we, &CacheCfg{CacheType: LongLiving}, "obscuro_config")
}
