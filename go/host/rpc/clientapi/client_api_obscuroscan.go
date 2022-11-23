package clientapi

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/obscuronet/go-obscuro/go/common/errutil"

	"github.com/obscuronet/go-obscuro/go/common/host"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
)

const txLimit = 100

// ObscuroScanAPI implements ObscuroScan-specific JSON RPC operations.
type ObscuroScanAPI struct {
	host host.Host
}

func NewObscuroScanAPI(host host.Host) *ObscuroScanAPI {
	return &ObscuroScanAPI{
		host: host,
	}
}

// GetBlockHeaderByHash returns the header for the block with the given hash.
func (api *ObscuroScanAPI) GetBlockHeaderByHash(blockHash gethcommon.Hash) (*types.Header, error) {
	blockHeader, err := api.host.DB().GetBlockHeader(blockHash)
	if err != nil {
		if errors.Is(err, errutil.ErrNotFound) {
			return nil, fmt.Errorf("no block with hash %s is stored", blockHash)
		}
		return nil, fmt.Errorf("could not retrieve block with hash %s. Cause: %w", blockHash, err)
	}
	return blockHeader, nil
}

// GetRollup returns the rollup with the given hash. Unlike `EthereumAPI.GetBlockByHash()`, returns the full
// `ExtRollup`, and not just the header.
// TODO - #718 - Switch to retrieving batch.
func (api *ObscuroScanAPI) GetRollup(hash gethcommon.Hash) (*common.ExtRollup, error) {
	return api.host.EnclaveClient().GetRollup(hash)
}

// GetRollupForTx returns the rollup containing a given transaction hash.
// TODO - #718 - Switch to retrieving batch based on transaction.
func (api *ObscuroScanAPI) GetRollupForTx(txHash gethcommon.Hash) (*common.ExtRollup, error) {
	rollupNumber, err := api.host.DB().GetRollupNumber(txHash)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve rollup containing a transaction with hash %s. Cause: %w", txHash, err)
	}

	rollupHash, err := api.host.DB().GetRollupHash(rollupNumber)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve rollup with number %d. Cause: %w", rollupNumber.Int64(), err)
	}

	rollup, err := api.host.EnclaveClient().GetRollup(*rollupHash)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve rollup with hash %s. Cause: %w", rollupNumber, err)
	}

	return rollup, nil
}

// GetLatestTransactions returns the hashes of the latest `num` transactions confirmed in batches (or all the
// transactions if there are less than `num` total transactions).
func (api *ObscuroScanAPI) GetLatestTransactions(num int) ([]gethcommon.Hash, error) {
	// We prevent someone from requesting an excessive amount of transactions.
	if num > txLimit {
		return nil, fmt.Errorf("cannot request more than 100 latest transactions")
	}

	headBatchHeader, err := api.host.DB().GetHeadBatchHeader()
	if err != nil {
		return nil, err
	}
	currentBatchHash := headBatchHeader.Hash()

	// We walk the chain until we've collected the requested number of transactions.
	var txHashes []gethcommon.Hash
	for {
		batchHeader, err := api.host.DB().GetBatchHeader(currentBatchHash)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve batch for hash %s. Cause: %w", currentBatchHash, err)
		}

		batchTxHashes, err := api.host.DB().GetBatchTxs(batchHeader.Hash())
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
	return api.host.DB().GetTotalTransactions()
}

// Attestation returns the node's attestation details.
func (api *ObscuroScanAPI) Attestation() (*common.AttestationReport, error) {
	return api.host.EnclaveClient().Attestation()
}
