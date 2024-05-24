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
