package evm

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/go-obscuro/go/obscuronode/nodecommon"
)

// Perform the conversion between an Obscuro header and an Ethereum header that the EVM understands
// in the first stage we just encode the obscuro header in the Extra field
// todo - find a better way
func convertToEthHeader(h *nodecommon.Header) *types.Header {
	obscuroHeader, err := rlp.EncodeToBytes(h)
	if err != nil {
		panic(err)
	}
	return &types.Header{
		ParentHash:  h.ParentHash,
		Root:        h.State,
		TxHash:      common.Hash{},
		ReceiptHash: common.Hash{},
		Bloom:       types.Bloom{},
		Difficulty:  common.Big0,
		Number:      big.NewInt(int64(h.Number)),
		GasLimit:    1_000_000_000,
		GasUsed:     0,
		Time:        uint64(time.Now().Unix()),
		Extra:       obscuroHeader,
		MixDigest:   common.Hash{},
		Nonce:       types.BlockNonce{},
		BaseFee:     common.Big0,
	}
}

func convertFromEthHeader(header *types.Header) *nodecommon.Header {
	h := new(nodecommon.Header)
	err := rlp.DecodeBytes(header.Extra, h)
	if err != nil {
		panic(err)
	}
	return h
}
