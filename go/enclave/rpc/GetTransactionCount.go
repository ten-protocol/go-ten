package rpc

import (
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
)

func GetTransactionCountValidate(reqParams []any, builder *CallBuilder[uint64, string], rpc *EncryptionManager) error {
	// Parameters are [Address, BlockHeader?]
	if len(reqParams) < 1 {
		builder.Err = fmt.Errorf("unexpected number of parameters")
		return nil
	}
	addressStr, ok := reqParams[0].(string)
	if !ok {
		builder.Err = fmt.Errorf("unexpected address parameter")
		return nil
	}

	address := gethcommon.HexToAddress(addressStr)

	seqNo := rpc.registry.HeadBatchSeq().Uint64()
	if len(reqParams) == 2 {
		tag, err := gethencoding.ExtractBlockNumber(reqParams[1])
		if err != nil {
			builder.Err = fmt.Errorf("unexpected tag parameter. Cause: %w", err)
			return nil
		}

		// todo - support BlockNumberOrHash
		b, err := rpc.registry.GetBatchAtHeight(builder.ctx, *tag.BlockNumber)
		if err != nil {
			builder.Err = fmt.Errorf("cant retrieve batch for tag. Cause: %w", err)
			return nil
		}
		seqNo = b.SeqNo().Uint64()
	}

	builder.From = &address
	builder.Param = &seqNo
	return nil
}

func GetTransactionCountExecute(builder *CallBuilder[uint64, string], rpc *EncryptionManager) error {
	err := authenticateFrom(builder.VK, builder.From)
	if err != nil {
		builder.Err = err
		return nil //nolint:nilerr
	}

	var nonce uint64
	l2Head, err := rpc.storage.FetchBatchHeaderBySeqNo(builder.ctx, *builder.Param)
	if err == nil {
		// todo - we should return an error when head state is not available, but for current test situations with race
		//  conditions we allow it to return zero while head state is uninitialized
		h := l2Head.Hash()
		s, err := rpc.registry.GetBatchState(builder.ctx, &h)
		if err != nil {
			return err
		}
		nonce = s.GetNonce(*builder.From)
	}

	enc := hexutil.EncodeUint64(nonce)
	builder.ReturnValue = &enc
	return nil
}
