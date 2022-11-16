package clientapi

import (
	"fmt"
	"math/big"

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

// GetBlockHeaderByHash returns the header for the block with the given number.
func (api *ObscuroScanAPI) GetBlockHeaderByHash(blockHash gethcommon.Hash) (*types.Header, error) {
	blockHeader, found := api.host.DB().GetBlockHeader(blockHash)
	if !found {
		return nil, fmt.Errorf("no block with hash %s is stored", blockHash)
	}
	return blockHeader, nil
}

// GetHeadRollupHeader returns the current head rollup's header.
// TODO - #718 - Switch to reading batch header.
func (api *ObscuroScanAPI) GetHeadRollupHeader() *common.Header {
	header, found := api.host.DB().GetHeadRollupHeader()
	if !found {
		return nil
	}
	return header
}

// GetRollup returns the rollup with the given hash.
func (api *ObscuroScanAPI) GetRollup(hash gethcommon.Hash) (*common.ExtRollup, error) {
	return api.host.EnclaveClient().GetRollup(hash)
}

// GetRollupForTx returns the rollup containing a given transaction hash. Required for ObscuroScan.
func (api *ObscuroScanAPI) GetRollupForTx(txHash gethcommon.Hash) (*common.ExtRollup, error) {
	rollupNumber, found := api.host.DB().GetRollupNumber(txHash)
	if !found {
		return nil, fmt.Errorf("no rollup containing a transaction with hash %s is stored", txHash)
	}

	rollupHash, found := api.host.DB().GetRollupHash(rollupNumber)
	if !found {
		return nil, fmt.Errorf("no rollup with number %d is stored", rollupNumber.Int64())
	}

	rollup, err := api.host.EnclaveClient().GetRollup(*rollupHash)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve rollup with hash %s. Cause: %w", rollupNumber, err)
	}

	return rollup, nil
}

// GetLatestTransactions returns the hashes of the latest `num` transactions, or as many as possible if less than `num`
// transactions exist.
// TODO - #718 - Switch to retrieving transactions from latest batch.
func (api *ObscuroScanAPI) GetLatestTransactions(num int) ([]gethcommon.Hash, error) {
	// We prevent someone from requesting an excessive amount of transactions.
	if num > txLimit {
		return nil, fmt.Errorf("cannot request more than 100 latest transactions")
	}

	headRollupHeader, found := api.host.DB().GetHeadRollupHeader()
	if !found {
		return nil, nil
	}
	currentRollupHash := headRollupHeader.Hash()

	// We walk the chain until we've collected the requested number of transactions.
	var txHashes []gethcommon.Hash
	for {
		rollupHeader, found := api.host.DB().GetRollupHeader(currentRollupHash)
		if !found {
			return nil, fmt.Errorf("could not retrieve rollup for hash %s", currentRollupHash)
		}

		rollupTxHashes, found := api.host.DB().GetRollupTxs(rollupHeader.Hash())
		if !found {
			return nil, fmt.Errorf("could not retrieve transaction hashes for rollup hash %s", currentRollupHash)
		}

		for _, txHash := range rollupTxHashes {
			txHashes = append(txHashes, txHash)
			if len(txHashes) >= num {
				break
			}
		}

		// If we've reached the top of the chain, we stop walking.
		if rollupHeader.Number.Uint64() == common.L2GenesisHeight {
			break
		}
		currentRollupHash = rollupHeader.ParentHash
	}

	return txHashes, nil
}

// GetTotalTransactions returns the number of recorded transactions on the network.
func (api *ObscuroScanAPI) GetTotalTransactions() *big.Int {
	totalTransactions := api.host.DB().GetTotalTransactions()
	return totalTransactions
}

// Attestation returns the node's attestation details.
func (api *ObscuroScanAPI) Attestation() (*common.AttestationReport, error) {
	return api.host.EnclaveClient().Attestation()
}
