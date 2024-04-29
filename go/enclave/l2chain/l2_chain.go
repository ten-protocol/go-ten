package l2chain

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ten-protocol/go-ten/go/config"

	"github.com/ten-protocol/go-ten/go/enclave/storage"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethcore "github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/status-im/keycard-go/hexutils"
	"github.com/ten-protocol/go-ten/go/common/gethapi"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/enclave/components"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"github.com/ten-protocol/go-ten/go/enclave/evm"
	"github.com/ten-protocol/go-ten/go/enclave/genesis"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

type obscuroChain struct {
	chainConfig         *params.ChainConfig
	config              config.EnclaveConfig
	storage             storage.Storage
	gethEncodingService gethencoding.EncodingService
	genesis             *genesis.Genesis

	logger gethlog.Logger

	Registry         components.BatchRegistry
	gasEstimationCap uint64
}

func NewChain(
	storage storage.Storage,
	config config.EnclaveConfig,
	gethEncodingService gethencoding.EncodingService,
	chainConfig *params.ChainConfig,
	genesis *genesis.Genesis,
	logger gethlog.Logger,
	registry components.BatchRegistry,
	gasEstimationCap uint64,
) ObscuroChain {
	return &obscuroChain{
		storage:             storage,
		config:              config,
		gethEncodingService: gethEncodingService,
		chainConfig:         chainConfig,
		logger:              logger,
		genesis:             genesis,
		Registry:            registry,
		gasEstimationCap:    gasEstimationCap,
	}
}

func (oc *obscuroChain) AccountOwner(ctx context.Context, address gethcommon.Address, blockNumber *gethrpc.BlockNumber) (*gethcommon.Address, error) {
	// check if account is a contract
	isContract, err := oc.isAccountContractAtBlock(ctx, address, blockNumber)
	if err != nil {
		return nil, err
	}
	if !isContract {
		return &address, nil
	}

	// If the address is a contract, find the signer of the deploy transaction
	txHash, err := oc.storage.GetContractCreationTx(ctx, address)
	if err != nil {
		return nil, err
	}
	transaction, _, _, _, err := oc.storage.GetTransaction(ctx, *txHash) //nolint:dogsled
	if err != nil {
		return nil, err
	}
	signer := types.NewLondonSigner(oc.chainConfig.ChainID)

	sender, err := signer.Sender(transaction)
	if err != nil {
		return nil, err
	}
	return &sender, nil
}

func (oc *obscuroChain) GetBalanceAtBlock(ctx context.Context, accountAddr gethcommon.Address, blockNumber *gethrpc.BlockNumber) (*hexutil.Big, error) {
	chainState, err := oc.Registry.GetBatchStateAtHeight(ctx, blockNumber, false)
	if err != nil {
		return nil, fmt.Errorf("unable to get blockchain state - %w", err)
	}

	return (*hexutil.Big)(chainState.GetBalance(accountAddr).ToBig()), nil
}

func (oc *obscuroChain) ObsCall(ctx context.Context, apiArgs *gethapi.TransactionArgs, blockNumber *gethrpc.BlockNumber) (*gethcore.ExecutionResult, error) {
	result, err := oc.ObsCallAtBlock(ctx, apiArgs, blockNumber)
	if err != nil {
		oc.logger.Info(fmt.Sprintf("Obs_Call: failed to execute contract %s.", apiArgs.To), log.CtrErrKey, err.Error())
		return nil, err
	}

	// the execution might have succeeded (err == nil) but the evm contract logic might have failed (result.Failed() == true)
	if result.Failed() {
		oc.logger.Debug(fmt.Sprintf("Obs_Call: Failed to execute contract %s.", apiArgs.To), log.CtrErrKey, result.Err)
		return nil, result.Err
	}

	if oc.logger.Enabled(context.Background(), gethlog.LevelTrace) {
		oc.logger.Trace("Obs_Call successful", "result", hexutils.BytesToHex(result.ReturnData))
	}
	return result, nil
}

func (oc *obscuroChain) ObsCallAtBlock(ctx context.Context, apiArgs *gethapi.TransactionArgs, blockNumber *gethrpc.BlockNumber) (*gethcore.ExecutionResult, error) {
	// fetch the chain state at given batch
	blockState, err := oc.Registry.GetBatchStateAtHeight(ctx, blockNumber, false)
	if err != nil {
		return nil, err
	}

	batch, err := oc.Registry.GetBatchAtHeight(ctx, *blockNumber)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch head state batch. Cause: %w", err)
	}

	callMsg, err := apiArgs.ToMessage(batch.Header.GasLimit-1, batch.Header.BaseFee)
	if err != nil {
		return nil, fmt.Errorf("unable to convert TransactionArgs to Message - %w", err)
	}

	if oc.logger.Enabled(context.Background(), gethlog.LevelTrace) {
		oc.logger.Trace("Obs_Call: Successful result", "result", fmt.Sprintf("contractAddress=%s, from=%s, data=%s, batch=%s, state=%s",
			callMsg.To,
			callMsg.From,
			hexutils.BytesToHex(callMsg.Data),
			batch.Hash(),
			batch.Header.Root.Hex()))
	}

	result, err := evm.ExecuteObsCall(ctx, callMsg, blockState, batch.Header, oc.storage, oc.gethEncodingService, oc.chainConfig, oc.gasEstimationCap, oc.config, oc.logger)
	if err != nil {
		// also return the result as the result can be evaluated on some errors like ErrIntrinsicGas
		return result, err
	}

	return result, nil
}

// GetChainStateAtTransaction Returns the state of the chain at certain block height after executing transactions up to the selected transaction
// TODO make this cacheable
func (oc *obscuroChain) GetChainStateAtTransaction(ctx context.Context, batch *core.Batch, txIndex int, _ uint64) (*gethcore.Message, vm.BlockContext, *state.StateDB, error) {
	// Short circuit if it's genesis batch.
	if batch.NumberU64() == 0 {
		return nil, vm.BlockContext{}, nil, errors.New("no transaction in genesis")
	}
	// Create the parent state database
	parent, err := oc.Registry.GetBatchAtHeight(ctx, gethrpc.BlockNumber(batch.NumberU64()-1))
	if err != nil {
		return nil, vm.BlockContext{}, nil, fmt.Errorf("unable to fetch parent batch - %w", err)
	}
	parentBlockNumber := gethrpc.BlockNumber(parent.NumberU64())

	// Lookup the statedb of parent batch from the live database,
	// otherwise regenerate it on the flight.
	statedb, err := oc.Registry.GetBatchStateAtHeight(ctx, &parentBlockNumber, false)
	if err != nil {
		return nil, vm.BlockContext{}, nil, err
	}
	if txIndex == 0 && len(batch.Transactions) == 0 {
		return nil, vm.BlockContext{}, statedb, nil
	}
	// Recompute transactions up to the target index.
	// TODO - Once the enclave's genesis.json is set, retrieve the signer type using `types.MakeSigner`.
	rules := oc.chainConfig.Rules(big.NewInt(0), true, 0)
	signer := types.LatestSigner(oc.chainConfig)
	for idx, tx := range batch.Transactions {
		// Assemble the transaction call message and return if the requested offset
		msg, err := gethcore.TransactionToMessage(tx, signer, big.NewInt(0))
		if err != nil {
			return nil, vm.BlockContext{}, nil, fmt.Errorf("unable to convert tx to message - %w", err)
		}
		txContext := gethcore.NewEVMTxContext(msg)

		chain := evm.NewObscuroChainContext(oc.storage, oc.gethEncodingService, oc.config, oc.logger)

		blockHeader, err := oc.gethEncodingService.CreateEthHeaderForBatch(ctx, batch.Header)
		if err != nil {
			return nil, vm.BlockContext{}, nil, fmt.Errorf("unable to convert batch header to eth header - %w", err)
		}
		context := gethcore.NewEVMBlockContext(blockHeader, chain, nil)
		if idx == txIndex {
			return msg, context, statedb, nil
		}
		// Not yet the searched for transaction, execute on top of the current state
		vmenv := vm.NewEVM(context, txContext, statedb, oc.chainConfig, vm.Config{})
		statedb.Prepare(rules, msg.From, gethcommon.Address{}, tx.To(), nil, nil)
		if _, err := gethcore.ApplyMessage(vmenv, msg, new(gethcore.GasPool).AddGas(tx.Gas())); err != nil {
			return nil, vm.BlockContext{}, nil, fmt.Errorf("transaction %#x failed: %w", tx.Hash(), err)
		}
		// Ensure any modifications are committed to the state
		// Only delete empty objects if EIP158/161 (a.k.a Spurious Dragon) is in effect
		statedb.Finalise(vmenv.ChainConfig().IsEIP158(batch.Number()))
	}
	return nil, vm.BlockContext{}, nil, fmt.Errorf("transaction index %d out of range for batch %#x", txIndex, batch.Hash())
}

// Returns whether the account is a contract
func (oc *obscuroChain) isAccountContractAtBlock(ctx context.Context, accountAddr gethcommon.Address, blockNumber *gethrpc.BlockNumber) (bool, error) {
	chainState, err := oc.Registry.GetBatchStateAtHeight(ctx, blockNumber, false)
	if err != nil {
		return false, fmt.Errorf("unable to get blockchain state - %w", err)
	}

	return len(chainState.GetCode(accountAddr)) > 0, nil
}
