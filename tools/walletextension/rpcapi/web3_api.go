package rpcapi

import (
	"context"
)

var _hardcodedClientVersion = "Geth/v10.0.0/drpc"

type Web3API struct {
	we *Services
}

func NewWeb3API(we *Services) *Web3API {
	return &Web3API{we}
}

func (api *Web3API) ClientVersion(_ context.Context) (*string, error) {
	// todo: have this return the Ten version from the node
	return &_hardcodedClientVersion, nil
}
