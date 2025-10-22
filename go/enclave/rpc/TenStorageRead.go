package rpc

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

type storageReadWithBlock struct {
	address     *common.Address
	storageSlot string
	block       *gethrpc.BlockNumberOrHash
}

func TenStorageReadValidate(reqParams []any, builder *CallBuilder[storageReadWithBlock, hexutil.Bytes], rpc *EncryptionManager) error {
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

func TenStorageReadExecute(builder *CallBuilder[storageReadWithBlock, hexutil.Bytes], rpc *EncryptionManager) error {
	state, _, err := rpc.registry.GetBatchState(builder.ctx, *builder.Param.block, false)
	if err != nil {
		builder.Err = fmt.Errorf("unable to read block number - %w", err)
		return nil
	}

	key, _, err := decodeHash(builder.Param.storageSlot)
	if err != nil {
		builder.Err = fmt.Errorf("unable to decode storage key: %s", err)
		return nil
	}

	res := state.GetState(*builder.Param.address, key)
	if state.Error() != nil {
		builder.Err = fmt.Errorf("unable to read storage: %s", state.Error())
		return nil
	}

	enc := (hexutil.Bytes)(res.Big().Bytes())
	builder.ReturnValue = &enc

	rpc.logger.Debug("TenStorageReadExecute",
		"address", builder.Param.address.Hex(),
		"slot", builder.Param.storageSlot,
		"slot decoded", key,
		"block", builder.Param.block.String(),
		"result", enc.String())

	return nil
}

// decodeHash parses a hex-encoded 32-byte hash. The input may optionally
// be prefixed by 0x and can have a byte length up to 32.
// from go-ethereum
func decodeHash(s string) (h common.Hash, inputLength int, err error) {
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		s = s[2:]
	}
	if (len(s) & 1) > 0 {
		s = "0" + s
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		return common.Hash{}, 0, errors.New("hex string invalid")
	}
	if len(b) > 32 {
		return common.Hash{}, len(b), errors.New("hex string too long, want at most 32 bytes")
	}
	return common.BytesToHash(b), len(b), nil
}
