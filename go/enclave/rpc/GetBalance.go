package rpc

import (
	"fmt"

	"github.com/status-im/keycard-go/hexutils"

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
		builder.Err = fmt.Errorf("unable to extract requested Block number - %w", err)
		return nil
	}
	builder.Param = &BalanceReq{
		Addr:  requestedAddress,
		Block: blockNumber,
	}
	return nil
}

func GetBalanceExecute(builder *CallBuilder[BalanceReq, hexutil.Big], rpc *EncryptionManager) error {
	acctOwner, err := rpc.chain.AccountOwner(builder.ctx, *builder.Param.Addr, builder.Param.Block.BlockNumber)
	if err != nil {
		return err
	}

	// authorise the call
	if acctOwner.Hex() != builder.VK.AccountAddress.Hex() {
		rpc.logger.Debug("Unauthorised call", "address", acctOwner, "vk", builder.VK.AccountAddress, "userId", hexutils.BytesToHex(builder.VK.UserID))
		builder.Status = NotAuthorised
		return nil
	}

	balance, err := rpc.chain.GetBalanceAtBlock(builder.ctx, *builder.Param.Addr, builder.Param.Block.BlockNumber)
	if err != nil {
		return fmt.Errorf("unable to get balance - %w", err)
	}
	builder.ReturnValue = balance
	return nil
}
