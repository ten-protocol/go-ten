package clientapi

import (
	"fmt"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
	"github.com/obscuronet/go-obscuro/go/host"
)

// ObscuroScanAPI implements ObscuroScan-specific JSON RPC operations.
type ObscuroScanAPI struct {
	host host.Host
}

func NewObscuroScanAPI(host host.Host) *ObscuroAPI {
	return &ObscuroAPI{
		host: host,
	}
}

// GetBlockHeaderByHash returns the header for the block with the given number.
func (api *ObscuroScanAPI) GetBlockHeaderByHash(blockHash gethcommon.Hash) (*types.Header, error) {
	blockHeader := api.host.DB().GetBlockHeader(blockHash)
	if blockHeader == nil {
		return nil, fmt.Errorf("no block with hash %s is stored", blockHash)
	}
	return blockHeader, nil
}

// GetRollupHeaderByNumber returns the header for the rollup with the given number.
func (api *ObscuroScanAPI) GetRollupHeaderByNumber(number *big.Int) (*common.Header, error) {
	rollupHash := api.host.DB().GetRollupHash(number)
	if rollupHash == nil {
		return nil, fmt.Errorf("no rollup with number %d is stored", number.Int64())
	}

	rollupHeader := api.host.DB().GetRollupHeader(*rollupHash)
	if rollupHeader == nil {
		return nil, fmt.Errorf("storage indicates that rollup %d has hash %s, but no such rollup is stored", number.Int64(), rollupHash)
	}

	return rollupHeader.Header, nil
}

// GetRollupForTx returns the rollup containing a given transaction hash. Required for ObscuroScan.
func (api *ObscuroScanAPI) GetRollupForTx(txHash gethcommon.Hash) (*common.ExtRollup, error) {
	rollupNumber := api.host.DB().GetRollupNumber(txHash)
	if rollupNumber == nil {
		return nil, fmt.Errorf("no rollup containing a transaction with hash %s is stored", txHash)
	}

	rollupHash := api.host.DB().GetRollupHash(rollupNumber)
	if rollupHash == nil {
		return nil, fmt.Errorf("no rollup with number %d is stored", rollupNumber.Int64())
	}

	rollup, err := api.host.EnclaveClient().GetRollup(*rollupHash)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve rollup with hash %s. Cause: %w", rollupNumber, err)
	}

	return rollup, nil
}

// GetLatestTransactions returns the hashes of the latest `num` transactions, or as many as possible if less than `num` transactions exist.
func (api *ObscuroAPI) GetLatestTransactions(num int) ([]gethcommon.Hash, error) {
	currentRollupHeaderWithHashes := api.host.DB().GetCurrentRollupHead()
	if currentRollupHeaderWithHashes == nil {
		return nil, nil
	}
	nextRollupHash := currentRollupHeaderWithHashes.Header.Hash()

	// We walk the chain until we've collected sufficient transactions.
	var txHashes []gethcommon.Hash
	for {
		rollupHeaderWithHashes := api.host.DB().GetRollupHeader(nextRollupHash)
		if rollupHeaderWithHashes == nil {
			return nil, fmt.Errorf("could not retrieve rollup for hash %s", nextRollupHash)
		}

		for _, txHash := range rollupHeaderWithHashes.TxHashes {
			txHashes = append(txHashes, txHash)
			if len(txHashes) >= num {
				return txHashes, nil
			}
		}

		// If we have reached the top of the chain (i.e. the current rollup's number is one), we stop walking.
		if rollupHeaderWithHashes.Header.Number.Cmp(big.NewInt(0)) == 0 {
			break
		}
		nextRollupHash = rollupHeaderWithHashes.Header.ParentHash
	}

	return txHashes, nil
}

// GetTotalTransactions returns the number of recorded transactions on the network.
func (api *ObscuroAPI) GetTotalTransactions() *big.Int {
	totalTransactions := api.host.DB().GetTotalTransactions()
	return totalTransactions
}

// Attestation returns the node's attestation details.
func (api *ObscuroAPI) Attestation() *common.AttestationReport {
	return api.host.EnclaveClient().Attestation()
}
