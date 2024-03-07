package rpcapi

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ten-protocol/go-ten/go/common/gethapi"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
	common2 "github.com/ten-protocol/go-ten/tools/walletextension/common"
)

type BlockChainAPI struct {
	we *Services
}

func NewBlockChainAPI(we *Services) *BlockChainAPI {
	return &BlockChainAPI{we}
}

func (api *BlockChainAPI) ChainId() *hexutil.Big {
	// chainid, _ := UnauthenticatedTenRPCCall[hexutil.Big](nil, api.we, &CacheCfg{TTL: longCacheTTL}, "eth_chainId")
	// return chainid
	return (*hexutil.Big)(big.NewInt(int64(api.we.TenChainID)))
}

func (s *BlockChainAPI) BlockNumber() hexutil.Uint64 {
	nr, err := UnauthenticatedTenRPCCall[hexutil.Uint64](nil, s.we, &CacheCfg{TTL: shortCacheTTL}, "eth_blockNumber")
	if err != nil {
		return hexutil.Uint64(0)
	}
	return *nr
}

func (s *BlockChainAPI) GetBalance(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Big, error) {
	// todo - how do you handle getBalance for contracts
	return ExecAuthRPC[hexutil.Big](
		ctx,
		s.we,
		&ExecCfg{
			cacheCfg: &CacheCfg{
				TTLCallback: func() time.Duration {
					if blockNrOrHash.BlockNumber != nil && blockNrOrHash.BlockNumber.Int64() <= 0 {
						return shortCacheTTL
					}
					return longCacheTTL
				},
			},
			account: &address,
		},
		"eth_getBalance",
		address,
		blockNrOrHash,
	)
}

// Result structs for GetProof
type AccountResult struct {
	Address      common.Address  `json:"address"`
	AccountProof []string        `json:"accountProof"`
	Balance      *hexutil.Big    `json:"balance"`
	CodeHash     common.Hash     `json:"codeHash"`
	Nonce        hexutil.Uint64  `json:"nonce"`
	StorageHash  common.Hash     `json:"storageHash"`
	StorageProof []StorageResult `json:"storageProof"`
}

type StorageResult struct {
	Key   string       `json:"key"`
	Value *hexutil.Big `json:"value"`
	Proof []string     `json:"proof"`
}

// proofList implements ethdb.KeyValueWriter and collects the proofs as
// hex-strings for delivery to rpc-caller.
type proofList []string

func (n *proofList) Put(key []byte, value []byte) error {
	*n = append(*n, hexutil.Encode(value))
	return nil
}

func (n *proofList) Delete(key []byte) error {
	panic("not supported")
}

func (s *BlockChainAPI) GetProof(ctx context.Context, address common.Address, storageKeys []string, blockNrOrHash rpc.BlockNumberOrHash) (*AccountResult, error) {
	// not implemented
	return nil, nil
}

func (s *BlockChainAPI) GetHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (map[string]interface{}, error) {
	resp, err := UnauthenticatedTenRPCCall[map[string]interface{}](ctx, s.we, &CacheCfg{TTLCallback: func() time.Duration {
		if number > 0 {
			return longCacheTTL
		}
		return shortCacheTTL
	}}, "eth_getHeaderByNumber", number)
	if resp == nil {
		return nil, err
	}
	return *resp, err
}

func (s *BlockChainAPI) GetHeaderByHash(ctx context.Context, hash common.Hash) map[string]interface{} {
	resp, _ := UnauthenticatedTenRPCCall[map[string]interface{}](ctx, s.we, &CacheCfg{TTL: longCacheTTL}, "eth_getHeaderByHash", hash)
	if resp == nil {
		return nil
	}
	return *resp
}

func (s *BlockChainAPI) GetBlockByNumber(ctx context.Context, number rpc.BlockNumber, fullTx bool) (map[string]interface{}, error) {
	resp, err := UnauthenticatedTenRPCCall[map[string]interface{}](
		ctx,
		s.we,
		&CacheCfg{
			TTLCallback: func() time.Duration {
				if number > 0 {
					return longCacheTTL
				}
				return shortCacheTTL
			},
		}, "eth_getBlockByNumber", number, fullTx)
	if resp == nil {
		return nil, err
	}
	return *resp, err
}

func (s *BlockChainAPI) GetBlockByHash(ctx context.Context, hash common.Hash, fullTx bool) (map[string]interface{}, error) {
	resp, err := UnauthenticatedTenRPCCall[map[string]interface{}](ctx, s.we, &CacheCfg{TTL: longCacheTTL}, "eth_getBlockByHash", hash, fullTx)
	if resp == nil {
		return nil, err
	}
	return *resp, err
}

func (s *BlockChainAPI) GetCode(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	resp, err := ExecAuthRPC[hexutil.Bytes](
		ctx,
		s.we,
		&ExecCfg{
			cacheCfg: &CacheCfg{
				TTLCallback: func() time.Duration {
					if blockNrOrHash.BlockNumber != nil && blockNrOrHash.BlockNumber.Int64() <= 0 {
						return shortCacheTTL
					}
					return longCacheTTL
				},
			},
			account: &address,
		},
		"eth_getCode",
		address,
		blockNrOrHash,
	)
	if resp == nil {
		return nil, err
	}
	return *resp, err
}

func (s *BlockChainAPI) GetStorageAt(ctx context.Context, address common.Address, hexKey string, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	// GetStorageAt is repurposed to return the userId
	if hexKey == common2.GetStorageAtUserIDRequestMethodName {
		userId, err := extractUserId(ctx)
		if err != nil {
			return nil, err
		}

		_, err = getUser(userId, s.we.Storage)
		if err != nil {
			return nil, err
		}
		return userId, nil
	}

	resp, err := ExecAuthRPC[hexutil.Bytes](ctx, s.we, &ExecCfg{account: &address}, "eth_getStorageAt", address, hexKey, blockNrOrHash)
	if resp == nil {
		return nil, err
	}
	return *resp, err
}

func (s *BlockChainAPI) GetBlockReceipts(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) ([]map[string]interface{}, error) {
	// not implemented
	return nil, nil
}

type ChainContext struct {
	we *Services
}

type OverrideAccount struct {
	Nonce     *hexutil.Uint64              `json:"nonce"`
	Code      *hexutil.Bytes               `json:"code"`
	Balance   **hexutil.Big                `json:"balance"`
	State     *map[common.Hash]common.Hash `json:"state"`
	StateDiff *map[common.Hash]common.Hash `json:"stateDiff"`
}
type (
	StateOverride  map[common.Address]OverrideAccount
	BlockOverrides struct {
		Number     *hexutil.Big
		Difficulty *hexutil.Big
		Time       *hexutil.Uint64
		GasLimit   *hexutil.Uint64
		Coinbase   *common.Address
		Random     *common.Hash
		BaseFee    *hexutil.Big
	}
)

func (s *BlockChainAPI) Call(ctx context.Context, args gethapi.TransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, overrides *StateOverride, blockOverrides *BlockOverrides) (hexutil.Bytes, error) {
	resp, err := ExecAuthRPC[hexutil.Bytes](ctx, s.we, &ExecCfg{
		cacheCfg: &CacheCfg{
			TTLCallback: func() time.Duration {
				if blockNrOrHash.BlockNumber != nil && blockNrOrHash.BlockNumber.Int64() <= 0 {
					return shortCacheTTL
				}
				return longCacheTTL
			},
		},
		computeFromCallback: func(user *GWUser) *common.Address {
			return searchFromAndData(user.GetAllAddresses(), args)
		},
		adjustArgs: func(acct *GWAccount) []any {
			// set the from
			args.From = acct.address
			return []any{args, blockNrOrHash, overrides, blockOverrides}
		},
	}, "eth_call", args, blockNrOrHash, overrides, blockOverrides)
	if resp == nil {
		return nil, err
	}
	return *resp, err
}

func (s *BlockChainAPI) EstimateGas(ctx context.Context, args gethapi.TransactionArgs, blockNrOrHash *rpc.BlockNumberOrHash, overrides *StateOverride) (hexutil.Uint64, error) {
	resp, err := ExecAuthRPC[hexutil.Uint64](ctx, s.we, &ExecCfg{
		cacheCfg: &CacheCfg{
			TTLCallback: func() time.Duration {
				if blockNrOrHash != nil && blockNrOrHash.BlockNumber != nil && blockNrOrHash.BlockNumber.Int64() <= 0 {
					return shortCacheTTL
				}
				return longCacheTTL
			},
		},
		computeFromCallback: func(user *GWUser) *common.Address {
			return searchFromAndData(user.GetAllAddresses(), args)
		},
	}, "eth_estimateGas", args, blockNrOrHash, overrides)
	if resp == nil {
		return 0, err
	}
	return *resp, err
}

type accessListResult struct {
	Accesslist *types.AccessList `json:"accessList"`
	Error      string            `json:"error,omitempty"`
	GasUsed    hexutil.Uint64    `json:"gasUsed"`
}

func (s *BlockChainAPI) CreateAccessList(ctx context.Context, args gethapi.TransactionArgs, blockNrOrHash *rpc.BlockNumberOrHash) (*accessListResult, error) {
	// not implemented
	return nil, nil
}
