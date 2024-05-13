package rpc

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ten-protocol/go-ten/go/enclave/evm/ethchainadapter"

	"github.com/ten-protocol/go-ten/go/enclave/core"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/log"
	"github.com/ten-protocol/go-ten/go/enclave/events"
)

func GetTransactionReceiptValidate(reqParams []any, builder *CallBuilder[gethcommon.Hash, map[string]interface{}], _ *EncryptionManager) error {
	// Parameters are [Hash]
	if len(reqParams) < 1 {
		builder.Err = fmt.Errorf("unexpected number of parameters")
		return nil
	}
	txHashStr, ok := reqParams[0].(string)
	if !ok {
		builder.Err = fmt.Errorf("invalid transaction hash")
		return nil
	}

	txHash := gethcommon.HexToHash(txHashStr)
	builder.Param = &txHash
	return nil
}

func GetTransactionReceiptExecute(builder *CallBuilder[gethcommon.Hash, map[string]interface{}], rpc *EncryptionManager) error {
	txHash := *builder.Param
	// todo - optimise these calls. This can be done with a single sql
	rpc.logger.Trace("Get receipt for ", log.TxKey, txHash)
	// We retrieve the transaction.
	tx, blockHash, number, txIndex, err := rpc.storage.GetTransaction(builder.ctx, txHash) //nolint:dogsled
	if err != nil {
		rpc.logger.Trace("error getting tx ", log.TxKey, txHash, log.ErrKey, err)
		if errors.Is(err, errutil.ErrNotFound) {
			builder.Status = NotFound
			return nil
		}
		return err
	}

	// We retrieve the txSigner's address.
	txSigner, err := core.GetTxSigner(tx)
	if err != nil {
		builder.Err = err
		return nil //nolint:nilerr
	}

	if txSigner.Hex() != builder.VK.AccountAddress.Hex() {
		builder.Status = NotAuthorised
		return nil
	}

	// We retrieve the transaction receipt.
	txReceipt, err := rpc.storage.GetTransactionReceipt(builder.ctx, txHash)
	if err != nil {
		rpc.logger.Trace("error getting tx receipt", log.TxKey, txHash, log.ErrKey, err)
		if errors.Is(err, errutil.ErrNotFound) {
			builder.Status = NotFound
			return nil
		}
		// this is a system error
		return fmt.Errorf("could not retrieve transaction receipt in eth_getTransactionReceipt request. Cause: %w", err)
	}

	// We filter out irrelevant logs.
	txReceipt.Logs, err = events.FilterLogsForReceipt(builder.ctx, txReceipt, &txSigner, rpc.registry)
	if err != nil {
		rpc.logger.Error("error filter logs ", log.TxKey, txHash, log.ErrKey, err)
		// this is a system error
		return err
	}

	rpc.logger.Trace("Successfully retrieved receipt for ", log.TxKey, txHash, "rec", txReceipt)
	signer := types.MakeSigner(ethchainadapter.ChainParams(big.NewInt(rpc.config.ObscuroChainID)), big.NewInt(int64(number)), 0)
	r := marshalReceipt(txReceipt, blockHash, number, signer, tx, int(txIndex))
	builder.ReturnValue = &r
	return nil
}

// marshalReceipt marshals a transaction receipt into a JSON object.
// taken from geth
func marshalReceipt(receipt *types.Receipt, blockHash gethcommon.Hash, blockNumber uint64, signer types.Signer, tx *types.Transaction, txIndex int) map[string]interface{} {
	from, _ := types.Sender(signer, tx)

	fields := map[string]interface{}{
		"blockHash":         blockHash,
		"blockNumber":       hexutil.Uint64(blockNumber),
		"transactionHash":   tx.Hash(),
		"transactionIndex":  hexutil.Uint64(txIndex),
		"from":              from,
		"to":                tx.To(),
		"gasUsed":           hexutil.Uint64(receipt.GasUsed),
		"cumulativeGasUsed": hexutil.Uint64(receipt.CumulativeGasUsed),
		"contractAddress":   nil,
		"logs":              receipt.Logs,
		"logsBloom":         receipt.Bloom,
		"type":              hexutil.Uint(tx.Type()),
		"effectiveGasPrice": (*hexutil.Big)(receipt.EffectiveGasPrice),
	}

	// Assign receipt status or post state.
	if len(receipt.PostState) > 0 {
		fields["root"] = hexutil.Bytes(receipt.PostState)
	} else {
		fields["status"] = hexutil.Uint(receipt.Status)
	}
	if receipt.Logs == nil {
		fields["logs"] = []*types.Log{}
	}

	if tx.Type() == types.BlobTxType {
		fields["blobGasUsed"] = hexutil.Uint64(receipt.BlobGasUsed)
		fields["blobGasPrice"] = (*hexutil.Big)(receipt.BlobGasPrice)
	}

	// If the ContractAddress is 20 0x0 bytes, assume it is not a contract creation
	if receipt.ContractAddress != (gethcommon.Address{}) {
		fields["contractAddress"] = receipt.ContractAddress
	}
	return fields
}
