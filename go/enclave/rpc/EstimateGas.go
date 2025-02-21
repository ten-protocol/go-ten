package rpc

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/ethereum/go-ethereum/params"

	"github.com/ten-protocol/go-ten/go/common/measure"
	"github.com/ten-protocol/go-ten/go/enclave/core"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethcore "github.com/ethereum/go-ethereum/core"
	"github.com/ten-protocol/go-ten/go/common/gethapi"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

var adjustPublishingGas = gethcommon.Big3

func EstimateGasValidate(reqParams []any, builder *CallBuilder[CallParamsWithBlock, hexutil.Uint64], _ *EncryptionManager) error {
	// Parameters are [callMsg, BlockHeader number (optional)]
	if len(reqParams) < 1 {
		builder.Err = fmt.Errorf("unexpected number of parameters")
		return nil
	}

	callMsg, err := gethencoding.ExtractEthCall(reqParams[0])
	if err != nil {
		builder.Err = fmt.Errorf("unable to decode EthCall Params - %w", err)
		return nil
	}

	if callMsg.From == nil {
		builder.Err = fmt.Errorf("no from Addr provided")
		return nil
	}

	// extract optional BlockHeader number - defaults to the latest BlockHeader if not avail
	blockNumber, err := gethencoding.ExtractOptionalBlockNumber(reqParams, 1)
	if err != nil {
		builder.Err = fmt.Errorf("unable to extract requested BlockHeader number - %w", err)
		return nil
	}

	builder.From = callMsg.From
	// todo
	builder.Param = &CallParamsWithBlock{callMsg, blockNumber.BlockNumber}
	return nil
}

// EstimateGasExecute - performs the gas estimation based on the provided parameters and the local environment configuration.
// Will accommodate l1 gas cost and stretch the final gas estimation.
func EstimateGasExecute(builder *CallBuilder[CallParamsWithBlock, hexutil.Uint64], rpc *EncryptionManager) error {
	err := authenticateFrom(builder.VK, builder.From)
	if err != nil {
		builder.Err = err
		return nil //nolint:nilerr
	}

	txArgs := builder.Param.callParams
	blockNumber := builder.Param.block
	block, err := rpc.l1BlockProcessor.GetHead(builder.ctx)
	if err != nil {
		return err
	}

	headBatchSeq := rpc.registry.HeadBatchSeq()
	batch, err := rpc.storage.FetchBatchHeaderBySeqNo(builder.ctx, headBatchSeq.Uint64())
	if err != nil {
		return err
	}

	// The message is run through the l1 publishing cost estimation for the current
	// known head BlockHeader.
	l1Cost, err := rpc.gasOracle.EstimateL1CostForMsg(txArgs, block, batch)
	if err != nil {
		return err
	}

	// We divide the total estimated l1 cost by the l2 fee per gas in order to convert
	// the expected cost into l2 gas based on current pricing.
	// todo @siliev - add overhead when the base fee becomes dynamic.
	publishingGas := big.NewInt(0).Div(l1Cost, batch.BaseFee)

	// Overestimate the publishing cost in case of spikes.
	// given that we publish in a blob, the amount will be very low.
	// Batch execution still deducts normally.
	// TODO: Change to fixed time period quotes, rather than this.
	publishingGas = publishingGas.Mul(publishingGas, adjustPublishingGas)

	// Run the execution simulation based on stateDB after head batch.
	// Notice that unfortunately, some slots might ve considered warm, which skews the estimation.
	// The single pass will run once at the highest gas cap and return gas used. Not completely reliable,
	// but is quick.
	executionGasEstimate, revert, gasPrice, err := estimateGasSinglePass(builder.ctx, rpc, txArgs, blockNumber, rpc.config.GasLocalExecutionCapFlag)
	if err != nil {
		if len(revert) > 0 {
			builder.Err = newRevertError(revert)
			rpc.logger.Debug("revert error", "err", builder.Err)
			return nil
		}

		// return EVM error
		builder.Err = err
		return nil
	}

	totalGasEstimateUint64 := publishingGas.Uint64() + uint64(executionGasEstimate)
	totalGasEstimate := hexutil.Uint64(totalGasEstimateUint64)
	balance, err := rpc.chain.GetBalanceAtBlock(builder.ctx, *txArgs.From, blockNumber)
	if err != nil {
		return err
	}

	if balance.ToInt().Cmp(big.NewInt(0).Mul(gasPrice, big.NewInt(0).SetUint64(totalGasEstimateUint64))) < 0 {
		return fmt.Errorf("insufficient funds for gas estimate")
	}
	rpc.logger.Debug("Estimation breakdown", "gasPrice", gasPrice, "executionGasEstimate", uint64(executionGasEstimate), "publishingGas", publishingGas, "totalGasEstimate", uint64(totalGasEstimate))
	builder.ReturnValue = &totalGasEstimate
	return nil
}

func calculateMaxGasCap(ctx context.Context, rpc *EncryptionManager, gasCap uint64, argsGas *hexutil.Uint64) uint64 {
	// Fetch the current batch header to get the batch gas limit
	batchHeader, err := rpc.storage.FetchHeadBatchHeader(ctx)
	if err != nil {
		rpc.logger.Error("Failed to fetch batch header", "error", err)
		return gasCap
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
			rpc.logger.Debug("Gas cap adjusted based on args.Gas",
				"argsGas", argsGasUint64,
				"previousGasCap", gasCap,
				"newGasCap", argsGasUint64,
			)
			gasCap = argsGasUint64
		}
	}

	return gasCap
}

// This adds a bit of an overhead to gas estimation. Fixes issues when calling proxies, but needs more investigation.
// Not sure why simulation is non consistent.
func calculateProxyOverhead(txArgs *gethapi.TransactionArgs) uint64 {
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

// estimateGasSinglePass - deduces the simulation params from the call parameters and the local environment configuration.
// will override the gas limit with one provided in transaction if lower. Furthermore figures out the gas cap and the allowance
// for the from address.
// In the binary search approach geth uses, the high of the range for gas limit is where our single pass runs.
// For example, if you estimate gas for a swap, the simulation EVM will be configured to run at the highest possible gas cap.
// This allows the maximum gas for running the call. Then we look at the gas used and return this with a couple modifications.
// The modifications are an overhead buffer and a 20% increase to account for warm storage slots. This is because the stateDB
// for the head batch might not be fully clean in terms of the running call. Cold storage slots cost far more than warm ones to
// read and write.
func estimateGasSinglePass(ctx context.Context, rpc *EncryptionManager, args *gethapi.TransactionArgs, blkNumber *gethrpc.BlockNumber, globalGasCap uint64) (hexutil.Uint64, []byte, *big.Int, error) {
	maxGasCap := calculateMaxGasCap(ctx, rpc, globalGasCap, args.Gas)
	// allowance will either be the maxGasCap or the balance allowance.
	// If the users funds are floaty, this might cause issues combined with the l1 pricing.
	allowance, feeCap, err := normalizeFeeCapAndAdjustGasLimit(ctx, rpc, args, blkNumber, maxGasCap)
	if err != nil {
		return 0, nil, nil, err
	}

	// Perform a single gas estimation pass using isGasEnough
	failed, result, err := isGasEnough(ctx, rpc, args, allowance, blkNumber)
	if err != nil {
		// Return zero values and the encountered error if estimation fails
		return 0, nil, nil, err
	}

	if failed {
		if result != nil && !errors.Is(result.Err, vm.ErrOutOfGas) {
			rpc.logger.Debug("Failed gas estimation", "error", result.Err)
			return 0, result.Revert(), nil, result.Err
		}
		// If the gas cap is insufficient, return an appropriate error
		return 0, nil, nil, fmt.Errorf("gas required exceeds allowance (%d)", globalGasCap)
	}

	if result == nil {
		// If there's no result, something went wrong
		return 0, nil, nil, fmt.Errorf("no execution result returned")
	}

	// Extract the gas used from the execution result.
	// Add an overhead buffer to account for the fact that the execution might not be able to be completed in the same batch.
	// There can be further discrepancies in the execution due to storage and other factors.
	gasUsedBig := big.NewInt(0).SetUint64(result.UsedGas)
	gasUsedBig.Add(gasUsedBig, big.NewInt(0).SetUint64(calculateProxyOverhead(args)))
	// Add 20% overhead to gas used - this is a rough accommodation for
	// warm storage slots.
	gasUsedBig.Mul(gasUsedBig, big.NewInt(120))
	gasUsedBig.Div(gasUsedBig, big.NewInt(100))
	gasUsed := hexutil.Uint64(gasUsedBig.Uint64())

	return gasUsed, nil, feeCap, nil
}

func normalizeFeeCapAndAdjustGasLimit(ctx context.Context, rpc *EncryptionManager, args *gethapi.TransactionArgs, blkNumber *gethrpc.BlockNumber, hi uint64) (uint64, *big.Int, error) {
	// Normalize the max fee per gas the call is willing to spend.
	var feeCap *big.Int
	if args.GasPrice != nil && (args.MaxFeePerGas != nil || args.MaxPriorityFeePerGas != nil) {
		return 0, gethcommon.Big0, errors.New("both gasPrice and (maxFeePerGas or maxPriorityFeePerGas) specified")
	} else if args.GasPrice != nil {
		feeCap = args.GasPrice.ToInt()
	} else if args.MaxFeePerGas != nil {
		feeCap = args.MaxFeePerGas.ToInt()
	} else {
		feeCap = gethcommon.Big0
	}

	// Recap the highest gas limit with account's available balance.
	if feeCap.BitLen() != 0 { //nolint:nestif
		balance, err := rpc.chain.GetBalanceAtBlock(ctx, *args.From, blkNumber)
		if err != nil {
			return 0, gethcommon.Big0, fmt.Errorf("unable to fetch account balance - %w", err)
		}

		available := new(big.Int).Set(balance.ToInt())
		if args.Value != nil {
			if args.Value.ToInt().Cmp(available) >= 0 {
				return 0, gethcommon.Big0, errors.New("insufficient funds for transfer")
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
			rpc.logger.Debug("Gas estimation capped by limited funds",
				"original", hi,
				"balance", balance,
				"sent", transfer.ToInt(),
				"maxFeePerGas", feeCap,
				"fundable", allowance,
			)
			hi = allowance.Uint64()
		}
	}

	return hi, feeCap, nil
}

// Create a helper to check if a gas allowance results in an executable transaction
// isGasEnough returns whether the gaslimit should be raised, lowered, or if it was impossible to execute the message
func isGasEnough(ctx context.Context, rpc *EncryptionManager, args *gethapi.TransactionArgs, gas uint64, blkNumber *gethrpc.BlockNumber) (bool, *gethcore.ExecutionResult, error) {
	defer core.LogMethodDuration(rpc.logger, measure.NewStopwatch(), "enclave.go:IsGasEnough")
	args.Gas = (*hexutil.Uint64)(&gas)
	result, err := rpc.chain.ObsCallAtBlock(ctx, args, blkNumber)
	if err != nil {
		// since we estimate gas in a single pass, any error is just returned
		return true, nil, err // Bail out
	}
	return result.Failed(), result, nil
}
