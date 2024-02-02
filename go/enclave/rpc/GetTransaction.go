package rpc

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ten-protocol/go-ten/go/enclave/core"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common/errutil"
)

type TxData struct {
	tx          *types.Transaction
	blockHash   gethcommon.Hash
	blockNumber uint64
	index       uint64
}

func ExtractGetTransactionRequest(reqParams []any, rpc *EncryptionManager) (*UserRPCRequest1[TxData], error) {
	// Parameters are [Hash]
	if len(reqParams) != 1 {
		return nil, fmt.Errorf("unexpected address parameter")
	}
	txHashStr, ok := reqParams[0].(string)
	if !ok {
		return nil, fmt.Errorf("unexpected tx hash parameter")
	}
	txHash := gethcommon.HexToHash(txHashStr)

	// Unlike in the Geth impl, we do not try and retrieve unconfirmed transactions from the mempool.
	tx, blockHash, blockNumber, index, err := rpc.storage.GetTransaction(txHash)
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			// like geth, return an empty response when a not found tx is requested
			return nil, nil
		}
		return nil, err
	}

	viewingKeyAddress, err := core.GetSender(tx)
	if err != nil {
		err = fmt.Errorf("could not recover viewing key address to encrypt eth_getTransactionByHash response. Cause: %w", err)
		return nil, err
	}

	return &UserRPCRequest1[TxData]{&viewingKeyAddress, &TxData{tx, blockHash, blockNumber, index}}, nil
}

func ExecuteGetTransaction(data *UserRPCRequest1[TxData], _ *EncryptionManager) (*UserResponse[RpcTransaction], error) {
	if data == nil {
		return nil, nil
	}
	// Unlike in the Geth impl, we hardcode the use of a London signer.
	// todo (#1553) - once the enclave's genesis.json is set, retrieve the signer type using `types.MakeSigner`
	signer := types.NewLondonSigner(data.Param1.tx.ChainId())
	rpcTx := newRPCTransaction(data.Param1.tx, data.Param1.blockHash, data.Param1.blockNumber, data.Param1.index, gethcommon.Big0, signer)
	return &UserResponse[RpcTransaction]{rpcTx, nil}, nil
}

// Lifted from Geth's internal `ethapi` package.
type RpcTransaction struct {
	BlockHash        *gethcommon.Hash    `json:"blockHash"`
	BlockNumber      *hexutil.Big        `json:"blockNumber"`
	From             gethcommon.Address  `json:"from"`
	Gas              hexutil.Uint64      `json:"gas"`
	GasPrice         *hexutil.Big        `json:"gasPrice"`
	GasFeeCap        *hexutil.Big        `json:"maxFeePerGas,omitempty"`
	GasTipCap        *hexutil.Big        `json:"maxPriorityFeePerGas,omitempty"`
	Hash             gethcommon.Hash     `json:"hash"`
	Input            hexutil.Bytes       `json:"input"`
	Nonce            hexutil.Uint64      `json:"nonce"`
	To               *gethcommon.Address `json:"to"`
	TransactionIndex *hexutil.Uint64     `json:"transactionIndex"`
	Value            *hexutil.Big        `json:"value"`
	Type             hexutil.Uint64      `json:"type"`
	Accesses         *types.AccessList   `json:"accessList,omitempty"`
	ChainID          *hexutil.Big        `json:"chainId,omitempty"`
	V                *hexutil.Big        `json:"v"`
	R                *hexutil.Big        `json:"r"`
	S                *hexutil.Big        `json:"s"`
}

// Lifted from Geth's internal `ethapi` package.
func newRPCTransaction(tx *types.Transaction, blockHash gethcommon.Hash, blockNumber uint64, index uint64, baseFee *big.Int, signer types.Signer) *RpcTransaction {
	from, _ := types.Sender(signer, tx)
	v, r, s := tx.RawSignatureValues()
	result := &RpcTransaction{
		Type:     hexutil.Uint64(tx.Type()),
		From:     from,
		Gas:      hexutil.Uint64(tx.Gas()),
		GasPrice: (*hexutil.Big)(tx.GasPrice()),
		Hash:     tx.Hash(),
		Input:    hexutil.Bytes(tx.Data()),
		Nonce:    hexutil.Uint64(tx.Nonce()),
		To:       tx.To(),
		Value:    (*hexutil.Big)(tx.Value()),
		V:        (*hexutil.Big)(v),
		R:        (*hexutil.Big)(r),
		S:        (*hexutil.Big)(s),
	}
	if blockHash != (gethcommon.Hash{}) {
		result.BlockHash = &blockHash
		result.BlockNumber = (*hexutil.Big)(new(big.Int).SetUint64(blockNumber))
		result.TransactionIndex = (*hexutil.Uint64)(&index)
	}
	switch tx.Type() {
	case types.AccessListTxType:
		al := tx.AccessList()
		result.Accesses = &al
		result.ChainID = (*hexutil.Big)(tx.ChainId())
	case types.DynamicFeeTxType:
		al := tx.AccessList()
		result.Accesses = &al
		result.ChainID = (*hexutil.Big)(tx.ChainId())
		result.GasFeeCap = (*hexutil.Big)(tx.GasFeeCap())
		result.GasTipCap = (*hexutil.Big)(tx.GasTipCap())
		// if the transaction has been mined, compute the effective gas price
		if baseFee != nil && blockHash != (gethcommon.Hash{}) {
			// price = min(tip, gasFeeCap - baseFee) + baseFee
			price := math.BigMin(new(big.Int).Add(tx.GasTipCap(), baseFee), tx.GasFeeCap())
			result.GasPrice = (*hexutil.Big)(price)
		} else {
			result.GasPrice = (*hexutil.Big)(tx.GasFeeCap())
		}
	}
	return result
}
