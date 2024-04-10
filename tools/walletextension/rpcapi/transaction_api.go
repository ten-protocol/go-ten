package rpcapi

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common/gethapi"
	"github.com/ten-protocol/go-ten/go/enclave/rpc"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

type TransactionAPI struct {
	we *Services
}

func NewTransactionAPI(we *Services) *TransactionAPI {
	return &TransactionAPI{we}
}

func (s *TransactionAPI) GetBlockTransactionCountByNumber(ctx context.Context, blockNr gethrpc.BlockNumber) *hexutil.Uint {
	count, err := UnauthenticatedTenRPCCall[hexutil.Uint](ctx, s.we, &CacheCfg{CacheTypeDynamic: func() CacheStrategy {
		return cacheTTLBlockNumber(blockNr)
	}}, "eth_getBlockTransactionCountByNumber", blockNr)
	if err != nil {
		return nil
	}
	return count
}

func (s *TransactionAPI) GetBlockTransactionCountByHash(ctx context.Context, blockHash common.Hash) *hexutil.Uint {
	count, err := UnauthenticatedTenRPCCall[hexutil.Uint](ctx, s.we, &CacheCfg{CacheType: LongLiving}, "eth_getBlockTransactionCountByHash", blockHash)
	if err != nil {
		return nil
	}
	return count
}

func (s *TransactionAPI) GetTransactionByBlockNumberAndIndex(ctx context.Context, blockNr gethrpc.BlockNumber, index hexutil.Uint) *rpc.RpcTransaction {
	// not implemented
	return nil
}

func (s *TransactionAPI) GetTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) *rpc.RpcTransaction {
	// not implemented
	return nil
}

func (s *TransactionAPI) GetRawTransactionByBlockNumberAndIndex(ctx context.Context, blockNr gethrpc.BlockNumber, index hexutil.Uint) hexutil.Bytes {
	// not implemented
	return nil
}

func (s *TransactionAPI) GetRawTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) hexutil.Bytes {
	// not implemented
	return nil
}

func (s *TransactionAPI) GetTransactionCount(ctx context.Context, address common.Address, blockNrOrHash gethrpc.BlockNumberOrHash) (*hexutil.Uint64, error) {
	return ExecAuthRPC[hexutil.Uint64](ctx, s.we, &ExecCfg{account: &address}, "eth_getTransactionCount", address, blockNrOrHash)
}

func (s *TransactionAPI) GetTransactionByHash(ctx context.Context, hash common.Hash) (*rpc.RpcTransaction, error) {
	return ExecAuthRPC[rpc.RpcTransaction](ctx, s.we, &ExecCfg{tryAll: true}, "eth_getTransactionByHash", hash)
}

func (s *TransactionAPI) GetRawTransactionByHash(ctx context.Context, hash common.Hash) (hexutil.Bytes, error) {
	tx, err := ExecAuthRPC[hexutil.Bytes](ctx, s.we, &ExecCfg{tryAll: true}, "eth_getRawTransactionByHash", hash)
	if tx != nil {
		return *tx, err
	}
	return nil, err
}

func (s *TransactionAPI) GetTransactionReceipt(ctx context.Context, hash common.Hash) (map[string]interface{}, error) {
	txRec, err := ExecAuthRPC[map[string]interface{}](ctx, s.we, &ExecCfg{tryUntilAuthorised: true}, "eth_getTransactionReceipt", hash)
	if err != nil {
		return nil, err
	}
	if txRec == nil {
		return nil, err
	}
	return *txRec, err
}

func (s *TransactionAPI) SendTransaction(ctx context.Context, args gethapi.TransactionArgs) (common.Hash, error) {
	txRec, err := ExecAuthRPC[common.Hash](ctx, s.we, &ExecCfg{account: args.From}, "eth_sendTransaction", args)
	if err != nil {
		return common.Hash{}, err
	}
	userIDBytes, _ := extractUserID(ctx, s.we)
	if s.we.Config.StoreIncomingTxs && len(userIDBytes) > 10 {
		tx, err := json.Marshal(args)
		if err != nil {
			s.we.Logger().Error("error marshalling transaction: %s", err)
			return *txRec, nil
		}
		err = s.we.Storage.StoreTransaction(string(tx), userIDBytes)
		if err != nil {
			s.we.Logger().Error("error storing transaction in the database: %s", err)
			return *txRec, nil
		}
	}
	return *txRec, err
}

type SignTransactionResult struct {
	Raw hexutil.Bytes      `json:"raw"`
	Tx  *types.Transaction `json:"tx"`
}

func (s *TransactionAPI) FillTransaction(ctx context.Context, args gethapi.TransactionArgs) (*SignTransactionResult, error) {
	return nil, rpcNotImplemented
}

func (s *TransactionAPI) SendRawTransaction(ctx context.Context, input hexutil.Bytes) (common.Hash, error) {
	txRec, err := ExecAuthRPC[common.Hash](ctx, s.we, &ExecCfg{tryAll: true}, "eth_sendRawTransaction", input)
	if err != nil {
		return common.Hash{}, err
	}
	userIDBytes, err := extractUserID(ctx, s.we)
	if s.we.Config.StoreIncomingTxs && len(userIDBytes) > 10 {
		err = s.we.Storage.StoreTransaction(input.String(), userIDBytes)
		if err != nil {
			s.we.Logger().Error(fmt.Errorf("error storing transaction in the database: %w", err).Error())
		}
	}
	return *txRec, err
}

func (s *TransactionAPI) PendingTransactions() ([]*rpc.RpcTransaction, error) {
	return nil, rpcNotImplemented
}

func (s *TransactionAPI) Resend(ctx context.Context, sendArgs gethapi.TransactionArgs, gasPrice *hexutil.Big, gasLimit *hexutil.Uint64) (common.Hash, error) {
	txRec, err := ExecAuthRPC[common.Hash](ctx, s.we, &ExecCfg{account: sendArgs.From}, "eth_resend", sendArgs, gasPrice, gasLimit)
	if txRec != nil {
		return *txRec, err
	}
	return common.Hash{}, err
}
