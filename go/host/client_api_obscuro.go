package host

import (
	"fmt"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
)

// TODO - Some methods return nil for an unfound block/rollup, while others return an error. Harmonise.

// ObscuroAPI implements Obscuro-specific JSON RPC operations.
type ObscuroAPI struct {
	host *Node
}

func NewObscuroAPI(host *Node) *ObscuroAPI {
	return &ObscuroAPI{
		host: host,
	}
}

// GetID returns the ID of the host.
func (api *ObscuroAPI) GetID() gethcommon.Address {
	return api.host.ID
}

// GetCurrentBlockHead returns the current head block's header.
func (api *ObscuroAPI) GetCurrentBlockHead() *types.Header {
	return api.host.nodeDB.GetCurrentBlockHead()
}

// GetBlockHeaderByHash returns the header for the block with the given number.
func (api *ObscuroAPI) GetBlockHeaderByHash(blockHash gethcommon.Hash) (*types.Header, error) {
	blockHeader := api.host.nodeDB.GetBlockHeader(blockHash)
	if blockHeader == nil {
		return nil, fmt.Errorf("no block with hash %s is stored", blockHash)
	}
	return blockHeader, nil
}

// GetCurrentRollupHead returns the current head rollup's header.
func (api *ObscuroAPI) GetCurrentRollupHead() *common.Header {
	headerWithHashes := api.host.nodeDB.GetCurrentRollupHead()
	if headerWithHashes == nil {
		return nil
	}
	return headerWithHashes.Header
}

// GetRollupHeader returns the header of the rollup with the given hash.
func (api *ObscuroAPI) GetRollupHeader(hash gethcommon.Hash) *common.Header {
	headerWithHashes := api.host.nodeDB.GetRollupHeader(hash)
	if headerWithHashes == nil {
		return nil
	}
	return headerWithHashes.Header
}

// GetRollupHeaderByNumber returns the header for the rollup with the given number.
func (api *ObscuroAPI) GetRollupHeaderByNumber(number *big.Int) (*common.Header, error) {
	rollupHash := api.host.nodeDB.GetRollupHash(number)
	if rollupHash == nil {
		return nil, fmt.Errorf("no rollup with number %d is stored", number.Int64())
	}

	rollupHeader := api.host.nodeDB.GetRollupHeader(*rollupHash)
	if rollupHeader == nil {
		return nil, fmt.Errorf("storage indicates that rollup %d has hash %s, but no such rollup is stored", number.Int64(), rollupHash)
	}

	return rollupHeader.Header, nil
}

// GetRollup returns the rollup with the given hash.
func (api *ObscuroAPI) GetRollup(hash gethcommon.Hash) (*common.ExtRollup, error) {
	return api.host.EnclaveClient.GetRollup(hash)
}

// Nonce returns the nonce of the wallet with the given address.
func (api *ObscuroAPI) Nonce(address gethcommon.Address) uint64 {
	return api.host.EnclaveClient.Nonce(address)
}

// AddViewingKey stores the viewing key on the enclave.
func (api *ObscuroAPI) AddViewingKey(viewingKeyBytes []byte, signature []byte) error {
	return api.host.EnclaveClient.AddViewingKey(viewingKeyBytes, signature)
}

// GetRollupForTx returns the rollup containing a given transaction hash. Required for ObscuroScan.
func (api *ObscuroAPI) GetRollupForTx(txHash gethcommon.Hash) (*common.ExtRollup, error) {
	rollupNumber := api.host.nodeDB.GetRollupNumber(txHash)
	if rollupNumber == nil {
		return nil, fmt.Errorf("no rollup containing a transaction with hash %s is stored", txHash)
	}

	rollupHash := api.host.nodeDB.GetRollupHash(rollupNumber)
	if rollupHash == nil {
		return nil, fmt.Errorf("no rollup with number %d is stored", rollupNumber.Int64())
	}

	rollup, err := api.host.EnclaveClient.GetRollup(*rollupHash)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve rollup with hash %s. Cause: %w", rollupNumber, err)
	}

	return rollup, nil
}

// GetLatestTransactions returns the hashes of the latest `num` transactions, or as many as possible if less than `num` transactions exist.
func (api *ObscuroAPI) GetLatestTransactions(num int) ([]gethcommon.Hash, error) {
	rollupHeaderWithHashes := api.host.nodeDB.GetCurrentRollupHead()
	if rollupHeaderWithHashes == nil {
		return nil, nil
	}
	nextRollupHash := rollupHeaderWithHashes.Header.Hash()

	// We walk the chain until we've collected sufficient transactions.
	var txHashes []gethcommon.Hash
	for {
		rollup, err := api.host.EnclaveClient.GetRollup(nextRollupHash)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve rollup for hash. Cause: %w", err)
		}

		for _, txHash := range rollup.TxHashes {
			txHashes = append(txHashes, txHash)
			if len(txHashes) >= num {
				return txHashes, nil
			}
		}

		// If we have reached the top of the chain (i.e. the current rollup's number is one), we stop walking.
		if rollup.Header.Number.Cmp(big.NewInt(0)) == 0 {
			break
		}
		nextRollupHash = rollup.Header.ParentHash
	}

	return txHashes, nil
}

// StopHost gracefully stops the host.
func (api *ObscuroAPI) StopHost() {
	go api.host.Stop()
}
