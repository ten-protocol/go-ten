package rpc

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/state"
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
	if !rpc.whitelist.AllowedStorageSlots[slot] && !contract.IsTransparent() {
		builder.Err = fmt.Errorf("eth_getStorageAt is not supported for this contract")
		return nil
	}

	blkNumber, err := gethencoding.ExtractBlockNumber(reqParams[2])
	if err != nil {
		builder.Err = fmt.Errorf("unable to extract requested block number - %w", err)
		return nil
	}

	builder.Param = &storageReadWithBlock{address, slot, blkNumber}

	return nil
}

func TenStorageReadExecute(builder *CallBuilder[storageReadWithBlock, string], rpc *EncryptionManager) error {
	var err error
	var stateDb *state.StateDB
	stateDb, err = rpc.registry.GetBatchState(builder.ctx, *builder.Param.block)
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

	account, err := stateDb.GetTrie().GetAccount(*builder.Param.address)
	if err != nil {
		builder.Err = fmt.Errorf("unable to get acct address - %w", err)
		return nil
	}

	trie, err := stateDb.Database().OpenTrie(account.Root)
	if err != nil {
		builder.Err = err
		return nil
	}

	value, err := trie.GetStorage(*builder.Param.address, storageSlot.Bytes())
	if err != nil {
		rpc.logger.Debug("Failed eth_getStorageAt.", log.ErrKey, err)

		// return system errors to the host
		if errors.Is(err, syserr.InternalError{}) {
			return fmt.Errorf("unable to get storage slot - %w", err)
		}

		builder.Err = fmt.Errorf("unable to get storage slot - %w", err)
		return nil
	}

	if len(value) == 0 {
		builder.ReturnValue = nil
		return nil
	}

	encodedResult := hexutil.Encode(value)
	builder.ReturnValue = &encodedResult
	return nil
}
