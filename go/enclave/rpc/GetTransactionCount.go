package rpc

import (
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/common/gethencoding"
	"github.com/ten-protocol/go-ten/go/responses"
)

func (rpc *EncryptionManager) GetTransactionCount(encryptedParams common.EncryptedParamsGetTxCount) (*responses.TxCount, common.SystemError) {
	return withVKEncryption1[uint64, string](
		rpc,
		rpc.config.ObscuroChainID,
		encryptedParams,
		// extract sender and arguments
		func(reqParams []any) (*UserRPCRequest1[uint64], error) {
			// Parameters are [Address, Block?]
			if len(reqParams) < 1 {
				return nil, fmt.Errorf("unexpected number of parameters")
			}
			addressStr, ok := reqParams[0].(string)
			if !ok {
				return nil, fmt.Errorf("unexpected address parameter")
			}

			address := gethcommon.HexToAddress(addressStr)

			seqNo := rpc.registry.HeadBatchSeq().Uint64()
			if len(reqParams) == 2 {
				tag, err := gethencoding.ExtractBlockNumber(reqParams[1])
				if err != nil {
					return nil, fmt.Errorf("unexpected tag parameter. Cause: %w", err)
				}

				b, err := rpc.registry.GetBatchAtHeight(*tag)
				if err != nil {
					return nil, fmt.Errorf("cant retrieve batch for tag. Cause: %w", err)
				}
				seqNo = b.SeqNo().Uint64()
			}

			return &UserRPCRequest1[uint64]{&address, &seqNo}, nil
		},
		// make call and return result
		func(decodedParams *UserRPCRequest1[uint64]) (*UserResponse[string], error) {
			var nonce uint64
			l2Head, err := rpc.storage.FetchBatchBySeqNo(*decodedParams.Param1)
			if err == nil {
				// todo - we should return an error when head state is not available, but for current test situations with race
				//  conditions we allow it to return zero while head state is uninitialized
				s, err := rpc.storage.CreateStateDB(l2Head.Hash())
				if err != nil {
					return nil, err
				}
				nonce = s.GetNonce(*decodedParams.Sender)
			}

			encoded := hexutil.EncodeUint64(nonce)
			return &UserResponse[string]{&encoded, nil}, nil
		})
}
