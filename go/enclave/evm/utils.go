package evm

import (
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/common"
)

// Perform the conversion between an Obscuro header and an Ethereum header that the EVM understands
// in the first stage we just encode the obscuro header in the Extra field
// todo - find a better way
func convertToEthHeader(h *common.Header) *types.Header {
	obscuroHeader, err := rlp.EncodeToBytes(h)
	if err != nil {
		panic(err)
	}
	return &types.Header{
		ParentHash:  h.ParentHash,
		Root:        h.Root,
		TxHash:      h.TxHash,
		ReceiptHash: h.ReceiptHash,
		Bloom:       h.Bloom,
		Difficulty:  gethcommon.Big0,
		Number:      h.Number,
		GasLimit:    1_000_000_000,
		GasUsed:     0,
		Time:        uint64(time.Now().Unix()),
		Extra:       obscuroHeader,
		MixDigest:   gethcommon.Hash{},
		Nonce:       types.BlockNonce{},
		BaseFee:     gethcommon.Big0,
	}
}

func convertFromEthHeader(header *types.Header) *common.Header {
	h := new(common.Header)
	err := rlp.DecodeBytes(header.Extra, h)
	if err != nil {
		panic(err)
	}
	return h
}
