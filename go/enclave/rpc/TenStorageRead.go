package rpc

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/privacy"
	"github.com/ten-protocol/go-ten/go/common/syserr"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

type StorageReadWithBlock struct {
	address     common.Address
	storageSlot string
	block       *gethrpc.BlockNumberOrHash
}

func TenStorageReadValidate(reqParams []any, builder *CallBuilder[StorageReadWithBlock, string], _ *EncryptionManager) error {
	// Parameters are [TransactionArgs, BlockNumber, 2 more which we don't support yet]
	if len(reqParams) < 2 || len(reqParams) > 3 {
		builder.Err = fmt.Errorf("unexpected number of parameters")
		return nil
	}

	addressHex, ok := reqParams[0].(string)
	if !ok {
		builder.Err = fmt.Errorf("address not provided in parameters")
		return nil
	}

	address := common.HexToAddress(addressHex)
	slot, ok := reqParams[1].(string)
	if !ok {
		builder.Err = fmt.Errorf("storage slot not provided in parameters")
		return nil
	}

	//todo: @siliev - this whitelist creation every time is bugging me
	whitelist := privacy.NewWhitelist()
	if !whitelist.AllowedStorageSlots[slot] {
		builder.Err = fmt.Errorf("storage slot not whitelisted")
		return nil
	}

	blkNumber, err := gethencoding.ExtractBlockNumber(reqParams[2])
	if err != nil {
		builder.Err = fmt.Errorf("unable to extract requested block number - %w", err)
		return nil
	}

	builder.Param = &StorageReadWithBlock{address, slot, blkNumber}

	return nil
}

func TenStorageReadExecute(builder *CallBuilder[StorageReadWithBlock, string], rpc *EncryptionManager) error {
	var err error
	var stateDb *state.StateDB
	blkNumber := builder.Param.block
	hash := blkNumber.BlockHash
	if hash != nil {
		stateDb, err = rpc.registry.GetBatchState(builder.ctx, hash)
	}

	number := blkNumber.BlockNumber
	if number != nil {
		stateDb, err = rpc.registry.GetBatchStateAtHeight(builder.ctx, number)
	}
	if err != nil {
		builder.Err = err
		return nil
	}

	value, err := stateDb.GetTrie().GetStorage(builder.Param.address, common.HexToHash(builder.Param.storageSlot).Bytes())
	if err != nil {
		rpc.logger.Debug("Failed eth_getStorageAt.", log.ErrKey, err)

		// return system errors to the host
		if errors.Is(err, syserr.InternalError{}) {
			return err
		}

		builder.Err = err
		return nil
	}

	encodedResult := hexutil.Encode(value)
	builder.ReturnValue = &encodedResult
	return nil
}
