package evm

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/obscuronet/obscuro-playground/go/obscuronode/nodecommon"
)

func convertToEthHeader(h *nodecommon.Header) *types.Header {
	// todo
	obscuroHeader, err := rlp.EncodeToBytes(h)
	if err != nil {
		panic(err)
	}
	return &types.Header{
		ParentHash: h.ParentHash,
		// UncleHash:   common.Hash{},
		// Coinbase:    common.Address{},
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
