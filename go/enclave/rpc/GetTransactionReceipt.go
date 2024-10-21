package rpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/ten-protocol/go-ten/go/enclave/storage"
	"github.com/ten-protocol/go-ten/go/enclave/storage/enclavedb"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/go/common/errutil"
	"github.com/ten-protocol/go-ten/go/common/log"
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
	requester := builder.VK.AccountAddress
	rpc.logger.Trace("Get receipt for ", log.TxKey, txHash, "requester", requester.Hex())

	rec, _ := rpc.cacheService.ReadReceipt(builder.ctx, txHash)
	if rec != nil {
		// receipt found in cache
		// we need to check whether the requester is the sender
		if rec.From == requester {
			logs := rec.Receipt.Logs
			if rec.To != nil && *rec.To != (gethcommon.Address{}) {
				// we need to filter the logs.
				ctr, err := rpc.storage.ReadContract(builder.ctx, *rec.To)
				if err != nil {
					return fmt.Errorf("could not read contract in eth_getTransactionReceipt request. Cause: %w", err)
				}
				logs, err = filterLogs(builder.ctx, rpc.storage, rec.Receipt.Logs, ctr, requester)
				if err != nil {
					return fmt.Errorf("could not filter cached logs in eth_getTransactionReceipt request. Cause: %w", err)
				}
			}
			r := marshalReceipt(rec.Receipt, logs, rec.From, rec.To)
			builder.ReturnValue = &r
			rpc.logger.Trace("Successfully retrieved receipt from cache for ", log.TxKey, txHash, "rec", rec)
			return nil
		}
	}

	exists, err := rpc.storage.ExistsTransactionReceipt(builder.ctx, txHash)
	if err != nil {
		return fmt.Errorf("could not retrieve transaction receipt in eth_getTransactionReceipt request. Cause: %w", err)
	}
	if !exists {
		builder.Status = NotFound
		return nil
	}

	// We retrieve the transaction receipt.
	receipt, err := rpc.storage.GetFilteredInternalReceipt(builder.ctx, txHash, requester, false)
	if err != nil {
		rpc.logger.Trace("error getting tx receipt", log.TxKey, txHash, log.ErrKey, err)
		if errors.Is(err, errutil.ErrNotFound) {
			builder.Status = NotAuthorised
			return nil
		}
		// this is a system error
		return fmt.Errorf("could not retrieve transaction receipt in eth_getTransactionReceipt request. Cause: %w", err)
	}

	rpc.logger.Trace("Successfully retrieved receipt for ", log.TxKey, txHash, "rec", receipt)
	r := receipt.MarshalToJson()
	builder.ReturnValue = &r
	return nil
}

// at this point the requester is the tx sender
func filterLogs(ctx context.Context, storage storage.Storage, logs []*types.Log, ctr *enclavedb.Contract, requester *gethcommon.Address) ([]*types.Log, error) {
	filtered := make([]*types.Log, 0)
	for _, l := range logs {

		eventSig := l.Topics[0]
		eventType, err := storage.ReadEventType(ctx, ctr.Address, eventSig)
		if err != nil {
			return nil, fmt.Errorf("could not read event type in eth_getTransactionReceipt request. Cause: %w", err)
		}
		// event visibility logic
		canView := eventType.ConfigPublic ||
			eventType.IsPublic() ||
			(eventType.AutoPublic != nil && *eventType.AutoPublic) ||
			(eventType.SenderCanView != nil && *eventType.SenderCanView) ||
			(eventType.Topic1CanView != nil && *eventType.Topic1CanView && isAddress(l.Topics, 1, requester)) ||
			(eventType.Topic2CanView != nil && *eventType.Topic2CanView && isAddress(l.Topics, 2, requester)) ||
			(eventType.Topic3CanView != nil && *eventType.Topic3CanView && isAddress(l.Topics, 3, requester)) ||
			(eventType.AutoVisibility && (isAddress(l.Topics, 1, requester) || isAddress(l.Topics, 2, requester) || isAddress(l.Topics, 3, requester)))

		if canView {
			filtered = append(filtered, l)
		}
	}
	return filtered, nil
}

func isAddress(topics []gethcommon.Hash, nr int, requester *gethcommon.Address) bool {
	if len(topics) < nr+1 {
		return false
	}
	topic := gethcommon.Address(topics[nr].Bytes())
	return topic == *requester
}

// marshalReceipt marshals a transaction receipt into a JSON object.
// taken from geth
func marshalReceipt(receipt *types.Receipt, logs []*types.Log, from, to *gethcommon.Address) map[string]interface{} {
	fields := map[string]interface{}{
		"blockHash":         receipt.BlockHash.Hex(),
		"blockNumber":       hexutil.Uint64(receipt.BlockNumber.Uint64()),
		"transactionHash":   receipt.TxHash.Hex(),
		"transactionIndex":  hexutil.Uint64(receipt.TransactionIndex),
		"from":              from,
		"to":                to,
		"gasUsed":           hexutil.Uint64(receipt.GasUsed),
		"cumulativeGasUsed": hexutil.Uint64(receipt.CumulativeGasUsed),
		"contractAddress":   nil,
		"logs":              logs,
		"logsBloom":         receipt.Bloom,
		"type":              hexutil.Uint(receipt.Type),
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

	// If the ContractAddress is 20 0x0 bytes, assume it is not a contract creation
	if receipt.ContractAddress != (gethcommon.Address{}) {
		fields["contractAddress"] = receipt.ContractAddress
	}
	return fields
}
