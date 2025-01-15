package rpcapi

import (
	"context"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/cache"
	"github.com/ten-protocol/go-ten/tools/walletextension/services"
)

type TenAPI struct {
	we *services.Services
}

func NewTenAPI(we *services.Services) *TenAPI {
	return &TenAPI{we}
}

type CrossChainProof struct {
	Proof []byte
	Root  gethcommon.Hash
}

func (api *TenAPI) GetCrossChainProof(ctx context.Context, messageType string, crossChainMessage gethcommon.Hash) (*CrossChainProof, error) {
	proof, err := UnauthenticatedTenRPCCall[CrossChainProof](ctx, api.we, &cache.Cfg{Type: cache.LatestBatch}, "ten_getCrossChainProof", messageType, crossChainMessage)
	if err != nil {
		return nil, err
	}
	return proof, nil
}
