package evm

import (
	"errors"

	"github.com/ten-protocol/go-ten/go/enclave/storage"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
	"github.com/ten-protocol/go-ten/go/common/log"
)

// ObscuroChainContext - basic implementation of the ChainContext needed for the EVM integration
type ObscuroChainContext struct {
	storage             storage.Storage
	gethEncodingService gethencoding.EncodingService
	logger              gethlog.Logger
}

// NewObscuroChainContext returns a new instance of the ObscuroChainContext given a storage ( and logger )
func NewObscuroChainContext(storage storage.Storage, gethEncodingService gethencoding.EncodingService, logger gethlog.Logger) *ObscuroChainContext {
	return &ObscuroChainContext{
		storage:             storage,
		gethEncodingService: gethEncodingService,
		logger:              logger,
	}
}

func (occ *ObscuroChainContext) Engine() consensus.Engine {
	return &ObscuroNoOpConsensusEngine{logger: occ.logger}
}

func (occ *ObscuroChainContext) GetHeader(hash common.Hash, _ uint64) *types.Header {
	batch, err := occ.storage.FetchBatch(hash)
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			return nil
		}
		occ.logger.Crit("Could not retrieve rollup", log.ErrKey, err)
	}

	h, err := occ.gethEncodingService.CreateEthHeaderForBatch(batch.Header)
	if err != nil {
		occ.logger.Crit("Could not convert to eth header", log.ErrKey, err)
		return nil
	}
	return h
}
