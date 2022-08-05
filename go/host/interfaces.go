package host

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/obscuronet/go-obscuro/go/common"
)

// P2PCallback -the glue between the P2p layer and the node. Notifies the node when rollups and transactions are received from peers
type P2PCallback interface {
	ReceiveRollup(r common.EncodedRollup)
	ReceiveTx(tx common.EncryptedTx)
}

// P2P is the layer responsible for sending and receiving messages to Obscuro network peers.
type P2P interface {
	StartListening(callback P2PCallback)
	StopListening() error
	UpdatePeerList([]string)
	BroadcastRollup(r common.EncodedRollup) error
	BroadcastTx(tx common.EncryptedTx) error
}

type StatsCollector interface {
	// L2Recalc - called when a node has to discard the speculative work built on top of the winner of the gossip round.
	L2Recalc(id gethcommon.Address)
	NewBlock(block *types.Block)
	NewRollup(node gethcommon.Address)
	RollupWithMoreRecentProof()
}
