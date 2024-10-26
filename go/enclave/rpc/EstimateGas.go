package rpc

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/measure"
	"github.com/ten-protocol/go-ten/go/enclave/core"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethcore "github.com/ethereum/go-ethereum/core"
	"github.com/ten-protocol/go-ten/go/common/gethapi"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
	"github.com/ten-protocol/go-ten/go/common/syserr"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

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

	// The message is run through the l1 publishing cost estimation for the current
	// known head BlockHeader.
	l1Cost, err := rpc.gasOracle.EstimateL1CostForMsg(txArgs, block)
	if err != nil {
		return err
	}

	headBatchSeq := rpc.registry.HeadBatchSeq()
	batch, err := rpc.storage.FetchBatchHeaderBySeqNo(builder.ctx, headBatchSeq.Uint64())
	if err != nil {
		return err
	}

	// We divide the total estimated l1 cost by the l2 fee per gas in order to convert
	// the expected cost into l2 gas based on current pricing.
	// todo @siliev - add overhead when the base fee becomes dynamic.
	publishingGas := big.NewInt(0).Div(l1Cost, batch.BaseFee)

	// The one additional gas captures the modulo leftover in some edge cases
	// where BaseFee is bigger than the l1cost.
	publishingGas = big.NewInt(0).Add(publishingGas, gethcommon.Big1)

	// Overestimate the publishing cost in case of spikes.
	// Batch execution still deducts normally.
	// TODO: Change to fixed time period quotes, rather than this.
	publishingGas = publishingGas.Mul(publishingGas, gethcommon.Big2)

	executionGasEstimate, gasPrice, err := rpc.estimateGasSinglePass(builder.ctx, txArgs, blockNumber, rpc.config.GasLocalExecutionCapFlag)
	if err != nil {
		err = fmt.Errorf("unable to estimate transaction - %w", err)

		if errors.Is(err, syserr.InternalError{}) {
			return err
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
	builder.ReturnValue = &totalGasEstimate
	return nil
}

func (rpc *EncryptionManager) calculateMaxGasCap(ctx context.Context, gasCap uint64, argsGas *hexutil.Uint64) uint64 {
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
		if argsGasUint64 < gasCap {
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
func (rpc *EncryptionManager) estimateGasSinglePass(ctx context.Context, args *gethapi.TransactionArgs, blkNumber *gethrpc.BlockNumber, gasCap uint64) (hexutil.Uint64, *big.Int, common.SystemError) {
	maxGasCap := rpc.calculateMaxGasCap(ctx, gasCap, args.Gas)
	// allowance will either be the maxGasCap or the balance allowance.
	// If the users funds are floaty, this might cause issues combined with the l1 pricing.
	allowance, feeCap, err := rpc.normalizeFeeCapAndAdjustGasLimit(ctx, args, blkNumber, maxGasCap)
	if err != nil {
		return 0, nil, err
	}

	// Set the gas limit to the provided gasCap
	args.Gas = (*hexutil.Uint64)(&allowance)

	// Perform a single gas estimation pass using isGasEnough
	failed, result, err := rpc.isGasEnough(ctx, args, gasCap, blkNumber)
	if err != nil {
		// Return zero values and the encountered error if estimation fails
		return 0, nil, err
	}

	if failed {
		if result != nil && result.Err != vm.ErrOutOfGas { //nolint: errorlint
			if len(result.Revert()) > 0 {
				return 0, gethcommon.Big0, newRevertError(result)
			}
			return 0, gethcommon.Big0, result.Err
		}
		// If the gas cap is insufficient, return an appropriate error
		return 0, nil, fmt.Errorf("gas required exceeds the provided gas cap (%d)", gasCap)
	}

	if result == nil {
		// If there's no result, something went wrong
		return 0, nil, fmt.Errorf("no execution result returned")
	}

	// Extract the gas used from the execution result
	gasUsed := hexutil.Uint64(result.UsedGas)

	return gasUsed, feeCap, nil
}

func (rpc *EncryptionManager) normalizeFeeCapAndAdjustGasLimit(ctx context.Context, args *gethapi.TransactionArgs, blkNumber *gethrpc.BlockNumber, hi uint64) (uint64, *big.Int, error) {
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
func (rpc *EncryptionManager) isGasEnough(ctx context.Context, args *gethapi.TransactionArgs, gas uint64, blkNumber *gethrpc.BlockNumber) (bool, *gethcore.ExecutionResult, error) {
	defer core.LogMethodDuration(rpc.logger, measure.NewStopwatch(), "enclave.go:IsGasEnough")
	args.Gas = (*hexutil.Uint64)(&gas)
	result, err := rpc.chain.ObsCallAtBlock(ctx, args, blkNumber)
	if err != nil {
		if errors.Is(err, gethcore.ErrIntrinsicGas) {
			return true, nil, nil // Special case, raise gas limit
		}
		return true, nil, err // Bail out
	}
	return result.Failed(), result, nil
}

func newRevertError(result *gethcore.ExecutionResult) *revertError {
	reason, errUnpack := abi.UnpackRevert(result.Revert())
	err := errors.New("execution reverted")
	if errUnpack == nil {
		err = fmt.Errorf("execution reverted: %v", reason)
	}
	return &revertError{
		error:  err,
		reason: hexutil.Encode(result.Revert()),
	}
}

// revertError is an API error that encompasses an EVM revertal with JSON error
// code and a binary data blob.
type revertError struct {
	error
	reason string // revert reason hex encoded
}

// ErrorCode returns the JSON error code for a revertal.
// See: https://github.com/ethereum/wiki/wiki/JSON-RPC-Error-Codes-Improvement-Proposal
func (e *revertError) ErrorCode() int {
	return 3
}

// ErrorData returns the hex encoded revert reason.
func (e *revertError) ErrorData() interface{} {
	return e.reason
}
