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
	err := fmt.Errorf("no rollup with number %d is stored", number.Int64())

	// TODO - Provide a more efficient method on node DB to retrieve a rollup by number.
	// We walk the chain back to the requested rollup.
	rollupHeader := api.host.nodeDB.GetCurrentRollupHead()
	for {
		cmp := rollupHeader.Number.Cmp(number)
		// We have found the rollup we're looking for.
		if cmp == 0 {
			return rollupHeader, nil
		}
		// The current rollup has a lower number than we are looking for. We stop walking the chain and return an error.
		if cmp == -1 {
			return nil, err
		}

		rollupHeader = api.host.nodeDB.GetRollupHeader(rollupHeader.ParentHash)
		if rollupHeader == nil {
			return nil, err
		}
	}
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

// StopHost gracefully stops the host.
func (api *ObscuroAPI) StopHost() {
	go api.host.Stop()
}
