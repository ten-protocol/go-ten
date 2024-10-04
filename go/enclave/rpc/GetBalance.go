package rpc

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
)

type BalanceReq struct {
	Addr  *common.Address
	Block *rpc.BlockNumberOrHash
}

func GetBalanceValidate(reqParams []any, builder *CallBuilder[BalanceReq, hexutil.Big], _ *EncryptionManager) error {
	// Parameters are [Address, BlockNumber]
	if len(reqParams) != 2 {
		builder.Err = fmt.Errorf("unexpected number of parameters")
		return nil
	}
	requestedAddress, err := gethencoding.ExtractAddress(reqParams[0])
	if err != nil {
		builder.Err = fmt.Errorf("unable to extract requested Addr - %w", err)
		return nil
	}

	blockNumber, err := gethencoding.ExtractBlockNumber(reqParams[1])
	if err != nil {
		builder.Err = fmt.Errorf("unable to extract requested BlockHeader number - %w", err)
		return nil
	}
	builder.Param = &BalanceReq{
		Addr:  requestedAddress,
		Block: blockNumber,
	}
	return nil
}

func GetBalanceExecute(builder *CallBuilder[BalanceReq, hexutil.Big], rpc *EncryptionManager) error {
	// anybody can query for the native balance of any address.
	// even if we added a check, it could be bypassed easily
	balance, err := rpc.chain.GetBalanceAtBlock(builder.ctx, *builder.Param.Addr, builder.Param.Block.BlockNumber)
	if err != nil {
		return fmt.Errorf("unable to get balance - %w", err)
	}
	builder.ReturnValue = balance
	return nil
}
