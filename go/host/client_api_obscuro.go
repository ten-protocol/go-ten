package host

import (
	"fmt"
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
)

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

// GetBlockHeaderByNumber returns the header for the block with the given number.
func (api *ObscuroAPI) GetBlockHeaderByNumber(number *big.Int) (*types.Header, error) {
	err := fmt.Errorf("no block with number %d is stored", number.Int64())

	// TODO - Provide a more efficient method on node DB to retrieve a block by number.
	// We walk the chain back to the requested block.
	blockHeader := api.host.nodeDB.GetCurrentBlockHead()
	for {
		cmp := blockHeader.Number.Cmp(number)
		// We have found the block we're looking for.
		if cmp == 0 {
			return blockHeader, nil
		}
		// The current block has a lower number than we are looking for. We stop walking the chain and return an error.
		if cmp == -1 {
			return nil, err
		}

		blockHeader = api.host.nodeDB.GetBlockHeader(blockHeader.ParentHash)
		if blockHeader == nil {
			return nil, err
		}
	}
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

// StopHost gracefully stops the host.
func (api *ObscuroAPI) StopHost() {
	go api.host.Stop()
}
