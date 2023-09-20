package l2chain

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/obscuronet/go-obscuro/go/enclave/storage"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethcore "github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/obscuronet/go-obscuro/go/common/gethapi"
	"github.com/obscuronet/go-obscuro/go/common/gethencoding"
	"github.com/obscuronet/go-obscuro/go/common/log"
	"github.com/obscuronet/go-obscuro/go/enclave/components"
	"github.com/obscuronet/go-obscuro/go/enclave/core"
	"github.com/obscuronet/go-obscuro/go/enclave/evm"
	"github.com/obscuronet/go-obscuro/go/enclave/genesis"
	"github.com/status-im/keycard-go/hexutils"
)

type obscuroChain struct {
	chainConfig *params.ChainConfig

	storage storage.Storage
	genesis *genesis.Genesis

	logger gethlog.Logger

	Registry components.BatchRegistry
}

func NewChain(
	storage storage.Storage,
	chainConfig *params.ChainConfig,
	genesis *genesis.Genesis,
	logger gethlog.Logger,
	registry components.BatchRegistry,
) ObscuroChain {
	return &obscuroChain{
		storage:     storage,
		chainConfig: chainConfig,
		logger:      logger,
		genesis:     genesis,
		Registry:    registry,
	}
}

func (oc *obscuroChain) GetBalance(accountAddress gethcommon.Address, blockNumber *gethrpc.BlockNumber) (*gethcommon.Address, *hexutil.Big, error) {
	// get account balance at certain block/height
	balance, err := oc.GetBalanceAtBlock(accountAddress, blockNumber)
	if err != nil {
		return nil, nil, err
	}

	// check if account is a contract
	isAddrContract, err := oc.isAccountContractAtBlock(accountAddress, blockNumber)
	if err != nil {
		return nil, nil, err
	}

	// Decide which address to encrypt the result with
	address := accountAddress
	// If the accountAddress is a contract, encrypt with the address of the contract owner
	if isAddrContract {
		txHash, err := oc.storage.GetContractCreationTx(accountAddress)
		if err != nil {
			return nil, nil, err
		}
		transaction, _, _, _, err := oc.storage.GetTransaction(*txHash)
		if err != nil {
			return nil, nil, err
		}
		signer := types.NewLondonSigner(oc.chainConfig.ChainID)

		sender, err := signer.Sender(transaction)
		if err != nil {
			return nil, nil, err
		}
		address = sender
	}

	return &address, balance, nil
}

func (oc *obscuroChain) GetBalanceAtBlock(accountAddr gethcommon.Address, blockNumber *gethrpc.BlockNumber) (*hexutil.Big, error) {
	chainState, err := oc.Registry.GetBatchStateAtHeight(blockNumber)
	if err != nil {
		return nil, fmt.Errorf("unable to get blockchain state - %w", err)
	}

	return (*hexutil.Big)(chainState.GetBalance(accountAddr)), nil
}

func (oc *obscuroChain) ObsCall(apiArgs *gethapi.TransactionArgs, blockNumber *gethrpc.BlockNumber) (*gethcore.ExecutionResult, error) {
	result, err := oc.ObsCallAtBlock(apiArgs, blockNumber)
	if err != nil {
		oc.logger.Info(fmt.Sprintf("Obs_Call: failed to execute contract %s.", apiArgs.To), log.CtrErrKey, err.Error())
		return nil, err
	}

	// the execution might have succeeded (err == nil) but the evm contract logic might have failed (result.Failed() == true)
	if result.Failed() {
		oc.logger.Debug(fmt.Sprintf("Obs_Call: Failed to execute contract %s.", apiArgs.To), log.CtrErrKey, result.Err)
		return nil, result.Err
	}

	oc.logger.Trace("Obs_Call successful", "result", gethlog.Lazy{Fn: func() string {
		return hexutils.BytesToHex(result.ReturnData)
	}})
	return result, nil
}

func (oc *obscuroChain) ObsCallAtBlock(apiArgs *gethapi.TransactionArgs, blockNumber *gethrpc.BlockNumber) (*gethcore.ExecutionResult, error) {
	// fetch the chain state at given batch
	blockState, err := oc.Registry.GetBatchStateAtHeight(blockNumber)
	if err != nil {
		return nil, err
	}

	batch, err := oc.Registry.GetBatchAtHeight(*blockNumber)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch head state batch. Cause: %w", err)
	}

	callMsg, err := apiArgs.ToMessage(batch.Header.GasLimit, batch.Header.BaseFee)
	if err != nil {
		return nil, fmt.Errorf("unable to convert TransactionArgs to Message - %w", err)
	}

	oc.logger.Trace("Obs_Call: Successful result", "result", gethlog.Lazy{Fn: func() string {
		return fmt.Sprintf("contractAddress=%s, from=%s, data=%s, batch=%s, state=%s",
			callMsg.To,
			callMsg.From,
			hexutils.BytesToHex(callMsg.Data),
			batch.Hash(),
			batch.Header.Root.Hex())
	}})

	result, err := evm.ExecuteObsCall(callMsg, blockState, batch.Header, oc.storage, oc.chainConfig, oc.logger)
	if err != nil {
		// also return the result as the result can be evaluated on some errors like ErrIntrinsicGas
		oc.logger.Info("Call failed with error", log.ErrKey, err)
		return result, err
	}

	return result, nil
}

// GetChainStateAtTransaction Returns the state of the chain at certain block height after executing transactions up to the selected transaction
// TODO make this cacheable
func (oc *obscuroChain) GetChainStateAtTransaction(batch *core.Batch, txIndex int, _ uint64) (*gethcore.Message, vm.BlockContext, *state.StateDB, error) {
	// Short circuit if it's genesis batch.
	if batch.NumberU64() == 0 {
		return nil, vm.BlockContext{}, nil, errors.New("no transaction in genesis")
	}
	// Create the parent state database
	parent, err := oc.Registry.GetBatchAtHeight(gethrpc.BlockNumber(batch.NumberU64() - 1))
	if err != nil {
		return nil, vm.BlockContext{}, nil, fmt.Errorf("unable to fetch parent batch - %w", err)
	}
	parentBlockNumber := gethrpc.BlockNumber(parent.NumberU64())

	// Lookup the statedb of parent batch from the live database,
	// otherwise regenerate it on the flight.
	statedb, err := oc.Registry.GetBatchStateAtHeight(&parentBlockNumber)
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

		chain := evm.NewObscuroChainContext(oc.storage, oc.logger)
		blockHeader, err := gethencoding.CreateEthHeaderForBatch(batch.Header, nil)
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

// Returns the whether the account is a contract or not at a certain height
func (oc *obscuroChain) isAccountContractAtBlock(accountAddr gethcommon.Address, blockNumber *gethrpc.BlockNumber) (bool, error) {
	chainState, err := oc.Registry.GetBatchStateAtHeight(blockNumber)
	if err != nil {
		return false, fmt.Errorf("unable to get blockchain state - %w", err)
	}

	return len(chainState.GetCode(accountAddr)) > 0, nil
}
