package evm

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common/errutil"
	"github.com/obscuronet/go-obscuro/go/common/gethencoding"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/enclave/db"

	gethlog "github.com/ethereum/go-ethereum/log"
)

// ObscuroChainContext - basic implementation of the ChainContext needed for the EVM integration
type ObscuroChainContext struct {
	storage db.Storage
	logger  gethlog.Logger
}

// NewObscuroChainContext returns a new instance of the ObscuroChainContext given a storage ( and logger )
func NewObscuroChainContext(storage db.Storage, logger gethlog.Logger) *ObscuroChainContext {
	return &ObscuroChainContext{
		storage: storage,
		logger:  logger,
	}
}

func (occ *ObscuroChainContext) Engine() consensus.Engine {
	return &ObscuroNoOpConsensusEngine{logger: occ.logger}
}

func (occ *ObscuroChainContext) GetHeader(hash common.Hash, height uint64) *types.Header {
	batch, err := occ.storage.FetchBatchHeader(hash)
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			return nil
		}
		occ.logger.Crit("Could not retrieve rollup", log.ErrKey, err)
	}

	h, err := gethencoding.ConvertToEthHeader(batch, secret(occ.storage))
	if err != nil {
		occ.logger.Crit("Could not convert to eth header", log.ErrKey, err)
		return nil
	}
	return h
}
