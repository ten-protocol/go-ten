package evm

import (
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ten-protocol/go-ten/go/common"
)

func MaskedSender(address gethcommon.Address) gethcommon.Address {
	return gethcommon.BigToAddress(big.NewInt(0).Sub(address.Big(), big.NewInt(1)))
}

func TransactionToMessageWithOverrides(
	tx *common.L2PricedTransaction,
	config *params.ChainConfig,
	header *types.Header,
) (*core.Message, error) {
	// Override from can be used for calling system contracts from underivable addresses like all zeroes
	if tx.FromSelf {
		msg := TransactionToMessageNoSender(tx.Tx, header.BaseFee)
		msg.From = MaskedSender(*msg.To)
		return msg, nil
	}
	return core.TransactionToMessage(tx.Tx, types.MakeSigner(config, header.Number, header.Time), header.BaseFee)
}

func TransactionToMessageNoSender(tx *types.Transaction, baseFee *big.Int) *core.Message {
	msg := &core.Message{
		Nonce:             tx.Nonce(),
		GasLimit:          tx.Gas(),
		GasPrice:          new(big.Int).Set(tx.GasPrice()),
		GasFeeCap:         new(big.Int).Set(tx.GasFeeCap()),
		GasTipCap:         new(big.Int).Set(tx.GasTipCap()),
		To:                tx.To(),
		Value:             tx.Value(),
		Data:              tx.Data(),
		AccessList:        tx.AccessList(),
		SkipAccountChecks: false,
		BlobHashes:        tx.BlobHashes(),
		BlobGasFeeCap:     tx.BlobGasFeeCap(),
	}

	if baseFee != nil {
		msg.GasPrice = math.BigMin(msg.GasPrice.Add(msg.GasTipCap, baseFee), msg.GasFeeCap)
	}

	return msg
}
