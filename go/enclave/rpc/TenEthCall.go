package rpc

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/gethapi"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/syserr"
	"github.com/ten-protocol/go-ten/go/responses"
)

// ObsCall handles param decryption, validation and encryption
// and requests the Rollup chain to execute the payload (eth_call)
func (rpc *EncryptionManager) ObsCall(encryptedParams common.EncryptedParamsCall) (*responses.Call, common.SystemError) {
	return withVKEncryption2[gethapi.TransactionArgs, gethrpc.BlockNumber](
		rpc,
		rpc.config.ObscuroChainID,
		encryptedParams,
		// extract sender and arguments
		func(reqParams []any) (*UserRPCRequest2[gethapi.TransactionArgs, gethrpc.BlockNumber], error) {
			// Parameters are [TransactionArgs, BlockNumber]
			if len(reqParams) != 2 {
				return nil, fmt.Errorf("unexpected number of parameters")
			}
			apiArgs, err := gethencoding.ExtractEthCall(reqParams[0])
			if err != nil {
				return nil, fmt.Errorf("unable to decode EthCall Params - %w", err)
			}

			// encryption will fail if no From address is provided
			if apiArgs.From == nil {
				return nil, fmt.Errorf("no from address provided")
			}

			blkNumber, err := gethencoding.ExtractBlockNumber(reqParams[1])
			if err != nil {
				return nil, fmt.Errorf("unable to extract requested block number - %w", err)
			}

			return &UserRPCRequest2[gethapi.TransactionArgs, gethrpc.BlockNumber]{apiArgs.From, apiArgs, blkNumber}, nil
		},
		// make call and return result
		func(decodedParams *UserRPCRequest2[gethapi.TransactionArgs, gethrpc.BlockNumber]) (any, error, error) {
			apiArgs := decodedParams.Param1
			blkNumber := decodedParams.Param2
			execResult, err := rpc.chain.ObsCall(apiArgs, blkNumber)
			if err != nil {
				rpc.logger.Debug("Failed eth_call.", log.ErrKey, err)

				// make sure it's not some internal error
				if errors.Is(err, syserr.InternalError{}) {
					return nil, nil, err
				}

				// make sure to serialize any possible EVM error
				evmErr, err := serializeEVMError(err)
				if err == nil {
					err = fmt.Errorf(string(evmErr))
				}
				return nil, err, nil
			}

			var encodedResult string
			if len(execResult.ReturnData) != 0 {
				encodedResult = hexutil.Encode(execResult.ReturnData)
			}

			return encodedResult, nil, nil
		})
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
