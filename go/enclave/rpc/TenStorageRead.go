package rpc

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/common/syserr"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

type storageReadWithBlock struct {
	address     *common.Address
	storageSlot string
	block       *gethrpc.BlockNumberOrHash
}

func TenStorageReadValidate(reqParams []any, builder *CallBuilder[storageReadWithBlock, string], rpc *EncryptionManager) error {
	if len(reqParams) < 2 || len(reqParams) > 3 {
		builder.Err = fmt.Errorf("unexpected number of parameters")
		return nil
	}

	address, err := gethencoding.ExtractAddress(reqParams[0])
	if err != nil {
		builder.Err = fmt.Errorf("error extracting address - %w", err)
		return nil
	}

	slot, ok := reqParams[1].(string)
	if !ok {
		builder.Err = fmt.Errorf("storage slot not provided in parameters")
		return nil
	}

	contract, err := rpc.storage.ReadContract(builder.ctx, *address)
	if err != nil {
		builder.Err = fmt.Errorf("eth_getStorageAt is not supported for this contract")
		return nil
	}

	// block the call for un-transparent contracts and non-whitelisted slots
	if !rpc.storageSlotWhitelist.AllowedStorageSlots[slot] && !contract.IsTransparent() {
		builder.Err = fmt.Errorf("eth_getStorageAt is not supported for this contract")
		return nil
	}

	var blk any = nil
	if len(reqParams) == 3 {
		blk = reqParams[2]
	}

	blkNumber, err := gethencoding.ExtractBlockNumber(blk)
	if err != nil {
		builder.Err = fmt.Errorf("unable to extract requested block number - %w", err)
		return nil
	}

	builder.Param = &storageReadWithBlock{address: address, storageSlot: slot, block: blkNumber}

	return nil
}

func TenStorageReadExecute(builder *CallBuilder[storageReadWithBlock, string], rpc *EncryptionManager) error {
	_, reader, err := rpc.registry.GetBatchState(builder.ctx, *builder.Param.block)
	if err != nil {
		builder.Err = fmt.Errorf("unable to read block number - %w", err)
		return nil
	}

	sl := new(big.Int)
	sl, ok := sl.SetString(builder.Param.storageSlot, 0)
	if !ok {
		builder.Err = fmt.Errorf("unable to parse storage slot (%s)", builder.Param.storageSlot)
		return nil
	}

	// the storage slot needs to be 32 bytes padded with 0s
	storageSlot := common.Hash{}
	storageSlot.SetBytes(sl.Bytes())

	value, err := reader.Storage(*builder.Param.address, storageSlot)
	if err != nil {
		rpc.logger.Debug("Failed eth_getStorageAt.", log.ErrKey, err)

		// return system errors to the host
		if errors.Is(err, syserr.InternalError{}) {
			return fmt.Errorf("unable to get storage slot - %w", err)
		}

		builder.Err = fmt.Errorf("unable to get storage slot - %w", err)
		return nil
	}

	encodedResult := value.Hex()
	builder.ReturnValue = &encodedResult
	return nil
}
