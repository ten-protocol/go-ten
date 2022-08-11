package clientapi

import (
	"github.com/obscuronet/go-obscuro/go/host"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
)

// TODO - Some methods return nil for an unfound block/rollup, while others return an error. Harmonise.

// ObscuroAPI implements Obscuro-specific JSON RPC operations.
type ObscuroAPI struct {
	host host.Host
}

func NewObscuroAPI(host host.Host) *ObscuroAPI {
	return &ObscuroAPI{
		host: host,
	}
}

// GetID returns the ID of the host.
func (api *ObscuroAPI) GetID() gethcommon.Address {
	return api.host.Config().ID
}

// GetCurrentBlockHead returns the current head block's header.
func (api *ObscuroAPI) GetCurrentBlockHead() *types.Header {
	return api.host.DB().GetCurrentBlockHead()
}

// GetRollupHeader returns the header of the rollup with the given hash.
func (api *ObscuroAPI) GetRollupHeader(hash gethcommon.Hash) *common.Header {
	headerWithHashes := api.host.DB().GetRollupHeader(hash)
	if headerWithHashes == nil {
		return nil
	}
	return headerWithHashes.Header
}

// AddViewingKey stores the viewing key on the enclave.
func (api *ObscuroAPI) AddViewingKey(viewingKeyBytes []byte, signature []byte) error {
	return api.host.EnclaveClient().AddViewingKey(viewingKeyBytes, signature)
}

// StopHost gracefully stops the host.
// TODO - Investigate how to authenticate this and other sensitive methods in production (Geth uses JWT).
func (api *ObscuroAPI) StopHost() {
	go api.host.Stop()
}
