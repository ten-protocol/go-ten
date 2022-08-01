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
	return api.host.nodeDB.GetCurrentRollupHead()
}

// GetRollupHeader returns the header of the rollup with the given hash.
func (api *ObscuroAPI) GetRollupHeader(hash gethcommon.Hash) *common.Header {
	return api.host.nodeDB.GetRollupHeader(hash)
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

	return rollupHeader, nil
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
	// TODO - Provide a more efficient method on node DB to retrieve a rollup by transaction hash.
	// We walk the chain back until we find the requested transaction hash.
	rollupHeader := api.host.nodeDB.GetCurrentRollupHead()
	if rollupHeader == nil {
		return nil, nil //nolint:nilnil
	}

	rollup, err := api.host.EnclaveClient.GetRollup(rollupHeader.Hash())
	if err != nil {
		return nil, fmt.Errorf("could not find rollup containing transaction. Cause: %w", err)
	}
	for {
		// We check whether the transaction is in the current rollup, and return it if so.
		for _, rollupTxHash := range rollup.TxHashes {
			if rollupTxHash == txHash {
				return rollup, nil
			}
		}

		// We get the next rollup in the chain, to be checked in turn.
		rollup, err = api.host.EnclaveClient.GetRollup(rollup.Header.ParentHash)
		if err != nil {
			return nil, fmt.Errorf("could not find rollup containing transaction. Cause: %w", err)
		}
	}
}

// GetLatestTransactions returns the hashes of the latest `num` transactions, or as many as possible if less than `num` transactions exist.
func (api *ObscuroAPI) GetLatestTransactions(num int) ([]gethcommon.Hash, error) {
	rollupHeader := api.host.nodeDB.GetCurrentRollupHead()
	if rollupHeader == nil {
		return nil, nil
	}
	nextRollupHash := rollupHeader.Hash()

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
