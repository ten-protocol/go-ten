package evm

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/enclave/db"
)

// Obscuro implementation of the ChainContext
type ObscuroChainContext struct {
	rollupResolver db.RollupResolver
}

func (*ObscuroChainContext) Engine() consensus.Engine {
	return &NoOpEngine{}
}

func (occ *ObscuroChainContext) GetHeader(hash common.Hash, height uint64) *types.Header {
	rol, f := occ.rollupResolver.FetchRollup(hash)
	if !f {
		return nil
	}
	return convertToEthHeader(rol.Header)
}
