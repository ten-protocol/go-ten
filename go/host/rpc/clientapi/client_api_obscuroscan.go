package clientapi

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/obscuronet/go-obscuro/go/host"

	"github.com/obscuronet/go-obscuro/go/common/errutil"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
)

const txLimit = 100

type obscuroscanAPIServiceLocator interface {
	host.DBLocator // todo (@matt) this api should depend on l1/l2 repos, not db directly
	host.EnclaveLocator
}

// ObscuroScanAPI implements ObscuroScan-specific JSON RPC operations.
type ObscuroScanAPI struct {
	sl obscuroscanAPIServiceLocator
}

func NewObscuroScanAPI(serviceLocator obscuroscanAPIServiceLocator) *ObscuroScanAPI {
	return &ObscuroScanAPI{
		sl: serviceLocator,
	}
}

// GetBlockHeaderByHash returns the header for the block with the given hash.
func (api *ObscuroScanAPI) GetBlockHeaderByHash(blockHash gethcommon.Hash) (*types.Header, error) {
	blockHeader, err := api.sl.DB().GetBlockHeader(blockHash)
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			return nil, fmt.Errorf("no block with hash %s is stored", blockHash)
		}
		return nil, fmt.Errorf("could not retrieve block with hash %s. Cause: %w", blockHash, err)
	}
	return blockHeader, nil
}

// GetBatch returns the batch with the given hash. Unlike `EthereumAPI.GetBlockByHash()`, returns the full
// `ExtBatch`, and not just the header.
func (api *ObscuroScanAPI) GetBatch(batchHash gethcommon.Hash) (*common.ExtBatch, error) {
	return api.sl.DB().GetBatch(batchHash)
}

// GetBatchForTx returns the batch containing a given transaction hash.
func (api *ObscuroScanAPI) GetBatchForTx(txHash gethcommon.Hash) (*common.ExtBatch, error) {
	batchNumber, err := api.sl.DB().GetBatchNumber(txHash)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve batch containing a transaction with hash %s. Cause: %w", txHash, err)
	}

	batchHash, err := api.sl.DB().GetBatchHash(batchNumber)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve batch with number %d. Cause: %w", batchNumber.Int64(), err)
	}

	return api.GetBatch(*batchHash)
}

// GetLatestTransactions returns the hashes of the latest `num` transactions confirmed in batches (or all the
// transactions if there are less than `num` total transactions).
func (api *ObscuroScanAPI) GetLatestTransactions(num int) ([]gethcommon.Hash, error) {
	// We prevent someone from requesting an excessive amount of transactions.
	if num > txLimit {
		return nil, fmt.Errorf("cannot request more than 100 latest transactions")
	}

	headBatchHeader, err := api.sl.DB().GetHeadBatchHeader()
	if err != nil {
		return nil, err
	}
	currentBatchHash := headBatchHeader.Hash()

	// We walk the chain until we've collected the requested number of transactions.
	var txHashes []gethcommon.Hash
	for {
		batchHeader, err := api.sl.DB().GetBatchHeader(currentBatchHash)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve batch for hash %s. Cause: %w", currentBatchHash, err)
		}

		batchTxHashes, err := api.sl.DB().GetBatchTxs(batchHeader.Hash())
		if err != nil {
			return nil, fmt.Errorf("could not retrieve transaction hashes for batch hash %s. Cause: %w", currentBatchHash, err)
		}

		for _, txHash := range batchTxHashes {
			txHashes = append(txHashes, txHash)
			if len(txHashes) >= num {
				return txHashes, nil
			}
		}

		// If we've reached the top of the chain, we stop walking.
		if batchHeader.Number.Uint64() == common.L2GenesisHeight {
			break
		}
		currentBatchHash = batchHeader.ParentHash
	}

	return txHashes, nil
}

// GetTotalTransactions returns the number of recorded transactions on the network.
func (api *ObscuroScanAPI) GetTotalTransactions() (*big.Int, error) {
	return api.sl.DB().GetTotalTransactions()
}

// Attestation returns the node's attestation details.
func (api *ObscuroScanAPI) Attestation() (*common.AttestationReport, error) {
	return api.sl.Enclave().GetEnclaveClient().Attestation()
}
