package obsclient

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ten-protocol/go-ten/go/common"
	"github.com/ten-protocol/go-ten/go/rpc"

	gethcommon "github.com/ethereum/go-ethereum/common"
	hostcommon "github.com/ten-protocol/go-ten/go/common/host"
)

// ObsClient provides access to general Obscuro functionality that doesn't require viewing keys.
//
// The methods in this client are analogous to the methods in geth's EthClient and should behave the same unless noted otherwise.
type ObsClient struct {
	rpcClient rpc.Client
}

func Dial(rawurl string) (*ObsClient, error) {
	rc, err := rpc.NewNetworkClient(rawurl)
	if err != nil {
		return nil, err
	}
	return NewObsClient(rc), nil
}

func NewObsClient(c rpc.Client) *ObsClient {
	return &ObsClient{c}
}

func (oc *ObsClient) Close() {
	oc.rpcClient.Stop()
}

// Blockchain Access

// ChainID retrieves the current chain ID for transaction replay protection.
func (oc *ObsClient) ChainID() (*big.Int, error) {
	var result hexutil.Big
	err := oc.rpcClient.Call(&result, rpc.ChainID)
	if err != nil {
		return nil, err
	}
	return (*big.Int)(&result), err
}

// BatchNumber returns the height of the head rollup
func (oc *ObsClient) BatchNumber() (uint64, error) {
	var result hexutil.Uint64
	err := oc.rpcClient.Call(&result, rpc.BatchNumber)
	return uint64(result), err
}

// GetBatchByHash returns the batch with the given hash.
func (oc *ObsClient) GetBatchByHash(hash gethcommon.Hash) (*common.ExtBatch, error) {
	var batch *common.ExtBatch
	err := oc.rpcClient.Call(&batch, rpc.GetBatch, hash)
	if err == nil && batch == nil {
		err = ethereum.NotFound
	}
	return batch, err
}

// GetBatchByHeight returns the batch with the given height.
func (oc *ObsClient) GetBatchByHeight(height *big.Int) (*common.PublicBatch, error) {
	var batch *common.PublicBatch

	err := oc.rpcClient.Call(&batch, rpc.GetBatchByHeight, height)
	if err == nil && batch == nil {
		err = ethereum.NotFound
	}
	return batch, err
}

// GetBatchBySeq returns the batch with the given sequence number.
func (oc *ObsClient) GetBatchBySeq(seq *big.Int) (*common.PublicBatch, error) {
	var batch *common.PublicBatch

	err := oc.rpcClient.Call(&batch, rpc.GetBatchBySeqNo, seq)
	if err == nil && batch == nil {
		err = ethereum.NotFound
	}
	return batch, err
}

// GetRollupBySeqNo returns the batch with the given height.
func (oc *ObsClient) GetRollupBySeqNo(seqNo uint64) (*common.PublicRollup, error) {
	var rollup *common.PublicRollup

	err := oc.rpcClient.Call(&rollup, rpc.GetRollupBySeqNo, seqNo)
	if err == nil && rollup == nil {
		err = ethereum.NotFound
	}
	return rollup, err
}

// GetBatchHeaderByNumber returns the header of the rollup with the given number
func (oc *ObsClient) GetBatchHeaderByNumber(number *big.Int) (*common.BatchHeader, error) {
	var batchHeader *common.BatchHeader
	err := oc.rpcClient.Call(&batchHeader, rpc.GetBatchByNumber, toBlockNumArg(number), false)
	if err == nil && batchHeader == nil {
		err = ethereum.NotFound
	}
	return batchHeader, err
}

// GetBatchHeaderByHash returns the block header with the given hash.
func (oc *ObsClient) GetBatchHeaderByHash(hash gethcommon.Hash) (*common.BatchHeader, error) {
	var batchHeader *common.BatchHeader
	err := oc.rpcClient.Call(&batchHeader, rpc.GetBatchByHash, hash, false)
	if err == nil && batchHeader == nil {
		err = ethereum.NotFound
	}
	return batchHeader, err
}

// GetTransaction returns the transaction.
func (oc *ObsClient) GetTransaction(hash gethcommon.Hash) (*common.PublicTransaction, error) {
	var tx *common.PublicTransaction
	err := oc.rpcClient.Call(&tx, rpc.GetTransaction, hash)
	if err == nil && tx == nil {
		err = ethereum.NotFound
	}
	return tx, err
}

// Health returns the Health status of the node.
func (oc *ObsClient) Health() (hostcommon.HealthCheck, error) {
	var healthy hostcommon.HealthCheck

	if err := oc.rpcClient.Call(&healthy, rpc.Health); err != nil {
		return hostcommon.HealthCheck{
			OverallHealth: false,
			Errors:        []string{fmt.Sprintf("RPC call failed: %v", err)},
		}, err
	}

	if !healthy.OverallHealth {
		if len(healthy.Errors) == 0 {
			healthy.Errors = []string{"Node reported unhealthy state without specific errors"}
		}
		return healthy, fmt.Errorf("node unhealthy: %s", strings.Join(healthy.Errors, ", "))
	}

	return healthy, nil
}

// GetTotalContractCount returns the total count of created contracts
func (oc *ObsClient) GetTotalContractCount() (int, error) {
	var count int
	err := oc.rpcClient.Call(&count, rpc.GetTotalContractCount)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetTotalTransactionCount returns the total count of executed transactions
func (oc *ObsClient) GetTotalTransactionCount() (int, error) {
	var count int
	err := oc.rpcClient.Call(&count, rpc.GetTotalTxCount)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetLatestRollupHeader returns the header of the latest rollup
func (oc *ObsClient) GetLatestRollupHeader() (*common.RollupHeader, error) {
	var header *common.RollupHeader
	err := oc.rpcClient.Call(&header, rpc.GetLatestRollupHeader)
	if err != nil {
		return nil, err
	}
	return header, nil
}

// GetLatestBatch returns the header of the latest rollup at tip
func (oc *ObsClient) GetLatestBatch() (*common.BatchHeader, error) {
	var header *common.BatchHeader
	err := oc.rpcClient.Call(&header, rpc.GetLatestBatch)
	if err != nil {
		return nil, err
	}
	return header, nil
}

// GetPublicTxListing returns a list of public transactions
func (oc *ObsClient) GetPublicTxListing(pagination *common.QueryPagination) (*common.TransactionListingResponse, error) {
	var result common.TransactionListingResponse
	err := oc.rpcClient.Call(&result, rpc.GetPublicTransactionData, pagination)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetBatchesListing returns a list of batches
func (oc *ObsClient) GetBatchesListing(pagination *common.QueryPagination) (*common.BatchListingResponse, error) {
	var result common.BatchListingResponse
	err := oc.rpcClient.Call(&result, rpc.GetBatchListing, pagination)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetBlockListing returns a list of block headers
func (oc *ObsClient) GetBlockListing(pagination *common.QueryPagination) (*common.BlockListingResponse, error) {
	var result common.BlockListingResponse
	err := oc.rpcClient.Call(&result, rpc.GetBlockListing, pagination)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetRollupListing returns a list of Rollups
func (oc *ObsClient) GetRollupListing(pagination *common.QueryPagination) (*common.RollupListingResponse, error) {
	var result common.RollupListingResponse
	err := oc.rpcClient.Call(&result, rpc.GetRollupListing, pagination)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetRollupByHash returns the public rollup data given its hash
func (oc *ObsClient) GetRollupByHash(hash gethcommon.Hash) (*common.PublicRollup, error) {
	var rollup *common.PublicRollup
	err := oc.rpcClient.Call(&rollup, rpc.GetRollupByHash, hash)
	if err == nil && rollup == nil {
		err = ethereum.NotFound
	}
	return rollup, err
}

// GetRollupBatches returns a list of public batch data within a given rollup hash
func (oc *ObsClient) GetRollupBatches(hash gethcommon.Hash, pagination *common.QueryPagination) (*common.BatchListingResponse, error) {
	var batchListing *common.BatchListingResponse
	err := oc.rpcClient.Call(&batchListing, rpc.GetRollupBatches, hash, pagination)
	if err == nil && batchListing == nil {
		err = ethereum.NotFound
	}
	return batchListing, err
}

// GetBatchTransactions returns a list of public transaction data within a given batch hash
func (oc *ObsClient) GetBatchTransactions(hash gethcommon.Hash, pagination *common.QueryPagination) (*common.TransactionListingResponse, error) {
	var txListing *common.TransactionListingResponse
	err := oc.rpcClient.Call(&txListing, rpc.GetBatchTransactions, hash, pagination)
	if err == nil && txListing == nil {
		err = ethereum.NotFound
	}
	return txListing, err
}

// GetConfig returns the network config for Ten
func (oc *ObsClient) GetConfig() (*common.TenNetworkInfo, error) {
	var result common.TenNetworkInfo
	err := oc.rpcClient.Call(&result, rpc.Config)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Search queries the host DB with the provided query string
func (oc *ObsClient) Search(query string) (*common.SearchResponse, error) {
	var result common.SearchResponse
	err := oc.rpcClient.Call(&result, rpc.Search, query)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
