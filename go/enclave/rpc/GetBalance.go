package rpc

import (
	"fmt"

	"github.com/ethereum/go-ethereum/rpc"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
)

func ExtractGetBalanceRequestParams(reqParams []any, builder *RPCCallBuilder[rpc.BlockNumber, hexutil.Big], _ *EncryptionManager) error {
	// Parameters are [Address, BlockNumber]
	if len(reqParams) != 2 {
		builder.Err = fmt.Errorf("unexpected number of parameters")
		return nil
	}
	requestedAddress, err := gethencoding.ExtractAddress(reqParams[0])
	if err != nil {
		builder.Err = fmt.Errorf("unable to extract requested address - %w", err)
		return nil
	}

	blockNumber, err := gethencoding.ExtractBlockNumber(reqParams[1])
	if err != nil {
		builder.Err = fmt.Errorf("unable to extract requested block number - %w", err)
		return nil
	}
	builder.From = requestedAddress
	builder.Param = blockNumber
	return nil
}

func ExecuteGetBalance(rpcBuilder *RPCCallBuilder[rpc.BlockNumber, hexutil.Big], rpc *EncryptionManager) error {
	encryptAddress, balance, err := rpc.chain.GetBalance(*rpcBuilder.From, rpcBuilder.Param)
	if err != nil {
		return fmt.Errorf("unable to get balance - %w", err)
	}
	rpcBuilder.ResourceOwner = encryptAddress
	rpcBuilder.ReturnValue = balance
	return nil
}
