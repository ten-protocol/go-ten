package rpcapi

import (
	"context"

	"github.com/ten-protocol/go-ten/tools/walletextension/services"
)

var _hardcodedClientVersion = "Geth/v10.0.0/ten"

type Web3API struct {
	we *services.Services
}

func NewWeb3API(we *services.Services) *Web3API {
	return &Web3API{we}
}

func (api *Web3API) ClientVersion(_ context.Context) (*string, error) {
	// todo: have this return the TEN version from the node
	return &_hardcodedClientVersion, nil
}
