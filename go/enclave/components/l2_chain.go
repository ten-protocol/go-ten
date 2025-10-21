package components

import (
	"context"
	"fmt"
	"sync"

	"github.com/ten-protocol/go-ten/go/common"

	enclaveconfig "github.com/ten-protocol/go-ten/go/enclave/config"
	"github.com/ten-protocol/go-ten/go/enclave/storage"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethcore "github.com/ethereum/go-ethereum/core"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/status-im/keycard-go/hexutils"
	"github.com/ten-protocol/go-ten/go/common/gethapi"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
	"github.com/ten-protocol/go-ten/go/enclave/evm"
	"github.com/ten-protocol/go-ten/go/enclave/genesis"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

type tenChain struct {
	chainConfig         *params.ChainConfig
	config              *enclaveconfig.EnclaveConfig
	evmFacade           evm.EVMFacade
	storage             storage.Storage
	gethEncodingService gethencoding.EncodingService
	genesis             *genesis.Genesis

	logger gethlog.Logger

	Registry     BatchRegistry
	statedbMutex *sync.Mutex
}

func NewChain(storage storage.Storage, config *enclaveconfig.EnclaveConfig, evmFacade evm.EVMFacade, gethEncodingService gethencoding.EncodingService, chainConfig *params.ChainConfig, genesis *genesis.Genesis, logger gethlog.Logger, registry BatchRegistry, mu *sync.Mutex) TENChain {
	return &tenChain{
		storage:             storage,
		evmFacade:           evmFacade,
		config:              config,
		gethEncodingService: gethEncodingService,
		chainConfig:         chainConfig,
		logger:              logger,
		genesis:             genesis,
		Registry:            registry,
		statedbMutex:        mu,
	}
}

func (oc *tenChain) GetBalanceAtBlock(ctx context.Context, accountAddr gethcommon.Address, blockNumber *gethrpc.BlockNumber) (*hexutil.Big, error) {
	chainState, err := oc.Registry.GetBatchStateAtHeight(ctx, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("unable to get blockchain state - %w", err)
	}

	return (*hexutil.Big)(chainState.GetBalance(accountAddr).ToBig()), nil
}

func (oc *tenChain) ObsCallAtBlock(ctx context.Context, apiArgs *gethapi.TransactionArgs, blockNumber *gethrpc.BlockNumber, isEstimateGas bool) (*gethcore.ExecutionResult, error, common.SystemError) {
	// fetch the chain state at given batch
	blockState, err := oc.Registry.GetBatchStateAtHeight(ctx, blockNumber)
	if err != nil {
		return nil, nil, err
	}

	batch, err := oc.Registry.GetBatchAtHeight(ctx, *blockNumber)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to fetch head state batch. Cause: %w", err)
	}

	callMsg, err := apiArgs.ToMessage(params.MaxTxGas, batch.Header.BaseFee)
	if err != nil {
		return nil, fmt.Errorf("unable to convert TransactionArgs to Message - %w", err), nil
	}

	if oc.logger.Enabled(context.Background(), gethlog.LevelTrace) {
		oc.logger.Trace("Obs_Call: Successful result", "result", fmt.Sprintf("to=%s, from=%s, data=%s, batch=%s, state=%s",
			callMsg.To,
			callMsg.From,
			hexutils.BytesToHex(callMsg.Data),
			batch.Hash(),
			batch.Header.Root.Hex()))
	}

	if *blockNumber == gethrpc.LatestBlockNumber {
		oc.statedbMutex.Lock()
		defer oc.statedbMutex.Unlock()
	}

	return oc.evmFacade.ExecuteCall(ctx, callMsg, blockState, batch.Header, isEstimateGas)
}
