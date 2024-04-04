package rpcapi

import (
	"context"
	"encoding/json"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ten-protocol/go-ten/go/common/gethapi"
	"github.com/ten-protocol/go-ten/lib/gethfork/rpc"
	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"
)

type BlockChainAPI struct {
	we *Services
}

func NewBlockChainAPI(we *Services) *BlockChainAPI {
	return &BlockChainAPI{we}
}

func (api *BlockChainAPI) ChainId() *hexutil.Big { //nolint:stylecheck
	chainID, _ := UnauthenticatedTenRPCCall[hexutil.Big](context.Background(), api.we, &CacheCfg{CacheType: LongLiving}, "eth_chainId")
	return chainID
}

func (api *BlockChainAPI) BlockNumber() hexutil.Uint64 {
	nr, err := UnauthenticatedTenRPCCall[hexutil.Uint64](context.Background(), api.we, &CacheCfg{CacheType: LatestBatch}, "eth_blockNumber")
	if err != nil {
		return hexutil.Uint64(0)
	}
	return *nr
}

func (api *BlockChainAPI) GetBalance(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Big, error) {
	return ExecAuthRPC[hexutil.Big](
		ctx,
		api.we,
		&ExecCfg{
			cacheCfg: &CacheCfg{
				CacheTypeDynamic: func() CacheStrategy {
					return cacheTTLBlockNumberOrHash(blockNrOrHash)
				},
			},
			account:            &address,
			tryUntilAuthorised: true, // the user can request the balance of a contract account
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

func (s *BlockChainAPI) GetProof(ctx context.Context, address common.Address, storageKeys []string, blockNrOrHash rpc.BlockNumberOrHash) (*AccountResult, error) {
	return nil, rpcNotImplemented
}

func (api *BlockChainAPI) GetHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (map[string]interface{}, error) {
	resp, err := UnauthenticatedTenRPCCall[map[string]interface{}](ctx, api.we, &CacheCfg{CacheTypeDynamic: func() CacheStrategy {
		return cacheTTLBlockNumber(number)
	}}, "eth_getHeaderByNumber", number)
	if resp == nil {
		return nil, err
	}
	return *resp, err
}

func (api *BlockChainAPI) GetHeaderByHash(ctx context.Context, hash common.Hash) map[string]interface{} {
	resp, _ := UnauthenticatedTenRPCCall[map[string]interface{}](ctx, api.we, &CacheCfg{CacheType: LongLiving}, "eth_getHeaderByHash", hash)
	if resp == nil {
		return nil
	}
	return *resp
}

func (api *BlockChainAPI) GetBlockByNumber(ctx context.Context, number rpc.BlockNumber, fullTx bool) (map[string]interface{}, error) {
	resp, err := UnauthenticatedTenRPCCall[map[string]interface{}](
		ctx,
		api.we,
		&CacheCfg{
			CacheTypeDynamic: func() CacheStrategy {
				return cacheTTLBlockNumber(number)
			},
		}, "eth_getBlockByNumber", number, fullTx)
	if resp == nil {
		return nil, err
	}
	return *resp, err
}

func (api *BlockChainAPI) GetBlockByHash(ctx context.Context, hash common.Hash, fullTx bool) (map[string]interface{}, error) {
	resp, err := UnauthenticatedTenRPCCall[map[string]interface{}](ctx, api.we, &CacheCfg{CacheType: LongLiving}, "eth_getBlockByHash", hash, fullTx)
	if resp == nil {
		return nil, err
	}
	return *resp, err
}

func (api *BlockChainAPI) GetCode(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	// todo - must be authenticated
	resp, err := UnauthenticatedTenRPCCall[hexutil.Bytes](
		ctx,
		api.we,
		&CacheCfg{
			CacheTypeDynamic: func() CacheStrategy {
				return cacheTTLBlockNumberOrHash(blockNrOrHash)
			},
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

func (api *BlockChainAPI) GetStorageAt(ctx context.Context, address common.Address, hexKey string, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	// GetStorageAt is repurposed to return the userID
	if address.Hex() == wecommon.GetStorageAtUserIDRequestMethodName {
		userID, err := extractUserID(ctx, api.we)
		if err != nil {
			return nil, err
		}

		_, err = getUser(userID, api.we)
		if err != nil {
			return nil, err
		}
		return userID, nil
	}

	resp, err := ExecAuthRPC[hexutil.Bytes](ctx, api.we, &ExecCfg{account: &address}, "eth_getStorageAt", address, hexKey, blockNrOrHash)
	if resp == nil {
		return nil, err
	}
	return *resp, err
}

func (s *BlockChainAPI) GetBlockReceipts(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) ([]map[string]interface{}, error) {
	return nil, rpcNotImplemented
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

func (api *BlockChainAPI) Call(ctx context.Context, args gethapi.TransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, overrides *StateOverride, blockOverrides *BlockOverrides) (hexutil.Bytes, error) {
	resp, err := ExecAuthRPC[hexutil.Bytes](ctx, api.we, &ExecCfg{
		cacheCfg: &CacheCfg{
			CacheTypeDynamic: func() CacheStrategy {
				return cacheTTLBlockNumberOrHash(blockNrOrHash)
			},
		},
		computeFromCallback: func(user *GWUser) *common.Address {
			return searchFromAndData(user.GetAllAddresses(), args)
		},
		adjustArgs: func(acct *GWAccount) []any {
			argsClone := populateFrom(acct, args)
			return []any{argsClone, blockNrOrHash, overrides, blockOverrides}
		},
		tryAll: true,
	}, "eth_call", args, blockNrOrHash, overrides, blockOverrides)
	if resp == nil {
		return nil, err
	}
	return *resp, err
}

func (api *BlockChainAPI) EstimateGas(ctx context.Context, args gethapi.TransactionArgs, blockNrOrHash *rpc.BlockNumberOrHash, overrides *StateOverride) (hexutil.Uint64, error) {
	resp, err := ExecAuthRPC[hexutil.Uint64](ctx, api.we, &ExecCfg{
		cacheCfg: &CacheCfg{
			CacheTypeDynamic: func() CacheStrategy {
				if blockNrOrHash != nil {
					return cacheTTLBlockNumberOrHash(*blockNrOrHash)
				}
				return LatestBatch
			},
		},
		computeFromCallback: func(user *GWUser) *common.Address {
			return searchFromAndData(user.GetAllAddresses(), args)
		},
		adjustArgs: func(acct *GWAccount) []any {
			argsClone := populateFrom(acct, args)
			return []any{argsClone, blockNrOrHash, overrides}
		},
		// is this a security risk?
		tryAll: true,
	}, "eth_estimateGas", args, blockNrOrHash, overrides)
	if resp == nil {
		return 0, err
	}
	return *resp, err
}

func populateFrom(acct *GWAccount, args gethapi.TransactionArgs) gethapi.TransactionArgs {
	// clone the args
	argsClone := cloneArgs(args)
	// set the from
	if args.From == nil || args.From.Hex() == (common.Address{}).Hex() {
		argsClone.From = acct.address
	}
	return argsClone
}

func cloneArgs(args gethapi.TransactionArgs) gethapi.TransactionArgs {
	serialised, _ := json.Marshal(args)
	var argsClone gethapi.TransactionArgs
	err := json.Unmarshal(serialised, &argsClone)
	if err != nil {
		return gethapi.TransactionArgs{}
	}
	return argsClone
}

type accessListResult struct {
	Accesslist *types.AccessList `json:"accessList"`
	Error      string            `json:"error,omitempty"`
	GasUsed    hexutil.Uint64    `json:"gasUsed"`
}

func (s *BlockChainAPI) CreateAccessList(ctx context.Context, args gethapi.TransactionArgs, blockNrOrHash *rpc.BlockNumberOrHash) (*accessListResult, error) {
	return nil, rpcNotImplemented
}
