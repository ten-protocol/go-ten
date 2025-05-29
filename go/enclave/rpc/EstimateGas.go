package rpc

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
	"github.com/ten-protocol/go-ten/go/enclave/components"
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

// EstimateGasExecute - performs the gas estimation based on the provided parameters and the local environment configuration.
// Will accommodate l1 gas cost and stretch the final gas estimation.
// Note that setting gas price on the external call affects behaviour - this is due to a change geth implemented; If the account has
// no balance with the gas price the estimation might zero out.
func EstimateGasExecute(builder *CallBuilder[CallParamsWithBlock, hexutil.Uint64], rpc *EncryptionManager) error {
	err := authenticateFrom(builder.VK, builder.From)
	if err != nil {
		builder.Err = err
		return nil
	}

	txArgs := builder.Param.callParams
	headBatchSeq := rpc.registry.HeadBatchSeq()
	batch, err := rpc.storage.FetchBatchHeaderBySeqNo(builder.ctx, headBatchSeq.Uint64())
	if err != nil {
		return fmt.Errorf("failed to fetch batch header: %w", err)
	}

	ge := components.NewGasEstimator(rpc.storage, rpc.chain, rpc.gasOracle, rpc.logger)
	totalCost, _, userErr, sysErr := ge.EstimateTotalGas(builder.ctx, txArgs, builder.Param.block, batch, rpc.config.GasLocalExecutionCapFlag)

	if sysErr != nil {
		return fmt.Errorf("system error during gas estimation: %w", sysErr)
	}

	if userErr != nil {
		builder.Err = userErr
		return nil
	}

	totalGasEstimate := hexutil.Uint64(totalCost)
	builder.ReturnValue = &totalGasEstimate
	return nil
}
