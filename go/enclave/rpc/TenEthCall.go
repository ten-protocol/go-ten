package rpc

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/syserr"
)

func ExtractObsCallRequest(reqParams []any, builder *RPCCallBuilder[CallParamsWithBlock, string], _ *EncryptionManager) error {
	// Parameters are [TransactionArgs, BlockNumber]
	if len(reqParams) != 2 {
		builder.Err = fmt.Errorf("unexpected number of parameters")
		return nil
	}
	apiArgs, err := gethencoding.ExtractEthCall(reqParams[0])
	if err != nil {
		builder.Err = fmt.Errorf("unable to decode EthCall Params - %w", err)
		return nil
	}

	// encryption will fail if no From address is provided
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
	builder.Param = &CallParamsWithBlock{apiArgs, blkNumber}

	return nil
}

func ExecuteObsCallGas(rpcBuilder *RPCCallBuilder[CallParamsWithBlock, string], rpc *EncryptionManager) error {
	apiArgs := rpcBuilder.Param.callParams
	blkNumber := rpcBuilder.Param.block
	execResult, err := rpc.chain.ObsCall(apiArgs, blkNumber)
	if err != nil {
		rpc.logger.Debug("Failed eth_call.", log.ErrKey, err)

		// return system errors to the host
		if errors.Is(err, syserr.InternalError{}) {
			return err
		}

		// extract the EVM error
		evmErr, err := serializeEVMError(err)
		if err == nil {
			err = fmt.Errorf(string(evmErr))
		}
		rpcBuilder.Err = err
		return nil
	}

	var encodedResult string
	if len(execResult.ReturnData) != 0 {
		encodedResult = hexutil.Encode(execResult.ReturnData)
	}
	rpcBuilder.ReturnValue = &encodedResult
	return nil
}

func serializeEVMError(err error) ([]byte, error) {
	var errReturn interface{}

	// check if it's a serialized error and handle any error wrapping that might have occurred
	var e *errutil.EVMSerialisableError
	if ok := errors.As(err, &e); ok {
		errReturn = e
	} else {
		// it's a generic error, serialise it
		errReturn = &errutil.EVMSerialisableError{Err: err.Error()}
	}

	// serialise the error object returned by the evm into a json
	errSerializedBytes, marshallErr := json.Marshal(errReturn)
	if marshallErr != nil {
		return nil, marshallErr
	}
	return errSerializedBytes, nil
}
