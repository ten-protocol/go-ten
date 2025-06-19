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

func GetTransactionReceiptValidate(reqParams []any, builder *CallBuilder[gethcommon.Hash, map[string]interface{}], rpc *EncryptionManager) error {
	if !storeTxEnabled(rpc, builder) {
		return nil
	}
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

var NotAuthorisedErr = errors.New("not authorised")

func GetTransactionReceiptExecute(builder *CallBuilder[gethcommon.Hash, map[string]interface{}], rpc *EncryptionManager) error {
	txHash := *builder.Param
	requester := builder.VK.AccountAddress
	rpc.logger.Trace("Get receipt for ", log.TxKey, txHash, "requester", requester.Hex())

	// first try the cache for recent transactions
	result, err := fetchFromCache(builder.ctx, rpc.storage, rpc.cacheService, txHash, requester)
	// there is an explicit entry in the cache that the tx was not found
	if err != nil && errors.Is(err, storage.ReceiptDoesNotExist) {
		builder.Status = NotFound
		return nil
	}
	if err != nil && errors.Is(err, NotAuthorisedErr) {
		builder.Status = NotAuthorised
		return nil
	}
	// unexpected error
	if err != nil && !errors.Is(err, errutil.ErrNotFound) {
		return err
	}

	if result != nil {
		rpc.logger.Debug("Cache hit for receipt", log.TxKey, txHash)
		builder.ReturnValue = &result
		return nil
	}

	rpc.logger.Debug("Cache miss for receipt", log.TxKey, txHash.String())
	exists, err := rpc.storage.ExistsTransactionReceipt(builder.ctx, txHash)
	if err != nil {
		return fmt.Errorf("could not retrieve transaction receipt in eth_getTransactionReceipt request. Cause: %w", err)
	}
	if !exists {
		rpc.cacheService.ReceiptDoesNotExist(txHash)
		builder.Status = NotFound
		return nil
	}

	// We retrieve the transaction receipt.
	receipt, err := rpc.storage.GetFilteredInternalReceipt(builder.ctx, txHash, requester, false)
	if err != nil {
		rpc.logger.Trace("Error getting tx receipt", log.TxKey, txHash, log.ErrKey, err)
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

func fetchFromCache(ctx context.Context, storage storage.Storage, cacheService *storage.CacheService, txHash gethcommon.Hash, requester *gethcommon.Address) (map[string]interface{}, error) {
	rec, err := cacheService.ReadReceipt(ctx, txHash)
	if err != nil {
		return nil, err
	}

	logs := rec.Receipt.Logs
	// filter out the logs that the sender can't read
	// doesn't apply to contract creation (when to=nil)
	if len(logs) > 0 && (rec.To != nil && *rec.To != (gethcommon.Address{})) {
		ctr, err := storage.ReadContract(ctx, *rec.To)
		if err != nil && !errors.Is(err, errutil.ErrNotFound) {
			return nil, fmt.Errorf("could not read contract in eth_getTransactionReceipt request. Cause: %w", err)
		}
		// only filter when the transaction calls a contract. Value transfers emit no events.
		if ctr != nil {
			logs, err = filterLogs(rec.Receipt.Logs, ctr, requester)
			if err != nil && !errors.Is(err, errutil.ErrNotFound) {
				return nil, fmt.Errorf("could not filter cached logs in eth_getTransactionReceipt request. Cause: %w", err)
			}
			// in case the contract or event type is not stored yet, we try a db lookup
			if errors.Is(err, errutil.ErrNotFound) {
				return nil, nil
			}
		}
	}

	// check whether the requester is the sender or the requester can see any logs
	if (*rec.From != *requester) && (len(logs) == 0) {
		return nil, NotAuthorisedErr
	}

	r := marshalReceipt(rec.Receipt, logs, rec.From, rec.To)
	return r, nil
}

func filterLogs(logs []*types.Log, ctr *enclavedb.Contract, requester *gethcommon.Address) ([]*types.Log, error) {
	filtered := make([]*types.Log, 0)
	for _, l := range logs {
		canView, err := senderCanViewLog(ctr, l, requester)
		if err != nil {
			return nil, err
		}
		if canView {
			filtered = append(filtered, l)
		}
	}
	return filtered, nil
}

func senderCanViewLog(ctr *enclavedb.Contract, l *types.Log, sender *gethcommon.Address) (bool, error) {
	eventSig := l.Topics[0]
	eventType := ctr.EventType(eventSig)
	if eventType == nil {
		return false, errutil.ErrNotFound
	}
	// event visibility logic
	canView := eventType.IsPublic() ||
		(eventType.SenderCanView != nil && *eventType.SenderCanView) ||
		(eventType.Topic1CanView != nil && *eventType.Topic1CanView && isAddress(l.Topics, 1, sender)) ||
		(eventType.Topic2CanView != nil && *eventType.Topic2CanView && isAddress(l.Topics, 2, sender)) ||
		(eventType.Topic3CanView != nil && *eventType.Topic3CanView && isAddress(l.Topics, 3, sender)) ||
		(eventType.AutoVisibility && (isAddress(l.Topics, 1, sender) || isAddress(l.Topics, 2, sender) || isAddress(l.Topics, 3, sender)))
	return canView, nil
}

func isAddress(topics []gethcommon.Hash, nr int, requester *gethcommon.Address) bool {
	if len(topics) < nr+1 {
		return false
	}
	addressFromTopic := gethcommon.BytesToAddress(topics[nr].Bytes())
	return addressFromTopic == *requester
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
