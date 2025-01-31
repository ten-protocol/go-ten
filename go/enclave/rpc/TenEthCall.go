package rpc

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
	"github.com/ten-protocol/go-ten/go/common/log"
)

func TenCallValidate(reqParams []any, builder *CallBuilder[CallParamsWithBlock, string], _ *EncryptionManager) error {
	// Parameters are [TransactionArgs, BlockNumber, 2 more which we don't support yet]
	if len(reqParams) < 2 && len(reqParams) > 4 {
		builder.Err = fmt.Errorf("unexpected number of parameters")
		return nil
	}
	apiArgs, err := gethencoding.ExtractEthCall(reqParams[0])
	if err != nil {
		builder.Err = fmt.Errorf("unable to decode EthCall Params - %w", err)
		return nil
	}

	if apiArgs.From == nil {
		builder.Err = fmt.Errorf("no from address provided")
		return nil
	}

	blkNumber, err := gethencoding.ExtractBlockNumber(reqParams[1])
	if err != nil {
		builder.Err = fmt.Errorf("unable to extract requested block number - %w", err)
		return nil
	}

	builder.From = apiArgs.From
	// todo - support BlockNumberOrHash
	builder.Param = &CallParamsWithBlock{apiArgs, blkNumber.BlockNumber}

	return nil
}

func TenCallExecute(builder *CallBuilder[CallParamsWithBlock, string], rpc *EncryptionManager) error {
	err := authenticateFrom(builder.VK, builder.From)
	if err != nil {
		builder.Err = err
		return nil //nolint:nilerr
	}

	apiArgs := builder.Param.callParams
	blkNumber := builder.Param.block
	execResult, err := rpc.chain.Call(builder.ctx, apiArgs, blkNumber)
	if err != nil {
		rpc.logger.Debug("Failed eth_call.", log.ErrKey, err)
		return err
	}
	// If the result contains a revert reason, try to unpack and return it.
	if len(execResult.Revert()) > 0 {
		builder.Err = newRevertError(execResult.Revert())
		return nil
	}

	builder.Err = execResult.Err

	var encodedResult string
	if len(execResult.ReturnData) != 0 {
		encodedResult = hexutil.Encode(execResult.ReturnData)
		builder.ReturnValue = &encodedResult
	} else {
		builder.ReturnValue = nil
	}
	return nil
}
