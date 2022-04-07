package ethclient

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/obscuro-playground/go/obscurocommon"
)

// Client provides a shared interface to access eth node RPC endpoints
type Client interface {
	BroadcastTx(t obscurocommon.EncodedL1Tx)
	FetchBlock(hash common.Hash) (*types.Block, bool)
}
