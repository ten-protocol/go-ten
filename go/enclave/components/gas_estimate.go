package components

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethcore "github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/gethapi"
	"github.com/ten-protocol/go-ten/go/common/measure"
	"github.com/ten-protocol/go-ten/go/enclave/core"
	"github.com/ten-protocol/go-ten/go/enclave/gas"
	"github.com/ten-protocol/go-ten/go/enclave/storage"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

var AdjustPublishingGas = gethcommon.Big2

type GasEstimator struct {
	storage   storage.Storage
	chain     TENChain
	logger    gethlog.Logger
	gasOracle gas.Oracle
	gasPricer *GasPricer
}

func NewGasEstimator(storage storage.Storage, chain TENChain, gasOracle gas.Oracle, gasPricer *GasPricer, logger gethlog.Logger) *GasEstimator {
	if gasPricer == nil {
		logger.Crit("gasPricer cannot be nil - this indicates a critical initialization failure")
	}
	return &GasEstimator{
		storage:   storage,
		chain:     chain,
		gasOracle: gasOracle,
		gasPricer: gasPricer,
		logger:    logger,
	}
}

func (ge *GasEstimator) EstimateTotalGas(ctx context.Context, args *gethapi.TransactionArgs, blockNumber *gethrpc.BlockNumber, batch *common.BatchHeader, globalGasCap uint64) (uint64, uint64, error, common.SystemError) {
	// The message is run through the l1 publishing cost estimation for the current
	// known head BlockHeader.
	l1Cost, err := ge.gasOracle.EstimateL1CostForMsg(ctx, args, batch)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("failed to estimate L1 cost: %w", err)
	}

	// We divide the total estimated l1 cost by the l2 fee per gas in order to convert
	// the expected cost into l2 gas based on current pricing.
	// todo @siliev - add overhead when the base fee becomes dynamic.
	divisor := ge.gasPricer.StaticL2BaseFee(common.ConvertBatchHeaderToHeader(batch))
	publishingGas := big.NewInt(0).Div(l1Cost, divisor)

	// Overestimate the publishing cost in case of spikes.
	// given that we publish in a blob, the amount will be very low.
	// Batch execution still deducts normally.
	// TODO: Change to fixed time period quotes, rather than this.
	publishingGas = publishingGas.Mul(publishingGas, AdjustPublishingGas)

	// Run the execution simulation based on stateDB after head batch.
	// Notice that unfortunately, some slots might ve considered warm, which skews the estimation.
	// The single pass will run once at the highest gas cap and return gas used. Not completely reliable,
	// but is quick.
	executionGasEstimate, revert, gasPrice, userErr, sysErr := ge.EstimateGasSinglePass(ctx, args, blockNumber, globalGasCap)
	if sysErr != nil {
		return 0, 0, nil, fmt.Errorf("system error during gas estimation: %w", sysErr)
	}

	if userErr != nil {
		return 0, 0, userErr, nil
	}

	if len(revert) > 0 {
		return 0, 0, newRevertError(revert), nil
	}

	totalGasEstimateUint64 := publishingGas.Uint64() + uint64(executionGasEstimate)
	balance, err := ge.chain.GetBalanceAtBlock(ctx, *args.From, blockNumber)
	if err != nil {
		return 0, 0, nil, fmt.Errorf("failed to get account balance: %w", err)
	}

	if balance.ToInt().Cmp(big.NewInt(0).Mul(gasPrice, big.NewInt(0).SetUint64(totalGasEstimateUint64))) < 0 {
		return 0, 0, errors.New("insufficient balance for transaction"), nil
	}

	return totalGasEstimateUint64, publishingGas.Uint64(), nil, nil
}

// EstimateGasSinglePass - deduces the simulation params from the call parameters and the local environment configuration.
// will override the gas limit with one provided in transaction if lower. Furthermore figures out the gas cap and the allowance
// for the from address.
// In the binary search approach geth uses, the high of the range for gas limit is where our single pass runs.
// For example, if you estimate gas for a swap, the simulation EVM will be configured to run at the highest possible gas cap.
// This allows the maximum gas for running the call. Then we look at the gas used and return this with a couple modifications.
// The modifications are an overhead buffer and a 20% increase to account for warm storage slots. This is because the stateDB
// for the head batch might not be fully clean in terms of the running call. Cold storage slots cost far more than warm ones to
// read and write.
func (ge *GasEstimator) EstimateGasSinglePass(ctx context.Context, args *gethapi.TransactionArgs, blkNumber *gethrpc.BlockNumber, globalGasCap uint64) (hexutil.Uint64, []byte, *big.Int, error, common.SystemError) {
	maxGasCap, err := ge.calculateMaxGasCap(ctx, globalGasCap, args.Gas)
	if err != nil {
		return 0, nil, nil, nil, err
	}

	// allowance will either be the maxGasCap or the balance allowance.
	// If the users funds are floaty, this might cause issues combined with the l1 pricing.
	allowance, feeCap, userErr, sysErr := ge.normalizeFeeCapAndAdjustGasLimit(ctx, args, blkNumber, maxGasCap)
	if sysErr != nil {
		return 0, nil, nil, nil, sysErr
	}

	if userErr != nil {
		return 0, nil, nil, userErr, nil
	}

	// Perform a single gas estimation pass using isNotEnoughGas
	failed, result, userErr, sysErr := ge.isNotEnoughGas(ctx, args, allowance, blkNumber)
	if sysErr != nil {
		// Return zero values and the encountered error if estimation fails
		return 0, nil, nil, nil, sysErr
	}

	if failed {
		if result != nil && !errors.Is(result.Err, vm.ErrOutOfGas) {
			ge.logger.Debug("Failed gas estimation", "error", result.Err)
			return 0, result.Revert(), nil, result.Err, nil
		}
		if userErr != nil {
			return 0, nil, nil, userErr, nil
		}
		// If the gas cap is insufficient, return an appropriate error
		return 0, nil, nil, fmt.Errorf("gas required exceeds allowance (%d)", globalGasCap), nil
	}

	if result == nil {
		// If there's no result, something went wrong
		return 0, nil, nil, nil, fmt.Errorf("no execution result returned")
	}

	// Extract the gas used from the execution result.
	// Add an overhead buffer to account for the fact that the execution might not be able to be completed in the same batch.
	// There can be further discrepancies in the execution due to storage and other factors.
	gasUsedBig := big.NewInt(0).SetUint64(result.UsedGas)
	gasUsedBig.Add(gasUsedBig, big.NewInt(0).SetUint64(ge.calculateProxyOverhead(args)))
	// Add 33% overhead to the single pass estimation
	gasUsedBig.Mul(gasUsedBig, big.NewInt(133))
	gasUsedBig.Div(gasUsedBig, big.NewInt(100))
	gasUsed := hexutil.Uint64(gasUsedBig.Uint64())

	return gasUsed, nil, feeCap, nil, nil
}

func (ge *GasEstimator) calculateMaxGasCap(ctx context.Context, gasCap uint64, argsGas *hexutil.Uint64) (uint64, error) {
	// Fetch the current batch header to get the batch gas limit
	batchHeader, err := ge.storage.FetchHeadBatchHeader(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch batch header: %w", err)
	}

	// Determine the gas limit based on the batch header
	batchGasLimit := batchHeader.GasLimit
	if batchGasLimit < gasCap {
		gasCap = batchGasLimit
	}

	// If args.Gas is specified, take the minimum of gasCap and args.Gas
	if argsGas != nil {
		argsGasUint64 := uint64(*argsGas)
		if argsGasUint64 < gasCap && argsGasUint64 >= params.TxGas {
			ge.logger.Debug("Gas cap adjusted based on args.Gas",
				"argsGas", argsGasUint64,
				"previousGasCap", gasCap,
				"newGasCap", argsGasUint64,
			)
			gasCap = argsGasUint64
		}
	}

	return gasCap, nil
}

// This adds a bit of an overhead to gas estimation. Fixes issues when calling proxies, but needs more investigation.
// Not sure why simulation is non consistent.
func (ge *GasEstimator) calculateProxyOverhead(txArgs *gethapi.TransactionArgs) uint64 {
	if txArgs == nil || txArgs.Data == nil {
		return 0
	}

	calldata := []byte(*txArgs.Data)

	// Base costs
	overhead := uint64(2200) // SLOAD (cold) + DELEGATECALL

	// Memory operations
	dataSize := uint64(len(calldata))
	memCost := (dataSize * 3) * 2 // calldatacopy in both contexts

	// Memory expansion
	words := (dataSize + 31) / 32
	memCost += words * 3

	return overhead + memCost
}

func (ge *GasEstimator) normalizeFeeCapAndAdjustGasLimit(ctx context.Context, args *gethapi.TransactionArgs, blkNumber *gethrpc.BlockNumber, hi uint64) (uint64, *big.Int, error, common.SystemError) {
	// Normalize the max fee per gas the call is willing to spend.
	var feeCap *big.Int
	if args.GasPrice != nil && (args.MaxFeePerGas != nil || args.MaxPriorityFeePerGas != nil) {
		return 0, gethcommon.Big0, errors.New("both gasPrice and (maxFeePerGas or maxPriorityFeePerGas) specified"), nil
	} else if args.GasPrice != nil {
		feeCap = args.GasPrice.ToInt()
	} else if args.MaxFeePerGas != nil {
		feeCap = args.MaxFeePerGas.ToInt()
	} else {
		feeCap = gethcommon.Big0
	}

	// Recap the highest gas limit with account's available balance.
	if feeCap.BitLen() != 0 { //nolint:nestif
		balance, err := ge.chain.GetBalanceAtBlock(ctx, *args.From, blkNumber)
		if err != nil {
			return 0, gethcommon.Big0, nil, fmt.Errorf("unable to fetch account balance: %w", err)
		}

		available := new(big.Int).Set(balance.ToInt())
		if args.Value != nil {
			if args.Value.ToInt().Cmp(available) >= 0 {
				return 0, gethcommon.Big0, errors.New("insufficient funds for transfer"), nil
			}
			available.Sub(available, args.Value.ToInt())
		}
		allowance := new(big.Int).Div(available, feeCap)

		// If the allowance is larger than maximum uint64, skip checking
		if allowance.IsUint64() && hi > allowance.Uint64() {
			transfer := args.Value
			if transfer == nil {
				transfer = new(hexutil.Big)
			}
			ge.logger.Debug("Gas estimation capped by limited funds",
				"original", hi,
				"balance", balance,
				"sent", transfer.ToInt(),
				"maxFeePerGas", feeCap,
				"fundable", allowance,
			)
			hi = allowance.Uint64()
		}
	}
	return hi, feeCap, nil, nil
}

// Create a helper to check if a gas allowance results in an executable transaction
// isNotEnoughGas returns whether the gaslimit should be raised, lowered, or if it was impossible to execute the message
func (ge *GasEstimator) isNotEnoughGas(ctx context.Context, args *gethapi.TransactionArgs, gas uint64, blkNumber *gethrpc.BlockNumber) (bool, *gethcore.ExecutionResult, error, common.SystemError) {
	defer core.LogMethodDuration(ge.logger, measure.NewStopwatch(), "enclave.go:IsGasEnough")
	args.Gas = (*hexutil.Uint64)(&gas)
	result, userErr, sysErr := ge.chain.ObsCallAtBlock(ctx, args, blkNumber)
	if sysErr != nil {
		return true, nil, nil, sysErr
	}

	if userErr != nil {
		// todo @siliev - do we need these?
		//if errors.Is(userErr, gethcore.ErrIntrinsicGas) {
		//	return true, nil, nil, nil // Special case, raise gas limit
		//}
		//if errors.Is(userErr, gethcore.ErrGasLimitTooHigh) {
		//	return true, nil, nil, nil // Special case, lower gas limit
		//}

		return true, nil, userErr, nil
	}

	return result.Failed(), result, nil, nil
}

// newRevertError creates a revertError instance with the provided revert data.
func newRevertError(revert []byte) *errutil.DataError {
	err := vm.ErrExecutionReverted

	reason, errUnpack := abi.UnpackRevert(revert)
	if errUnpack == nil {
		err = fmt.Errorf("%w: %v", vm.ErrExecutionReverted, reason)
	}
	return &errutil.DataError{
		Err:    err.Error(),
		Code:   3, // See: https://github.com/ethereum/wiki/wiki/JSON-RPC-Error-Codes-Improvement-Proposal
		Reason: hexutil.Encode(revert),
	}
}
