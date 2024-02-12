package rpc

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/rpc"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
)

type BalanceReq struct {
	Addr  *common.Address
	Block *rpc.BlockNumber
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
	acctOwner, err := rpc.chain.AccountOwner(*builder.Param.Addr, builder.Param.Block)
	if err != nil {
		return err
	}

	// authorise the call
	if acctOwner.Hex() != builder.VK.AccountAddress.Hex() {
		rpc.logger.Debug("Unauthorised call", "address", acctOwner, "vk", builder.VK.AccountAddress, "userId", builder.VK.UserId)
		// return a default value
		builder.ReturnValue = (*hexutil.Big)(big.NewInt(0))
		return nil
	}

	balance, err := rpc.chain.GetBalanceAtBlock(*builder.Param.Addr, builder.Param.Block)
	if err != nil {
		return fmt.Errorf("unable to get balance - %w", err)
	}
	builder.ReturnValue = balance
	return nil
}
