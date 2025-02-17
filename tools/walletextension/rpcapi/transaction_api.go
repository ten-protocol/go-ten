package rpcapi

import (
	"context"
	"fmt"

	tenrpc "github.com/ten-protocol/go-ten/go/common/rpc"

	"github.com/ten-protocol/go-ten/tools/walletextension/cache"

	"github.com/ten-protocol/go-ten/tools/walletextension/services"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common/gethapi"
	"github.com/ten-protocol/go-ten/go/enclave/rpc"
	gethrpc "github.com/ten-protocol/go-ten/lib/gethfork/rpc"
)

type TransactionAPI struct {
	we *services.Services
}

func NewTransactionAPI(we *services.Services) *TransactionAPI {
	return &TransactionAPI{we}
}

func (s *TransactionAPI) GetBlockTransactionCountByNumber(ctx context.Context, blockNr gethrpc.BlockNumber) *hexutil.Uint {
	count, err := UnauthenticatedTenRPCCall[hexutil.Uint](ctx, s.we, &cache.Cfg{DynamicType: func() cache.Strategy {
		return cacheBlockNumber(blockNr)
	}}, "ten_getBlockTransactionCountByNumber", blockNr)
	if err != nil {
		return nil
	}
	return count
}

func (s *TransactionAPI) GetBlockTransactionCountByHash(ctx context.Context, blockHash common.Hash) *hexutil.Uint {
	count, err := UnauthenticatedTenRPCCall[hexutil.Uint](ctx, s.we, &cache.Cfg{Type: cache.LongLiving}, "eth_getBlockTransactionCountByHash", blockHash)
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
	return ExecAuthRPC[hexutil.Uint64](
		ctx,
		s.we,
		&AuthExecCfg{
			account: &address,
			cacheCfg: &cache.Cfg{
				DynamicType: func() cache.Strategy {
					return cacheBlockNumberOrHash(blockNrOrHash)
				},
			},
		},
		"ten_getTransactionCount",
		address,
		blockNrOrHash,
	)
}

func (s *TransactionAPI) GetTransactionByHash(ctx context.Context, hash common.Hash) (*rpc.RpcTransaction, error) {
	return ExecAuthRPC[rpc.RpcTransaction](ctx, s.we, &AuthExecCfg{tryAll: true, cacheCfg: &cache.Cfg{Type: cache.LongLiving}}, tenrpc.ERPCGetTransactionByHash, hash)
}

func (s *TransactionAPI) GetRawTransactionByHash(ctx context.Context, hash common.Hash) (hexutil.Bytes, error) {
	tx, err := ExecAuthRPC[hexutil.Bytes](ctx, s.we, &AuthExecCfg{tryAll: true, cacheCfg: &cache.Cfg{Type: cache.LongLiving}}, tenrpc.ERPCGetRawTransactionByHash, hash)
	if tx != nil {
		return *tx, err
	}
	return nil, err
}

func (s *TransactionAPI) GetTransactionReceipt(ctx context.Context, hash common.Hash) (map[string]interface{}, error) {
	txRec, err := ExecAuthRPC[map[string]interface{}](ctx, s.we, &AuthExecCfg{tryUntilAuthorised: true, cacheCfg: &cache.Cfg{Type: cache.LongLiving}}, tenrpc.ERPCGetTransactionReceipt, hash)
	if err != nil {
		return nil, err
	}
	if txRec == nil {
		return nil, err
	}
	return *txRec, err
}

func (s *TransactionAPI) SendTransaction(ctx context.Context, args gethapi.TransactionArgs) (common.Hash, error) {
	user, err := extractUserForRequest(ctx, s.we)
	if err != nil {
		return common.Hash{}, err
	}
	if !user.ActiveSK {
		return common.Hash{}, fmt.Errorf("please activate session key")
	}

	// when there is an active Session Key, sign all incoming transactions with that SK
	signedTx, err := s.we.SKManager.SignTx(ctx, user, args.ToTransaction())
	if err != nil {
		return common.Hash{}, err
	}

	blob, err := signedTx.MarshalBinary()
	if err != nil {
		return common.Hash{}, err
	}

	return s.sendRawTx(ctx, blob)
}

type SignTransactionResult struct {
	Raw hexutil.Bytes      `json:"raw"`
	Tx  *types.Transaction `json:"tx"`
}

func (s *TransactionAPI) FillTransaction(ctx context.Context, args gethapi.TransactionArgs) (*SignTransactionResult, error) {
	return nil, rpcNotImplemented
}

func (s *TransactionAPI) SendRawTransaction(ctx context.Context, input hexutil.Bytes) (common.Hash, error) {
	user, err := extractUserForRequest(ctx, s.we)
	if err != nil {
		return common.Hash{}, err
	}

	signedTxBlob := input
	// when there is an active Session Key, sign all incoming transactions with that SK
	if user.ActiveSK && user.SessionKey != nil {
		tx := new(types.Transaction)
		if err = tx.UnmarshalBinary(input); err != nil {
			return common.Hash{}, err
		}
		signedTx, err := s.we.SKManager.SignTx(ctx, user, tx)
		if err != nil {
			return common.Hash{}, err
		}
		signedTxBlob, err = signedTx.MarshalBinary()
		if err != nil {
			return common.Hash{}, err
		}
	}

	return s.sendRawTx(ctx, signedTxBlob)
}

func (s *TransactionAPI) sendRawTx(ctx context.Context, input hexutil.Bytes) (common.Hash, error) {
	txRec, err := ExecAuthRPC[common.Hash](ctx, s.we, &AuthExecCfg{tryAll: true, timeout: sendTransactionDuration}, tenrpc.ERPCSendRawTransaction, input)
	if err != nil {
		return common.Hash{}, err
	}
	return *txRec, err
}

func (s *TransactionAPI) PendingTransactions() ([]*rpc.RpcTransaction, error) {
	return nil, rpcNotImplemented
}

func (s *TransactionAPI) Resend(ctx context.Context, sendArgs gethapi.TransactionArgs, gasPrice *hexutil.Big, gasLimit *hexutil.Uint64) (common.Hash, error) {
	txRec, err := ExecAuthRPC[common.Hash](ctx, s.we, &AuthExecCfg{account: sendArgs.From}, tenrpc.ERPCResend, sendArgs, gasPrice, gasLimit)
	if txRec != nil {
		return *txRec, err
	}
	return common.Hash{}, err
}
