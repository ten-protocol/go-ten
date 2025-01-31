package evm

import (
	"context"
	"errors"

	"github.com/ten-protocol/go-ten/go/enclave/config"
	"github.com/ten-protocol/go-ten/go/enclave/storage"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
	"github.com/ten-protocol/go-ten/go/common/log"
)

// TenChainContext - basic implementation of the ChainContext needed for the EVM integration
type TenChainContext struct {
	storage             storage.Storage
	config              config.EnclaveConfig
	gethEncodingService gethencoding.EncodingService
	logger              gethlog.Logger
}

// NewTenChainContext returns a new instance of the TenChainContext given a storage ( and logger )
func NewTenChainContext(storage storage.Storage, gethEncodingService gethencoding.EncodingService, config config.EnclaveConfig, logger gethlog.Logger) *TenChainContext {
	return &TenChainContext{
		storage:             storage,
		config:              config,
		gethEncodingService: gethEncodingService,
		logger:              logger,
	}
}

func (occ *TenChainContext) Engine() consensus.Engine {
	return &NoOpConsensusEngine{logger: occ.logger}
}

func (occ *TenChainContext) GetHeader(hash common.Hash, _ uint64) *types.Header {
	ctx, cancelCtx := context.WithTimeout(context.Background(), occ.config.RPCTimeout)
	defer cancelCtx()

	batch, err := occ.storage.FetchBatchHeader(ctx, hash)
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			return nil
		}
		occ.logger.Crit("Could not retrieve rollup", log.ErrKey, err)
	}

	h, err := occ.gethEncodingService.CreateEthHeaderForBatch(ctx, batch)
	if err != nil {
		occ.logger.Crit("Could not convert to eth header", log.ErrKey, err)
		return nil
	}
	return h
}
